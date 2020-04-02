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
	}, {
		name: "parse a function definition",
		input: `
			a = (x, y) => { x + y }
			a(1, 2)
		`,
		want: block(
			fcall("a", iconst(1), iconst(2)),
			assign("a", fdef(block(fcall("+", id("x"), id("y"))), "x", "y"))),
	}, {
		name: "parse a nested function definition",
		input: `
			a = (x, y) => {
				b = (x) => { x + 1 }
				x + b(y)
			}
			a(1, 2)
		`,
		want: block(
			fcall("a", iconst(1), iconst(2)),
			assign("a", fdef(
				block(
					fcall("+", id("x"), fcall("b", id("y"))),
					assign("b", fdef(block(fcall("+", id("x"), iconst(1))), "x"))),
				"x", "y"))),
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			got := prog.ToAst()
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

func assign(name string, value *ast.Exp) *ast.Assign {
	return &ast.Assign{
		Name:  name,
		Value: value,
	}
}

func block(e *ast.Exp, assigns ...*ast.Assign) *ast.Exp {
	as := assigns
	if as == nil {
		as = []*ast.Assign{}
	}
	return &ast.Exp{
		Block: &ast.Block{Value: e, Assignments: as},
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

func fdef(body *ast.Exp, args ...string) *ast.Exp {
	params := make([]*ast.FParam, len(args))
	for i, p := range args {
		params[i] = &ast.FParam{Name: p}
	}
	fdef := &ast.FDef{
		Body:   body,
		Params: params,
	}
	return &ast.Exp{
		Def: fdef,
	}
}
