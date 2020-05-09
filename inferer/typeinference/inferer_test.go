package typeinference

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/test"
)

func TestInfer(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		typ  string
		err  error
	}{{
		name: "infer constant int correctly",
		exp:  IConst(5),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer constant bool correctly",
		exp:  BConst(false),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer assigments in blocks",
		exp:  Block(Id("a"), Assign("a", IConst(5))),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer integer comparisons as boolean",
		exp:  Block(Fcall(">", IConst(1), IConst(2))),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  Block(Fcall("if", BConst(true), IConst(1), IConst(2))),
		typ:  "int",
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  Block(Fcall("if", BConst(true), IConst(1), BConst(false))),
		typ:  "",
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "fail when adding booleans together",
		exp:  Fcall("+", BConst(true), BConst(false)),
		typ:  "",
		err:  errors.New("can not unify bool with V1[int|real]"),
	}, {
		name: "infer recursive functions",
		exp: Block(
			Fcall("a", BConst(false)),
			Assign("a", Fdef(Block(
				Fcall("if", BConst(false), Id("x"), Fcall("a", BConst(true)))),
				"x",
			))),
		typ: "bool",
		err: nil,
	}, {
		name: "infer function calls",
		exp: Block(
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Block(Fcall("+", IConst(1), Id("x"))), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer function parameters",
		exp: Block(
			Fcall("a", IConst(1), BConst(true)),
			Assign("a", Fdef(Block(Fcall("if", Id("b"), Id("x"), IConst(0))), "x", "b"))),
		typ: "int",
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: Block(
			Fcall("a", BConst(true), BConst(true)),
			Assign("a", Fdef(Block(Fcall("if", Id("b"), Id("x"), IConst(0))), "x", "b"))),
		typ: "",
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  Fdef(Block(Fcall("if", BConst(true), Id("x"), Id("x"))), "x"),
		typ:  "(V1)=>V1",
		err:  nil,
	}, {
		name: "fail on recursive values",
		exp:  Block(Id("a"), Assign("a", Id("b")), Assign("b", Id("a"))),
		typ:  "",
		err:  errors.New("recursive value: a -> b -> a"),
	}, {
		name: "unify one function multiple ways",
		exp: Block(
			Fcall("if", Fcall("a", BConst(true)), Fcall("a", IConst(1)), IConst(2)),
			Assign("a", Fdef(Block(Fcall("if", BConst(true), Id("x"), Id("x"))), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer parameters in block values",
		exp: Block(
			Fdef(Fcall("if", BConst(true), Block(Id("x")), IConst(2)), "x"),
		),
		typ: "(int)=>int",
		err: nil,
	}, {
		name: "infer functions as arguments",
		exp: Block(
			Fdef(Fcall("+", Fcall("x", BConst(true), IConst(2)), IConst(1)), "x"),
		),
		typ: "((bool,int)=>int)=>int",
		err: nil,
	}, {
		name: "fail to unify functions with wrong number of arguments",
		exp: Block(
			Fcall("a", Id("b")),
			Assign("a", Fdef(Fcall("x", IConst(2), IConst(2)), "x")),
			Assign("b", Fdef(Id("x"), "x")),
		),
		typ: "",
		err: errors.New("wrong number of function arguments: 3 != 2"),
	}, {
		name: "infer multiple function arguments",
		exp: Block(
			Fcall("+", Fcall("op", Id("a")), Fcall("op", Id("b"))),
			Assign("a", Fdef(Fcall("+", Id("x"), Id("y")), "x", "y")),
			Assign("b", Fdef(Fcall("-", Id("x"), Id("y")), "x", "y")),
			Assign("op", Fdef(Fcall("x", IConst(1), IConst(2)), "x")),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer functions as return values",
		exp: Block(
			Fcall("r", Fcall("sw", BConst(true))),
			Assign("a", Fdef(Fcall("+", Id("x"), Id("y")), "x", "y")),
			Assign("b", Fdef(Fcall("-", Id("x"), Id("y")), "x", "y")),
			Assign("sw", Fdef(Fcall("if", Id("x"), Id("a"), Id("b")), "x")),
			Assign("r", Fdef(Fcall("f", RConst(1.0), RConst(2.0)), "f")),
		),
		typ: "real",
		err: nil,
	}, {
		name: "infer return values based on Closure",
		exp: Block(
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Block(
				Fcall("b", BConst(true)),
				Assign("b", Fdef(Fcall("if", Id("b"), Id("x"), IConst(2)), "b")),
			), "x")),
		),
		typ: "int",
		err: nil,
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := Infer(tt.exp)
			if err != nil {
				if !reflect.DeepEqual(err, tt.err) {
					t.Errorf("Infer() error = %v, want %v", err, tt.err)
				}
			} else {
				if (!tt.exp.Type().IsDefined()) && tt.typ != "" {
					t.Errorf("Infer() wrong type = nil, want %v", tt.typ)
				} else if tt.exp.Type().IsDefined() && tt.exp.Type().Signature() != tt.typ {
					t.Errorf("Infer() wrong type = %v, want %v", tt.exp.Type().Signature(), tt.typ)
				}
			}
		})
	}
}
