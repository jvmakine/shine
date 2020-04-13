package types

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/test"
)

func TestInfer(tes *testing.T) {
	tests := []struct {
		name    string
		exp     *ast.Exp
		typ     *Type
		wantErr bool
	}{{
		name:    "infer constant int correctly",
		exp:     t.IConst(5),
		typ:     Int,
		wantErr: false,
	}, {
		name:    "infer constant bool correctly",
		exp:     t.BConst(false),
		typ:     Bool,
		wantErr: false,
	}, {
		name:    "infer assigments in blocks",
		exp:     t.Block(t.Id("a"), t.Assign("a", t.IConst(5))),
		typ:     Int,
		wantErr: false,
	}, {
		name:    "infer integer comparisons as boolean",
		exp:     t.Block(t.Fcall(">", t.IConst(1), t.IConst(2))),
		typ:     Bool,
		wantErr: false,
	}, {
		name:    "infer if expressions",
		exp:     t.Block(t.Fcall("if", t.BConst(true), t.IConst(1), t.IConst(2))),
		typ:     Int,
		wantErr: false,
	}, {
		name: "infer recursive functions",
		exp: t.Block(
			t.Fcall("a", t.BConst(false)),
			t.Assign("a", t.Fdef(t.Block(
				t.Fcall("if", t.BConst(false), t.Id("x"), t.Fcall("a", t.BConst(false)))),
				"x",
			))),
		typ:     Bool,
		wantErr: false,
	}, {
		name: "infer function calls",
		exp: t.Block(
			t.Fcall("a", t.IConst(1)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("+", t.IConst(1), t.Id("x"))), "x"))),
		typ:     Int,
		wantErr: false,
	}, {
		name: "infer function parameters",
		exp: t.Block(
			t.Fcall("a", t.IConst(1), t.BConst(true)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("if", t.Id("b"), t.Id("x"), t.IConst(0))), "x", "b"))),
		typ:     Int,
		wantErr: false,
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			if err := Infer(tt.exp); (err != nil) != tt.wantErr {
				t.Errorf("Infer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.exp.Type == nil {
				t.Errorf("Infer() wrong type = nil, want %v", tt.typ)
			} else if !reflect.DeepEqual(tt.exp.Type.(*Type), tt.typ) {
				t.Errorf("Infer() wrong type = %v, want %v", tt.exp.Type, tt.typ)
			}
		})
	}
}
