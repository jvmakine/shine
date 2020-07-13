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
	defin := ctx.DefinitionOf(id.Name)
	if ctx.Path()[id.Name] {
		id.IdType = NewVariable()
	} else if defin != nil && ctx.DefinitionOf(id.Name).Assignments[id.Name] != nil {
		ref := ctx.DefinitionOf(id.Name).Assignments[id.Name]
		if _, ok := ref.(*FDef); ok {
			id.IdType = ref.Type().Copy(NewTypeCopyCtx())
		} else {
			id.IdType = ref.Type()
		}
	} else if p := ctx.ParamOf(id.Name); p != nil {
		id.IdType = p.ParamType
	} else {
		if id.Name == "$" {
			inter := ctx.Interface()
			if inter == nil {
				panic("$ id outside of an interface")
			}
			id.IdType = inter.InterfaceType
		} else {
			return errors.New("undefined id " + id.Name)
		}
	}
	return nil
}

func typeCall(call *FCall, unifier Substitutions) error {
	call.CallType = NewVariable()
	ftype := call.MakeFunType()
	s, err := Unifier(ftype, call.Function.Type())
	if err != nil {
		return err
	}
	call.CallType = s.Apply(call.CallType)
	for _, p := range call.Params {
		ConvertTypes(p, s)
	}
	return unifier.Combine(s)
}

func initialiseVariables(exp Expression) error {
	return VisitBefore(exp, func(v Ast, ctx *VisitContext) error {
		if d, ok := v.(*FDef); ok {
			for _, p := range d.Params {
				name := p.Name
				if ctx.DefinitionOf(name) != nil || ctx.ParamOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
				if p.ParamType == nil {
					p.ParamType = NewVariable()
				}
			}
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

func buildInInterfaceaFor(name string) []*Interface {
	result := []*Interface{}
	for _, i := range global {
		if i.Definitions.Assignments[name] != nil {
			result = append(result, i)
		}
	}
	return result
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
	return buildInInterfaceaFor(name)
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
			typeConstant(c)
		} else if b, ok := v.(*Branch); ok {
			u, err := Unifier(b.Condition.Type(), Bool)
			if err != nil {
				return err
			}
			err = unifier.Combine(u)
			if err != nil {
				return err
			}
			u, err = Unifier(b.True.Type(), b.False.Type())
			if err != nil {
				return err
			}
			err = unifier.Combine(u)
			if err != nil {
				return err
			}
		} else if b, ok := v.(*Block); ok {
			ConvertTypes(b.Value, unifier)
			for _, a := range b.Def.Assignments {
				_, isDef := a.(*FDef)
				if !isDef {
					ConvertTypes(a, unifier)
				}
			}
			newinfers := []*Interface{}
			for _, i := range b.Def.Interfaces {
				ConvertTypes(i.Definitions, unifier)
				i.InterfaceType = unifier.Apply(i.InterfaceType)
				newinfers = append(newinfers, i)
			}
			b.Def.Interfaces = newinfers
		} else if i, ok := v.(*Id); ok {
			if err := typeId(i, ctx, unifier); err != nil {
				return err
			}
		} else if o, ok := v.(*Op); ok {
			inters := interfacesFor(o.Name, ctx)
			wantFun := NewFunction(o.Left.Type(), o.Right.Type(), o.OpType)
			match := false
			for _, in := range inters {
				it := in.Definitions.Assignments[o.Name].Type()
				s, err := Unifier(it, wantFun)
				if err == nil {
					err = unifier.Combine(s)
					if err == nil {
						match = true
						break
					}
				}
			}
			if !match {
				return errors.New("no operation \"" + o.Name + "\" found of type " + Signature(wantFun))
			}
		} else if c, ok := v.(*FCall); ok {
			if err := typeCall(c, unifier); err != nil {
				return err
			}
		} else if d, ok := v.(*FDef); ok {
			ConvertTypes(d, unifier)
		} else if t, ok := v.(*TypeDecl); ok {
			uni, err := Unifier(t.DeclType, t.Exp.Type())
			if err != nil {
				return err
			}
			err = unifier.Combine(uni)
			if err != nil {
				return err
			}
		} else if a, ok := v.(*FieldAccessor); ok {
			inters := interfacesFor(a.Field, ctx)
			if len(inters) == 0 {
				return errors.New("no interface with field \"" + a.Field + "\" found")
			}
			panic("NOT IMPLEMENTED")
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
	// TODO: skip defs / interfaces / assignments for performance
	ConvertTypes(exp, unifier)
	return err
}
