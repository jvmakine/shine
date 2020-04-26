package inferer

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/test"
)

func TestResolveSignatureGeneration(tes *testing.T) {
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
		want: []string{"a%%1%%bool", "a%%1%%int"},
	}}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := Infer(tt.exp)
			if err != nil {
				panic(err)
			}
			Resolve(tt.exp)
			result := collectResolved(tt.exp)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Resolve() = %v, want %v", result, tt.want)
			}
		})
	}
}

func collectResolved(exp *ast.Exp) []string {
	res := []string{}
	if exp.Resolved != "" {
		res = []string{exp.Resolved}
	}
	if exp.Block != nil {
		for _, a := range exp.Block.Assignments {
			res = append(res, collectResolved(a.Value)...)
		}
		res = append(res, collectResolved(exp.Block.Value)...)
	}
	if exp.Call != nil {
		for _, p := range exp.Call.Params {
			res = append(res, collectResolved(p)...)
		}
	}
	if exp.Def != nil {
		res = append(res, collectResolved(exp.Def.Body)...)
	}
	return res
}
