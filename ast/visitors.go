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
	a.visit(f, func(_ *Exp, _ *VisitContext) {}, &VisitContext{})
}

func (a *Exp) VisitAfter(f func(p *Exp, ctx *VisitContext)) {
	a.visit(func(_ *Exp, _ *VisitContext) {}, f, &VisitContext{})
}

func (a *Exp) Crawl(f func(p *Exp, ctx *VisitContext)) {
	visited := map[*Exp]bool{}
	a.crawl(f, func(_ *Exp, _ *VisitContext) {}, &VisitContext{}, &visited)
}

func (a *Exp) CrawlAfter(f func(p *Exp, ctx *VisitContext)) {
	visited := map[*Exp]bool{}
	a.crawl(func(_ *Exp, _ *VisitContext) {}, f, &VisitContext{}, &visited)
}

func (a *Exp) crawl(f func(p *Exp, ctx *VisitContext), l func(p *Exp, ctx *VisitContext), ctx *VisitContext, visited *map[*Exp]bool) {
	if (*visited)[a] {
		return
	}
	(*visited)[a] = true
	f(a, ctx)
	if a.Block != nil {
		sub := &VisitContext{block: a.Block, parent: ctx}
		a.Block.Value.crawl(f, l, sub, visited)
	} else if a.Def != nil {
		a.Def.Body.crawl(f, l, ctx, visited)
	} else if a.Call != nil {
		a.Call.Function.crawl(f, l, ctx, visited)
		for _, p := range a.Call.Params {
			p.crawl(f, l, ctx, visited)
		}
	} else if a.Id != nil {
		if r, c := ctx.resolve(a.Id.Name); r != nil {
			r.crawl(f, l, c, visited)
		}
	}
	l(a, ctx)
}

func (a *Exp) visit(f func(p *Exp, ctx *VisitContext), l func(p *Exp, ctx *VisitContext), ctx *VisitContext) {
	f(a, ctx)
	if a.Block != nil {
		sub := &VisitContext{block: a.Block, parent: ctx}
		for _, a := range a.Block.Assignments {
			a.visit(f, l, sub)
		}
		a.Block.Value.visit(f, l, sub)
	} else if a.Def != nil {
		a.Def.Body.visit(f, l, ctx)
	} else if a.Call != nil {
		a.Call.Function.visit(f, l, ctx)
		for _, p := range a.Call.Params {
			p.visit(f, l, ctx)
		}
	}
	l(a, ctx)
}
