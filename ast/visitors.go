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
	ids    map[string]*Exp
	blocks map[string]*Block
}

func (c *CrawlContext) BlockOf(id string) *Block {
	if c.ids[id] != nil {
		return c.blocks[id]
	} else if c.parent != nil {
		return c.parent.BlockOf(id)
	}
	return nil
}

func (c *CrawlContext) resolve(id string) (*Exp, *CrawlContext) {
	if c.ids[id] != nil {
		return c.ids[id], c
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
		sub := &CrawlContext{ids: map[string]*Exp{}, blocks: map[string]*Block{}, parent: ctx}
		for k, as := range a.Block.Assignments {
			sub.ids[k] = as
			sub.blocks[k] = a.Block
		}
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
	a.crawl(f, &CrawlContext{ids: map[string]*Exp{}, blocks: map[string]*Block{}}, &visited)
}
