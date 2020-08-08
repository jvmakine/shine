package closureresolver

import (
	"testing"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/callresolver"
	"github.com/jvmakine/shine/passes/optimisation"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
)

func TestResolveFunctionDef(tes *testing.T) {
	tests := []struct {
		name string
		exp  Expression
		want map[string]map[string]Type
	}{{
		name: "resolves empty closure for function without closure",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(true), NewConst(true), NewConst(false))).
			WithAssignment("a", NewFDef(NewBranch(NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")).
			WithID(1),
		want: map[string]map[string]Type{
			"a%%1%%(bool,bool,bool)=>bool": map[string]Type{},
		},
	}, {
		name: "resolve closure parameters for function referring to outer ids",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1), NewConst(true))).WithID(3).
			WithAssignment("a", NewFDef(NewBlock(NewFCall(NewId("b"))).WithID(2).
				WithAssignment("b", NewFDef(NewBlock(NewBranch(NewId("c"), NewId("x"), NewConst(2))).WithID(1).
					WithAssignment("c", NewId("y")),
				)),
				"x", "y")),
		want: map[string]map[string]Type{
			"a%%3%%(int,bool)=>int": map[string]Type{},
			"b%%2%%()=>int": map[string]Type{
				"y": Bool,
				"x": Int,
			},
		},
	}, {
		name: "not include static function references in the closure",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("a", NewFDef(NewFCall(NewId("b"), NewId("x"), NewId("s")), "x")).
			WithAssignment("b", NewFDef(NewOp("+", NewFCall(NewId("f"), NewId("y")), NewConst(2)), "y", "f")).
			WithAssignment("s", NewFDef(NewOp("+", NewId("y"), NewConst(3)), "y")).
			WithID(1),
		want: map[string]map[string]Type{
			"a%%1%%(int)=>int":            map[string]Type{},
			"b%%1%%(int,(int)=>int)=>int": map[string]Type{},
			"s%%1%%(int)=>int":            map[string]Type{},
		},
	}, {
		name: "resolves closures for sequential functions",
		exp: NewBlock(NewFCall(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2)), NewConst(3))).
			WithAssignment("a", NewFDef(NewFDef(NewFDef(NewOp("+", NewOp("+", NewId("x"), NewId("y")), NewId("z")), "z"), "y"), "x")).
			WithID(1),
		want: map[string]map[string]Type{
			"a%%1%%(int)=>(int)=>(int)=>int": map[string]Type{},
			"<anon1>%%1%%(int)=>(int)=>int": map[string]Type{
				"x": Int,
			},
			"<anon2>%%1%%(int)=>int": map[string]Type{
				"x": Int,
				"y": Int,
			},
		},
	}, /*{
		name: "resolves structures as closures",
		exp: NewBlock(NewFCall(NewId("f"), NewConst(2))).
			WithAssignment("S", NewStruct(ast.StructField{"x", Int})).
			WithAssignment("a", NewFCall(NewId("S"), NewConst(1))).
			WithAssignment("f", NewFDef(NewOp("+", NewId("y"), NewFieldAccessor("x", NewId("a"))), "y")).
			WithID(1),
		want: map[string]map[string]Type{
			"f%%1%%(int)=>int": map[string]Type{
				"a": types.NewNamed("S", types.NewStructure(types.NewNamed("x", Int))),
			},
		},
	}*/}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := typeinference.Infer(tt.exp)
			if err != nil {
				panic(err)
			}
			callresolver.ResolveFunctions(tt.exp)
			optimisation.DeadCodeElimination(tt.exp)
			CollectClosures(tt.exp)
			fcat := callresolver.Collect(tt.exp)
			result := collectClosures(&fcat)
			ok, err := deepdiff.DeepDiff(result, tt.want)
			if !ok {
				t.Error(err)
			}
		})
	}
}

func collectClosures(cat *callresolver.FCat) map[string]map[string]Type {
	res := map[string]map[string]Type{}
	for k, v := range *cat {
		if v.Def != nil {
			if v.Def.Closure != nil {
				res[k] = map[string]Type{}
				for _, c := range v.Def.Closure.Fields {
					r := c.Type
					if co, isC := r.(Contextual); isC {
						r = co.WithContext(nil).(Type)
					}
					res[k][c.Name] = r
				}
			}
		}
	}
	return res
}
