package typeinference

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/grammar"
	"github.com/jvmakine/shine/types"

	. "github.com/jvmakine/shine/test"
	"github.com/stretchr/testify/require"
)

func TestInfer(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		prg  string
		typ  string
		err  error
	}{{
		name: "infer constant int correctly",
		exp:  IConst(5),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer constant bool correctly",
		exp:  BConst(false),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer assigments in blocks",
		exp:  Block(Assgs{"a": IConst(5)}, Typedefs{}, Bindings{}, Id("a")),
		typ:  "int",
		err:  nil,
	}, {
		name: "infer integer comparisons as boolean",
		exp:  Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op(">"), IConst(1), IConst(2))),
		typ:  "bool",
		err:  nil,
	}, {
		name: "infer if expressions",
		exp:  Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op("if"), BConst(true), IConst(1), IConst(2))),
		typ:  "int",
		err:  nil,
	}, {
		name: "fail on mismatching if expression branches",
		exp:  Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op("if"), BConst(true), IConst(1), BConst(false))),
		typ:  "",
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "fail when adding booleans together",
		exp:  Fcall(Op("+"), BConst(true), BConst(false)),
		typ:  "",
		err:  errors.New("can not unify bool with V1[int|real|string]"),
	}, {
		name: "infer recursive functions",
		exp: Block(
			Assgs{"a": Fdef(Block(
				Assgs{},
				Typedefs{},
				Bindings{},
				Fcall(Op("if"), BConst(false), Id("x"), Fcall(Id("a"), BConst(true)))),
				"x",
			)},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), BConst(false)),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer deeply nested recursive functions",
		exp: Block(
			Assgs{"a": Fdef(Block(
				Assgs{"b": Fdef(Fcall(Id("a"), Id("y")), "y")},
				Typedefs{}, Bindings{},
				Fcall(Op("if"), BConst(false), Id("x"), Fcall(Id("b"), BConst(true)))),
				"x",
			)},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), BConst(false)),
		),
		typ: "bool",
		err: nil,
	}, {
		name: "infer function calls",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op("+"), IConst(1), Id("x"))), "x")},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), IConst(1)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer function parameters",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op("if"), Id("b"), Id("x"), IConst(0))), "x", "b")},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), IConst(1), BConst(true)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on inferred function parameter mismatch",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op("if"), Id("b"), Id("x"), IConst(0))), "x", "b")},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), BConst(true), BConst(true)),
		),
		typ: "",
		err: errors.New("can not unify bool with int"),
	}, {
		name: "unify function return values",
		exp:  Fdef(Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op("if"), BConst(true), Id("x"), Id("x"))), "x"),
		typ:  "(V1)=>V1",
		err:  nil,
	}, {
		name: "fail on recursive values",
		exp:  Block(Assgs{"a": Id("b"), "b": Id("a")}, Typedefs{}, Bindings{}, Id("a")),
		typ:  "",
		err:  errors.New("recursive value: a -> b -> a"),
	}, {
		name: "work on non-recursive values",
		exp: Block(
			Assgs{
				"a": Fcall(Op("+"), Id("b"), Id("c")),
				"b": Fcall(Op("+"), Id("c"), Id("c")),
				"c": Fcall(Op("+"), IConst(1), IConst(2)),
			},
			Typedefs{}, Bindings{},
			Id("a"),
		),
		typ: "int",
		err: nil,
	}, {
		name: "unify one function multiple ways",
		exp: Block(
			Assgs{"a": Fdef(Block(Assgs{}, Typedefs{}, Bindings{}, Fcall(Op("if"), BConst(true), Id("x"), Id("x"))), "x")},
			Typedefs{}, Bindings{},
			Fcall(Op("if"), Fcall(Id("a"), BConst(true)), Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer parameters in block values",
		exp: Block(
			Assgs{},
			Typedefs{}, Bindings{},
			Fdef(Fcall(Op("if"), BConst(true), Block(Assgs{}, Typedefs{}, Bindings{}, Id("x")), IConst(2)), "x"),
		),
		typ: "(int)=>int",
		err: nil,
	}, {
		name: "infer functions as arguments",
		exp: Block(
			Assgs{},
			Typedefs{}, Bindings{},
			Fdef(Fcall(Op("+"), Fcall(Id("x"), BConst(true), IConst(2)), IConst(1)), "x"),
		),
		typ: "((bool,int)=>int)=>int",
		err: nil,
	}, {
		name: "fail to unify functions with wrong number of arguments",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall(Id("x"), IConst(2), IConst(2)), "x"),
				"b": Fdef(Id("x"), "x"),
			},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), Id("b")),
		),
		typ: "",
		err: errors.New("can not unify (V1)=>V1 with (int,int)=>V1"),
	}, {
		name: "infer multiple function arguments",
		exp: Block(
			Assgs{
				"a":  Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
				"b":  Fdef(Fcall(Op("-"), Id("x"), Id("y")), "x", "y"),
				"op": Fdef(Fcall(Id("x"), IConst(1), IConst(2)), "x"),
			},
			Typedefs{}, Bindings{},
			Fcall(Op("+"), Fcall(Id("op"), Id("a")), Fcall(Id("op"), Id("b"))),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer functions as return values",
		exp: Block(
			Assgs{
				"a":  Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
				"b":  Fdef(Fcall(Op("-"), Id("x"), Id("y")), "x", "y"),
				"sw": Fdef(Fcall(Op("if"), Id("x"), Id("a"), Id("b")), "x"),
				"r":  Fdef(Fcall(Id("f"), RConst(1.0), RConst(2.0)), "f"),
			},
			Typedefs{}, Bindings{},
			Fcall(Id("r"), Fcall(Id("sw"), BConst(true))),
		),
		typ: "real",
		err: nil,
	}, {
		name: "infer return values based on Closure",
		exp: Block(
			Assgs{
				"a": Fdef(Block(
					Assgs{"b": Fdef(Fcall(Op("if"), Id("bo"), Id("x"), IConst(2)), "bo")},
					Typedefs{}, Bindings{},
					Fcall(Id("b"), BConst(true)),
				), "x"),
			},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), IConst(1)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on type errors in unused code",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall(Op("if"), BConst(true), Id("x"), IConst(2)), "x"),
				"b": Fdef(Fcall(Op("if"), Id("bo"), RConst(1.0), IConst(2)), "bo"),
			},
			Typedefs{}, Bindings{},
			Fcall(Id("a"), IConst(3)),
		),
		typ: "",
		err: errors.New("can not unify int with real"),
	}, {
		name: "leave free variables to functions",
		exp: Block(
			Assgs{"a": Fdef(Fcall(Op("+"), Id("x"), Id("x")), "x")},
			Typedefs{}, Bindings{},
			Fcall(Op("if"), Fcall(Op("<"), Fcall(Id("a"), RConst(1.0)), RConst(2.0)), Fcall(Id("a"), IConst(1)), IConst(3)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "infer sequential function definitions",
		exp: Block(
			Assgs{"a": Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("y")), "y"), "x")},
			Typedefs{}, Bindings{},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "int",
		err: nil,
	}, {
		name: "fail on parameter redefinitions",
		exp: Block(
			Assgs{"a": Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("x")), "x"), "x")},
			Typedefs{}, Bindings{},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "",
		err: errors.New("redefinition of x"),
	}, {
		name: "fails when function return type contradicts explicit type",
		exp: Block(
			Assgs{"a": Fdef(TDecl(Fcall(Op("+"), Id("x"), Id("x")), types.BoolP), "x")},
			Typedefs{}, Bindings{},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		typ: "",
		err: errors.New("can not unify bool with V1[int|real|string]"),
	}, {
		name: "fail to unify two different named types",
		prg: `
				ai = a(1)
				bi = b(1)
				a :: (a1: int)
				b :: (a1: int)
				if (true) ai else bi
			`,
		err: errors.New("can not unify a with b"),
	}, {
		name: "fail on unknown named type",
		prg:  `a: t`,
		err:  errors.New("type t is undefined"),
	}, {
		name: "fail on unknown named type argument",
		exp: Block(
			Assgs{},
			Typedefs{"A": Struct(ast.StructField{Name: "x", Type: types.MakeNamed("X")}).WithFreeVars("X")},
			Bindings{},
			TDecl(Id("a"), types.MakeNamed("A", types.MakeNamed("Z"))),
		),
		err: errors.New("type Z is undefined"),
	}, {
		name: "work on known named type",
		prg: `
			t :: (x: int)
			a: t
		`,
		typ: "t",
	}, {
		name: "unify recursive types",
		prg: `
				a[X] :: (a1: X)
				a(a(0))
			`,
		typ: "a[a]",
		err: nil,
	}, {
		name: "infer function types from structure fields",
		exp: Block(
			Assgs{},
			Typedefs{},
			Bindings{},
			Fdef(Fcall(Op("+"), Faccess(Id("x"), "a"), IConst(1)), "x"),
		),
		typ: "(V1{a:int})=>int",
		err: nil,
	}, {
		name: "fail to unify free type to two different types",
		exp: Block(
			Assgs{},
			Typedefs{"A": Struct(
				ast.StructField{"a1", types.MakeNamed("X")},
				ast.StructField{"a2", types.MakeNamed("X")},
			).WithFreeVars("X")},
			Bindings{},
			Fcall(Id("A"), IConst(1), RConst(2.0)),
		),
		err: errors.New("can not unify int with real"),
	}, {
		name: "fail on reusing defined type as free type",
		exp: Block(
			Assgs{},
			Typedefs{
				"X": Struct(ast.StructField{"a1", types.IntP}),
				"A": Struct(ast.StructField{"a1", types.MakeNamed("X")}).WithFreeVars("X"),
			},
			Bindings{},
			Fcall(Id("A"), IConst(1)),
		),
		err: errors.New("redefinition of X"),
	}, {
		name: "fail on unused free types",
		exp: Block(
			Assgs{},
			Typedefs{"A": Struct(ast.StructField{"a1", types.MakeNamed("X")}).WithFreeVars("X", "Y")},
			Bindings{},
			Fcall(Id("A"), IConst(1)),
		),
		err: errors.New("unused free type Y"),
	}, {
		name: "fail on incorrect type variable",
		prg: `
					f = (a: A[real]) => a.x
					A[X] :: (x: X)
					f(A(1))
				`,
		err: errors.New("can not unify int with real"),
	}, {
		name: "fail on redefinitions",
		prg: `
					a = 1
					{
						a = 1
						a
					}
				`,
		err: errors.New("redefinition of a"),
	}, {
		name: "infers type parameters in functions",
		prg: `
					S[A] :: (a: (A)=>A)
					S((x:int) => 1.0)
				`,
		err: errors.New("can not unify int with real"),
	}, {
		name: "infer types based on functions with type arguments",
		prg: `
					f: F[int] = (x) => x
					F[A] :: (A) => A
					f
				`,
		typ: "(int)=>int",
	}, {
		name: "infer type variables as arguments",
		prg: `
					f: G[int] = (x) => x
					F[A] :: (A) => A
					G[X] :: F[X]
					f
				`,
		typ: "(int)=>int",
	}, {
		name: "fails on redefinition of type class function",
		prg: `
					Foo[A] :: { f :: (A) => A }
					f = (x) => x
					f
				`,
		err: errors.New("redefinition of f"),
	}, {
		name: "infers type class functions",
		prg: `
					Foo[A] :: { f[B] :: (A,B) => A }
					f
				`,
		typ: "(Foo[V1],V2)=>Foo[V1]",
	}, {
		name: "fails if binding for given type is not found",
		prg: `
					Foo[A] :: { f[B] :: (A,B) => A }
					a = (x) => f(x,5)
					a(1)
				`,
		err: errors.New("can not unify Foo[V1] with int"),
	}, {
		name: "succeeds if binding for given type is found",
		prg: `
					Foo[A] :: { f :: (A) => A }
					Foo[int] -> { f = (x) => x + 1 }
					a = (x) => f(x)
					a(1)
				`,
		typ: "int",
	}, {
		name: "resolve same functions using multiple bindings",
		prg: `
				Foo[A] :: { f :: (A) => A }
				Foo[int] -> { f = (x) => x + 1 }
				Foo[bool] -> { f = (x) => x }
				if (f(true)) f(1) else f(2)
			`,
		typ: "int",
	}, {
		name: "fail to unify structures with incorrect named variables",
		prg: `
				S[A] :: (value: A)
				N[X] :: S[X]
				S(1): N[real]
			`,
		err: errors.New("can not unify int with real"),
	}, {
		name: "infer parametrised named types in TC definitions",
		prg: `
			Functor[F] :: { map[A,B] :: (F[A], (A) => B) => F[B] }
			map
		`,
		typ: "(Functor[V1[V2]],(V2)=>V3)=>Functor[V1[V3]]",
	}, {
		name: "resolve a second order type class",
		prg: `
				S[A] :: (value: A)
				Functor[F] :: { map[A,B] :: (F[A], (A) => B) => F[B] }
				Functor[S] -> { map = (s, f) => S(f(s.value)) }
				map(S(1), (x) => 1.0)
			`,
		typ: "S[real]",
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			exp := tt.exp
			if tt.prg != "" {
				p, err := grammar.Parse(tt.prg)
				require.NoError(t, err)
				e, err := p.ToAst()
				require.NoError(t, err)
				exp = e
			}
			err := Infer(exp)
			if err != nil {
				if !reflect.DeepEqual(err, tt.err) {
					t.Errorf("Infer() error = %v, want %v", err, tt.err)
				}
			} else {
				require.NoError(t, err)
				res := exp.Type()
				if (!res.IsDefined()) && tt.typ != "" {
					t.Errorf("Infer() wrong type = nil, want %v", tt.typ)
				} else if res.IsDefined() && res.Signature() != tt.typ {
					t.Errorf("Infer() wrong type = %v, want %v", res.Signature(), tt.typ)
				}
			}
		})
	}
}
