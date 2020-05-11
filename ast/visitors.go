package ast

type VisitContext struct {
	parent *VisitContext
	block  *Block
}

func (c *VisitContext) Block() *Block {
	return c.block
}

func (c *VisitContext) BlockOf(id string) *Block {
	if c.block == nil {
		return nil
	} else if c.block.Assignments[id] != nil {
		return c.block
	} else if c.parent != nil {
		return c.parent.BlockOf(id)
	}
	return nil
}

func (c *VisitContext) NameOf(exp *Exp) string {
	if c.block == nil {
		return ""
	}
	for n, a := range c.block.Assignments {
		if a == exp {
			return n
		}
	}
	if c.parent != nil {
		return c.parent.NameOf(exp)
	}
	return ""
}

func (c *VisitContext) resolve(id string) (*Exp, *VisitContext) {
	if c.block != nil && c.block.Assignments[id] != nil {
		return c.block.Assignments[id], c
	} else if c.parent != nil {
		return c.parent.resolve(id)
	}
	return nil, nil
}

func (a *Exp) Visit(f func(p *Exp, ctx *VisitContext)) {
	a.visit(f, &VisitContext{})
}

func (a *Exp) Crawl(f func(p *Exp, ctx *VisitContext)) {
	visited := map[*Exp]bool{}
	a.crawl(f, &VisitContext{}, &visited)
}

func (a *Exp) crawl(f func(p *Exp, ctx *VisitContext), ctx *VisitContext, visited *map[*Exp]bool) {
	if (*visited)[a] {
		return
	}
	(*visited)[a] = true
	f(a, ctx)
	if a.Block != nil {
		sub := &VisitContext{block: a.Block, parent: ctx}
		a.Block.Value.crawl(f, sub, visited)
	} else if a.Def != nil {
		a.Def.Body.crawl(f, ctx, visited)
	} else if a.Call != nil {
		for _, p := range a.Call.Params {
			p.crawl(f, ctx, visited)
		}
		if r, c := ctx.resolve(a.Call.Name); r != nil {
			r.crawl(f, c, visited)
		}
	} else if a.Id != nil {
		if r, c := ctx.resolve(a.Id.Name); r != nil {
			r.crawl(f, c, visited)
		}
	}
}

func (a *Exp) visit(f func(p *Exp, ctx *VisitContext), ctx *VisitContext) {
	f(a, ctx)
	if a.Block != nil {
		sub := &VisitContext{block: a.Block, parent: ctx}
		for _, a := range a.Block.Assignments {
			a.visit(f, sub)
		}
		a.Block.Value.visit(f, sub)
	} else if a.Def != nil {
		a.Def.Body.visit(f, ctx)
	} else if a.Call != nil {
		for _, p := range a.Call.Params {
			p.visit(f, ctx)
		}
	}
}
