package optimisation

import "github.com/jvmakine/shine/ast"

func Optimise(exp *ast.Exp) {
	visited := map[*ast.Exp]bool{}
	exp.Crawl(func(v *ast.Exp, _ *ast.CrawlContext) {
		visited[v] = true
	})
	exp.Visit(func(v *ast.Exp) {
		if v.Block != nil {
			seen := map[string]*ast.Exp{}
			for k, a := range v.Block.Assignments {
				if visited[a] {
					seen[k] = a
				}
			}
			v.Block.Assignments = seen
		}
	})
}
