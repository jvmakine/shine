package inferer

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/resolved"
	. "github.com/jvmakine/shine/resolved"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
)

func TestResolveFunctionCall(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		want []string
	}{{
		name: "resolves function signatures based on the call type",
		exp: Block(
			Fcall("if",
				Fcall("a", BConst(true), BConst(true), BConst(false)),
				Fcall("a", BConst(true), IConst(5), IConst(6)),
				IConst(7)),
			Assign("a", Fdef(Fcall("if", Id("b"), Id("y"), Id("x")), "b", "y", "x")),
		),
		want: []string{"a%%1%%(bool,bool,bool)=>bool", "a%%1%%(bool,int,int)=>int"},
	}, {
		name: "resolves functions as arguments",
		exp: Block(
			Fcall("a", Id("b")),
			Assign("a", Fdef(Fcall("f", IConst(1), IConst(2)), "f")),
			Assign("b", Fdef(Fcall("+", Id("x"), Id("y")), "x", "y")),
		),
		want: []string{"a%%1%%((int,int)=>int)=>int", "b%%1%%(int,int)=>int"},
	}, {
		name: "resolves anonymous functions",
		exp: Block(
			Fcall("a", Fdef(Fcall("+", Id("x"), Id("y")), "x", "y")),
			Assign("a", Fdef(Fcall("f", IConst(1), IConst(2)), "f")),
		),
		want: []string{"a%%1%%((int,int)=>int)=>int", "<anon1>%%1%%(int,int)=>int"},
	}}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := Infer(tt.exp)
			if err != nil {
				panic(err)
			}
			Resolve(tt.exp)
			result := collectResolvedCalls(tt.exp)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Resolve() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestResolveFunctionDef(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		want []resolved.Closure
	}{{
		name: "resolves empty Closure for function without closure",
		exp: Block(
			Fcall("a", BConst(true), BConst(true), BConst(false)),
			Assign("a", Fdef(Fcall("if", Id("b"), Id("y"), Id("x")), "b", "y", "x")),
		),
		want: []Closure{Closure{}},
	}, {
		name: "resolve closure parameters for function referring to outer ids",
		exp: Block(
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Block(
				Fcall("b", BConst(true)),
				Assign("b", Fdef(Fcall("if", Id("bo"), Id("x"), IConst(2)), "bo")),
			), "x")),
		),
		want: []Closure{Closure{}, Closure{ClosureParam{Name: "x", Type: types.IntP}}},
	}, {
		name: "not include static function references in the closure",
		exp: Block(
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Fcall("b", Id("x"), Id("s")), "x")),
			Assign("b", Fdef(Fcall("+", Fcall("f", Id("y")), IConst(2)), "y", "f")),
			Assign("s", Fdef(Fcall("+", Id("y"), IConst(3)), "y")),
		),
		want: []Closure{Closure{}, Closure{}, Closure{}},
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := Infer(tt.exp)
			if err != nil {
				panic(err)
			}
			fcat := Resolve(tt.exp)
			result := collectResolvedDefs(fcat)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Resolve() = %v, want %v", result, tt.want)
			}
		})
	}
}

func collectResolvedCalls(exp *ast.Exp) []string {
	res := []string{}
	if exp.Resolved != nil {
		res = []string{exp.Resolved.ID}
	}
	if exp.Block != nil {
		for _, a := range exp.Block.Assignments {
			res = append(res, collectResolvedCalls(a.Value)...)
		}
		res = append(res, collectResolvedCalls(exp.Block.Value)...)
	}
	if exp.Call != nil {
		for _, p := range exp.Call.Params {
			res = append(res, collectResolvedCalls(p)...)
		}
	}
	if exp.Def != nil {
		res = append(res, collectResolvedCalls(exp.Def.Body)...)
	}
	return res
}

func collectResolvedDefs(cat *FCat) []Closure {
	res := []Closure{}
	for _, v := range *cat {
		if v.Resolved != nil {
			res = append(res, v.Resolved.Closure)
		}
	}
	return res
}
