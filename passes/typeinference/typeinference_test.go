package typeinference

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/grammar"
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
	}, {
		name: "infer constant bool correctly",
		exp:  NewConst(false),
		typ:  "bool",
	}, {
		name: "infer assigments in blocks",
		exp:  NewBlock(NewId("a")).WithAssignment("a", NewConst(5)),
		typ:  "int",
	}, {
		name: "infer integer comparisons as boolean",
		exp:  NewBlock(NewOp(">", NewConst(1), NewConst(2))),
		typ:  "bool",
	}, {
		name: "infer if expressions",
		exp:  NewBlock(NewBranch(NewConst(true), NewConst(1), NewConst(2))),
		typ:  "int",
	}, {
		name: "fail on mismatching if expression branches",
		exp:  NewBlock(NewBranch(NewConst(true), NewConst(1), NewConst(false))),
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "fail when adding booleans together",
		exp:  NewOp("+", NewConst(true), NewConst(false)),
		err:  errors.New("can not unify V1{+:(bool)=>V2} with bool"),
	}, {
		name: "infer recursive functions",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(false))).WithAssignment(
			"a", NewFDef(NewBlock(
				NewBranch(NewConst(false), NewId("x"), NewFCall(NewId("a"), NewConst(true)))),
				"x",
			),
		),
		typ: "bool",
	}, {
		name: "infer deeply nested recursive functions",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(false))).WithAssignment(
			"a", NewFDef(NewBlock(
				NewBranch(NewConst(false), NewId("x"), NewFCall(NewId("b"), NewConst(true))),
			).WithAssignment(
				"b", NewFDef(NewFCall(NewId("a"), NewId("y")), "y"),
			), "x",
			),
		),
		typ: "bool",
	}, {
		name: "infer function calls",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1))).WithAssignment(
			"a", NewFDef(NewBlock(NewOp("+", NewConst(1), NewId("x"))), "x"),
		),
		typ: "int",
	}, {
		name: "infer function parameters",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1), NewConst(true))).WithAssignment(
			"a", NewFDef(NewBlock(NewBranch(NewId("b"), NewId("x"), NewConst(0))), "x", "b"),
		),
		typ: "int",
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(true), NewConst(true))).WithAssignment(
			"a", NewFDef(NewBlock(NewBranch(NewId("b"), NewId("x"), NewConst(0))), "x", "b"),
		),
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  NewFDef(NewBlock(NewBranch(NewConst(true), NewId("x"), NewId("x"))), "x"),
		typ:  "(V1)=>V1",
	}, {
		name: "fail on recursive values",
		exp: NewBlock(NewId("a")).
			WithAssignment("a", NewId("b")).
			WithAssignment("b", NewId("a")),
		err: errors.New("recursive value: a -> b -> a"),
	}, {
		name: "work on non-recursive values",
		exp: NewBlock(NewId("a")).
			WithAssignment("a", NewOp("+", NewId("b"), NewId("c"))).
			WithAssignment("b", NewOp("+", NewId("c"), NewId("c"))).
			WithAssignment("c", NewOp("+", NewConst(1), NewConst(2))),
		typ: "int",
	}, {
		name: "unify one function multiple ways",
		exp: NewBlock(NewBranch(NewFCall(NewId("a"), NewConst(true)), NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewBlock(NewBranch(NewConst(true), NewId("x"), NewId("x"))), "x")),
		typ: "int",
	}, {
		name: "infer parameters in block values",
		exp: NewBlock(
			NewFDef(NewBranch(NewConst(true), NewBlock(NewId("x")), NewConst(2)), "x"),
		),
		typ: "(int)=>int",
	}, {
		name: "infer functions as arguments",
		exp: NewBlock(
			NewFDef(NewOp("+", NewConst(1), NewFCall(NewId("x"), NewConst(true), NewConst(2))), "x"),
		),
		typ: "((bool,int)=>int)=>int",
	}, {
		name: "fail to unify functions with wrong number of arguments",
		exp: NewBlock(NewFCall(NewId("a"), NewId("b"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("x"), NewConst(2), NewConst(2)), "x")).
			WithAssignment("b", NewFDef(NewId("x"), "x")),
		err: errors.New("can not unify (V1)=>V1 with (int,int)=>V1"),
	}, {
		name: "infer multiple function arguments",
		exp: NewBlock(NewOp("+", NewFCall(NewId("op"), NewId("a")), NewFCall(NewId("op"), NewId("b")))).
			WithAssignment("a", NewFDef(NewOp("+", NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("b", NewFDef(NewOp("-", NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("op", NewFDef(NewFCall(NewId("x"), NewConst(1), NewConst(2)), "x")),
		typ: "int",
	}, {
		name: "infer functions as return values",
		exp: NewBlock(NewFCall(NewId("r"), NewFCall(NewId("sw"), NewConst(true)))).
			WithAssignment("a", NewFDef(NewOp("+", NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("b", NewFDef(NewOp("-", NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("sw", NewFDef(NewBranch(NewId("x"), NewId("a"), NewId("b")), "x")).
			WithAssignment("r", NewFDef(NewFCall(NewId("f"), NewConst(1.0), NewConst(2.0)), "f")),
		typ: "real",
	}, {
		name: "infer return values based on Closure",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("a", NewFDef(NewBlock(NewFCall(NewId("b"), NewConst(true))).
				WithAssignment("b", NewFDef(NewBranch(NewId("bo"), NewId("x"), NewConst(2)), "bo")),
				"x")),
		typ: "int",
	}, {
		name: "fail on type errors in unused code",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(3))).
			WithAssignment("a", NewFDef(NewBranch(NewConst(true), NewId("x"), NewConst(2)), "x")).
			WithAssignment("b", NewFDef(NewBranch(NewId("bo"), NewConst(1.0), NewConst(2)), "bo")),
		err: errors.New("can not unify int with real"),
	}, {
		name: "leave free variables to functions",
		exp: NewBlock(NewBranch(NewOp("<", NewFCall(NewId("a"), NewConst(1.0)), NewConst(2.0)), NewFCall(NewId("a"), NewConst(1)), NewConst(3))).
			WithAssignment("a", NewFDef(NewOp("+", NewId("x"), NewId("x")), "x")),
		typ: "int",
	}, {
		name: "infer sequential function definitions",
		exp: NewBlock(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewFDef(NewOp("+", NewId("x"), NewId("y")), "y"), "x")),
		typ: "int",
	}, {
		name: "fail on parameter redefinitions",
		exp: NewBlock(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewFDef(NewOp("+", NewId("x"), NewId("x")), "x"), "x")),
		err: errors.New("redefinition of x"),
	}, {
		name: "fails when function return type contradicts explicit type",
		exp: NewBlock(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewTypeDecl(types.Bool, NewOp("+", NewId("x"), NewId("x"))), "x")),
		err: errors.New("can not unify V1{+:(V2)=>bool} with int"),
	}, {
		name: "fail to unify two different named types",
		exp: NewBlock(NewBranch(NewConst(true), NewId("ai"), NewId("bi"))).
			WithAssignment("a", NewStruct(StructField{"a1", types.Int})).
			WithAssignment("b", NewStruct(StructField{"a1", types.Int})).
			WithAssignment("ai", NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("bi", NewFCall(NewId("b"), NewConst(1))),
		err: errors.New("can not unify a with b"),
	}, {
		name: "fail on invalid field access",
		exp:  NewBlock(NewFieldAccessor("xx", NewConst(1))),
		err:  errors.New("can not unify V1{xx:V2} with int"),
	}, {
		name: "fail on unknown named type",
		exp:  NewBlock(NewTypeDecl(types.NewNamed("t", nil), NewId("a"))),
		err:  errors.New("type t is undefined"),
	}, {
		name: "work on known named type",
		exp: NewBlock(NewTypeDecl(types.NewNamed("t", nil), NewFCall(NewId("t"), NewConst(1)))).
			WithAssignment("t", NewStruct(StructField{"x", types.Int})),
		typ: "t",
	}, {
		name: "unify recursive types",
		exp: NewBlock(NewFCall(NewId("a"), NewFCall(NewId("a"), NewConst(0)))).
			WithAssignment("a", NewStruct(StructField{"a1", nil})),
		typ: "a",
	}, {
		name: "infer function types from structure fields",
		exp:  NewBlock(NewFDef(NewOp("+", NewFieldAccessor("a", NewId("x")), NewConst(1)), "x")),
		typ:  "(V1{a:V2{+:(int)=>V3}})=>V3",
	}, {
		name: "infer typed interface function calls",
		exp: NewBlock(NewFCall(NewFieldAccessor("add", NewConst(1)), NewConst(2))).
			WithInterface(types.Int, NewDefinitions(0).WithAssignment("add", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer untyped interface function calls",
		exp: NewBlock(NewFCall(NewFieldAccessor("add", NewConst(1)), NewConst(2))).
			WithInterface(nil, NewDefinitions(0).WithAssignment("add", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer multiple interface invocations",
		exp: NewBlock(NewFCall(NewFieldAccessor("add", NewFCall(NewFieldAccessor("add", NewConst(1)), NewConst(2))), NewConst(3))).
			WithInterface(nil, NewDefinitions(0).WithAssignment("add", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"))),
		typ: "int",
		err: nil,
	}, {
		name: "infer interface usage in functions",
		exp: NewBlock(NewFCall(NewId("f"), NewConst(2))).
			WithInterface(nil, NewDefinitions(0).WithAssignment("isOdd", NewFDef(NewOp("==", NewConst(0), NewOp("%", NewId("$"), NewConst(2)))))).
			WithAssignment("f", NewFDef(NewFCall(NewFieldAccessor("isOdd", NewId("x"))), &FParam{"x", types.Int})),
		typ: "bool",
	}, {
		name: "fail on invalid interface call",
		exp: NewBlock(NewFCall(NewFCall(NewFieldAccessor("add", NewConst("a")), NewConst("b")))).
			WithInterface(nil, NewDefinitions(0).WithAssignment("add", NewFDef(NewOp("-", NewId("$"), NewId("x")), "x"))),
		err: errors.New("can not unify V1{add:V2} with string"),
	}, {
		name: "infers aggregate type from all methods of an interface",
		exp: NewBlock(NewFCall(NewFCall(NewFieldAccessor("identity", NewConst("str"))))).
			WithInterface(nil, NewDefinitions(0).
				WithAssignment("sub", NewFDef(NewOp("-", NewId("$"), NewId("x")), "x")).
				WithAssignment("identity", NewFDef(NewId("$"))),
			),
		err: errors.New("can not unify V1{identity:V2} with string"),
	}, {
		name: "support functions summing over same argument",
		exp: NewBlock(NewFCall(NewId("a"), NewConst(0))).
			WithAssignment("a", NewFDef(NewOp("+", NewId("x"), NewId("x")), "x")),
		typ: "int",
	}, {
		name: "infers the interface method return type from the interface type",
		exp: NewBlock(NewFCall(NewFieldAccessor("a", NewConst(1)))).
			WithInterface(nil, NewDefinitions(0).
				WithAssignment("a", NewFDef(NewOp("+", NewId("$"), NewId("$")))),
			),
		typ: "int",
	}, {
		name: "infers the interface method return type from multiple interfaces with different types",
		exp: NewBlock(NewFCall(NewFieldAccessor("a", NewConst(1)), NewConst(1))).
			WithInterface(types.Int, NewDefinitions(0).
				WithAssignment("a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x")),
			).
			WithInterface(types.Real, NewDefinitions(0).
				WithAssignment("a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x")),
			),
		typ: "int",
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
				if (tt.exp.Type() == nil) && tt.typ != "" {
					t.Errorf("Infer() wrong type = nil, want %v", tt.typ)
				} else if tt.exp.Type() != nil && types.Signature(tt.exp.Type()) != tt.typ {
					t.Errorf("Infer() wrong type = %v, want %v", types.Signature(tt.exp.Type()), tt.typ)
				}
			}
		})
	}
}

func TestComplexInferences(t *testing.T) {
	p, _ := grammar.Parse(`
		operate = (x, y,  f) => { f(x, y) }
		add = (x, y) => { x + y }
		sub = (x, y) => { x - y }
		pick = (b) => { if (b) sub else add }
		operate(3, 1, pick(true)) + operate(5, 1, pick(false)) + operate(1, 1, (x, y) => { x + y })
	`)
	a := p.ToAst()
	err := Infer(a)
	if err != nil {
		t.Error(err)
	}
	sign := types.Signature(a.Type())
	if sign != "int" {
		t.Error("got " + sign + ", expected int")
	}
	ast.VisitAfter(a.(*Block).Value, func(v ast.Ast, ctx *ast.VisitContext) error {
		if e, ok := v.(ast.Expression); ok {
			if types.HasFreeVars(e.Type()) {
				t.Error("free variables in the root expression: " + ast.Stringify(e))
			}
		}
		return nil
	})
}
