package ast

import "github.com/jvmakine/shine/types"

type GlobalVCtx struct {
	visited map[Expression]bool
}

type VisitContext struct {
	parent     *VisitContext
	block      *Block
	def        *FDef
	assignment string
	global     *GlobalVCtx
}

type VisitFunc = func(p Ast, ctx *VisitContext) error

func (c *VisitContext) WithBlock(b *Block) *VisitContext {
	return &VisitContext{parent: c, block: b, def: c.def, assignment: c.assignment, global: c.global}
}

func (c *VisitContext) WithDef(d *FDef) *VisitContext {
	return &VisitContext{parent: c, block: c.block, def: d, assignment: c.assignment, global: c.global}
}

func (c *VisitContext) WithAssignment(a string) *VisitContext {
	return &VisitContext{parent: c, block: c.block, def: c.def, assignment: a, global: c.global}
}

func NewVisitCtx() *VisitContext {
	global := &GlobalVCtx{visited: map[Expression]bool{}}
	return &VisitContext{parent: nil, block: nil, def: nil, assignment: "", global: global}
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
	} else if c.block.Def.Assignments[id] != nil {
		return c.block
	} else if c.block.Def.TypeDefs[id] != nil {
		return c.block
	} else if c.parent != nil {
		return c.parent.BlockOf(id)
	}
	return nil
}

func (c *VisitContext) NameOf(exp Expression) string {
	if c.block == nil {
		return ""
	}
	for n, a := range c.block.Def.Assignments {
		if a == exp {
			return n
		}
	}
	if c.parent != nil {
		return c.parent.NameOf(exp)
	}
	return ""
}

func (c *VisitContext) Resolve(id string) (Expression, *VisitContext) {
	if c.block != nil && c.block.Def.Assignments[id] != nil {
		return c.block.Def.Assignments[id], c
	} else if c.parent != nil {
		return c.parent.Resolve(id)
	}
	return nil, nil
}

func NullFun(_ Ast, _ *VisitContext) error {
	return nil
}

func VisitBefore(a Ast, f VisitFunc) error {
	return a.Visit(f, NullFun, false, NewVisitCtx())
}

func VisitAfter(a Ast, f VisitFunc) error {
	return a.Visit(NullFun, f, false, NewVisitCtx())
}

func CrawlAfter(a Ast, f VisitFunc) (map[Ast]bool, error) {
	res := map[Ast]bool{}
	err := a.Visit(NullFun, func(v Ast, ctx *VisitContext) error {
		res[v] = true
		return f(v, ctx)
	}, true, NewVisitCtx())
	return res, err
}

func RewriteTypes(a Ast, f func(t types.Type, ctx *VisitContext) (types.Type, error)) error {
	return a.Visit(NullFun, func(v Ast, ctx *VisitContext) error {
		if op, ok := v.(*Op); ok {
			t, err := f(op.Type(), ctx)
			if err != nil {
				return err
			}
			op.OpType = t
		} else if id, ok := v.(*Id); ok {
			t, err := f(id.Type(), ctx)
			if err != nil {
				return err
			}
			id.IdType = t
		} else if c, ok := v.(*Const); ok {
			t, err := f(c.ConstType, ctx)
			if err != nil {
				return err
			}
			c.ConstType = t
		} else if d, ok := v.(*TypeDecl); ok {
			t, err := f(d.DeclType, ctx)
			if err != nil {
				return err
			}
			d.DeclType = t
		} else if a, ok := v.(*FieldAccessor); ok {
			t, err := f(a.FAType, ctx)
			if err != nil {
				return err
			}
			a.FAType = t
		} else if c, ok := v.(*FCall); ok {
			t, err := f(c.CallType, ctx)
			if err != nil {
				return err
			}
			c.CallType = t
		} else if d, ok := v.(*FDef); ok {
			for _, p := range d.Params {
				t, err := f(p.ParamType, ctx)
				if err != nil {
					return err
				}
				p.ParamType = t
			}
		} else if s, ok := v.(*Struct); ok {
			t, err := f(s.StructType, ctx)
			if err != nil {
				return err
			}
			s.StructType = t
			for _, p := range s.Fields {
				t, err := f(p.Type, ctx)
				if err != nil {
					return err
				}
				p.Type = t
			}
		}
		return nil
	}, false, NewVisitCtx())
}
