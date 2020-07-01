package optimisation

import (
	. "github.com/jvmakine/shine/ast"
)

func DeadCodeElimination(exp Expression) {
	visited, _ := CrawlBefore(exp, func(v Ast, _ *VisitContext) error {
		return nil
	})
	VisitBefore(exp, func(v Ast, _ *VisitContext) error {
		if b, ok := v.(*Block); ok {
			newDef := NewDefinitions()
			for k, a := range b.Def.Assignments {
				if visited[a] {
					newDef.Assignments[k] = a
				}
			}
			for k, a := range b.Def.TypeDefs {
				if visited[a] {
					newDef.TypeDefs[k] = a
				}
			}
			b.Def = newDef
		}
		return nil
	})
}
