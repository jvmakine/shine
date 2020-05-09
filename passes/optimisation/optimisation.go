package optimisation

import "github.com/jvmakine/shine/ast"

func Optimise(exp *ast.Exp) {
	visited := map[*ast.Exp]bool{}
	exp.Crawl(func(v *ast.Exp, _ *ast.CrawlContext) {
		visited[v] = true
	})
	exp.Visit(func(v *ast.Exp) {
		if v.Block != nil {
			seen := []*ast.Assign{}
			for _, a := range v.Block.Assignments {
				if visited[a.Value] {
					seen = append(seen, a)
				}
			}
			v.Block.Assignments = seen
		}
	})
}
