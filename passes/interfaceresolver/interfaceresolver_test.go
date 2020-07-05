package interfaceresolver

import (
	"testing"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	"github.com/jvmakine/shine/types"
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
			WithInterface(types.IntP, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewFCall(NewOp("+"), NewId("$"), NewId("x")), "x"),
			)),
		after: NewBlock(NewFCall(NewFCall(NewId("a%interface%0%int"), NewConst(0)), NewConst(1))).
			WithAssignment("a%interface%0%int", NewFDef(NewFDef(NewFCall(NewOp("+"), NewId("$"), NewId("x")), "x"), "$")).
			WithInterface(types.IntP, NewDefinitions(0).WithAssignment(
				"a", NewFDef(NewFCall(NewOp("+"), NewId("$"), NewId("x")), "x"),
			)),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeinference.Infer(tt.before)
			Resolve(tt.before)
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
			b.Def.ID = 0
		}
		return nil
	})
}
