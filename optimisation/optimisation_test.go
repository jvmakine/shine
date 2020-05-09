package optimisation

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/typeinference"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/test"
)

func TestOptimise(t *testing.T) {
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
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Fcall("+", Id("x"), Id("y")), "x")),
			Assign("y", IConst(5)),
			Assign("z", IConst(4)),
		),
		after: Block(
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Fcall("+", Id("x"), Id("y")), "x")),
			Assign("y", IConst(5)),
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
			Optimise(tt.before)
			if !reflect.DeepEqual(tt.before, tt.after) {
				t.Errorf("Resolve() = %v, want %v", tt.before, tt.after)
			}
		})
	}
}
