package optimisation

import (
	"testing"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	"github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
)

func TestSeqFN(t *testing.T) {
	tests := []struct {
		name   string
		before Expression
		after  Expression
	}{{
		name: "combines sequential functions when possible",
		before: NewBlock(NewFCall(NewFCall(NewId("a"), NewConst(1)), NewConst(2))).
			WithAssignment("a", NewFDef(NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "y"), "x")),
		after: NewBlock(NewFCall(NewId("a%c"), NewConst(1), NewConst(2))).
			WithAssignment("a", NewFDef(NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "y"), "x")).
			WithAssignment("a%c", NewFDef(NewFCall(NewOp("+"), NewId("x"), NewId("y")), "x", "y")),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeinference.Infer(tt.before)
			SequentialFunctionPass(tt.before)

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
	VisitAfter(e, func(a Ast, ctx *VisitContext) error {
		if b, ok := a.(*Block); ok {
			b.Def.ID = 0
		}
		return nil
	})
}
