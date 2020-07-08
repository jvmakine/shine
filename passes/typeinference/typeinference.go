package typeinference

import (
	"errors"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
)

func fun(ts ...Type) Expression {
	return &Op{OpType: function(ts...)}
}

func union(un ...Type) Type {
	return Type{Variable: &TypeVar{Union: un}}
}

func function(ts ...Type) Type {
	return MakeFunction(ts...)
}

func withVar(v Type, f func(t Type) Expression) Expression {
	return f(v)
}

var global map[string]Expression = map[string]Expression{
	"+":  withVar(union(IntP, RealP, StringP), func(t Type) Expression { return fun(t, t, t) }),
	"-":  withVar(union(IntP, RealP), func(t Type) Expression { return fun(t, t, t) }),
	"*":  withVar(union(IntP, RealP), func(t Type) Expression { return fun(t, t, t) }),
	"%":  fun(IntP, IntP, IntP),
	"/":  withVar(union(IntP, RealP), func(t Type) Expression { return fun(t, t, t) }),
	"<":  withVar(union(IntP, RealP), func(t Type) Expression { return fun(t, t, BoolP) }),
	">":  withVar(union(IntP, RealP), func(t Type) Expression { return fun(t, t, BoolP) }),
	">=": withVar(union(IntP, RealP), func(t Type) Expression { return fun(t, t, BoolP) }),
	"<=": withVar(union(IntP, RealP), func(t Type) Expression { return fun(t, t, BoolP) }),
	"||": fun(BoolP, BoolP, BoolP),
	"&&": fun(BoolP, BoolP, BoolP),
	"==": withVar(union(IntP, BoolP, StringP), func(t Type) Expression { return fun(t, t, BoolP) }),
	"!=": withVar(union(IntP, BoolP), func(t Type) Expression { return fun(t, t, BoolP) }),
	"if": withVar(MakeVariable(), func(t Type) Expression { return fun(BoolP, t, t, t) }),
}

func typeConstant(constant *Const) {
	if constant.Int != nil {
		constant.ConstType = IntP
	} else if constant.Bool != nil {
		constant.ConstType = BoolP
	} else if constant.Real != nil {
		constant.ConstType = RealP
	} else if constant.String != nil {
		constant.ConstType = StringP
	} else {
		panic("invalid const")
	}
}

func typeId(id *Id, ctx *VisitContext, unifier Substitutions) error {
	defin := ctx.DefinitionOf(id.Name)
	if ctx.Path()[id.Name] {
		id.IdType = MakeVariable()
	} else if defin != nil && ctx.DefinitionOf(id.Name).Assignments[id.Name] != nil {
		ref := ctx.DefinitionOf(id.Name).Assignments[id.Name]
		id.IdType = ref.Type().Copy(NewTypeCopyCtx())
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

func typeOp(op *Op, ctx *VisitContext) error {
	g := global[op.Name]
	if g == nil {
		panic("invalid op " + op.Name)
	}
	op.OpType = g.Type().Copy(NewTypeCopyCtx())
	return nil
}

func typeCall(call *FCall, unifier Substitutions) error {
	call.CallType = MakeVariable()
	ftype := call.MakeFunType()
	s, err := ftype.Unifier(call.Function.Type())
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
				if !p.ParamType.IsDefined() {
					p.ParamType = MakeVariable()
				}
			}
		} else if b, ok := v.(*Block); ok {
			err := b.CheckValueCycles()
			if err != nil {
				return err
			}
		} else if d, ok := v.(*Definitions); ok {
			for name, value := range d.Assignments {
				if s, ok := value.(*Struct); ok {
					ts := make([]types.Type, len(s.Fields)+1)
					sf := make([]types.SField, len(s.Fields))
					for i, v := range s.Fields {
						typ := v.Type
						if !typ.IsDefined() {
							typ = MakeVariable()
						}
						v.Type = typ

						ts[i] = typ
						sf[i] = SField{
							Name: v.Name,
							Type: typ,
						}
					}

					stru := types.MakeStructure(name, sf...)
					ts[len(s.Fields)] = stru

					s.StructType = stru
				}
			}
			for _, in := range d.Interfaces {
				if !in.InterfaceType.IsDefined() {
					varit := MakeVariable()
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
			fa.FAType = MakeVariable()
		}
		return nil
	})
}

func resolveNamed(name string, ctx *VisitContext) (Type, error) {
	defin := ctx.DefinitionOf(name)
	if defin == nil {
		return Type{}, errors.New("type " + name + " is undefined")
	}
	value := defin.Assignments[name]
	if value == nil {
		return Type{}, errors.New(name + " is not a type")
	}
	struc, ok := value.(*Struct)
	if !ok {
		return Type{}, errors.New(name + " is not a type")
	}
	fs := make([]types.SField, len(struc.Fields))
	for i, f := range struc.Fields {
		if !f.Type.IsDefined() {
			fs[i] = types.SField{
				Name: f.Name,
				Type: types.MakeVariable(),
			}
		} else {
			fs[i] = types.SField{
				Name: f.Name,
				Type: f.Type,
			}
		}
	}
	return types.MakeStructure(name, fs...), nil
}

func rewriteNamed(exp Expression) error {
	return RewriteTypes(exp, func(t Type, ctx *VisitContext) (Type, error) {
		if t.IsNamed() {
			return resolveNamed(*t.Named, ctx)
		}
		return t, nil
	})
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
			if err := typeOp(o, ctx); err != nil {
				return err
			}
		} else if c, ok := v.(*FCall); ok {
			if err := typeCall(c, unifier); err != nil {
				return err
			}
		} else if d, ok := v.(*FDef); ok {
			ConvertTypes(d, unifier)
		} else if t, ok := v.(*TypeDecl); ok {
			uni, err := t.DeclType.Unifier(t.Exp.Type())
			if err != nil {
				return err
			}
			err = unifier.Combine(uni)
			if err != nil {
				return err
			}
		} else if a, ok := v.(*FieldAccessor); ok {
			expType := MakeStructuralVar(map[string]Type{a.Field: a.FAType})
			fieldType := Type{}
			inters := ctx.InterfacesWith(a.Field)
			tctx := types.NewTypeCopyCtx()
			for _, in := range inters {
				ift := in.Interf.InterfaceType
				funt := in.Interf.Definitions.Assignments[a.Field].Type()
				expType = expType.AddToUnion(ift.Copy(tctx))
				fat := funt.Copy(tctx)
				if !fieldType.IsDefined() {
					fieldType = fat
				} else {
					fieldType = fieldType.AddToUnion(fat)
				}
			}
			uni, err := a.Exp.Type().Unifier(expType)
			if err != nil {
				return err
			}
			err = unifier.Combine(uni)
			if err != nil {
				return err
			}
			if fieldType.IsDefined() {
				uni, err := a.FAType.Unifier(fieldType)
				if err != nil {
					return err
				}
				err = unifier.Combine(uni)
				if err != nil {
					return err
				}
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
