package closureresolver

import (
	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func combine(a map[string]bool, b map[string]bool) map[string]bool {
	res := map[string]bool{}
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}

func CollectClosures(exp *ast.Exp) {
	closureAt := map[*ast.Exp]map[string]bool{}
	exp.VisitAfter(func(v *ast.Exp, ctx *ast.VisitContext) error {
		closureAt[v] = map[string]bool{}
		if v.Id != nil {
			closureAt[v] = map[string]bool{v.Id.Name: true}
		} else if v.Call != nil {
			closureAt[v] = map[string]bool{}
			for _, p := range v.Call.Params {
				closureAt[v] = combine(closureAt[v], closureAt[p])
			}
			closureAt[v] = combine(closureAt[v], closureAt[v.Call.Function])
		} else if v.Const != nil {
			closureAt[v] = map[string]bool{}
		} else if v.Block != nil {
			closureAt[v] = map[string]bool{}
			assigns := map[string]bool{}
			for n := range v.Block.Assignments {
				assigns[n] = true
			}
			for _, a := range v.Block.Assignments {
				for k, b := range closureAt[a] {
					if !assigns[k] {
						closureAt[v][k] = b
					}
				}
			}
			for k, b := range closureAt[v.Block.Value] {
				if !assigns[k] {
					closureAt[v][k] = b
				}
			}
		} else if v.Def != nil {
			params := map[string]bool{}
			for _, n := range v.Def.Params {
				params[n.Name] = true
			}
			for k, b := range closureAt[v.Def.Body] {
				if !params[k] {
					closureAt[v][k] = b
				}
			}
			result := Closure{}
			for n := range closureAt[v] {
				if block := ctx.BlockOf(n); block == nil || !block.Assignments[n].Type().IsFunction() {
					t := ctx.TypeOf(n)
					result = append(result, ClosureParam{Name: n, Type: t})
				}
			}
			v.Def.Closure = &result
		}
		return nil
	})
}
