package closureresolver

import (
	"reflect"
	"testing"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/callresolver"
	"github.com/jvmakine/shine/passes/optimisation"
	"github.com/jvmakine/shine/passes/typeinference"
	. "github.com/jvmakine/shine/test"
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
)

func TestResolveFunctionDef(tes *testing.T) {
	tests := []struct {
		name string
		exp  *ast.Exp
		want []Closure
	}{{
		name: "resolves empty Closure for function without closure",
		exp: Block(
			Assgs{
				"a": Fdef(Fcall(Op("if"), Id("b"), Id("y"), Id("x")), "b", "y", "x"),
			},
			Fcall(Id("a"), BConst(true), BConst(true), BConst(false)),
		),
		want: []Closure{Closure{}},
	}, {
		name: "resolve closure parameters for function referring to outer ids",
		exp: Block(
			Assgs{
				"a": Fdef(Block(
					Assgs{
						"b": Fdef(Block(
							Assgs{"c": Id("y")},
							Fcall(Op("if"), Id("c"), Id("x"), IConst(2)),
						)),
					},
					Fcall(Id("b")),
				), "x", "y"),
			},
			Fcall(Id("a"), IConst(1), BConst(true)),
		),
		want: []Closure{Closure{}, Closure{
			ClosureParam{Name: "y", Type: types.BoolP},
			ClosureParam{Name: "x", Type: types.IntP}},
		},
	}, /*{
		name: "not include static function references in the closure",
		exp: Block(
			Fcall("a", IConst(1)),
			Assign("a", Fdef(Fcall("b", Id("x"), Id("s")), "x")),
			Assign("b", Fdef(Fcall("+", Fcall("f", Id("y")), IConst(2)), "y", "f")),
			Assign("s", Fdef(Fcall("+", Id("y"), IConst(3)), "y")),
		),
		want: []Closure{Closure{}, Closure{}, Closure{}},
	},*/
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			err := typeinference.Infer(tt.exp)
			if err != nil {
				panic(err)
			}
			callresolver.ResolveFunctions(tt.exp)
			optimisation.Optimise(tt.exp)
			CollectClosures(tt.exp)
			fcat := callresolver.Collect(tt.exp)
			result := collectClosures(&fcat)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Resolve() = %v, want %v", result, tt.want)
			}
		})
	}
}

func collectClosures(cat *callresolver.FCat) []Closure {
	res := []Closure{}
	for _, v := range *cat {
		if v.Closure != nil {
			res = append(res, *v.Closure)
		}
	}
	return res
}
