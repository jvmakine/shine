package optimisation

import "github.com/jvmakine/shine/ast"

func DeadCodeElimination(exp *ast.Exp) {
	visited := map[*ast.Exp]bool{}
	exp.Crawl(func(v *ast.Exp, _ *ast.VisitContext) error {
		visited[v] = true
		return nil
	})
	exp.Visit(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Block != nil {
			seen := map[string]*ast.Exp{}
			for k, a := range v.Block.Defin.Assignments {
				if visited[a] {
					seen[k] = a
				}
			}
			v.Block.Defin.Assignments = seen
		}
		return nil
	})
}
