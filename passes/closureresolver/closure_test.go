package closureresolver

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/callresolver"
	"github.com/jvmakine/shine/passes/optimisation"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
)

func TestResolveFunctionDef(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		want map[string]map[string]Type
	}{{
		name: "resolves empty Closure for function without closure",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
			},
			Fcall(Id("a"), BConst(true), BConst(true), BConst(false)),
		),
		want: map[string]map[string]Type{
			"a%%1%%(bool,bool,bool)=>bool": map[string]Type{},
		},
	}, {
		name: "resolve closure parameters for function referring to outer ids",
		exp: Block(
			Assgs{
				"a": Fdef(Block(
					Assgs{
						"b": Fdef(Block(
							Assgs{"c": Id("y")},
							Fcall(Op("if"), Id("c"), Id("x"), IConst(2)),
						)),
					},
					Fcall(Id("b")),
				), "x", "y"),
			},
			Fcall(Id("a"), IConst(1), BConst(true)),
		),
		want: map[string]map[string]Type{
			"a%%3%%(int,bool)=>int": map[string]Type{},
			"b%%2%%()=>int": map[string]Type{
				"y": types.BoolP,
				"x": types.IntP,
			},
		},
	}, {
		name: "not include static function references in the closure",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall(Id("b"), Id("x"), Id("s")), "x"),
				"b": Fdef(Fcall(Op("+"), Fcall(Id("f"), Id("y")), IConst(2)), "y", "f"),
				"s": Fdef(Fcall(Op("+"), Id("y"), IConst(3)), "y"),
			},
			Fcall(Id("a"), IConst(1)),
		),
		want: map[string]map[string]Type{
			"a%%1%%(int)=>int":            map[string]Type{},
			"b%%1%%(int,(int)=>int)=>int": map[string]Type{},
			"s%%1%%(int)=>int":            map[string]Type{},
		},
	}, {
		name: "resolves closures for sequential functions",
		exp: Block(
			Assgs{"a": Fdef(Fdef(Fdef(Fcall(Op("+"), Fcall(Op("+"), Id("x"), Id("y")), Id("z")), "z"), "y"), "x")},
			Fcall(Fcall(Fcall(Id("a"), IConst(1)), IConst(2)), IConst(3)),
		),
		want: map[string]map[string]Type{
			"a%%1%%(int)=>(int)=>(int)=>int": map[string]Type{},
			"<anon1>%%1%%(int)=>(int)=>int": map[string]Type{
				"x": types.IntP,
			},
			"<anon2>%%1%%(int)=>int": map[string]Type{
				"x": types.IntP,
				"y": types.IntP,
			},
		},
	},
	}
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
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Resolve() = %v, want %v", result, tt.want)
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
				for _, c := range *v.Def.Closure {
					res[k][c.Name] = c.Type
				}
			}
		}
	}
	return res
}
