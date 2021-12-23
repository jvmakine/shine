package typeinference

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
)

func fun(ts ...Type) *ast.Exp {
	return &ast.Exp{Op: &ast.Op{Type: function(ts...)}}
}

func union(un ...Primitive) Type {
	return Type{Variable: &TypeVar{Union: un}}
}

func function(ts ...Type) Type {
	return MakeFunction(ts...)
}

func withVar(v Type, f func(t Type) *ast.Exp) *ast.Exp {
	return f(v)
}

var global map[string]*ast.Exp = map[string]*ast.Exp{
	"+":  withVar(union(Int, Real, String), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"-":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"*":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"%":  fun(IntP, IntP, IntP),
	"/":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"<":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	">":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	">=": withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"<=": withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"||": fun(BoolP, BoolP, BoolP),
	"&&": fun(BoolP, BoolP, BoolP),
	"==": withVar(union(Int, Bool, String), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"!=": withVar(union(Int, Bool), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"if": withVar(MakeVariable(), func(t Type) *ast.Exp { return fun(BoolP, t, t, t) }),
}

func typeConstant(constant *ast.Const) {
	if constant.Int != nil {
		constant.Type = IntP
	} else if constant.Bool != nil {
		constant.Type = BoolP
	} else if constant.Real != nil {
		constant.Type = RealP
	} else if constant.String != nil {
		constant.Type = StringP
	} else {
		panic("invalid const")
	}
}

func typeId(id *ast.Id, ctx *ast.VisitContext) error {
	block := ctx.BlockOf(id.Name)
	if ctx.Path()[id.Name] {
		id.Type = MakeVariable()
	} else if block != nil {
		b := ctx.BlockOf(id.Name)
		ref := b.Assignments[id.Name]
		if ref != nil {
			id.Type = ref.Type().Copy(NewTypeCopyCtx())
			return nil
		}
		tdef := b.TypeDefs[id.Name]
		if tdef == nil {
			panic("no id found: " + id.Name)
		}
		if tdef.Struct != nil {
			id.Type = tdef.Type().Copy(NewTypeCopyCtx())
		} else {
			return errors.New("invalid type def")
		}
	} else if p := ctx.ParamOf(id.Name); p != nil {
		id.Type = p.Type
	} else {
		return errors.New("undefined id " + id.Name)
	}
	return nil
}

func typeOp(op *ast.Op, ctx *ast.VisitContext) error {
	g := global[op.Name]
	if g == nil {
		panic("invalid op " + op.Name)
	}
	op.Type = g.Type().Copy(NewTypeCopyCtx())
	return nil
}

func typeCall(call *ast.FCall, unifier Substitutions) error {
	call.Type = MakeVariable()
	ftype := call.MakeFunType()
	s, err := ftype.Unifier(call.Function.Type())
	if err != nil {
		return err
	}
	call.Type = s.Apply(call.Type)
	for _, p := range call.Params {
		p.Convert(s)
	}
	unifier.Combine(s)
	return nil
}

func initialiseVariables(exp *ast.Exp) error {
	return exp.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Def != nil {
			for _, p := range v.Def.Params {
				name := p.Name
				if ctx.BlockOf(name) != nil || ctx.ParamOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
				if !p.Type.IsDefined() {
					p.Type = MakeVariable()
				}
			}
		} else if v.Block != nil {
			err := v.Block.CheckValueCycles()
			if err != nil {
				return err
			}
			for name, value := range v.Block.TypeDefs {
				if value.Struct != nil {
					ts := make([]types.Type, len(value.Struct.Fields)+1)
					sf := make([]types.SField, len(value.Struct.Fields))
					for i, v := range value.Struct.Fields {
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
					ts[len(value.Struct.Fields)] = stru

					value.Struct.Type = types.MakeFunction(ts...)
				}
			}
		}
		return nil
	})
}

func resolveNamed(name string, ctx *ast.VisitContext) (Type, error) {
	block := ctx.BlockOf(name)
	if block == nil {
		return Type{}, errors.New("type " + name + " is undefined")
	}
	exp := block.TypeDefs[name]
	if exp == nil || exp.Struct == nil {
		return Type{}, errors.New(name + " is not a type")
	}
	fs := make([]types.SField, len(exp.Struct.Fields))
	for i, f := range exp.Struct.Fields {
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

func rewriteNamed(exp *ast.Exp) error {
	return exp.RewriteTypes(func(t Type, ctx *ast.VisitContext) (Type, error) {
		if t.IsNamed() {
			return resolveNamed(*t.Named, ctx)
		}
		return t, nil
	})
}

func Infer(exp *ast.Exp) error {
	blockCount := 0
	if err := initialiseVariables(exp); err != nil {
		return err
	}
	if err := rewriteNamed(exp); err != nil {
		return err
	}
	unifier := MakeSubstitutions()
	crawler := func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Const != nil {
			typeConstant(v.Const)
		} else if v.Block != nil {
			blockCount++
			v.Block.ID = blockCount
		} else if v.Id != nil {
			if err := typeId(v.Id, ctx); err != nil {
				return err
			}
		} else if v.Op != nil {
			if err := typeOp(v.Op, ctx); err != nil {
				return err
			}
		} else if v.Call != nil {
			if err := typeCall(v.Call, unifier); err != nil {
				return err
			}
		} else if v.Def != nil {
			v.Convert(unifier)
		} else if v.TDecl != nil {
			uni, err := v.TDecl.Type.Unifier(v.TDecl.Exp.Type())
			if err != nil {
				return err
			}
			unifier.Combine(uni)
		} else if v.FAccess != nil {
			vari := MakeVariable()
			typ := MakeStructuralVar(map[string]Type{v.FAccess.Field: vari})
			uni, err := v.FAccess.Exp.Type().Unifier(typ)
			if err != nil {
				return err
			}
			unifier.Combine(uni)
			v.FAccess.Type = vari
		}
		return nil
	}
	// infer used code
	visited, err := exp.CrawlAfter(crawler)
	if err != nil {
		return err
	}
	// infer unused code
	err = exp.VisitAfter(func(v *ast.Exp, ctx *ast.VisitContext) error {
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
