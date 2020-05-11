package ast

func (a *Exp) Visit(f func(p *Exp)) {
	f(a)
	if a.Def != nil {
		a.Def.Body.Visit(f)
	} else if a.Call != nil {
		for _, p := range a.Call.Params {
			p.Visit(f)
		}
	} else if a.Block != nil {
		for _, a := range a.Block.Assignments {
			a.Visit(f)
		}
		a.Block.Value.Visit(f)
	}
}

type CrawlContext struct {
	parent *CrawlContext
	block  *Block
}

func (c *CrawlContext) Block() *Block {
	return c.block
}

func (c *CrawlContext) BlockOf(id string) *Block {
	if c.block == nil {
		return nil
	} else if c.block.Assignments[id] != nil {
		return c.block
	} else if c.parent != nil {
		return c.parent.BlockOf(id)
	}
	return nil
}

func (c *CrawlContext) NameOf(exp *Exp) string {
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

func (c *CrawlContext) resolve(id string) (*Exp, *CrawlContext) {
	if c.block != nil && c.block.Assignments[id] != nil {
		return c.block.Assignments[id], c
	} else if c.parent != nil {
		return c.parent.resolve(id)
	}
	return nil, nil
}

func (a *Exp) crawl(f func(p *Exp, ctx *CrawlContext), ctx *CrawlContext, visited *map[*Exp]bool) {
	if (*visited)[a] {
		return
	}
	(*visited)[a] = true
	f(a, ctx)
	if a.Block != nil {
		sub := &CrawlContext{block: a.Block, parent: ctx}
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

func (a *Exp) Crawl(f func(p *Exp, ctx *CrawlContext)) {
	visited := map[*Exp]bool{}
	a.crawl(f, &CrawlContext{}, &visited)
}
