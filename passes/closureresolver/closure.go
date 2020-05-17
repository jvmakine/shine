package closureresolver

import (
	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func combine(a map[string]Type, b map[string]Type) map[string]Type {
	res := map[string]Type{}
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}

func CollectClosures(exp *ast.Exp) {
	closureAt := map[*ast.Exp]map[string]Type{}
	exp.VisitAfter(func(v *ast.Exp, ctx *ast.VisitContext) error {
		closureAt[v] = map[string]Type{}
		if v.Id != nil {
			closureAt[v] = map[string]Type{v.Id.Name: v.Id.Type}
		} else if v.Call != nil {
			closureAt[v] = map[string]Type{}
			for _, p := range v.Call.Params {
				closureAt[v] = combine(closureAt[v], closureAt[p])
			}
			closureAt[v] = combine(closureAt[v], closureAt[v.Call.Function])
		} else if v.Const != nil {
			closureAt[v] = map[string]Type{}
		} else if v.Block != nil {
			closureAt[v] = map[string]Type{}
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
			for n, t := range closureAt[v] {
				if block := ctx.BlockOf(n); block == nil || !block.Assignments[n].Type().IsFunction() {
					result = append(result, ClosureParam{Name: n, Type: t})
				}
			}
			v.Def.Closure = &result
		}
		return nil
	})
}
