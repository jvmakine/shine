package typeinference

import (
	"errors"

	. "github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func typeConstant(constant *Const) {
	if constant.Int != nil {
		constant.ConstType = Int
	} else if constant.Bool != nil {
		constant.ConstType = Bool
	} else if constant.Real != nil {
		constant.ConstType = Real
	} else if constant.String != nil {
		constant.ConstType = String
	} else {
		panic("invalid const")
	}
}

func typeId(id *Id, ctx *VisitContext, unifier Substitutions) error {
	name := id.Name
	defin := ctx.DefinitionOf(name)
	if ctx.Path()[name] {
		id.IdType = NewVariable()
	} else if defin != nil && ctx.DefinitionOf(name).Assignments[name] != nil {
		ref := ctx.DefinitionOf(name).Assignments[name]
		if _, ok := ref.(*FDef); ok {
			id.IdType = ref.Type().Copy(NewTypeCopyCtx())
		} else if _, ok := ref.(*Struct); ok {
			id.IdType = ref.Type().Copy(NewTypeCopyCtx())
		} else {
			id.IdType = ref.Type()
		}
	} else if p := ctx.ParamOf(name); p != nil {
		id.IdType = p.ParamType
	} else {
		if name == "$" {
			inter := ctx.Interface()
			if inter == nil {
				panic("$ id outside of an interface")
			}
			id.IdType = inter.InterfaceType
		} else {
			return errors.New("undefined id " + name)
		}
	}
	return nil
}

func initialiseVariables(exp Expression) error {
	return VisitBefore(exp, func(v Ast, ctx *VisitContext) error {
		if c, ok := v.(*Const); ok {
			typeConstant(c)
		} else if d, ok := v.(*FDef); ok {
			for _, p := range d.Params {
				name := p.Name
				if ctx.DefinitionOf(name) != nil || ctx.ParamOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
				if p.ParamType == nil {
					p.ParamType = NewVariable()
				}
			}
		} else if c, ok := v.(*FCall); ok {
			c.CallType = NewVariable()
		} else if o, ok := v.(*Op); ok {
			o.OpType = NewVariable()
		} else if b, ok := v.(*Block); ok {
			err := b.CheckValueCycles()
			if err != nil {
				return err
			}
		} else if d, ok := v.(*Definitions); ok {
			for name, value := range d.Assignments {
				if s, ok := value.(*Struct); ok {
					ts := make([]Type, len(s.Fields)+1)
					sf := make([]Named, len(s.Fields))
					for i, v := range s.Fields {
						typ := v.Type
						if typ == nil {
							typ = NewVariable()
						}
						v.Type = typ

						ts[i] = typ
						sf[i] = NewNamed(v.Name, typ)
					}

					stru := NewNamed(name, NewStructure(sf...))
					ts[len(s.Fields)] = stru

					s.StructType = stru
				}
			}
			for _, in := range d.Interfaces {
				if in.InterfaceType == nil {
					varit := NewVariable()
					in.InterfaceType = varit
					for _, d := range in.Definitions.Assignments {
						_, isDef := d.(*FDef)
						if !isDef {
							return errors.New("only methods are supported in interfaces")
						}
					}
				}
			}
		} else if fa, ok := v.(*FieldAccessor); ok {
			fa.FAType = NewVariable()
		}
		return nil
	})
}

func resolveNamed(name string, ctx *VisitContext) (Type, error) {
	defin := ctx.DefinitionOf(name)
	if defin == nil {
		return nil, errors.New("type " + name + " is undefined")
	}
	value := defin.Assignments[name]
	if value == nil {
		return nil, errors.New(name + " is not a type")
	}
	struc, ok := value.(*Struct)
	if !ok {
		return nil, errors.New(name + " is not a type")
	}
	fs := make([]Named, len(struc.Fields))
	for i, f := range struc.Fields {
		if f.Type == nil {
			fs[i] = NewNamed(f.Name, NewVariable())
		} else {
			fs[i] = NewNamed(f.Name, f.Type)
		}
	}
	return NewStructure(fs...), nil
}

func rewriteNamed(exp Expression) error {
	return RewriteTypes(exp, func(t Type, ctx *VisitContext) (Type, error) {
		if n, ok := t.(Named); ok {
			typ, err := resolveNamed(n.Name, ctx)
			if err != nil {
				return nil, err
			}
			n.Type = typ
			return n, nil
		}
		return t, nil
	})
}

func interfacesFor(name string, ctx *VisitContext) []*Interface {
	inters := ctx.InterfacesWith(name)
	if len(inters) > 0 {
		res := []*Interface{}
		for _, i := range inters {
			res = append(res, i.Interf)
		}
		return res
	}
	return []*Interface{}
}

func Infer(exp Expression) error {
	if err := initialiseVariables(exp); err != nil {
		return err
	}
	if err := rewriteNamed(exp); err != nil {
		return err
	}
	visited := map[Ast]bool{}
	unifier := MakeSubstitutions()
	crawler := func(v Ast, ctx *VisitContext) error {
		if visited[v] {
			return nil
		}
		visited[v] = true
		if b, ok := v.(*Branch); ok {
			if err := unifier.Add(b.Condition.Type(), Bool, ctx); err != nil {
				return err
			}
			if err := unifier.Add(b.True.Type(), b.False.Type(), ctx); err != nil {
				return err
			}
		} else if b, ok := v.(*Block); ok {
			ConvertTypes(b.Value, unifier)
			for _, a := range b.Def.Assignments {
				ConvertTypes(a, unifier)
			}
		} else if i, ok := v.(*Id); ok {
			if err := typeId(i, ctx, unifier); err != nil {
				return err
			}
		} else if o, ok := v.(*Op); ok {
			wantFun := NewFunction(o.Right.Type(), o.OpType)
			strct := NewVariable(NewNamed(o.Name, wantFun))
			if err := unifier.Add(o.Left.Type(), strct, ctx); err != nil {
				return err
			}
		} else if c, ok := v.(*FCall); ok {
			ftype1 := c.MakeFunType()
			ftype2 := c.Function.Type()
			if err := unifier.Add(ftype1, ftype2, ctx); err != nil {
				return err
			}
		} else if d, ok := v.(*FDef); ok {
			ConvertTypes(d, unifier)
		} else if t, ok := v.(*TypeDecl); ok {
			if err := unifier.Add(t.DeclType, t.Exp.Type(), ctx); err != nil {
				return err
			}
		} else if a, ok := v.(*FieldAccessor); ok {
			et := a.Exp.Type()
			strct := NewVariable(NewNamed(a.Field, a.FAType))
			if err := unifier.Add(et, strct, ctx); err != nil {
				return err
			}
		}
		return nil
	}
	// infer interfaces
	err := VisitAfter(exp, func(v Ast, ctx *VisitContext) error {
		if e, ok := v.(*Block); ok {
			for _, in := range e.Def.Interfaces {
				err := in.Definitions.Visit(NullFun, crawler, true, IdRewrite, ctx.WithBlock(e).WithInterface(in))
				if err != nil {
					return err
				}
				in.InterfaceType = unifier.Apply(in.InterfaceType)
				ConvertTypes(in.Definitions, unifier)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// infer used code
	_, err = CrawlAfter(exp, crawler)
	if err != nil {
		return err
	}
	// infer unused code
	err = VisitAfter(exp, func(v Ast, ctx *VisitContext) error {
		if !visited[v] {
			err := crawler(v, ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
