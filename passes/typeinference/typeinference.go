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

func union(un ...Primitive) Type {
	return Type{Variable: &TypeVar{Union: un}}
}

func function(ts ...Type) Type {
	return MakeFunction(ts...)
}

func withVar(v Type, f func(t Type) Expression) Expression {
	return f(v)
}

var global map[string]Expression = map[string]Expression{
	"+":  withVar(union(Int, Real, String), func(t Type) Expression { return fun(t, t, t) }),
	"-":  withVar(union(Int, Real), func(t Type) Expression { return fun(t, t, t) }),
	"*":  withVar(union(Int, Real), func(t Type) Expression { return fun(t, t, t) }),
	"%":  fun(IntP, IntP, IntP),
	"/":  withVar(union(Int, Real), func(t Type) Expression { return fun(t, t, t) }),
	"<":  withVar(union(Int, Real), func(t Type) Expression { return fun(t, t, BoolP) }),
	">":  withVar(union(Int, Real), func(t Type) Expression { return fun(t, t, BoolP) }),
	">=": withVar(union(Int, Real), func(t Type) Expression { return fun(t, t, BoolP) }),
	"<=": withVar(union(Int, Real), func(t Type) Expression { return fun(t, t, BoolP) }),
	"||": fun(BoolP, BoolP, BoolP),
	"&&": fun(BoolP, BoolP, BoolP),
	"==": withVar(union(Int, Bool, String), func(t Type) Expression { return fun(t, t, BoolP) }),
	"!=": withVar(union(Int, Bool), func(t Type) Expression { return fun(t, t, BoolP) }),
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

func typeId(id *Id, ctx *VisitContext) error {
	block := ctx.BlockOf(id.Name)
	if ctx.Path()[id.Name] {
		id.IdType = MakeVariable()
	} else if block != nil && ctx.BlockOf(id.Name).Def.Assignments[id.Name] != nil {
		ref := ctx.BlockOf(id.Name).Def.Assignments[id.Name]
		id.IdType = ref.Type().Copy(NewTypeCopyCtx())
	} else if block != nil && ctx.BlockOf(id.Name).Def.TypeDefs[id.Name] != nil {
		ref := ctx.BlockOf(id.Name).Def.TypeDefs[id.Name]
		id.IdType = ref.StructType.Copy(NewTypeCopyCtx())
	} else if p := ctx.ParamOf(id.Name); p != nil {
		id.IdType = p.ParamType
	} else {
		return errors.New("undefined id " + id.Name)
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

func convert(e Expression, s Substitutions) {
	RewriteTypes(e, func(t Type, ctx *VisitContext) (Type, error) {
		return s.Apply(t), nil
	})
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
		convert(p, s)
	}
	unifier.Combine(s)
	return nil
}

func initialiseVariables(exp Expression) error {
	return VisitBefore(exp, func(v Ast, ctx *VisitContext) error {
		if d, ok := v.(*FDef); ok {
			for _, p := range d.Params {
				name := p.Name
				if ctx.BlockOf(name) != nil || ctx.ParamOf(name) != nil {
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
			for name, value := range b.Def.TypeDefs {
				ts := make([]types.Type, len(value.Fields)+1)
				sf := make([]types.SField, len(value.Fields))
				for i, v := range value.Fields {
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
				ts[len(value.Fields)] = stru

				value.StructType = types.MakeFunction(ts...)
			}
		}
		return nil
	})
}

func resolveNamed(name string, ctx *VisitContext) (Type, error) {
	block := ctx.BlockOf(name)
	if block == nil {
		return Type{}, errors.New("type " + name + " is undefined")
	}
	struc := block.Def.TypeDefs[name]
	if struc == nil {
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
	blockCount := 0
	if err := initialiseVariables(exp); err != nil {
		return err
	}
	if err := rewriteNamed(exp); err != nil {
		return err
	}
	unifier := MakeSubstitutions()
	crawler := func(v Ast, ctx *VisitContext) error {
		if c, ok := v.(*Const); ok {
			typeConstant(c)
		} else if b, ok := v.(*Block); ok {
			blockCount++
			b.ID = blockCount
		} else if i, ok := v.(*Id); ok {
			if err := typeId(i, ctx); err != nil {
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
			convert(d, unifier)
		} else if t, ok := v.(*TypeDecl); ok {
			uni, err := t.DeclType.Unifier(t.Exp.Type())
			if err != nil {
				return err
			}
			unifier.Combine(uni)
		} else if a, ok := v.(*FieldAccessor); ok {
			vari := MakeVariable()
			typ := MakeStructuralVar(map[string]Type{a.Field: vari})
			uni, err := a.Exp.Type().Unifier(typ)
			if err != nil {
				return err
			}
			unifier.Combine(uni)
			a.FAType = vari
		}
		return nil
	}
	// infer used code
	visited, err := CrawlAfter(exp, crawler)
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
