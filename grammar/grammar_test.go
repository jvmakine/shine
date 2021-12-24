package grammar

import (
	"errors"
	"testing"

	"github.com/jvmakine/shine/ast"
	a "github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
	"github.com/stretchr/testify/require"
)

func TestExpressionParsing(tes *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *a.Exp
		err   error
	}{{
		name:  "parse an int const",
		input: "42",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.IConst(42)),
	}, {
		name:  "parse a real const",
		input: "0.1",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.RConst(0.1)),
	}, {
		name:  "parse a bool const",
		input: "true",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.BConst(true)),
	}, {
		name:  "parse an identifier",
		input: "abc",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Id("abc")),
	}, {
		name:  "parse + term expression",
		input: "1 + 2",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("+"), t.IConst(1), t.IConst(2))),
	}, {
		name:  "parse - term expression",
		input: "1 - 2",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("-"), t.IConst(1), t.IConst(2))),
	}, {
		name:  "parse * factor expression",
		input: "2 * 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("*"), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse / factor expression",
		input: "2 / 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("/"), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse % factor expression",
		input: "2 % 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("%"), t.IConst(2), t.IConst(3))),
	}, {
		name:  "maintain right precedence with + and *",
		input: "2 + 3 * 4",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("+"), t.IConst(2), t.Fcall(t.Op("*"), t.IConst(3), t.IConst(4)))),
	}, {
		name:  "parse numeric expressions with brackets",
		input: "(2 + 4) * 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("*"), t.Fcall(t.Op("+"), t.IConst(2), t.IConst(4)), t.IConst(3))),
	}, {
		name:  "parse id expressions with brackets",
		input: "(c % 2) == 0",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("=="), t.Fcall(t.Op("%"), t.Id("c"), t.IConst(2)), t.IConst(0))),
	}, {
		name:  "parse == operator",
		input: "2 == 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("=="), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse != operator",
		input: "a != 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("!="), t.Id("a"), t.IConst(3))),
	}, {
		name:  "parse < operator",
		input: "2 < 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("<"), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse > operator",
		input: "2 > 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op(">"), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse >= operator",
		input: "2 >= 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op(">="), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse <= operator",
		input: "2 <= 3",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("<="), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse || operator",
		input: "true || false",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("||"), t.BConst(true), t.BConst(false))),
	}, {
		name:  "parse && operator",
		input: "true && false",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("&&"), t.BConst(true), t.BConst(false))),
	}, {
		name:  "parse if expression",
		input: "if(2 > 3) 1 else 2",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("if"), t.Fcall(t.Op(">"), t.IConst(2), t.IConst(3)), t.IConst(1), t.IConst(2))),
	}, {
		name:  "parse if expressions with blocks",
		input: "if(2 > 3) { 1 } else { 2 }",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("if"), t.Fcall(t.Op(">"), t.IConst(2), t.IConst(3)), t.Block(t.Assgs{}, t.Typedefs{}, t.IConst(1)), t.Block(t.Assgs{}, t.Typedefs{}, t.IConst(2)))),
	}, {
		name:  "parse if else if expression",
		input: "if (2 > 3) 1 else if (3 > 4) 2 else 4",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("if"), t.Fcall(t.Op(">"), t.IConst(2), t.IConst(3)), t.IConst(1), t.Fcall(t.Op("if"), t.Fcall(t.Op(">"), t.IConst(3), t.IConst(4)), t.IConst(2), t.IConst(4)))),
	}, {
		name:  "parse a function call",
		input: "f(1, x, y)",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Id("f"), t.IConst(1), t.Id("x"), t.Id("y"))),
	}, {
		name:  "parse a function calls of returned function values",
		input: "f(1, x, y)(2, 3)",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Fcall(t.Id("f"), t.IConst(1), t.Id("x"), t.Id("y")), t.IConst(2), t.IConst(3))),
	}, {
		name:  "parse functions as values",
		input: "f((x) => {x + 2}, (y) => {y + 1})",
		want:  t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Id("f"), t.Fdef(t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("+"), t.Id("x"), t.IConst(2))), "x"), t.Fdef(t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("+"), t.Id("y"), t.IConst(1))), "y"))),
	}, {
		name: "parse several assignments",
		input: `
			a = 1 + 2
			b = 2 + 3
			c = 3 + 4
			a + b + c
		`,
		want: t.Block(
			t.Assgs{
				"a": t.Fcall(t.Op("+"), t.IConst(1), t.IConst(2)),
				"b": t.Fcall(t.Op("+"), t.IConst(2), t.IConst(3)),
				"c": t.Fcall(t.Op("+"), t.IConst(3), t.IConst(4)),
			},
			t.Typedefs{},
			t.Fcall(t.Op("+"), t.Fcall(t.Op("+"), t.Id("a"), t.Id("b")), t.Id("c")),
		),
	}, {
		name: "parse a function definition",
		input: `
			a = (x, y) => { x + y }
			a(1, 2)
		`,
		want: t.Block(
			t.Assgs{"a": t.Fdef(t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("+"), t.Id("x"), t.Id("y"))), "x", "y")},
			t.Typedefs{},
			t.Fcall(t.Id("a"), t.IConst(1), t.IConst(2)),
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
							"b": t.Fdef(t.Block(t.Assgs{}, t.Typedefs{}, t.Fcall(t.Op("+"), t.Id("x"), t.IConst(1))), "x"),
						},
						t.Typedefs{},
						t.Fcall(t.Op("+"), t.Id("x"), t.Fcall(t.Id("b"), t.Id("y"))),
					),
					"x", "y"),
			},
			t.Typedefs{},
			t.Fcall(t.Id("a"), t.IConst(1), t.IConst(2)),
		),
	}, {
		name: "parse sequential function definitions",
		input: `
			a = (x) => (y) => x + y
			a(1)(2)
		`,
		want: t.Block(
			t.Assgs{"a": t.Fdef(t.Fdef(t.Fcall(t.Op("+"), t.Id("x"), t.Id("y")), "y"), "x")},
			t.Typedefs{},
			t.Fcall(t.Fcall(t.Id("a"), t.IConst(1)), t.IConst(2)),
		),
	}, {
		name: "parse explicit type definitions on functions",
		input: `
			a = (x:int, y:real, z:bool): real => if (b && y > 1.0) x else 0
			a(1, 2.0, true)
		`,
		want: t.Block(
			t.Assgs{"a": t.Fdef(t.TDecl(t.Fcall(
				t.Op("if"),
				t.Fcall(t.Op("&&"), t.Id("b"), t.Fcall(t.Op(">"), t.Id("y"), t.RConst(1.0))),
				t.Id("x"),
				t.IConst(0),
			), types.RealP), t.Param("x", types.IntP), t.Param("y", types.RealP), t.Param("z", types.BoolP))},
			t.Typedefs{},
			t.Fcall(t.Id("a"), t.IConst(1), t.RConst(2.0), t.BConst(true)),
		),
	}, {
		name:  "parse explicit type definitions on generic expression",
		input: `((1:int) + (2:bool)):real`,
		want: t.Block(t.Assgs{}, t.Typedefs{}, t.TDecl(
			t.Fcall(
				t.Op("+"),
				t.TDecl(t.IConst(1), types.IntP),
				t.TDecl(t.IConst(2), types.BoolP),
			),
			types.RealP,
		)),
	}, {
		name: "parse function type definitions",
		input: `
			a = (x:int, f:(int)=>bool) => if (f(x)) x else 0
			a(1, b)
		`,
		want: t.Block(
			t.Assgs{"a": t.Fdef(t.Fcall(
				t.Op("if"),
				t.Fcall(t.Id("f"), t.Id("x")),
				t.Id("x"),
				t.IConst(0),
			), t.Param("x", types.IntP), t.Param("f", types.MakeFunction(types.IntP, types.BoolP)))},
			t.Typedefs{},
			t.Fcall(t.Id("a"), t.IConst(1), t.Id("b")),
		),
	}, {
		name: "parse structure definitions",
		input: `
			a :: (x:int, c)
			a(1, true)
		`,
		want: t.Block(
			t.Assgs{},
			t.Typedefs{"a": t.Struct(ast.StructField{"x", types.IntP}, ast.StructField{"c", types.Type{}})},
			t.Fcall(t.Id("a"), t.IConst(1), t.BConst(true)),
		),
	}, {
		name:  "parse simple field accessors",
		input: `a.foo.bar`,
		want: t.Block(
			t.Assgs{},
			t.Typedefs{},
			t.Faccess(t.Faccess(t.Id("a"), "foo"), "bar"),
		),
	}, {
		name:  "parse method calls",
		input: `a.foo(1)`,
		want: t.Block(
			t.Assgs{},
			t.Typedefs{},
			t.Fcall(t.Faccess(t.Id("a"), "foo"), t.IConst(1)),
		),
	}, {
		name:  "parse sequential method calls",
		input: `a.foo(1).bar(2)`,
		want: t.Block(
			t.Assgs{},
			t.Typedefs{},
			t.Fcall(t.Faccess(t.Fcall(t.Faccess(t.Id("a"), "foo"), t.IConst(1)), "bar"), t.IConst(2)),
		),
	}, {
		name: "parse custom types in functions",
		input: `
			a = (x: A) => x
			a(b)
		`,
		want: t.Block(
			t.Assgs{"a": t.Fdef(t.Id("x"), t.Param("x", types.MakeNamed("A")))},
			t.Typedefs{},
			t.Fcall(t.Id("a"), t.Id("b")),
		),
	}, {
		name: "parse type parameters in functions",
		input: `
			a = (x: A[int]) => x
			a(b)
		`,
		want: t.Block(
			t.Assgs{"a": t.Fdef(t.Id("x"), t.Param("x", types.MakeNamed("A", types.IntP)))},
			t.Typedefs{},
			t.Fcall(t.Id("a"), t.Id("b")),
		),
	}, {
		name: "parse typed constants",
		input: `a:int = 5
				a
			`,
		want: t.Block(
			t.Assgs{"a": t.TDecl(t.IConst(5), types.IntP)},
			t.Typedefs{},
			t.Id("a"),
		),
	}, {
		name: "parse type variables",
		input: `A[X] :: (l:X, r:X) 
				A(1,2)
			`,
		want: t.Block(
			t.Assgs{},
			t.Typedefs{"A": &a.TypeDefinition{
				FreeVariables: []string{"X"},
				Struct: &a.Struct{Fields: []*a.StructField{{
					Name: "l",
					Type: types.MakeNamed("X"),
				}, {
					Name: "r",
					Type: types.MakeNamed("X"),
				}}},
			}},
			t.Fcall(t.Id("A"), t.IConst(1), t.IConst(2)),
		),
	}, {
		name: "fails on duplicate definitions",
		input: `a = 1 
				a = 2
				a
			`,
		err: errors.New("redefinition of a"),
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			prog, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			got, err := prog.ToAst()
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, got)
		})
	}
}
