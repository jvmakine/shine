package inferer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/test"
)

func TestInfer(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		typ  string
		err  error
	}{{
		name: "infer constant int correctly",
		exp:  t.IConst(5),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer constant bool correctly",
		exp:  t.BConst(false),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer assigments in blocks",
		exp:  t.Block(t.Id("a"), t.Assign("a", t.IConst(5))),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer integer comparisons as boolean",
		exp:  t.Block(t.Fcall(">", t.IConst(1), t.IConst(2))),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  t.Block(t.Fcall("if", t.BConst(true), t.IConst(1), t.IConst(2))),
		typ:  "int",
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  t.Block(t.Fcall("if", t.BConst(true), t.IConst(1), t.BConst(false))),
		typ:  "",
		err:  errors.New("can not unify int with bool"),
	}, {
		name: "infer recursive functions",
		exp: t.Block(
			t.Fcall("a", t.BConst(false)),
			t.Assign("a", t.Fdef(t.Block(
				t.Fcall("if", t.BConst(false), t.Id("x"), t.Fcall("a", t.BConst(true)))),
				"x",
			))),
		typ: "bool",
		err: nil,
	}, {
		name: "infer function calls",
		exp: t.Block(
			t.Fcall("a", t.IConst(1)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("+", t.IConst(1), t.Id("x"))), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer function parameters",
		exp: t.Block(
			t.Fcall("a", t.IConst(1), t.BConst(true)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("if", t.Id("b"), t.Id("x"), t.IConst(0))), "x", "b"))),
		typ: "int",
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: t.Block(
			t.Fcall("a", t.BConst(true), t.BConst(true)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("if", t.Id("b"), t.Id("x"), t.IConst(0))), "x", "b"))),
		typ: "",
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  t.Fdef(t.Block(t.Fcall("if", t.BConst(true), t.Id("x"), t.Id("x"))), "x"),
		typ:  "(V1,V1)",
		err:  nil,
	}, {
		name: "fail on recursive values",
		exp:  t.Block(t.Id("a"), t.Assign("a", t.Id("b")), t.Assign("b", t.Id("a"))),
		typ:  "",
		err:  errors.New("recursive value: b -> a -> b"),
	}, {
		name: "unify one function multiple ways",
		exp: t.Block(
			t.Fcall("if", t.Fcall("a", t.BConst(true)), t.Fcall("a", t.IConst(1)), t.IConst(2)),
			t.Assign("a", t.Fdef(t.Block(t.Fcall("if", t.BConst(true), t.Id("x"), t.Id("x"))), "x"))),
		typ: "int",
		err: nil,
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := Infer(tt.exp)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Infer() error = %v, want %v", err, tt.err)
			}
			if (!tt.exp.Type.IsDefined()) && tt.typ != "" {
				t.Errorf("Infer() wrong type = nil, want %v", tt.typ)
			} else if tt.exp.Type.IsDefined() && tt.exp.Type.Signature() != tt.typ {
				t.Errorf("Infer() wrong type = %v, want %v", tt.exp.Type.Signature(), tt.typ)
			}
		})
	}
}
