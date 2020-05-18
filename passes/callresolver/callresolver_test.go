package callresolver

import (
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
	"github.com/roamz/deepdiff"
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
		name: "combines sequential functions when possible",
		before: Block(
			Assgs{"a": Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("y")), "y"), "x")},
			Fcall(Fcall(Id("a"), IConst(1)), IConst(2)),
		),
		after: Block(
			Assgs{
				"a":                      Fdef(Fdef(Fcall(Op("+"), Id("x"), Id("y")), "y"), "x"),
				"a%c":                    Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
				"a%c%%1%%(int,int)=>int": Fdef(Fcall(Op("+"), Id("x"), Id("y")), "x", "y"),
			},
			Fcall(Id("a%c%%1%%(int,int)=>int"), IConst(2), IConst(1)),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeinference.Infer(tt.before)
			ResolveFunctions(tt.before)

			eraseType(tt.after)
			eraseType(tt.before)
			ok, err := deepdiff.DeepDiff(tt.before, tt.after)
			if !ok {
				t.Error(err)
			}
		})
	}
}

func eraseType(e *ast.Exp) {
	e.Visit(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Id != nil {
			v.Id.Type = types.IntP
		} else if v.Op != nil {
			v.Op.Type = types.IntP
		} else if v.Const != nil {
			v.Const.Type = types.IntP
		} else if v.Call != nil {
			v.Call.Type = types.IntP
		} else if v.Def != nil {
			for _, p := range v.Def.Params {
				p.Type = types.IntP
			}
			v.Def.Closure = nil
		} else if v.Block != nil {
			v.Block.ID = 0
		}
		return nil
	})
}
