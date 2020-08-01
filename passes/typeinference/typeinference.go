package typeinference

import (
	"errors"

	. "github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func typeConstant(constant *Const, unifier Substitutions, ctx UnificationCtx) error {
	if constant.Int != nil {
		return unifier.Add(constant.ConstType, Int, ctx)
	} else if constant.Bool != nil {
		return unifier.Add(constant.ConstType, Bool, ctx)
	} else if constant.Real != nil {
		return unifier.Add(constant.ConstType, Real, ctx)
	} else if constant.String != nil {
		return unifier.Add(constant.ConstType, String, ctx)
	}
	panic("invalid const")
}

func typeId(id *Id, ctx *VisitContext, unifier Substitutions) error {
	name := id.Name
	defin := ctx.DefinitionOf(name)
	if ctx.Path()[name] {
		return nil
	} else if defin != nil && ctx.DefinitionOf(name).Assignments[name] != nil {
		ref := ctx.DefinitionOf(name).Assignments[name]
		rtyp := ref.Type()
		if _, ok := ref.(*FDef); ok {
			return unifier.Add(rtyp.Copy(NewTypeCopyCtx()), id.IdType, ctx)
		} else if _, ok := ref.(*Struct); ok {
			return unifier.Add(rtyp.Copy(NewTypeCopyCtx()), id.IdType, ctx)
		} else {
			return unifier.Add(rtyp, id.IdType, ctx)
		}
	} else if p := ctx.ParamOf(name); p != nil {
		return unifier.Add(p.ParamType, id.IdType, ctx)
	} else {
		if name == "$" {
			inter := ctx.Interface()
			if inter == nil {
				panic("$ id outside of an interface")
			}
			return unifier.Add(inter.InterfaceType, id.IdType, ctx)
		} else {
			return errors.New("undefined id " + name)
		}
	}
	return nil
}

func initialiseVariables(exp Expression) error {
	return VisitBefore(exp, func(v Ast, ctx *VisitContext) error {
		if id, ok := v.(*Id); ok {
			id.IdType = NewVariable()
		} else if c, ok := v.(*Const); ok {
			c.ConstType = NewVariable()
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
		if c, ok := v.(*Const); ok {
			err := typeConstant(c, unifier, ctx)
			if err != nil {
				return err
			}
		} else if b, ok := v.(*Branch); ok {
			if err := unifier.Add(b.Condition.Type(), Bool, ctx); err != nil {
				return err
			}
			if err := unifier.Add(b.True.Type(), b.False.Type(), ctx); err != nil {
				return err
			}
		} else if i, ok := v.(*Id); ok {
			if err := typeId(i, ctx, unifier); err != nil {
				return err
			}
		} else if o, ok := v.(*Op); ok {
			nv1 := NewVariable()
			nv2 := NewVariable()
			wantFun := NewFunction(nv1, o.OpType)
			strct := NewVariable(NewNamed(o.Name, wantFun))
			if err := unifier.Add(o.Left.Type(), strct, ctx); err != nil {
				return err
			}
			if err := unifier.Add(nv1, o.Right.Type(), ctx); err != nil {
				return err
			}
			if err := unifier.Add(nv2, o.Left.Type(), ctx); err != nil {
				return err
			}
		} else if c, ok := v.(*FCall); ok {
			ftype1 := c.MakeFunType()
			ftype2 := c.Function.Type()
			if err := unifier.Add(ftype1, ftype2, ctx); err != nil {
				return err
			}
		} else if t, ok := v.(*TypeDecl); ok {
			if err := unifier.Add(t.DeclType, t.Exp.Type(), ctx); err != nil {
				return err
			}
		} else if a, ok := v.(*FieldAccessor); ok {
			nv := NewVariable()
			et := a.Exp.Type()
			strct := NewVariable(NewNamed(a.Field, a.FAType))
			if err := unifier.Add(nv, strct, ctx); err != nil {
				return err
			}
			if err := unifier.Add(nv, et, ctx); err != nil {
				return err
			}
		}
		return nil
	}

	crawlerWithDefRewrite := func(v Ast, ctx *VisitContext) error {
		err := crawler(v, ctx)
		if err != nil {
			return err
		}
		// When crawling through the AST, we need to convert the function types
		// so that when they are copied to the caller, they have the right structure
		if def, ok := v.(*FDef); ok && ctx.IsActiveAssignment(def) {
			ConvertTypes(def, unifier)
		}
		return nil
	}

	// infer interfaces and function definitions
	err := VisitAfter(exp, func(v Ast, ctx *VisitContext) error {
		if e, ok := v.(*Block); ok {
			for _, in := range e.Def.Interfaces {
				err := in.Definitions.Visit(NullFun, crawlerWithDefRewrite, true, IdRewrite, ctx.WithBlock(e).WithInterface(in))
				if err != nil {
					return err
				}
				in.InterfaceType = unifier.Apply(in.InterfaceType)
			}
			for name, ass := range e.Def.Assignments {
				if def, ok := ass.(*FDef); ok {
					err := ass.Visit(NullFun, crawlerWithDefRewrite, true, IdRewrite, ctx.WithBlock(e).WithAssignment(name).WithDef(def))
					if err != nil {
						return err
					}
				}
			}
			ConvertTypes(e, unifier)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Use a separate unifier for actual values
	unifier = MakeSubstitutions()

	// infer used code
	_, err = CrawlAfter(exp, crawler)
	if err != nil {
		return err
	}

	// infer unused code
	err = VisitAfter(exp, crawler)
	if err != nil {
		return err
	}

	ConvertTypes(exp, unifier)
	return nil
}
