package callresolver

import (
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/grammar"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
	"github.com/stretchr/testify/require"
)

func TestResolveFunctions(t *testing.T) {
	tests := []struct {
		name  string
		prg   string
		after *ast.Exp
	}{{
		name: "resolves function signatures based on the call type",
		prg: `
			a = (b, y, x) => if (b) y else x
			if (a(true, true, false)) a(true, 5, 6) else 7
		`,
		after: Block(
			Assgs{
				"a":                            Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
				"a%%1%%(bool,bool,bool)=>bool": Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
				"a%%1%%(bool,int,int)=>int":    Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
			},
			Typedefs{}, Bindings{},
			Fcall(Op("if"),
				Fcall(Id("a%%1%%(bool,bool,bool)=>bool"), BConst(true), BConst(true), BConst(false)),
				Fcall(Id("a%%1%%(bool,int,int)=>int"), BConst(true), IConst(5), IConst(6)),
				IConst(7)),
		),
	}, {
		name: "resolves functions as arguments",
		prg: `
			a = (f) => f(1, 2)
			b = (x, y) => x + y
			a(b)
		`,
		after: Block(
			Assgs{
				"a":                           Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"b":                           Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
				"a%%1%%((int,int)=>int)=>int": Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"b%%1%%(int,int)=>int":        Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
			},
			Typedefs{}, Bindings{},
			Fcall(Id("a%%1%%((int,int)=>int)=>int"), Id("b%%1%%(int,int)=>int")),
		),
	}, {
		name: "resolves anonymous functions",
		prg: `
			a = (f) => f(1,2)
			a((x,y) => x + y)
		`,
		after: Block(
			Assgs{
				"a":                           Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"a%%1%%((int,int)=>int)=>int": Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"<anon1>%%1%%(int,int)=>int":  Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
			},
			Typedefs{}, Bindings{},
			Fcall(Id("a%%1%%((int,int)=>int)=>int"), Id("<anon1>%%1%%(int,int)=>int")),
		),
	}, {
		name: "resolves simple structures",
		prg: `
			a :: (x: int)
			a(1)
		`,
		after: Block(
			Assgs{},
			Typedefs{
				"a":              Struct(ast.StructField{"x", types.IntP}),
				"a%%1%%(int)=>a": Struct(ast.StructField{"x", types.IntP}),
			},
			Bindings{},
			Fcall(Id("a%%1%%(int)=>a"), IConst(1)),
		),
	}, {
		name: "resolves multitype structures",
		prg: `
			a[X] :: (x: X)
			a(a(1))
		`,
		after: Block(
			Assgs{},
			Typedefs{
				"a":                   Struct(ast.StructField{"x", types.MakeVariable()}).WithFreeVars("X"),
				"a%%1%%(int)=>a[int]": Struct(ast.StructField{"x", types.IntP}),
				"a%%1%%(a[int])=>a":   Struct(ast.StructField{"x", types.MakeNamed("a")}),
			},
			Bindings{},
			Fcall(Id("a%%1%%(a[int])=>a"), Fcall(Id("a%%1%%(int)=>a[int]"), IConst(1))),
		),
	}, {
		name: "resolves typeclass references",
		prg: `
			A[X] :: { add :: (X,X) => X }
			A[int] -> { add = (a,b) => a + b }
			A[real] -> { add = (a,b) => a + b }
			if (add(1.0,2.0) > 2.0) add(1,2) else 4
		`,
		after: Block(
			Assgs{
				"add%%1%%(int,int)=>int":    Fdef(Fcall(Op("+"), Id("a"), Id("b")), "a", "b"),
				"add%%1%%(real,real)=>real": Fdef(Fcall(Op("+"), Id("a"), Id("b")), "a", "b"),
			},
			Typedefs{"A": &ast.TypeDefinition{
				FreeVariables: []string{"X"},
				VaribleMap:    map[string]types.Type{"X": types.MakeVariable()},
				TypeClass: &ast.TypeClass{
					Functions: map[string]*ast.TypeDefinition{
						"add": {TypeDecl: types.MakeFunction(types.MakeNamed("X"), types.MakeNamed("X"), types.MakeNamed("X"))},
					},
				},
			}},
			Bindings{&ast.TypeBinding{
				Name:       "A",
				Parameters: []types.Type{types.IntP},
				Bindings: map[string]*ast.FDef{
					"add": Fdef(Fcall(Op("+"), Id("a"), Id("b")), "a", "b").Def,
				},
			}, &ast.TypeBinding{
				Name:       "A",
				Parameters: []types.Type{types.RealP},
				Bindings: map[string]*ast.FDef{
					"add": Fdef(Fcall(Op("+"), Id("a"), Id("b")), "a", "b").Def,
				},
			}},
			Fcall(Op("if"), Fcall(Op(">"), Fcall(Id("add%%1%%(real,real)=>real"), RConst(1.0), RConst(2.0)), RConst(2.0)), Fcall(Id("add%%1%%(int,int)=>int"), IConst(1), IConst(2)), IConst(4)),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := grammar.Parse(tt.prg)
			require.NoError(t, err)
			ast, err := p.ToAst()
			require.NoError(t, err)
			typeinference.Infer(ast)
			ResolveFunctions(ast)
			eraseType(tt.after)
			eraseType(ast)
			require.Equal(t, tt.after, ast)
		})
	}
}

func eraseType(e *ast.Exp) {
	e.RewriteTypes(func(t types.Type, ctx *ast.VisitContext) (types.Type, error) {
		return types.IntP, nil
	})
	e.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Block != nil {
			v.Block.ID = 0
			for _, t := range v.Block.TypeDefs {
				t.VaribleMap = nil
				if t.Struct != nil {
					t.Struct.Type = types.IntP
					for _, f := range t.Struct.Fields {
						f.Type = types.IntP
					}
				}
			}
			v.Block.TCFunctions = nil
		}
		return nil
	})
}
