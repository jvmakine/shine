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
			a.Value.Visit(f)
		}
		a.Block.Value.Visit(f)
	}
}

type ccontext struct {
	parent *ccontext
	ids    map[string]*Exp
}

func (c *ccontext) resolve(id string) (*Exp, *ccontext) {
	if c.ids[id] != nil {
		return c.ids[id], c
	} else if c.parent != nil {
		return c.parent.resolve(id)
	}
	return nil, nil
}

func (a *Exp) crawl(f func(p *Exp), ctx *ccontext, visited *map[*Exp]bool) {
	if (*visited)[a] {
		return
	}
	(*visited)[a] = true
	f(a)
	if a.Block != nil {
		sub := &ccontext{ids: map[string]*Exp{}, parent: ctx}
		for _, a := range a.Block.Assignments {
			sub.ids[a.Name] = a.Value
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

func (a *Exp) Crawl(f func(p *Exp)) {
	visited := map[*Exp]bool{}
	a.crawl(f, &ccontext{ids: map[string]*Exp{}}, &visited)
}
