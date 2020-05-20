package optimisation

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/passes/typeinference"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/test"
)

func TestDCE(t *testing.T) {
	type args struct {
		exp *ast.Exp
	}
	tests := []struct {
		name   string
		before *ast.Exp
		after  *ast.Exp
	}{{
		name: "removes unused assignments from blocks",
		before: Block(
			Assgs{
				"a": Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x"),
				"y": IConst(5),
				"z": IConst(4),
			},
			Fcall(Id("a"), IConst(1)),
		),
		after: Block(
			Assgs{
				"a": Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x"),
				"y": IConst(5),
			},
			Fcall(Id("a"), IConst(1)),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := typeinference.Infer(tt.before)
			if err != nil {
				panic(err)
			}
			err = typeinference.Infer(tt.after)
			if err != nil {
				panic(err)
			}
			DeadCodeElimination(tt.before)
			if !reflect.DeepEqual(tt.before, tt.after) {
				t.Errorf("Resolve() = %v, want %v", tt.before, tt.after)
			}
		})
	}
}
