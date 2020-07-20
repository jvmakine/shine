package interfaceresolver

import (
	"errors"
	"testing"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name   string
		before Expression
		after  Expression
		err    error
	}{{
		name: "converts interface call into a function",
		before: NewBlock(NewFCall(NewFieldAccessor("a", NewConst(0)), NewConst(1))).
			WithInterface(Int, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"),
			)),
		after: NewBlock(NewFCall(NewFieldAccessor("a", NewConst(0)), NewConst(1))).
			WithAssignment("a%interface%0%int", NewFDef(NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"), "$")).
			WithInterface(Int, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"),
			)),
	}, {
		name: "returns an error if a method is declared twice for the same type",
		before: NewBlock(NewFCall(NewFieldAccessor("a", NewConst(0)), NewConst(1))).
			WithInterface(Int, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"),
			)).
			WithInterface(Int, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"),
			)),
		after: nil,
		err:   errors.New("a declared twice for unifiable types: int, int"),
	}, {
		name: "returns an error if a method could unify to two different implementations",
		before: NewBlock(NewFCall(NewFieldAccessor("a", NewConst(0)), NewConst(1))).
			WithInterface(Int, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"),
			)).
			WithInterface(nil, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewOp("+", NewId("$"), NewId("x")), "x"),
			)),
		after: nil,
		err:   errors.New("a declared twice for unifiable types: V1{+:(V2)=>V3}, int"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeinference.Infer(tt.before)
			Resolve(tt.before)
			if tt.after != nil {
				eraseType(tt.after)
			}
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
		return nil, nil
	})
	VisitBefore(e, func(v Ast, ctx *VisitContext) error {
		if b, ok := v.(*Block); ok {
			b.Def.ID = 0
		}
		return nil
	})
}
