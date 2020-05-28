package ast

type VisitContext struct {
	parent     *VisitContext
	block      *Block
	def        *FDef
	assignment string
}

func (c *VisitContext) Path() map[string]bool {
	p := map[string]bool{}
	if c.parent != nil {
		p = c.parent.Path()
	}
	if c.assignment != "" {
		p[c.assignment] = true
	}
	return p
}

func (c *VisitContext) Def() *FDef {
	return c.def
}

func (c *VisitContext) Block() *Block {
	return c.block
}

func (c *VisitContext) ParamOf(id string) *FParam {
	if c.def != nil {
		if p := c.def.ParamOf(id); p != nil {
			return p
		}
	}
	if c.parent != nil {
		return c.parent.ParamOf(id)
	}
	return nil
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

type VisitFunc = func(p *Exp, ctx *VisitContext) error

func nullVisitFun(_ *Exp, _ *VisitContext) error {
	return nil
}

func (a *Exp) Visit(f VisitFunc) error {
	return a.visit(f, nullVisitFun, &VisitContext{})
}

func (a *Exp) VisitAfter(f VisitFunc) error {
	return a.visit(nullVisitFun, f, &VisitContext{})
}

func (a *Exp) Crawl(f VisitFunc) (map[*Exp]bool, error) {
	visited := map[*Exp]bool{}
	return visited, a.crawl(f, nullVisitFun, &VisitContext{}, &visited)
}

func (a *Exp) CrawlAfter(f VisitFunc) (map[*Exp]bool, error) {
	visited := map[*Exp]bool{}
	return visited, a.crawl(nullVisitFun, f, &VisitContext{}, &visited)
}

func (a *Exp) crawl(f VisitFunc, l VisitFunc, ctx *VisitContext, visited *map[*Exp]bool) error {
	if (*visited)[a] {
		return nil
	}
	(*visited)[a] = true
	if err := f(a, ctx); err != nil {
		return err
	}
	if a.Block != nil {
		sub := &VisitContext{block: a.Block, def: ctx.def, parent: ctx}
		if err := a.Block.Value.crawl(f, l, sub, visited); err != nil {
			return err
		}
	} else if a.Def != nil {
		sub := &VisitContext{block: ctx.block, parent: ctx, def: a.Def}
		if err := a.Def.Body.crawl(f, l, sub, visited); err != nil {
			return err
		}
	} else if a.Call != nil {
		if err := a.Call.Function.crawl(f, l, ctx, visited); err != nil {
			return err
		}
		for _, p := range a.Call.Params {
			if err := p.crawl(f, l, ctx, visited); err != nil {
				return err
			}
		}
	} else if a.Id != nil {
		if r, c := ctx.resolve(a.Id.Name); r != nil {
			sub := &VisitContext{assignment: a.Id.Name, block: c.block, def: c.def, parent: c}
			if err := r.crawl(f, l, sub, visited); err != nil {
				return err
			}
		}
	} else if a.TDecl != nil {
		a.TDecl.Exp.crawl(f, l, ctx, visited)
	}
	if err := l(a, ctx); err != nil {
		return err
	}
	return nil
}

func (a *Exp) visit(f VisitFunc, l VisitFunc, ctx *VisitContext) error {
	if err := f(a, ctx); err != nil {
		return err
	}
	if a.Block != nil {
		sub := &VisitContext{block: a.Block, parent: ctx, def: ctx.def}
		for n, a := range a.Block.Assignments {
			ssub := &VisitContext{assignment: n, block: sub.block, def: sub.def, parent: sub}
			if err := a.visit(f, l, ssub); err != nil {
				return err
			}
		}
		if err := a.Block.Value.visit(f, l, sub); err != nil {
			return err
		}
	} else if a.Def != nil {
		sub := &VisitContext{block: ctx.block, parent: ctx, def: a.Def}
		if err := a.Def.Body.visit(f, l, sub); err != nil {
			return err
		}
	} else if a.Call != nil {
		if err := a.Call.Function.visit(f, l, ctx); err != nil {
			return err
		}
		for _, p := range a.Call.Params {
			if err := p.visit(f, l, ctx); err != nil {
				return err
			}
		}
	} else if a.TDecl != nil {
		a.TDecl.Exp.visit(f, l, ctx)
	}
	if err := l(a, ctx); err != nil {
		return err
	}
	return nil
}
