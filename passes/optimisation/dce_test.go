package optimisation

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/passes/typeinference"

	. "github.com/jvmakine/shine/ast"
)

func TestDCE(t *testing.T) {
	type args struct {
		exp Expression
	}
	tests := []struct {
		name   string
		before Expression
		after  Expression
	}{{
		name: "removes unused assignments from blocks",
		before: NewBlock(NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("a", NewFDef(NewOp("+", NewId("x"), NewId("y")), "x")).
			WithAssignment("y", NewConst(5)).
			WithAssignment("z", NewConst(4)),
		after: NewBlock(NewFCall(NewId("a"), NewConst(1))).
			WithAssignment("a", NewFDef(NewOp("+", NewId("x"), NewId("y")), "x")).
			WithAssignment("y", NewConst(5)),
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
			eraseType(tt.after)
			eraseType(tt.before)
			if !reflect.DeepEqual(tt.before, tt.after) {
				t.Errorf("Resolve() = %v, want %v", tt.before, tt.after)
			}
		})
	}
}
