package grammar

import (
	"reflect"
	"testing"

	a "github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/test"
)

func TestExpressionParsing(tes *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *a.Exp
	}{{
		name:  "parse an int const",
		input: "42",
		want:  t.Block(t.Assgs{}, t.IConst(42)),
	}, {
		name:  "parse a real const",
		input: "0.1",
		want:  t.Block(t.Assgs{}, t.RConst(0.1)),
	}, {
		name:  "parse a bool const",
		input: "true",
		want:  t.Block(t.Assgs{}, t.BConst(true)),
	}, {
		name:  "parse an identifier",
		input: "abc",
		want:  t.Block(t.Assgs{}, t.Id("abc")),
	}, {
		name:  "parse + term expression",
		input: "1 + 2",
		want:  t.Block(t.Assgs{}, t.Fcall("+", t.IConst(1), t.IConst(2))),
	}, {
		name:  "parse - term expression",
		input: "1 - 2",
		want:  t.Block(t.Assgs{}, t.Fcall("-", t.IConst(1), t.IConst(2))),
	}, {
		name:  "parse * factor expression",
		input: "2 * 3",
		want:  t.Block(t.Assgs{}, t.Fcall("*", t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse / factor expression",
		input: "2 / 3",
		want:  t.Block(t.Assgs{}, t.Fcall("/", t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse % factor expression",
		input: "2 % 3",
		want:  t.Block(t.Assgs{}, t.Fcall("%", t.IConst(2), t.IConst(3))),
	}, {
		name:  "maintain right precedence with + and *",
		input: "2 + 3 * 4",
		want:  t.Block(t.Assgs{}, t.Fcall("+", t.IConst(2), t.Fcall("*", t.IConst(3), t.IConst(4)))),
	}, {
		name:  "parse == operator",
		input: "2 == 3",
		want:  t.Block(t.Assgs{}, t.Fcall("==", t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse < operator",
		input: "2 < 3",
		want:  t.Block(t.Assgs{}, t.Fcall("<", t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse > operator",
		input: "2 > 3",
		want:  t.Block(t.Assgs{}, t.Fcall(">", t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse >= operator",
		input: "2 >= 3",
		want:  t.Block(t.Assgs{}, t.Fcall(">=", t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse <= operator",
		input: "2 <= 3",
		want:  t.Block(t.Assgs{}, t.Fcall("<=", t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse || operator",
		input: "true || false",
		want:  t.Block(t.Assgs{}, t.Fcall("||", t.BConst(true), t.BConst(false))),
	}, {
		name:  "parse && operator",
		input: "true && false",
		want:  t.Block(t.Assgs{}, t.Fcall("&&", t.BConst(true), t.BConst(false))),
	}, {
		name:  "parse if expression",
		input: "if(2 > 3) 1 else 2",
		want:  t.Block(t.Assgs{}, t.Fcall("if", t.Fcall(">", t.IConst(2), t.IConst(3)), t.IConst(1), t.IConst(2))),
	}, {
		name:  "parse a function call",
		input: "f(1, x, y)",
		want:  t.Block(t.Assgs{}, t.Fcall("f", t.IConst(1), t.Id("x"), t.Id("y"))),
	}, {
		name: "parse a function definition",
		input: `
			a = (x, y) => { x + y }
			a(1, 2)
		`,
		want: t.Block(
			t.Assgs{"a": t.Fdef(t.Block(t.Assgs{}, t.Fcall("+", t.Id("x"), t.Id("y"))), "x", "y")},
			t.Fcall("a", t.IConst(1), t.IConst(2)),
		),
	}, {
		name: "parse a nested function definition",
		input: `
			a = (x, y) => {
				b = (x) => { x + 1 }
				x + b(y)
			}
			a(1, 2)
		`,
		want: t.Block(
			t.Assgs{
				"a": t.Fdef(
					t.Block(
						t.Assgs{
							"b": t.Fdef(t.Block(t.Assgs{}, t.Fcall("+", t.Id("x"), t.IConst(1))), "x"),
						},
						t.Fcall("+", t.Id("x"), t.Fcall("b", t.Id("y"))),
					),
					"x", "y"),
			},
			t.Fcall("a", t.IConst(1), t.IConst(2)),
		),
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
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
