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
		exp:  Block(Assgs{"a": IConst(5)}, Id("a")),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer integer comparisons as boolean",
		exp:  Block(Assgs{}, Fcall(">", IConst(1), IConst(2))),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  Block(Assgs{}, Fcall("if", BConst(true), IConst(1), IConst(2))),
		typ:  "int",
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  Block(Assgs{}, Fcall("if", BConst(true), IConst(1), BConst(false))),
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
			Assgs{"a": Fdef(Block(
				Assgs{},
				Fcall("if", BConst(false), Id("x"), Fcall("a", BConst(true)))),
				"x",
			)},
			Fcall("a", BConst(false)),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer function calls",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall("+", IConst(1), Id("x"))), "x")},
			Fcall("a", IConst(1)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer function parameters",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall("if", Id("b"), Id("x"), IConst(0))), "x", "b")},
			Fcall("a", IConst(1), BConst(true)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall("if", Id("b"), Id("x"), IConst(0))), "x", "b")},
			Fcall("a", BConst(true), BConst(true)),
		),
		typ: "",
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  Fdef(Block(Assgs{}, Fcall("if", BConst(true), Id("x"), Id("x"))), "x"),
		typ:  "(V1)=>V1",
		err:  nil,
	}, {
		name: "fail on recursive values",
		exp:  Block(Assgs{"a": Id("b"), "b": Id("a")}, Id("a")),
		typ:  "",
		err:  errors.New("recursive value: a -> b -> a"),
	}, {
		name: "unify one function multiple ways",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall("if", BConst(true), Id("x"), Id("x"))), "x")},
			Fcall("if", Fcall("a", BConst(true)), Fcall("a", IConst(1)), IConst(2)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer parameters in block values",
		exp: Block(
			Assgs{},
			Fdef(Fcall("if", BConst(true), Block(Assgs{}, Id("x")), IConst(2)), "x"),
		),
		typ: "(int)=>int",
		err: nil,
	}, {
		name: "infer functions as arguments",
		exp: Block(
			Assgs{},
			Fdef(Fcall("+", Fcall("x", BConst(true), IConst(2)), IConst(1)), "x"),
		),
		typ: "((bool,int)=>int)=>int",
		err: nil,
	}, {
		name: "fail to unify functions with wrong number of arguments",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall("x", IConst(2), IConst(2)), "x"),
				"b": Fdef(Id("x"), "x"),
			},
			Fcall("a", Id("b")),
		),
		typ: "",
		err: errors.New("wrong number of function arguments: 3 != 2"),
	}, {
		name: "infer multiple function arguments",
		exp: Block(
			Assgs{
				"a":  Fdef(Fcall("+", Id("x"), Id("y")), "x", "y"),
				"b":  Fdef(Fcall("-", Id("x"), Id("y")), "x", "y"),
				"op": Fdef(Fcall("x", IConst(1), IConst(2)), "x"),
			},
			Fcall("+", Fcall("op", Id("a")), Fcall("op", Id("b"))),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer functions as return values",
		exp: Block(
			Assgs{
				"a":  Fdef(Fcall("+", Id("x"), Id("y")), "x", "y"),
				"b":  Fdef(Fcall("-", Id("x"), Id("y")), "x", "y"),
				"sw": Fdef(Fcall("if", Id("x"), Id("a"), Id("b")), "x"),
				"r":  Fdef(Fcall("f", RConst(1.0), RConst(2.0)), "f"),
			},
			Fcall("r", Fcall("sw", BConst(true))),
		),
		typ: "real",
		err: nil,
	}, {
		name: "infer return values based on Closure",
		exp: Block(
			Assgs{
				"a": Fdef(Block(
					Assgs{
						"b": Fdef(Fcall("if", Id("b"), Id("x"), IConst(2)), "b"),
					},
					Fcall("b", BConst(true)),
				), "x"),
			},
			Fcall("a", IConst(1)),
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
