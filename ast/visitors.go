package ast

import (
	"github.com/jvmakine/shine/types"
)

type VisitContext struct {
	parent     *VisitContext
	block      *Block
	def        *FDef
	binding    *TypeBinding
	assignment string
}

func (ctx *VisitContext) SubBlock(b *Block) *VisitContext {
	return &VisitContext{parent: ctx, block: b, def: ctx.def, binding: ctx.binding}
}

func (ctx *VisitContext) SubBinding(b *TypeBinding) *VisitContext {
	return &VisitContext{parent: ctx, block: ctx.block, def: ctx.def, binding: b}
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

func (c *VisitContext) ParentBlock() *VisitContext {
	currentBlock := c.block
	current := c
	if currentBlock == nil {
		return nil
	}
	for currentBlock == current.block && current != nil {
		current = current.parent
	}
	return current
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
	} else if c.block.Assignments[id] != nil || c.block.TypeDefs[id] != nil || c.block.TCFunctions[id] != nil {
		return c.block
	} else if c.parent != nil {
		return c.parent.BlockOf(id)
	}
	return nil
}

func (c *VisitContext) TypeDef(id string) *TypeDefinition {
	if c.block == nil {
		return nil
	} else if c.block.TypeDefs[id] != nil {
		return c.block.TypeDefs[id]
	} else if c.parent != nil {
		return c.parent.TypeDef(id)
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

func (c *VisitContext) Binding() *TypeBinding {
	return c.binding
}

func (c *VisitContext) resolve(id string) (*Exp, *VisitContext) {
	if c.block != nil && c.block.Assignments[id] != nil {
		return c.block.Assignments[id], c
	} else if c.parent != nil {
		return c.parent.resolve(id)
	}
	return nil, nil
}

func (c *VisitContext) GetBindings() []types.TCBinding {
	result := []types.TCBinding{}
	block := c.Block()
	if block == nil {
		return result
	}
	for _, b := range block.TypeBindings {
		result = append(result, types.TCBinding{
			Name:    b.Name,
			Args:    b.Parameters,
			BlockID: block.ID,
		})
	}
	parent := c.ParentBlock()
	result = append(result, parent.GetBindings()...)
	return result
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
		sub := ctx.SubBlock(a.Block)
		for _, b := range a.Block.TypeBindings {
			for _, def := range b.Bindings {
				exp := &Exp{Def: def}
				if err := exp.crawl(f, l, sub.SubBinding(b), visited); err != nil {
					return err
				}
			}
		}
		if err := a.Block.Value.crawl(f, l, sub, visited); err != nil {
			return err
		}
	} else if a.Def != nil {
		sub := &VisitContext{block: ctx.block, parent: ctx, def: a.Def, binding: ctx.binding}
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
		if err := a.TDecl.Exp.crawl(f, l, ctx, visited); err != nil {
			return err
		}
	} else if a.FAccess != nil {
		if err := a.FAccess.Exp.crawl(f, l, ctx, visited); err != nil {
			return err
		}
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
		for _, b := range a.Block.TypeBindings {
			for _, def := range b.Bindings {
				e := &Exp{Def: def}
				if err := e.visit(f, l, sub); err != nil {
					return err
				}
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
	} else if a.FAccess != nil {
		a.FAccess.Exp.visit(f, l, ctx)
	}
	if err := l(a, ctx); err != nil {
		return err
	}
	return nil
}

func (a *Exp) RewriteTypes(f func(t types.Type, ctx *VisitContext) (types.Type, error)) error {
	return a.VisitAfter(func(v *Exp, ctx *VisitContext) error {
		if v.Op != nil {
			t, err := f(v.Op.Type, ctx)
			if err != nil {
				return err
			}
			v.Op.Type = t
		} else if v.Id != nil {
			t, err := f(v.Id.Type, ctx)
			if err != nil {
				return err
			}
			v.Id.Type = t
		} else if v.Const != nil {
			t, err := f(v.Const.Type, ctx)
			if err != nil {
				return err
			}
			v.Const.Type = t
		} else if v.TDecl != nil {
			t, err := f(v.TDecl.Type, ctx)
			if err != nil {
				return err
			}
			v.TDecl.Type = t
		} else if v.FAccess != nil {
			t, err := f(v.FAccess.Type, ctx)
			if err != nil {
				return err
			}
			v.FAccess.Type = t
		} else if v.Call != nil {
			t, err := f(v.Call.Type, ctx)
			if err != nil {
				return err
			}
			v.Call.Type = t
		} else if v.Def != nil {
			for _, p := range v.Def.Params {
				t, err := f(p.Type, ctx)
				if err != nil {
					return err
				}
				p.Type = t
			}
		} else if v.Block != nil {
			for _, bind := range v.Block.TypeBindings {
				nt := make([]types.Type, len(bind.Parameters))
				for i, t := range bind.Parameters {
					x, err := f(t, ctx.SubBlock(v.Block))
					if err != nil {
						return err
					}
					nt[i] = x
				}
				bind.Parameters = nt
			}
			for _, td := range v.Block.TypeDefs {
				if err := rewriteTypeDef(td, f, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func rewriteTypeDef(td *TypeDefinition, f func(t types.Type, ctx *VisitContext) (types.Type, error), ctx *VisitContext) error {
	if s := td.Struct; s != nil {
		nt, err := f(s.Type, ctx)
		if err != nil {
			return err
		}
		s.Type = nt
	} else if c := td.TypeClass; c != nil {
		for _, d := range c.Functions {
			if err := rewriteTypeDef(d, f, ctx); err != nil {
				return err
			}
		}
	} else if td.TypeDecl.IsDefined() {
		nt, err := f(td.TypeDecl, ctx)
		if err != nil {
			return err
		}
		td.TypeDecl = nt
	}
	return nil
}
