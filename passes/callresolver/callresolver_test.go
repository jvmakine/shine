package callresolver

import (
	"testing"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	"github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
)

func TestResolveFunctions(t *testing.T) {
	tests := []struct {
		name   string
		before Expression
		after  Expression
	}{{
		name: "resolves function signatures based on the call type",
		before: NewBlock(NewFCall(NewOp("if"),
			NewFCall(NewId("a"), NewConst(true), NewConst(true), NewConst(false)),
			NewFCall(NewId("a"), NewConst(true), NewConst(5), NewConst(6)),
			NewConst(7)),
		).WithAssignment("a", NewFDef(NewFCall(NewOp("if"), NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")),
		after: NewBlock(NewFCall(NewOp("if"),
			NewFCall(NewId("a%%1%%(bool,bool,bool)=>bool"), NewConst(true), NewConst(true), NewConst(false)),
			NewFCall(NewId("a%%1%%(bool,int,int)=>int"), NewConst(true), NewConst(5), NewConst(6)),
			NewConst(7)),
		).WithAssignment("a", NewFDef(NewFCall(NewOp("if"), NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")).
			WithAssignment("a%%1%%(bool,bool,bool)=>bool", NewFDef(NewFCall(NewOp("if"), NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")).
			WithAssignment("a%%1%%(bool,int,int)=>int", NewFDef(NewFCall(NewOp("if"), NewId("b"), NewId("y"), NewId("x")), "b", "y", "x")),
	}, {
		name: "resolves functions as arguments",
		before: NewBlock(NewFCall(NewId("a"), NewId("b"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("b", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y")),
		after: NewBlock(NewFCall(NewId("a%%1%%((int,int)=>int)=>int"), NewId("b%%1%%(int,int)=>int"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("b", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y")).
			WithAssignment("a%%1%%((int,int)=>int)=>int", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("b%%1%%(int,int)=>int", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y")),
	}, {
		name: "resolves anonymous functions",
		before: NewBlock(NewFCall(NewId("a"), NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")),
		after: NewBlock(NewFCall(NewId("a%%1%%((int,int)=>int)=>int"), NewId("<anon1>%%1%%(int,int)=>int"))).
			WithAssignment("a", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("a%%1%%((int,int)=>int)=>int", NewFDef(NewFCall(NewId("f"), NewConst(1), NewConst(2)), "f")).
			WithAssignment("<anon1>%%1%%(int,int)=>int", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y")),
	}, {
		name: "resolves simple structures",
		before: NewBlock(NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("a", NewStruct(ast.StructField{"x", types.IntP})),
		after: NewBlock(NewFCall(NewId("a%%1%%(int)=>a{x:int}"), NewConst(1))).
			WithAssignment("a", NewStruct(ast.StructField{"x", types.IntP})).
			WithAssignment("a%%1%%(int)=>a{x:int}", NewStruct(ast.StructField{"x", types.IntP})),
	}, {
		name: "resolves multitype structures",
		before: NewBlock(NewFCall(NewId("a"), NewFCall(NewId("a"), NewConst(1)))).
			WithAssignment("a", NewStruct(ast.StructField{"x", types.MakeVariable()})),
		after: NewBlock(NewFCall(NewId("a%%1%%(a{x:int})=>a"), NewFCall(NewId("a%%1%%(int)=>a{x:int}"), NewConst(1)))).
			WithAssignment("a", NewStruct(ast.StructField{"x", types.MakeVariable()})).
			WithAssignment("a%%1%%(int)=>a{x:int}", NewStruct(ast.StructField{"x", types.IntP})).
			WithAssignment("a%%1%%(a{x:int})=>a", NewStruct(ast.StructField{"x", types.MakeNamed("a")})),
	}}
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
	RewriteTypes(e, func(t types.Type, ctx *VisitContext) (types.Type, error) {
		return types.IntP, nil
	})
	VisitBefore(e, func(v Ast, ctx *VisitContext) error {
		if b, ok := v.(*Block); ok {
			b.ID = 0
		}
		return nil
	})
}
