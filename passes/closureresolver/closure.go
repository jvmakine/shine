package closureresolver

import (
	. "github.com/jvmakine/shine/ast"
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

func CollectClosures(exp Expression) {
	closureAt := map[Ast]map[string]Type{}
	CrawlAfter(exp, func(v Ast, ctx *VisitContext) error {
		closureAt[v] = map[string]Type{}
		switch a := v.(type) {
		case *Id:
			closureAt[v] = map[string]Type{a.Name: a.IdType}
			if b := ctx.BlockOf(a.Name); b != nil {
				def := b.Def.Assignments[a.Name]
				if d, ok := def.(*FDef); ok && def.Type().IsFunction() && d.Closure != nil {
					for _, c := range d.Closure.Fields {
						closureAt[v][c.Name] = c.Type
					}
				}
			}
		case *FCall:
			closureAt[v] = map[string]Type{}
			for _, p := range a.Params {
				closureAt[v] = combine(closureAt[v], closureAt[p])
			}
			closureAt[v] = combine(closureAt[v], closureAt[a.Function])
		case *Const:
			closureAt[v] = map[string]Type{}
		case *Block:
			closureAt[v] = map[string]Type{}
			assigns := map[string]bool{}
			for n := range a.Def.Assignments {
				assigns[n] = true
			}
			for _, a := range a.Def.Assignments {
				for k, b := range closureAt[a] {
					if !assigns[k] {
						closureAt[v][k] = b
					}
				}
			}
			for k, b := range closureAt[a.Value] {
				if !assigns[k] {
					closureAt[v][k] = b
				}
			}
		case *FDef:
			params := map[string]bool{}
			for _, n := range a.Params {
				params[n.Name] = true
			}
			for k, b := range closureAt[a.Body] {
				if !params[k] {
					closureAt[v][k] = b
				}
			}
			result := Structure{Name: "", Fields: []SField{}}
			for n, t := range closureAt[v] {
				if block := ctx.BlockOf(n); block == nil || !block.Def.Assignments[n].Type().IsFunction() {
					result.Fields = append(result.Fields, SField{Name: n, Type: t})
				}
			}
			a.Closure = &result
		case *FieldAccessor:
			closureAt[v] = closureAt[a.Exp]
		}
		return nil
	})
}
