package closure

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/inferer/callresolver"
	"github.com/jvmakine/shine/inferer/typeinference"
	"github.com/jvmakine/shine/resolved"
	. "github.com/jvmakine/shine/resolved"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
)

func TestResolveFunctionDef(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		want []resolved.Closure
	}{{
		name: "resolves empty Closure for function without closure",
		exp: Block(
			Fcall("a", BConst(true), BConst(true), BConst(false)),
			Assign("a", Fdef(Fcall("if", Id("b"), Id("y"), Id("x")), "b", "y", "x")),
		),
		want: []Closure{Closure{}},
	}, {
		name: "resolve closure parameters for function referring to outer ids",
		exp: Block(
			Fcall("a", IConst(1), BConst(true)),
			Assign("a", Fdef(Block(
				Fcall("b"),
				Assign("b", Fdef(Block(
					Fcall("if", Id("c"), Id("x"), IConst(2)),
					Assign("c", Id("y")),
				))),
			), "x", "y")),
		),
		want: []Closure{Closure{
			ClosureParam{Name: "y", Type: types.BoolP},
			ClosureParam{Name: "x", Type: types.IntP}},
			Closure{},
		},
	}, {
		name: "not include static function references in the closure",
		exp: Block(
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Fcall("b", Id("x"), Id("s")), "x")),
			Assign("b", Fdef(Fcall("+", Fcall("f", Id("y")), IConst(2)), "y", "f")),
			Assign("s", Fdef(Fcall("+", Id("y"), IConst(3)), "y")),
		),
		want: []Closure{Closure{}, Closure{}, Closure{}},
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := typeinference.Infer(tt.exp)
			if err != nil {
				panic(err)
			}
			fcat := callresolver.Resolve(tt.exp)
			result := collectClosures(fcat)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Resolve() = %v, want %v", result, tt.want)
			}
		})
	}
}

func collectClosures(cat *callresolver.FCat) []Closure {
	res := []Closure{}
	for _, v := range *cat {
		if v.Resolved != nil {
			res = append(res, v.Resolved.Closure)
		}
	}
	return res
}
