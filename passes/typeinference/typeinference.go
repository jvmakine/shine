package typeinference

import (
	"errors"
	"strings"

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
		tc := b.TCFunctions[id.Name]
		if tc != nil {
			id.Type = tc.TypeClass.Functions[id.Name].TypeDecl.Copy(NewTypeCopyCtx())
			return nil
		}
		tdef := b.TypeDefs[id.Name]
		if tdef == nil {
			panic("no id found: " + id.Name)
		}
		if tdef.Struct != nil {
			id.Type = tdef.Struct.Constructor().Copy(types.NewTypeCopyCtx())
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

func typeCall(call *ast.FCall, unifier Substitutions, ctx *ast.VisitContext) error {
	call.Type = MakeVariable()
	ftype := call.MakeFunType()

	s, err := ftype.Unifier(call.Function.Type(), ctx)
	if err != nil {
		return err
	}
	call.Type = s.Apply(call.Type, ctx)

	for _, p := range call.Params {
		p.Convert(s, ctx)
	}

	err = unifier.Combine(s, ctx)
	if err != nil {
		return err
	}

	return nil
}

func Infer(exp *ast.Exp) error {
	blockCount := 0
	if err := initialiseVariables(exp); err != nil {
		return err
	}
	unifier := MakeSubstitutions()
	crawler := func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Const != nil {
			typeConstant(v.Const)
		} else if v.Block != nil {
			blockCount++
			v.Block.ID = blockCount
			nctx := ctx.SubBlock(v.Block)

			for name := range v.Block.Assignments {
				if ctx.BlockOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
			}
			for name := range v.Block.TypeDefs {
				if ctx.BlockOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
			}
			for _, binding := range v.Block.TypeBindings {
				name := binding.Name
				class := nctx.TypeDef(name).TypeClass
				found := map[string]bool{}

				for fname, fdef := range binding.Bindings {
					found[fname] = true
					fun := class.Functions[fname]
					if fun == nil {
						return errors.New("function " + fname + " not defined in " + name)
					}

					exp := &ast.Exp{Def: fdef}
					expType := exp.Type()

					s, err := fun.Type().Unifier(expType, nctx)
					if err != nil {
						return err
					}

					exp.Convert(s, nctx)
					binding.Bindings[fname] = exp.Def
				}

				for fname := range class.Functions {
					if !found[fname] {
						strs := make([]string, len(binding.Parameters))
						for i, t := range binding.Parameters {
							strs[i] = t.Signature()
						}
						sig := binding.Name + "[" + strings.Join(strs, ",") + "]"
						return errors.New("following functions were not defined for " + sig + ": " + fname)
					}
				}
			}
		} else if v.Id != nil {
			if err := typeId(v.Id, ctx); err != nil {
				return err
			}
		} else if v.Op != nil {
			if err := typeOp(v.Op, ctx); err != nil {
				return err
			}
		} else if call := v.Call; call != nil {
			if err := typeCall(call, unifier, ctx); err != nil {
				return err
			}
		} else if v.Def != nil {
			v.Convert(unifier, ctx)
		} else if v.TDecl != nil {
			uni, err := v.TDecl.Type.Unifier(v.TDecl.Exp.Type(), ctx)
			if err != nil {
				return err
			}
			unifier.Combine(uni, ctx)
			v.TDecl.Exp.Convert(unifier, ctx)
		} else if v.FAccess != nil {
			vari := MakeVariable()
			typ := MakeStructuralVar(map[string]Type{v.FAccess.Field: vari})
			uni, err := v.FAccess.Exp.Type().Unifier(typ, ctx)
			if err != nil {
				return err
			}
			err = unifier.Combine(uni, ctx)
			if err != nil {
				return err
			}
			v.FAccess.Type = unifier.Apply(vari, ctx)
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
