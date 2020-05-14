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
	return &ast.Exp{Id: &ast.Id{Type: function(result...)}}
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

func variable() Type {
	return Type{Variable: &TypeVar{}}
}

func withVar(v Type, f func(t Type) *ast.Exp) *ast.Exp {
	return f(v)
}

var (
	integer = base(Int)
	real    = base(Real)
	boolean = base(Bool)
)

var global map[string]*ast.Exp = map[string]*ast.Exp{
	"+":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"-":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"*":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"%":  fun(integer, integer, integer),
	"/":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"<":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, boolean) }),
	">":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, boolean) }),
	">=": withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, boolean) }),
	"<=": withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, boolean) }),
	"||": fun(boolean, boolean, boolean),
	"&&": fun(boolean, boolean, boolean),
	"==": withVar(union(Int, Bool), func(t Type) *ast.Exp { return fun(t, t, boolean) }),
	"if": withVar(variable(), func(t Type) *ast.Exp { return fun(boolean, t, t, t) }),
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
	if block != nil {
		ref := ctx.BlockOf(id.Name).Assignments[id.Name]
		id.Type = ref.Type().Copy(NewTypeCopyCtx())
	} else if g := global[id.Name]; g != nil {
		id.Type = g.Type().Copy(NewTypeCopyCtx())
	} else if p := ctx.ParamOf(id.Name); p != nil {
		id.Type = p.Type
	} else {
		return errors.New("undefined id " + id.Name)
	}
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

func Infer(exp *ast.Exp) error {
	blockCount := 0
	// set function parameters as variables
	exp.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Def != nil {
			for _, p := range v.Def.Params {
				if ctx.BlockOf(p.Name) != nil || ctx.ParamOf(p.Name) != nil {
					return errors.New("redefinition of " + p.Name)
				}
				p.Type = MakeVariable()
			}
		}
		return nil
	})
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
