package callresolver

import (
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
	"github.com/stretchr/testify/require"
)

func TestResolveFunctions(t *testing.T) {
	tests := []struct {
		name   string
		before *ast.Exp
		after  *ast.Exp
	}{{
		name: "resolves function signatures based on the call type",
		before: Block(
			Assgs{
				"a": Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
			},
			Fcall(Op("if"),
				Fcall(Id("a"), BConst(true), BConst(true), BConst(false)),
				Fcall(Id("a"), BConst(true), IConst(5), IConst(6)),
				IConst(7)),
		),
		after: Block(
			Assgs{
				"a":                            Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
				"a%%1%%(bool,bool,bool)=>bool": Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
				"a%%1%%(bool,int,int)=>int":    Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
			},
			Fcall(Op("if"),
				Fcall(Id("a%%1%%(bool,bool,bool)=>bool"), BConst(true), BConst(true), BConst(false)),
				Fcall(Id("a%%1%%(bool,int,int)=>int"), BConst(true), IConst(5), IConst(6)),
				IConst(7)),
		),
	}, {
		name: "resolves functions as arguments",
		before: Block(
			Assgs{
				"a": Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"b": Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
			},
			Fcall(Id("a"), Id("b")),
		),
		after: Block(
			Assgs{
				"a":                           Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"b":                           Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
				"a%%1%%((int,int)=>int)=>int": Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"b%%1%%(int,int)=>int":        Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
			},
			Fcall(Id("a%%1%%((int,int)=>int)=>int"), Id("b%%1%%(int,int)=>int")),
		),
	}, {
		name: "resolves anonymous functions",
		before: Block(
			Assgs{
				"a": Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
			},
			Fcall(Id("a"), Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y")),
		),
		after: Block(
			Assgs{
				"a":                           Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"a%%1%%((int,int)=>int)=>int": Fdef(Fcall(Id("f"), IConst(1), IConst(2)), "f"),
				"<anon1>%%1%%(int,int)=>int":  Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
			},
			Fcall(Id("a%%1%%((int,int)=>int)=>int"), Id("<anon1>%%1%%(int,int)=>int")),
		),
	}, {
		name: "resolves simple structures",
		before: Block(
			Assgs{"a": Struct(ast.StructField{"x", types.IntP})},
			Fcall(Id("a"), IConst(1)),
		),
		after: Block(
			Assgs{
				"a":                     Struct(ast.StructField{"x", types.IntP}),
				"a%%1%%(int)=>a{x:int}": Struct(ast.StructField{"x", types.IntP}),
			},
			Fcall(Id("a%%1%%(int)=>a{x:int}"), IConst(1)),
		),
	}, {
		name: "resolves multitype structures",
		before: Block(
			Assgs{"a": Struct(ast.StructField{"x", types.MakeVariable()})},
			Fcall(Id("a"), Fcall(Id("a"), IConst(1))),
		),
		after: Block(
			Assgs{
				"a":                     Struct(ast.StructField{"x", types.MakeVariable()}),
				"a%%1%%(int)=>a{x:int}": Struct(ast.StructField{"x", types.IntP}),
				"a%%1%%(a{x:int})=>a":   Struct(ast.StructField{"x", types.MakeNamed("a")}),
			},
			Fcall(Id("a%%1%%(a{x:int})=>a"), Fcall(Id("a%%1%%(int)=>a{x:int}"), IConst(1))),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeinference.Infer(tt.before)
			ResolveFunctions(tt.before)
			eraseType(tt.after)
			eraseType(tt.before)
			require.Equal(t, tt.before, tt.after)
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
		}
		return nil
	})
}
