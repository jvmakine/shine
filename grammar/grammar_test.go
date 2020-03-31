package grammar

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
)

func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *ast.Exp
	}{{
		name:  "parse a simple numeric expression",
		input: "42",
		want:  block(iconst(42)),
	}, {
		name:  "parse an identifier",
		input: "abc",
		want:  block(id("abc")),
	}, {
		name:  "parse + term expression",
		input: "1 + 2",
		want:  block(fcall("+", iconst(1), iconst(2))),
	}, {
		name:  "parse - term expression",
		input: "1 - 2",
		want:  block(fcall("-", iconst(1), iconst(2))),
	}, {
		name:  "parse * factor expression",
		input: "2 * 3",
		want:  block(fcall("*", iconst(2), iconst(3))),
	}, {
		name:  "parse / factor expression",
		input: "2 / 3",
		want:  block(fcall("/", iconst(2), iconst(3))),
	}, {
		name:  "maintain right precedence with + and *",
		input: "2 + 3 * 4",
		want:  block(fcall("+", iconst(2), fcall("*", iconst(3), iconst(4)))),
	}, {
		name:  "parse a function call",
		input: "f(1, x, y)",
		want:  block(fcall("f", iconst(1), id("x"), id("y"))),
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := Parse(tt.input)
			got := prog.ToAst()
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func iconst(v int) *ast.Exp {
	return &ast.Exp{
		Const: &ast.Const{Int: &v},
	}
}

func id(name string) *ast.Exp {
	return &ast.Exp{
		Id: &name,
	}
}

func block(e *ast.Exp) *ast.Exp {
	return &ast.Exp{
		Block: &ast.Block{Value: e, Assignments: []*ast.Assign{}},
	}
}

func fcall(name string, args ...*ast.Exp) *ast.Exp {
	call := &ast.FCall{
		Name:   name,
		Params: args,
	}
	return &ast.Exp{
		Call: call,
	}
}
