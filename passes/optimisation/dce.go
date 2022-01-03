package optimisation

import "github.com/jvmakine/shine/ast"

func DeadCodeElimination(exp *ast.Exp) {
	visited := map[*ast.Exp]bool{}
	visitedT := map[*ast.TypeDefinition]bool{}
	exp.Visit(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Block != nil {
			v.Block.TypeBindings = nil
		}
		return nil
	})
	exp.Crawl(func(v *ast.Exp, c *ast.VisitContext) error {
		if v.Id != nil {
			n := v.Id.Name
			b := c.BlockOf(n)
			if b != nil && b.TypeDefs[n] != nil {
				visitedT[b.TypeDefs[n]] = true
			}
		}
		visited[v] = true
		return nil
	})
	exp.Visit(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Block != nil {
			seen := map[string]*ast.Exp{}
			seenT := map[string]*ast.TypeDefinition{}
			for k, a := range v.Block.Assignments {
				if visited[a] {
					seen[k] = a
				}
			}
			for k, a := range v.Block.TypeDefs {
				if visitedT[a] {
					seenT[k] = a
				}
			}
			v.Block.Assignments = seen
			v.Block.TypeDefs = seenT
		}
		return nil
	})
}
