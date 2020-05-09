package callresolver

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/test"
)

func TestResolveFunctionCall(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		want []string
	}{{
		name: "resolves function signatures based on the call type",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall("if", Id("b"), Id("y"), Id("x")), "b", "y", "x"),
			},
			Fcall("if",
				Fcall("a", BConst(true), BConst(true), BConst(false)),
				Fcall("a", BConst(true), IConst(5), IConst(6)),
				IConst(7)),
		),
		want: []string{"a%%1%%(bool,bool,bool)=>bool", "a%%1%%(bool,int,int)=>int"},
	}, {
		name: "resolves functions as arguments",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall("f", IConst(1), IConst(2)), "f"),
				"b": Fdef(Fcall("+", Id("x"), Id("y")), "x", "y"),
			},
			Fcall("a", Id("b")),
		),
		want: []string{"b%%1%%(int,int)=>int", "a%%1%%((int,int)=>int)=>int", "b%%1%%(int,int)=>int"},
	}, {
		name: "resolves anonymous functions",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall("f", IConst(1), IConst(2)), "f"),
			},
			Fcall("a", Fdef(Fcall("+", Id("x"), Id("y")), "x", "y")),
		),
		want: []string{"a%%1%%((int,int)=>int)=>int", "<anon1>%%1%%(int,int)=>int"},
	}}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := typeinference.Infer(tt.exp)
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

func collectResolvedCalls(exp *ast.Exp) []string {
	res := []string{}
	if exp.Resolved != nil {
		res = []string{exp.Resolved.ID}
	}
	if exp.Block != nil {
		for _, a := range exp.Block.Assignments {
			res = append(res, collectResolvedCalls(a)...)
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
