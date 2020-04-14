package typeinferer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/typedef"
)

func TestInfer(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		typ  *Type
		err  error
	}{{
		name: "infer constant int correctly",
		exp:  t.IConst(5),
		typ:  base(typedef.Int),
		err:  nil,
	}, {
		name: "infer constant bool correctly",
		exp:  t.BConst(false),
		typ:  base(typedef.Bool),
		err:  nil,
	}, {
		name: "infer assigments in blocks",
		exp:  t.Block(t.Id("a"), t.Assign("a", t.IConst(5))),
		typ:  base(typedef.Int),
		err:  nil,
	}, {
		name: "infer integer comparisons as boolean",
		exp:  t.Block(t.Fcall(">", t.IConst(1), t.IConst(2))),
		typ:  base(typedef.Bool),
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  t.Block(t.Fcall("if", t.BConst(true), t.IConst(1), t.IConst(2))),
		typ:  base(typedef.Int),
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  t.Block(t.Fcall("if", t.BConst(true), t.IConst(1), t.BConst(false))),
		typ:  nil,
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "infer recursive functions",
		exp: t.Block(
			t.Fcall("a", t.BConst(false)),
			t.Assign("a", t.Fdef(t.Block(
				t.Fcall("if", t.BConst(false), t.Id("x"), t.Fcall("a", t.BConst(true)))),
				"x",
			))),
		typ: base(typedef.Bool),
		err: nil,
	}, {
		name: "infer function calls",
		exp: t.Block(
			t.Fcall("a", t.IConst(1)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("+", t.IConst(1), t.Id("x"))), "x"))),
		typ: base(typedef.Int),
		err: nil,
	}, {
		name: "infer function parameters",
		exp: t.Block(
			t.Fcall("a", t.IConst(1), t.BConst(true)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("if", t.Id("b"), t.Id("x"), t.IConst(0))), "x", "b"))),
		typ: base(typedef.Int),
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: t.Block(
			t.Fcall("a", t.BConst(true), t.BConst(true)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("if", t.Id("b"), t.Id("x"), t.IConst(0))), "x", "b"))),
		typ: nil,
		err: errors.New("can not unify bool with int"),
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := Infer(tt.exp)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Infer() error = %v, want %v", err, tt.err)
			}
			if tt.exp.Type == nil {
				t.Errorf("Infer() wrong type = nil, want %v", tt.typ)
			} else if !reflect.DeepEqual(tt.exp.Type.(*Type), tt.typ) {
				t.Errorf("Infer() wrong type = %v, want %v", tt.exp.Type, tt.typ)
			}
		})
	}
}
