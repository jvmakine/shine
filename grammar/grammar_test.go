package grammar

import (
	"testing"

	a "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
)

func TestExpressionParsing(tes *testing.T) {
	tests := []struct {
		name  string
		input string
		want  a.Expression
	}{{
		name:  "parse an int const",
		input: "42",
		want:  a.NewBlock(a.NewConst(42)),
	}, {
		name:  "parse a real const",
		input: "0.1",
		want:  a.NewBlock(a.NewConst(0.1)),
	}, {
		name:  "parse a bool const",
		input: "true",
		want:  a.NewBlock(a.NewConst(true)),
	}, {
		name:  "parse an identifier",
		input: "abc",
		want:  a.NewBlock(a.NewId("abc")),
	}, {
		name:  "parse + term expression",
		input: "1 + 2",
		want:  a.NewBlock(a.NewFCall(a.NewOp("+"), a.NewConst(1), a.NewConst(2))),
	}, {
		name:  "parse - term expression",
		input: "1 - 2",
		want:  a.NewBlock(a.NewFCall(a.NewOp("-"), a.NewConst(1), a.NewConst(2))),
	}, {
		name:  "parse * factor expression",
		input: "2 * 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp("*"), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse / factor expression",
		input: "2 / 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp("/"), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse % factor expression",
		input: "2 % 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp("%"), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "maintain right precedence with + and *",
		input: "2 + 3 * 4",
		want: a.NewBlock(
			a.NewFCall(a.NewOp("+"),
				a.NewConst(2),
				a.NewFCall(a.NewOp("*"), a.NewConst(3), a.NewConst(4)),
			),
		),
	}, {
		name:  "parse numeric expressions with brackets",
		input: "(2 + 4) * 3",
		want: a.NewBlock(
			a.NewFCall(a.NewOp("*"),
				a.NewFCall(a.NewOp("+"), a.NewConst(2), a.NewConst(4)),
				a.NewConst(3),
			),
		)}, {
		name:  "parse id expressions with brackets",
		input: "(c % 2) == 0",
		want: a.NewBlock(
			a.NewFCall(a.NewOp("=="),
				a.NewFCall(a.NewOp("%"), a.NewId("c"), a.NewConst(2)),
				a.NewConst(0),
			),
		)}, {
		name:  "parse == operator",
		input: "2 == 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp("=="), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse != operator",
		input: "a != 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp("!="), a.NewId("a"), a.NewConst(3))),
	}, {
		name:  "parse < operator",
		input: "2 < 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp("<"), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse > operator",
		input: "2 > 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp(">"), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse >= operator",
		input: "2 >= 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp(">="), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse <= operator",
		input: "2 <= 3",
		want:  a.NewBlock(a.NewFCall(a.NewOp("<="), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse || operator",
		input: "true || false",
		want:  a.NewBlock(a.NewFCall(a.NewOp("||"), a.NewConst(true), a.NewConst(false))),
	}, {
		name:  "parse && operator",
		input: "true && false",
		want:  a.NewBlock(a.NewFCall(a.NewOp("&&"), a.NewConst(true), a.NewConst(false))),
	}, {
		name:  "parse if expression",
		input: "if(2 > 3) 1 else 2",
		want: a.NewBlock(a.NewFCall(a.NewOp("if"),
			a.NewFCall(a.NewOp(">"), a.NewConst(2), a.NewConst(3)),
			a.NewConst(1),
			a.NewConst(2),
		)),
	}, {
		name:  "parse if expressions with blocks",
		input: "if(2 > 3) { 1 } else { 2 }",
		want: a.NewBlock(a.NewFCall(a.NewOp("if"),
			a.NewFCall(a.NewOp(">"), a.NewConst(2), a.NewConst(3)),
			a.NewBlock(a.NewConst(1)),
			a.NewBlock(a.NewConst(2)),
		)),
	}, {
		name:  "parse if else if expression",
		input: "if (2 > 3) 1 else if (3 > 4) 2 else 4",
		want: a.NewBlock(a.NewFCall(a.NewOp("if"),
			a.NewFCall(a.NewOp(">"), a.NewConst(2), a.NewConst(3)),
			a.NewConst(1),
			a.NewFCall(a.NewOp("if"),
				a.NewFCall(a.NewOp(">"), a.NewConst(3), a.NewConst(4)),
				a.NewConst(2),
				a.NewConst(4),
			),
		)),
	}, {
		name:  "parse a function call",
		input: "f(1, x, y)",
		want:  a.NewBlock(a.NewFCall(a.NewId("f"), a.NewConst(1), a.NewId("x"), a.NewId("y"))),
	}, {
		name:  "parse a function calls of returned function values",
		input: "f(1, x, y)(2, 3)",
		want:  a.NewBlock(a.NewFCall(a.NewFCall(a.NewId("f"), a.NewConst(1), a.NewId("x"), a.NewId("y")), a.NewConst(2), a.NewConst(3))),
	}, {
		name:  "parse functions as values",
		input: "f((x) => {x + 2}, (y) => {y + 1})",
		want: a.NewBlock(a.NewFCall(a.NewId("f"),
			a.NewFDef(a.NewBlock(a.NewFCall(a.NewOp("+"), a.NewId("x"), a.NewConst(2))), "x"),
			a.NewFDef(a.NewBlock(a.NewFCall(a.NewOp("+"), a.NewId("y"), a.NewConst(1))), "y"),
		)),
	}, {
		name: "parse several assignments",
		input: `
				a = 1 + 2
				b = 2 + 3
				c = 3 + 4
				a + b + c
			`,
		want: a.
			NewBlock(a.NewFCall(a.NewOp("+"), a.NewFCall(a.NewOp("+"), a.NewId("a"), a.NewId("b")), a.NewId("c"))).
			WithAssignment("a", a.NewFCall(a.NewOp("+"), a.NewConst(1), a.NewConst(2))).
			WithAssignment("b", a.NewFCall(a.NewOp("+"), a.NewConst(2), a.NewConst(3))).
			WithAssignment("c", a.NewFCall(a.NewOp("+"), a.NewConst(3), a.NewConst(4))),
	}, {
		name: "parse a function definition",
		input: `
				a = (x, y) => { x + y }
				a(1, 2)
			`,
		want: a.
			NewBlock(a.NewFCall(a.NewId("a"), a.NewConst(1), a.NewConst(2))).
			WithAssignment("a", a.NewFDef(a.NewBlock(a.NewFCall(a.NewOp("+"), a.NewId("x"), a.NewId("y"))), "x", "y")),
	}, {
		name: "parse a nested function definition",
		input: `
			a = (x, y) => {
				b = (z) => { z + 1 }
				x + b(y)
			}
			a(1, 2)
		`,
		want: a.
			NewBlock(a.NewFCall(a.NewId("a"), a.NewConst(1), a.NewConst(2))).
			WithAssignment("a", a.NewFDef(a.
				NewBlock(a.NewFCall(a.NewOp("+"), a.NewId("x"), a.NewFCall(a.NewId("b"), a.NewId("y")))).
				WithAssignment("b", a.NewFDef(a.NewBlock(a.NewFCall(a.NewOp("+"), a.NewId("z"), a.NewConst(1))), "z")),
				"x", "y",
			)),
	}, {
		name: "parse sequential function definitions",
		input: `
			a = (x) => (y) => x + y
			a(1)(2)
		`,
		want: a.
			NewBlock(a.NewFCall(a.NewFCall(a.NewId("a"), a.NewConst(1)), a.NewConst(2))).
			WithAssignment("a", a.NewFDef(a.NewFDef(a.NewFCall(a.NewOp("+"), a.NewId("x"), a.NewId("y")), "y"), "x")),
	}, {
		name: "parse explicit type definitions on functions",
		input: `
			a = (x:int, y:real, z:bool): real => if (b && y > 1.0) x else 0
			a(1, 2.0, true)
		`,
		want: a.
			NewBlock(a.NewFCall(a.NewId("a"), a.NewConst(1), a.NewConst(2.0), a.NewConst(true))).
			WithAssignment("a", a.NewFDef(
				a.NewTypeDecl(types.RealP, a.NewFCall(a.NewOp("if"),
					a.NewFCall(a.NewOp("&&"), a.NewId("b"), a.NewFCall(a.NewOp(">"), a.NewId("y"), a.NewConst(1.0))),
					a.NewId("x"),
					a.NewConst(0))),
				&a.FParam{"x", types.IntP}, &a.FParam{"y", types.RealP}, &a.FParam{"z", types.BoolP},
			)),
	}, {
		name:  "parse explicit type definitions on generic expression",
		input: `((1:int) + (2:bool)):real`,
		want: a.NewBlock(a.NewTypeDecl(types.RealP,
			a.NewFCall(a.NewOp("+"),
				a.NewTypeDecl(types.IntP, a.NewConst(1)),
				a.NewTypeDecl(types.BoolP, a.NewConst(2)),
			),
		)),
	}, {
		name: "parse function type definitions",
		input: `
			a = (x:int, f:(int)=>bool) => if (f(x)) x else 0
			a(1, b)
		`,
		want: a.
			NewBlock(a.NewFCall(a.NewId("a"), a.NewConst(1), a.NewId("b"))).
			WithAssignment("a", a.NewFDef(
				a.NewFCall(a.NewOp("if"), a.NewFCall(a.NewId("f"), a.NewId("x")), a.NewId("x"), a.NewConst(0)),
				&a.FParam{"x", types.IntP}, &a.FParam{"f", types.MakeFunction(types.IntP, types.BoolP)},
			)),
	}, {
		name: "parse structure definitions",
		input: `
			a = (x:int, c)
			a(1, true)
		`,
		want: a.
			NewBlock(a.NewFCall(a.NewId("a"), a.NewConst(1), a.NewConst(true))).
			WithAssignment("a", a.NewStruct(a.StructField{"x", types.IntP}, a.StructField{"c", types.Type{}})),
	}, {
		name:  "parse simple field accessors",
		input: `a.foo.bar`,
		want:  a.NewBlock(a.NewFieldAccessor("bar", a.NewFieldAccessor("foo", a.NewId("a")))),
	}, {
		name:  "parse method calls",
		input: `a.foo(1)`,
		want:  a.NewBlock(a.NewFCall(a.NewFieldAccessor("foo", a.NewId("a")), a.NewConst(1))),
	}, {
		name:  "parse sequential method calls",
		input: `a.foo(1).bar(2)`,
		want: a.NewBlock(a.NewFCall(
			a.NewFieldAccessor("bar", a.NewFCall(a.NewFieldAccessor("foo", a.NewId("a")), a.NewConst(1))),
			a.NewConst(2))),
	}, {
		name: "parse named types in functions",
		input: `
			a = (x: A) => x
			a(b)
		`,
		want: a.
			NewBlock(a.NewFCall(a.NewId("a"), a.NewId("b"))).
			WithAssignment("a", a.NewFDef(a.NewId("x"), &a.FParam{"x", types.MakeNamed("A")})),
	}, {
		name: "parse typed constants",
		input: `a:int = 5
				a
			`,
		want: a.
			NewBlock(a.NewId("a")).
			WithAssignment("a", a.NewTypeDecl(types.IntP, a.NewConst(5))),
	}, {
		name: "parse interface bindings to a primitive",
		input: `a:int ~> {
					add = (b) => a + b
				}
				a.add(4)
		`,
		want: a.NewBlock(a.NewFCall(a.NewFieldAccessor("add", a.NewId("a")), a.NewConst(4))).
			WithInterface("a", types.IntP, a.NewDefinitions().WithAssignment(
				"add", a.NewFDef(a.NewFCall(a.NewOp("+"), a.NewId("a"), a.NewId("b")), "b"),
			)),
	}}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			prog, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			got := prog.ToAst()
			ok, err := deepdiff.DeepDiff(got, tt.want)
			if !ok {
				t.Error(err)
			}
		})
	}
}
