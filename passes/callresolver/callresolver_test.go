package callresolver

import (
	"testing"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
)

func TestResolveFunctions(t *testing.T) {
	tests := []struct {
		name   string
		before Expression
		after  Expression
	}{{
		name: "resolves function signatures based on the call type",
		before: NewBlock(NewBranch(
			NewFCall(NewId("a"), NewConst(true), NewConst(true), NewConst(false)),
			NewFCall(NewId("a"), NewConst(true), NewConst(5), NewConst(6)),
			NewConst(7)),
		).WithID(1).WithAssignment("a",
			NewFDef(NewBranch(NewId("b"), NewId("y"), NewId("x")), "b", "y", "x"),
		),
		after: NewBlock(NewBranch(
			NewFCall(NewId("a%%1%%(bool,bool,bool)=>bool"), NewConst(true), NewConst(true), NewConst(false)),
			NewFCall(NewId("a%%1%%(bool,int,int)=>int"), NewConst(true), NewConst(5), NewConst(6)),
			NewConst(7)),
		).WithAssignment("a", NewFDef(NewBranch(NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")).
			WithAssignment("a%%1%%(bool,bool,bool)=>bool", NewFDef(NewBranch(NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")).
			WithAssignment("a%%1%%(bool,int,int)=>int", NewFDef(NewBranch(NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")),
	}, {
		name: "resolves operations in assignments",
		before: NewBlock(NewFCall(NewId("a"), NewConst(6), NewConst(7))).
			WithID(1).WithAssignment("a",
			NewFDef(
				NewBlock(NewId("added")).WithAssignment("added", NewOp("+", NewId("x"), NewId("y"))),
				"y", "x",
			),
		),
		after: NewBlock(NewFCall(NewId("a%%1%%(int,int)=>int"), NewConst(6), NewConst(7))).
			WithID(1).WithAssignment("a%%1%%(int,int)=>int",
			NewFDef(
				NewBlock(NewId("added")).WithAssignment("added", NewPrimitiveOp("int_+", Int, NewId("x"), NewId("y"))),
				"y", "x",
			)).WithAssignment("a",
			NewFDef(
				NewBlock(NewId("added")).WithAssignment("added", NewOp("+", NewId("x"), NewId("y"))),
				"y", "x",
			)),
	}, {
		name: "resolves functions as arguments",
		before: NewBlock(NewFCall(NewId("a"), NewId("b"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("b", NewFDef(NewOp("+", NewId("x"), NewId("y")), "x", "y")).
			WithID(1),
		after: NewBlock(NewFCall(NewId("a%%1%%((int,int)=>int)=>int"), NewId("b%%1%%(int,int)=>int"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("b", NewFDef(NewOp("+", NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("a%%1%%((int,int)=>int)=>int", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("b%%1%%(int,int)=>int", NewFDef(NewPrimitiveOp("int_+", Int, NewId("x"), NewId("y")), "x", "y")).
			WithID(1),
	}, {
		name: "resolves anonymous functions",
		before: NewBlock(NewFCall(NewId("a"), NewFDef(NewOp("+", NewId("x"), NewId("y")), "x", "y"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithID(1),
		after: NewBlock(NewFCall(NewId("a%%1%%((int,int)=>int)=>int"), NewId("<anon1>%%1%%(int,int)=>int"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("a%%1%%((int,int)=>int)=>int", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("<anon1>%%1%%(int,int)=>int", NewFDef(NewPrimitiveOp("int_+", Int, NewId("x"), NewId("y")), "x", "y")).
			WithID(1),
	}, /* {
		name: "resolves simple structures",
		before: NewBlock(NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("a", NewStruct(ast.StructField{"x", Int})).
			WithID(1),
		after: NewBlock(NewFCall(NewId("a%%1%%(int)=>a{x:int}"), NewConst(1))).
			WithAssignment("a", NewStruct(ast.StructField{"x", Int})).
			WithAssignment("a%%1%%(int)=>a{x:int}", NewStruct(ast.StructField{"x", Int})).
			WithID(1),
	}, {
		name: "resolves multitype structures",
		before: NewBlock(NewFCall(NewId("a"), NewFCall(NewId("a"), NewConst(1)))).
			WithAssignment("a", NewStruct(ast.StructField{"x", NewVariable()})).
			WithID(1),
		after: NewBlock(NewFCall(NewId("a%%1%%(a{x:int})=>a"), NewFCall(NewId("a%%1%%(int)=>a{x:int}"), NewConst(1)))).
			WithAssignment("a", NewStruct(ast.StructField{"x", NewVariable()})).
			WithAssignment("a%%1%%(int)=>a{x:int}", NewStruct(ast.StructField{"x", Int})).
			WithAssignment("a%%1%%(a{x:int})=>a", NewStruct(ast.StructField{"x", NewNamed("a", nil)})).
			WithID(1),
	}*/}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeinference.Infer(tt.before)
			ResolveFunctions(tt.before)
			eraseType(tt.after)
			eraseType(tt.before)
			ok, err := deepdiff.DeepDiff(tt.before, tt.after)
			if !ok {
				t.Error(err)
			}
		})
	}
}

func eraseType(e Expression) {
	RewriteTypes(e, func(t Type, ctx *VisitContext) (Type, error) {
		return Int, nil
	})
	VisitBefore(e, func(v Ast, ctx *VisitContext) error {
		if b, ok := v.(*Block); ok {
			b.Def.ID = 0
		}
		return nil
	})
}
