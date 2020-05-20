package optimisation

import (
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
)

func TestSeqFN(t *testing.T) {
	tests := []struct {
		name   string
		before *ast.Exp
		after  *ast.Exp
	}{{
		name: "combines sequential functions when possible",
		before: Block(
			Assgs{"a": Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("y")), "y"), "x")},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		after: Block(
			Assgs{
				"a":   Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("y")), "y"), "x"),
				"a%c": Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
			},
			Fcall(Id("a%c"), IConst(1), IConst(2)),
		),
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

func eraseType(e *ast.Exp) {
	e.Visit(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Id != nil {
			v.Id.Type = types.IntP
		} else if v.Op != nil {
			v.Op.Type = types.IntP
		} else if v.Const != nil {
			v.Const.Type = types.IntP
		} else if v.Call != nil {
			v.Call.Type = types.IntP
		} else if v.Def != nil {
			for _, p := range v.Def.Params {
				p.Type = types.IntP
			}
			v.Def.Closure = nil
		} else if v.Block != nil {
			v.Block.ID = 0
		}
		return nil
	})
}
