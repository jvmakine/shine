package ast

import (
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
)

var global = []*Interface{
	&Interface{InterfaceType: Int, Definitions: &Definitions{Assignments: map[string]Expression{
		">":  NewPrimitiveOp("int_>", Bool, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"<":  NewPrimitiveOp("int_<", Bool, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"==": NewPrimitiveOp("int_==", Bool, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"!=": NewPrimitiveOp("int_!=", Bool, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"+":  NewPrimitiveOp("int_+", Int, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"-":  NewPrimitiveOp("int_-", Int, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"*":  NewPrimitiveOp("int_*", Int, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"/":  NewPrimitiveOp("int_/", Int, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
		"%":  NewPrimitiveOp("int_%", Int, NewId("$").WithType(Int), NewId("$2").WithType(Int)),
	}}},
	&Interface{InterfaceType: Real, Definitions: &Definitions{Assignments: map[string]Expression{
		">":  NewPrimitiveOp("real_>", Bool, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
		"<":  NewPrimitiveOp("real_<", Bool, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
		"==": NewPrimitiveOp("real_==", Bool, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
		"!=": NewPrimitiveOp("real_!=", Bool, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
		"+":  NewPrimitiveOp("real_+", Real, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
		"-":  NewPrimitiveOp("real_-", Real, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
		"*":  NewPrimitiveOp("real_*", Real, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
		"/":  NewPrimitiveOp("real_/", Real, NewId("$").WithType(Real), NewId("$2").WithType(Real)),
	}}},
	&Interface{InterfaceType: String, Definitions: &Definitions{Assignments: map[string]Expression{
		"==": NewPrimitiveOp("string_==", Bool, NewId("$").WithType(String), NewId("$2").WithType(String)),
		"!=": NewPrimitiveOp("string_!=", Bool, NewId("$").WithType(String), NewId("$2").WithType(String)),
		"+":  NewPrimitiveOp("string_+", String, NewId("$").WithType(String), NewId("$2").WithType(String)),
	}}},
	&Interface{InterfaceType: Bool, Definitions: &Definitions{Assignments: map[string]Expression{
		"==": NewPrimitiveOp("bool_==", Bool, NewId("$").WithType(Bool), NewId("$2").WithType(Bool)),
		"!=": NewPrimitiveOp("bool_!=", Bool, NewId("$").WithType(Bool), NewId("$2").WithType(Bool)),
	}}},
}

type GlobalVCtx struct {
	visited map[Ast]bool
}

type VisitContext struct {
	parent     *VisitContext
	defin      *Definitions
	interf     *Interface
	def        *FDef
	assignment string
	global     *GlobalVCtx
}

type VisitFunc = func(p Ast, ctx *VisitContext) error
type RewriteFunc = func(from Ast, ctx *VisitContext) Ast

func IdRewrite(from Ast, ctx *VisitContext) Ast {
	return from
}

func (c *VisitContext) WithBlock(b *Block) *VisitContext {
	return &VisitContext{parent: c, defin: b.Def, def: c.def, assignment: c.assignment, interf: c.interf, global: c.global}
}

func (c *VisitContext) WithDefinitions(d *Definitions) *VisitContext {
	return &VisitContext{parent: c, defin: d, def: c.def, assignment: c.assignment, interf: c.interf, global: c.global}
}

func (c *VisitContext) WithDef(d *FDef) *VisitContext {
	return &VisitContext{parent: c, defin: c.defin, def: d, assignment: c.assignment, interf: c.interf, global: c.global}
}

func (c *VisitContext) WithAssignment(a string) *VisitContext {
	return &VisitContext{parent: c, defin: c.defin, def: c.def, assignment: a, interf: c.interf, global: c.global}
}

func (c *VisitContext) WithInterface(i *Interface) *VisitContext {
	return &VisitContext{parent: c, defin: c.defin, def: c.def, assignment: c.assignment, interf: i, global: c.global}
}

func NewVisitCtx() *VisitContext {
	gctx := &GlobalVCtx{visited: map[Ast]bool{}}
	defin := &Definitions{
		Interfaces: global,
	}
	return &VisitContext{parent: nil, defin: defin, def: nil, assignment: "", interf: nil, global: gctx}
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

func (c *VisitContext) Definitions() *Definitions {
	return c.defin
}

func (c *VisitContext) Interface() *Interface {
	return c.interf
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

func (c *VisitContext) DefinitionOf(id string) *Definitions {
	if c.defin == nil {
		return nil
	} else if c.defin.Assignments[id] != nil {
		return c.defin
	} else if c.parent != nil {
		return c.parent.DefinitionOf(id)
	}
	return nil
}

func (c *VisitContext) NameOf(exp Expression) string {
	if c.defin == nil {
		return ""
	}
	for n, a := range c.defin.Assignments {
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
	if c.defin != nil && c.defin.Assignments[id] != nil {
		return c.defin.Assignments[id], c
	} else if c.parent != nil {
		return c.parent.Resolve(id)
	}
	return nil, nil
}

type IResult struct {
	Interf *Interface
	Ctx    *VisitContext
}

func (c *VisitContext) InterfacesWith(id string) []IResult {
	seen := map[*Interface]bool{}
	res := []IResult{}
	if c.defin != nil {
		for _, is := range c.defin.Interfaces {
			seen[is] = true
			if is.Definitions.Assignments[id] != nil {
				res = append(res, IResult{is, c.WithInterface(is)})
			}
		}
	}
	if c.parent != nil {
		pres := c.parent.InterfacesWith(id)
		for _, r := range pres {
			if !seen[r.Interf] {
				res = append(res, r)
			}
		}
	}
	return res
}

func (c VisitContext) StructuralTypeFor(name string, typ types.Type) types.Type {
	ifs := c.InterfacesWith(name)
	if len(ifs) == 0 {
		return nil
	}
	for _, in := range ifs {
		if types.UnifiesWith(in.Interf.InterfaceType, typ, c) {
			return NewFunction(in.Interf.InterfaceType, in.Interf.Definitions.Assignments[name].Type())
		}
	}
	return nil
}

func NullFun(_ Ast, _ *VisitContext) error {
	return nil
}

func VisitBefore(a Ast, f VisitFunc) error {
	return a.Visit(f, NullFun, false, IdRewrite, NewVisitCtx())
}

func VisitAfter(a Ast, f VisitFunc) error {
	return a.Visit(NullFun, f, false, IdRewrite, NewVisitCtx())
}

func CrawlBefore(a Ast, f VisitFunc) (map[Ast]bool, error) {
	res := map[Ast]bool{}
	err := a.Visit(func(v Ast, ctx *VisitContext) error {
		res[v] = true
		return f(v, ctx)
	}, NullFun, true, IdRewrite, NewVisitCtx())
	return res, err
}

func CrawlAfter(a Ast, f VisitFunc) (map[Ast]bool, error) {
	res := map[Ast]bool{}
	err := a.Visit(NullFun, func(v Ast, ctx *VisitContext) error {
		res[v] = true
		return f(v, ctx)
	}, true, IdRewrite, NewVisitCtx())
	return res, err
}

func RewriteTypes(a Ast, f func(t types.Type, ctx *VisitContext) (types.Type, error)) error {
	return a.Visit(NullFun, func(v Ast, ctx *VisitContext) error {
		if op, ok := v.(*Op); ok {
			t, err := f(op.OpType, ctx)
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
	}, false, IdRewrite, NewVisitCtx())
}

func ConvertTypes(e Ast, s types.Substitutions) {
	RewriteTypes(e, func(t types.Type, ctx *VisitContext) (types.Type, error) {
		return s.Apply(t), nil
	})
}
