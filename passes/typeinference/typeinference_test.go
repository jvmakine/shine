package typeinference

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

func TestInfer(tes *testing.T) {
	tests := []struct {
		name string
		exp  Expression
		typ  string
		err  error
	}{{
		name: "infer constant int correctly",
		exp:  NewConst(5),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer constant bool correctly",
		exp:  NewConst(false),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer assigments in blocks",
		exp:  NewBlock(NewId("a")).WithAssignment("a", NewConst(5)),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer integer comparisons as boolean",
		exp:  NewBlock(NewFCall(NewOp(">"), NewConst(1), NewConst(2))),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  NewBlock(NewFCall(NewOp("if"), NewConst(true), NewConst(1), NewConst(2))),
		typ:  "int",
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  NewBlock(NewFCall(NewOp("if"), NewConst(true), NewConst(1), NewConst(false))),
		typ:  "",
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "fail when adding booleans together",
		exp:  NewFCall(NewOp("+"), NewConst(true), NewConst(false)),
		typ:  "",
		err:  errors.New("can not unify V1[int|real|string] with bool"),
	}, {
		name: "infer recursive functions",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(false))).WithAssignment(
			"a", NewFDef(NewBlock(
				NewFCall(NewOp("if"), NewConst(false), NewId("x"), NewFCall(NewId("a"), NewConst(true)))),
				"x",
			),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer deeply nested recursive functions",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(false))).WithAssignment(
			"a", NewFDef(NewBlock(
				NewFCall(NewOp("if"), NewConst(false), NewId("x"), NewFCall(NewId("b"), NewConst(true))),
			).WithAssignment(
				"b", NewFDef(NewFCall(NewId("a"), NewId("y")), "y"),
			), "x",
			),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer function calls",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1))).WithAssignment(
			"a", NewFDef(NewBlock(NewFCall(NewOp("+"), NewConst(1), NewId("x"))), "x"),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer function parameters",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1), NewConst(true))).WithAssignment(
			"a", NewFDef(NewBlock(NewFCall(NewOp("if"), NewId("b"), NewId("x"), NewConst(0))), "x", "b"),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(true), NewConst(true))).WithAssignment(
			"a", NewFDef(NewBlock(NewFCall(NewOp("if"), NewId("b"), NewId("x"), NewConst(0))), "x", "b"),
		),
		typ: "",
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  NewFDef(NewBlock(NewFCall(NewOp("if"), NewConst(true), NewId("x"), NewId("x"))), "x"),
		typ:  "(V1)=>V1",
		err:  nil,
	}, {
		name: "fail on recursive values",
		exp: NewBlock(NewId("a")).
			WithAssignment("a", NewId("b")).
			WithAssignment("b", NewId("a")),
		typ: "",
		err: errors.New("recursive value: a -> b -> a"),
	}, {
		name: "work on non-recursive values",
		exp: NewBlock(NewId("a")).
			WithAssignment("a", NewFCall(NewOp("+"), NewId("b"), NewId("c"))).
			WithAssignment("b", NewFCall(NewOp("+"), NewId("c"), NewId("c"))).
			WithAssignment("c", NewFCall(NewOp("+"), NewConst(1), NewConst(2))),
		typ: "int",
		err: nil,
	}, {
		name: "unify one function multiple ways",
		exp: NewBlock(NewFCall(NewOp("if"), NewFCall(NewId("a"), NewConst(true)), NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewBlock(NewFCall(NewOp("if"), NewConst(true), NewId("x"), NewId("x"))), "x")),
		typ: "int",
		err: nil,
	}, {
		name: "infer parameters in block values",
		exp: NewBlock(
			NewFDef(NewFCall(NewOp("if"), NewConst(true), NewBlock(NewId("x")), NewConst(2)), "x"),
		),
		typ: "(int)=>int",
		err: nil,
	}, {
		name: "infer functions as arguments",
		exp: NewBlock(
			NewFDef(NewFCall(NewOp("+"), NewFCall(NewId("x"), NewConst(true), NewConst(2)), NewConst(1)), "x"),
		),
		typ: "((bool,int)=>int)=>int",
		err: nil,
	}, {
		name: "fail to unify functions with wrong number of arguments",
		exp: NewBlock(NewFCall(NewId("a"), NewId("b"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("x"), NewConst(2), NewConst(2)), "x")).
			WithAssignment("b", NewFDef(NewId("x"), "x")),
		typ: "",
		err: errors.New("can not unify (V1)=>V1 with (int,int)=>V1"),
	}, {
		name: "infer multiple function arguments",
		exp: NewBlock(NewFCall(NewOp("+"), NewFCall(NewId("op"), NewId("a")), NewFCall(NewId("op"), NewId("b")))).
			WithAssignment("a", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("b", NewFDef(NewFCall(NewOp("-"), NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("op", NewFDef(NewFCall(NewId("x"), NewConst(1), NewConst(2)), "x")),
		typ: "int",
		err: nil,
	}, {
		name: "infer functions as return values",
		exp: NewBlock(NewFCall(NewId("r"), NewFCall(NewId("sw"), NewConst(true)))).
			WithAssignment("a", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("b", NewFDef(NewFCall(NewOp("-"), NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("sw", NewFDef(NewFCall(NewOp("if"), NewId("x"), NewId("a"), NewId("b")), "x")).
			WithAssignment("r", NewFDef(NewFCall(NewId("f"), NewConst(1.0), NewConst(2.0)), "f")),
		typ: "real",
		err: nil,
	}, {
		name: "infer return values based on Closure",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("a", NewFDef(NewBlock(NewFCall(NewId("b"), NewConst(true))).
				WithAssignment("b", NewFDef(NewFCall(NewOp("if"), NewId("bo"), NewId("x"), NewConst(2)), "bo")),
				"x")),
		typ: "int",
		err: nil,
	}, {
		name: "fail on type errors in unused code",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(3))).
			WithAssignment("a", NewFDef(NewFCall(NewOp("if"), NewConst(true), NewId("x"), NewConst(2)), "x")).
			WithAssignment("b", NewFDef(NewFCall(NewOp("if"), NewId("bo"), NewConst(1.0), NewConst(2)), "bo")),
		typ: "",
		err: errors.New("can not unify int with real"),
	}, {
		name: "leave free variables to functions",
		exp: NewBlock(NewFCall(NewOp("if"), NewFCall(NewOp("<"), NewFCall(NewId("a"), NewConst(1.0)), NewConst(2.0)), NewFCall(NewId("a"), NewConst(1)), NewConst(3))).
			WithAssignment("a", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("x")), "x")),
		typ: "int",
		err: nil,
	}, {
		name: "infer sequential function definitions",
		exp: NewBlock(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "y"), "x")),
		typ: "int",
		err: nil,
	}, {
		name: "fail on parameter redefinitions",
		exp: NewBlock(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("x")), "x"), "x")),
		typ: "",
		err: errors.New("redefinition of x"),
	}, {
		name: "fails when function return type contradicts explicit type",
		exp: NewBlock(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewTypeDecl(types.BoolP, NewFCall(NewOp("+"), NewId("x"), NewId("x"))), "x")),
		typ: "",
		err: errors.New("can not unify V1[int|real|string] with bool"),
	}, {
		name: "fail to unify two different named types",
		exp: NewBlock(NewFCall(NewOp("if"), NewConst(true), NewId("ai"), NewId("bi"))).
			WithAssignment("a", NewStruct(StructField{"a1", types.IntP})).
			WithAssignment("b", NewStruct(StructField{"a1", types.IntP})).
			WithAssignment("ai", NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("bi", NewFCall(NewId("b"), NewConst(1))),
		typ: "",
		err: errors.New("can not unify a{a1:int} with b{a1:int}"),
	}, {
		name: "fail on invalid field access",
		exp:  NewBlock(NewFieldAccessor("xx", NewConst(1))),
		typ:  "",
		err:  errors.New("can not unify V1{xx:V2} with int"),
	}, {
		name: "fail on unknown named type",
		exp:  NewBlock(NewTypeDecl(types.MakeNamed("t"), NewId("a"))),
		typ:  "",
		err:  errors.New("type t is undefined"),
	}, {
		name: "work on known named type",
		exp: NewBlock(NewTypeDecl(types.MakeNamed("t"), NewFCall(NewId("t"), NewConst(1)))).
			WithAssignment("t", NewStruct(StructField{"x", types.IntP})),
		typ: "t{x:int}",
		err: nil,
	}, {
		name: "unify recursive types",
		exp: NewBlock(NewFCall(NewId("a"), NewFCall(NewId("a"), NewConst(0)))).
			WithAssignment("a", NewStruct(StructField{"a1", types.Type{}})),
		typ: "a{a1:a}",
		err: nil,
	}, {
		name: "infer function types from structure fields",
		exp:  NewBlock(NewFDef(NewFCall(NewOp("+"), NewFieldAccessor("a", NewId("x")), NewConst(1)), "x")),
		typ:  "(V1{a:int})=>int",
		err:  nil,
	}, {
		name: "infer typed interface function calls",
		exp: NewBlock(NewFCall(NewFieldAccessor("add", NewConst(1)), NewConst(2))).
			WithInterface(types.IntP, NewDefinitions(0).WithAssignment("add", NewFDef(NewFCall(NewOp("+"), NewId("$"), NewId("x")), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer untyped interface function calls",
		exp: NewBlock(NewFCall(NewFieldAccessor("add", NewConst(1)), NewConst(2))).
			WithInterface(types.Type{}, NewDefinitions(0).WithAssignment("add", NewFDef(NewFCall(NewOp("+"), NewId("$"), NewId("x")), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer multiple interface invocations",
		exp: NewBlock(NewFCall(NewFieldAccessor("add", NewFCall(NewFieldAccessor("add", NewConst(1)), NewConst(2))), NewConst(3))).
			WithInterface(types.Type{}, NewDefinitions(0).WithAssignment("add", NewFDef(NewFCall(NewOp("+"), NewId("$"), NewId("x")), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer interface usage in functions",
		exp: NewBlock(NewFCall(NewId("f"), NewConst(2))).
			WithInterface(types.Type{}, NewDefinitions(0).WithAssignment("isOdd", NewFDef(NewFCall(NewOp("=="), NewConst(0), NewFCall(NewOp("%"), NewId("$"), NewConst(2)))))).
			WithAssignment("f", NewFDef(NewFCall(NewFieldAccessor("isOdd", NewId("x"))), &FParam{"x", types.IntP})),
		typ: "bool",
	}, {
		name: "fail on invalid interface call",
		exp: NewBlock(NewFCall(NewFCall(NewFieldAccessor("add", NewConst("a")), NewConst("b")))).
			WithInterface(types.Type{}, NewDefinitions(0).WithAssignment("add", NewFDef(NewFCall(NewOp("-"), NewId("$"), NewId("x")), "x"))),
		err: errors.New("can not unify V1[int|real] with string"),
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
