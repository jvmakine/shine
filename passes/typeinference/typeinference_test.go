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
		exp:  Block(Assgs{}, Fcall(Op(">"), IConst(1), IConst(2))),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  Block(Assgs{}, Fcall(Op("if"), BConst(true), IConst(1), IConst(2))),
		typ:  "int",
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  Block(Assgs{}, Fcall(Op("if"), BConst(true), IConst(1), BConst(false))),
		typ:  "",
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "fail when adding booleans together",
		exp:  Fcall(Op("+"), BConst(true), BConst(false)),
		typ:  "",
		err:  errors.New("can not unify bool with V1[int|real]"),
	}, {
		name: "infer recursive functions",
		exp: Block(
			Assgs{"a": Fdef(Block(
				Assgs{},
				Fcall(Op("if"), BConst(false), Id("x"), Fcall(Id("a"), BConst(true)))),
				"x",
			)},
			Fcall(Id("a"), BConst(false)),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer deeply nested recursive functions",
		exp: Block(
			Assgs{"a": Fdef(Block(
				Assgs{"b": Fdef(Fcall(Id("a"), Id("y")), "y")},
				Fcall(Op("if"), BConst(false), Id("x"), Fcall(Id("b"), BConst(true)))),
				"x",
			)},
			Fcall(Id("a"), BConst(false)),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer function calls",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall(Op("+"), IConst(1), Id("x"))), "x")},
			Fcall(Id("a"), IConst(1)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer function parameters",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall(Op("if"), Id("b"), Id("x"), IConst(0))), "x", "b")},
			Fcall(Id("a"), IConst(1), BConst(true)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall(Op("if"), Id("b"), Id("x"), IConst(0))), "x", "b")},
			Fcall(Id("a"), BConst(true), BConst(true)),
		),
		typ: "",
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  Fdef(Block(Assgs{}, Fcall(Op("if"), BConst(true), Id("x"), Id("x"))), "x"),
		typ:  "(V1)=>V1",
		err:  nil,
	}, {
		name: "fail on recursive values",
		exp:  Block(Assgs{"a": Id("b"), "b": Id("a")}, Id("a")),
		typ:  "",
		err:  errors.New("recursive value: a -> b -> a"),
	}, {
		name: "work on non-recursive values",
		exp: Block(
			Assgs{
				"a": Fcall(Op("+"), Id("b"), Id("b")),
				"b": Fcall(Op("+"), Id("c"), Id("c")),
				"c": Fcall(Op("+"), IConst(1), IConst(2)),
			},
			Id("a"),
		),
		typ: "int",
		err: nil,
	}, {
		name: "unify one function multiple ways",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Fcall(Op("if"), BConst(true), Id("x"), Id("x"))), "x")},
			Fcall(Op("if"), Fcall(Id("a"), BConst(true)), Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer parameters in block values",
		exp: Block(
			Assgs{},
			Fdef(Fcall(Op("if"), BConst(true), Block(Assgs{}, Id("x")), IConst(2)), "x"),
		),
		typ: "(int)=>int",
		err: nil,
	}, {
		name: "infer functions as arguments",
		exp: Block(
			Assgs{},
			Fdef(Fcall(Op("+"), Fcall(Id("x"), BConst(true), IConst(2)), IConst(1)), "x"),
		),
		typ: "((bool,int)=>int)=>int",
		err: nil,
	}, {
		name: "fail to unify functions with wrong number of arguments",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall(Id("x"), IConst(2), IConst(2)), "x"),
				"b": Fdef(Id("x"), "x"),
			},
			Fcall(Id("a"), Id("b")),
		),
		typ: "",
		err: errors.New("can not unify (V1)=>V1 with (int,int)=>V1"),
	}, {
		name: "infer multiple function arguments",
		exp: Block(
			Assgs{
				"a":  Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
				"b":  Fdef(Fcall(Op("-"), Id("x"), Id("y")), "x", "y"),
				"op": Fdef(Fcall(Id("x"), IConst(1), IConst(2)), "x"),
			},
			Fcall(Op("+"), Fcall(Id("op"), Id("a")), Fcall(Id("op"), Id("b"))),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer functions as return values",
		exp: Block(
			Assgs{
				"a":  Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
				"b":  Fdef(Fcall(Op("-"), Id("x"), Id("y")), "x", "y"),
				"sw": Fdef(Fcall(Op("if"), Id("x"), Id("a"), Id("b")), "x"),
				"r":  Fdef(Fcall(Id("f"), RConst(1.0), RConst(2.0)), "f"),
			},
			Fcall(Id("r"), Fcall(Id("sw"), BConst(true))),
		),
		typ: "real",
		err: nil,
	}, {
		name: "infer return values based on Closure",
		exp: Block(
			Assgs{
				"a": Fdef(Block(
					Assgs{"b": Fdef(Fcall(Op("if"), Id("bo"), Id("x"), IConst(2)), "bo")},
					Fcall(Id("b"), BConst(true)),
				), "x"),
			},
			Fcall(Id("a"), IConst(1)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on type errors in unused code",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall(Op("if"), BConst(true), Id("x"), IConst(2)), "x"),
				"b": Fdef(Fcall(Op("if"), Id("bo"), RConst(1.0), IConst(2)), "bo"),
			},
			Fcall(Id("a"), IConst(3)),
		),
		typ: "",
		err: errors.New("can not unify int with real"),
	}, {
		name: "leave free variables to functions",
		exp: Block(
			Assgs{"a": Fdef(Fcall(Op("+"), Id("x"), Id("x")), "x")},
			Fcall(Op("if"), Fcall(Op("<"), Fcall(Id("a"), RConst(1.0)), RConst(2.0)), Fcall(Id("a"), IConst(1)), IConst(3)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer sequential function definitions",
		exp: Block(
			Assgs{"a": Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("y")), "y"), "x")},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
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
