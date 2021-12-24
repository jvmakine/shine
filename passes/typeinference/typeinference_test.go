package typeinference

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
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
		exp:  Block(Assgs{"a": IConst(5)}, Typedefs{}, Id("a")),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer integer comparisons as boolean",
		exp:  Block(Assgs{}, Typedefs{}, Fcall(Op(">"), IConst(1), IConst(2))),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  Block(Assgs{}, Typedefs{}, Fcall(Op("if"), BConst(true), IConst(1), IConst(2))),
		typ:  "int",
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  Block(Assgs{}, Typedefs{}, Fcall(Op("if"), BConst(true), IConst(1), BConst(false))),
		typ:  "",
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "fail when adding booleans together",
		exp:  Fcall(Op("+"), BConst(true), BConst(false)),
		typ:  "",
		err:  errors.New("can not unify bool with V1[int|real|string]"),
	}, {
		name: "infer recursive functions",
		exp: Block(
			Assgs{"a": Fdef(Block(
				Assgs{},
				Typedefs{},
				Fcall(Op("if"), BConst(false), Id("x"), Fcall(Id("a"), BConst(true)))),
				"x",
			)},
			Typedefs{},
			Fcall(Id("a"), BConst(false)),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer deeply nested recursive functions",
		exp: Block(
			Assgs{"a": Fdef(Block(
				Assgs{"b": Fdef(Fcall(Id("a"), Id("y")), "y")},
				Typedefs{},
				Fcall(Op("if"), BConst(false), Id("x"), Fcall(Id("b"), BConst(true)))),
				"x",
			)},
			Typedefs{},
			Fcall(Id("a"), BConst(false)),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer function calls",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Fcall(Op("+"), IConst(1), Id("x"))), "x")},
			Typedefs{},
			Fcall(Id("a"), IConst(1)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer function parameters",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Fcall(Op("if"), Id("b"), Id("x"), IConst(0))), "x", "b")},
			Typedefs{},
			Fcall(Id("a"), IConst(1), BConst(true)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Fcall(Op("if"), Id("b"), Id("x"), IConst(0))), "x", "b")},
			Typedefs{},
			Fcall(Id("a"), BConst(true), BConst(true)),
		),
		typ: "",
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  Fdef(Block(Assgs{}, Typedefs{}, Fcall(Op("if"), BConst(true), Id("x"), Id("x"))), "x"),
		typ:  "(V1)=>V1",
		err:  nil,
	}, {
		name: "fail on recursive values",
		exp:  Block(Assgs{"a": Id("b"), "b": Id("a")}, Typedefs{}, Id("a")),
		typ:  "",
		err:  errors.New("recursive value: a -> b -> a"),
	}, {
		name: "work on non-recursive values",
		exp: Block(
			Assgs{
				"a": Fcall(Op("+"), Id("b"), Id("c")),
				"b": Fcall(Op("+"), Id("c"), Id("c")),
				"c": Fcall(Op("+"), IConst(1), IConst(2)),
			},
			Typedefs{},
			Id("a"),
		),
		typ: "int",
		err: nil,
	}, {
		name: "unify one function multiple ways",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Fcall(Op("if"), BConst(true), Id("x"), Id("x"))), "x")},
			Typedefs{},
			Fcall(Op("if"), Fcall(Id("a"), BConst(true)), Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer parameters in block values",
		exp: Block(
			Assgs{},
			Typedefs{},
			Fdef(Fcall(Op("if"), BConst(true), Block(Assgs{}, Typedefs{}, Id("x")), IConst(2)), "x"),
		),
		typ: "(int)=>int",
		err: nil,
	}, {
		name: "infer functions as arguments",
		exp: Block(
			Assgs{},
			Typedefs{},
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
			Typedefs{},
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
			Typedefs{},
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
			Typedefs{},
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
					Typedefs{},
					Fcall(Id("b"), BConst(true)),
				), "x"),
			},
			Typedefs{},
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
			Typedefs{},
			Fcall(Id("a"), IConst(3)),
		),
		typ: "",
		err: errors.New("can not unify int with real"),
	}, {
		name: "leave free variables to functions",
		exp: Block(
			Assgs{"a": Fdef(Fcall(Op("+"), Id("x"), Id("x")), "x")},
			Typedefs{},
			Fcall(Op("if"), Fcall(Op("<"), Fcall(Id("a"), RConst(1.0)), RConst(2.0)), Fcall(Id("a"), IConst(1)), IConst(3)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer sequential function definitions",
		exp: Block(
			Assgs{"a": Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("y")), "y"), "x")},
			Typedefs{},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on parameter redefinitions",
		exp: Block(
			Assgs{"a": Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("x")), "x"), "x")},
			Typedefs{},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "",
		err: errors.New("redefinition of x"),
	}, {
		name: "fails when function return type contradicts explicit type",
		exp: Block(
			Assgs{"a": Fdef(TDecl(Fcall(Op("+"), Id("x"), Id("x")), types.BoolP), "x")},
			Typedefs{},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "",
		err: errors.New("can not unify bool with V1[int|real|string]"),
	}, {
		name: "fail to unify two different named types",
		exp: Block(
			Assgs{
				"ai": Fcall(Id("a"), IConst(1)),
				"bi": Fcall(Id("b"), IConst(1)),
			},
			Typedefs{
				"a": Struct(ast.StructField{"a1", types.IntP}),
				"b": Struct(ast.StructField{"a1", types.IntP}),
			},
			Fcall(Op("if"), BConst(true), Id("ai"), Id("bi")),
		),
		err: errors.New("can not unify a{a1:int} with b{a1:int}"),
	}, {
		name: "fail on unknown named type",
		exp:  Block(Assgs{}, Typedefs{}, TDecl(Id("a"), types.MakeNamed("t"))),
		err:  errors.New("type t is undefined"),
	}, {
		name: "fail on unknown named type argument",
		exp: Block(
			Assgs{},
			Typedefs{"A": Struct(ast.StructField{Name: "x", Type: types.MakeNamed("X")}).WithFreeVars("X")},
			TDecl(Id("a"), types.MakeNamed("A", types.MakeNamed("Z"))),
		),
		err: errors.New("type Z is undefined"),
	}, {
		name: "work on known named type",
		exp: Block(
			Assgs{},
			Typedefs{"t": Struct(ast.StructField{"x", types.IntP})},
			TDecl(Id("a"), types.MakeNamed("t")),
		),
		typ: "t{x:int}",
	}, {
		name: "unify recursive types",
		exp: Block(
			Assgs{},
			Typedefs{"a": Struct(ast.StructField{"a1", types.Type{}})},
			Fcall(Id("a"), Fcall(Id("a"), IConst(0))),
		),
		typ: "a{a1:a}",
		err: nil,
	}, {
		name: "infer function types from structure fields",
		exp: Block(
			Assgs{},
			Typedefs{},
			Fdef(Fcall(Op("+"), Faccess(Id("x"), "a"), IConst(1)), "x"),
		),
		typ: "(V1{a:int})=>int",
		err: nil,
	}, {
		name: "fail to unify free type to two different types",
		exp: Block(
			Assgs{},
			Typedefs{"A": Struct(
				ast.StructField{"a1", types.MakeNamed("X")},
				ast.StructField{"a2", types.MakeNamed("X")},
			).WithFreeVars("X")},
			Fcall(Id("A"), IConst(1), RConst(2.0)),
		),
		err: errors.New("can not unify int with real"),
	}, {
		name: "fail on reusing defined type as free type",
		exp: Block(
			Assgs{},
			Typedefs{
				"X": Struct(ast.StructField{"a1", types.IntP}),
				"A": Struct(ast.StructField{"a1", types.MakeNamed("X")}).WithFreeVars("X"),
			},
			Fcall(Id("A"), IConst(1)),
		),
		err: errors.New("redefinition of X"),
	}, {
		name: "fail on unused free types",
		exp: Block(
			Assgs{},
			Typedefs{"A": Struct(ast.StructField{"a1", types.MakeNamed("X")}).WithFreeVars("X", "Y")},
			Fcall(Id("A"), IConst(1)),
		),
		err: errors.New("unused free type Y"),
	}, {
		name: "fail on incorrect type variable",
		exp: Block(
			Assgs{"f": Fdef(Faccess(Id("a"), "x"), &ast.FParam{Name: "a", Type: types.MakeNamed("A", types.RealP)})},
			Typedefs{"A": Struct(ast.StructField{"x", types.MakeNamed("X")}).WithFreeVars("X")},
			Fcall(Id("f"), Fcall(Id("A"), IConst(1))),
		),
		err: errors.New("can not unify int with real"),
	}, {
		name: "fail on redefinitions",
		exp: Block(
			Assgs{"a": IConst(1)}, Typedefs{},
			Block(
				Assgs{"a": IConst(1)}, Typedefs{},
				Id("a"),
			),
		),
		err: errors.New("redefinition of a"),
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
