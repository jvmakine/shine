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
			if b := ctx.DefinitionOf(a.Name); b != nil {
				def := b.Assignments[a.Name]
				if d, ok := def.(*FDef); ok {
					_, isFun := def.Type().(Function)
					if isFun && d.Closure != nil {
						for _, c := range d.Closure.Fields {
							closureAt[v][c.Name] = c.Type
						}
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
			result := NewStructure()
			for n, t := range closureAt[v] {
				d := ctx.DefinitionOf(n)
				isFun := false
				if d != nil {
					_, isFun = d.Assignments[n].Type().(Function)
				}
				if !isFun {
					result.Fields = append(result.Fields, NewNamed(n, t))
				}
			}
			a.Closure = &result
		case *FieldAccessor:
			closureAt[v] = closureAt[a.Exp]
		}
		return nil
	})
}
