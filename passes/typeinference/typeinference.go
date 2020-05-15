package typeinference

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func fun(ts ...interface{}) *ast.Exp {
	result := make([]Type, len(ts))
	var variables map[string]*TypeVar = map[string]*TypeVar{}
	for _, t := range ts {
		switch v := t.(type) {
		case string:
			if variables[v] == nil {
				variables[v] = &TypeVar{}
			}
		}
	}

	for i, t := range ts {
		switch v := t.(type) {
		case Type:
			result[i] = v
		case string:
			result[i] = Type{Variable: variables[v]}
		}
	}
	return &ast.Exp{Op: &ast.Op{Type: function(result...)}}
}

func base(t Primitive) Type {
	return Type{Primitive: &t}
}

func union(un ...Primitive) Type {
	return Type{Variable: &TypeVar{Restrictions: un}}
}

func function(ts ...Type) Type {
	return MakeFunction(ts...)
}

func withVar(v Type, f func(t Type) *ast.Exp) *ast.Exp {
	return f(v)
}

var global map[string]*ast.Exp = map[string]*ast.Exp{
	"+":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
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
	"==": withVar(union(Int, Bool), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
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
	} else {
		panic("invalid const")
	}
}

func typeId(id *ast.Id, ctx *ast.VisitContext) error {
	block := ctx.BlockOf(id.Name)
	if ctx.Path()[id.Name] {
		id.Type = MakeVariable()
	} else if block != nil {
		ref := ctx.BlockOf(id.Name).Assignments[id.Name]
		id.Type = ref.Type().Copy(NewTypeCopyCtx())
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
				p.Type = MakeVariable()
			}
		} else if v.Block != nil {
			err := v.Block.CheckValueCycles()
			if err != nil {
				return err
			}
		}
		return nil
	})
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
		}
		return nil
	}
	// infer used code
	visited, err := exp.CrawlAfter(crawler)
	if err != nil {
		return err
	}
	// infer unused code
	err = exp.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if !visited[v] {
			v, err := v.CrawlAfter(crawler)
			if err != nil {
				return err
			}
			for k, v := range v {
				visited[k] = v
			}
		}
		return nil
	})
	return err
}
