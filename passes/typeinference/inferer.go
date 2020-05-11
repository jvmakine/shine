package typeinference

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/typeinference/internal/graph"
	. "github.com/jvmakine/shine/types"
)

var blockCount = 0

func fun(ts ...interface{}) *excon {
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
	return &excon{
		&ast.Exp{Id: &ast.Id{Type: function(result...)}},
		&context{},
	}
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

func withVar(v Type, f func(t Type) *excon) *excon {
	return f(v)
}

var (
	integer = base(Int)
	real    = base(Real)
	boolean = base(Bool)
)

var global map[string]*excon = map[string]*excon{
	"+":  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, t) }),
	"-":  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, t) }),
	"*":  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, t) }),
	"%":  fun(integer, integer, integer),
	"/":  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, t) }),
	"<":  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, boolean) }),
	">":  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, boolean) }),
	">=": withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, boolean) }),
	"<=": withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, boolean) }),
	"||": fun(boolean, boolean, boolean),
	"&&": fun(boolean, boolean, boolean),
	"==": withVar(union(Int, Bool), func(t Type) *excon { return fun(t, t, boolean) }),
	"if": withVar(variable(), func(t Type) *excon { return fun(boolean, t, t, t) }),
}

func (ctx *context) getId(id string) *excon {
	if ctx.ids[id] != nil {
		return ctx.ids[id]
	} else if ctx.parent != nil {
		return ctx.parent.getId(id)
	}
	return nil
}

func Infer(exp *ast.Exp) error {
	blockCount = 0 // TODO: Remove global var
	root := &context{ids: global}
	initialise(exp, root)
	tgraph := graph.MakeTypeGraph()
	if err := inferExp(exp, root, &tgraph); err != nil {
		return err
	}
	sub, err := tgraph.Substitutions()
	if err != nil {
		return err
	}
	sub.Convert(exp)
	return nil
}

func Unify(a Type, b Type) (graph.Substitutions, error) {
	tgraph := graph.MakeTypeGraph()
	if err := tgraph.Add(a, b); err != nil {
		return nil, err
	}
	return tgraph.Substitutions()
}

func initialise(exp *ast.Exp, ctx *context) {
	if exp.Const != nil {
		if exp.Const.Int != nil {
			exp.Const.Type = IntP
		} else if exp.Const.Bool != nil {
			exp.Const.Type = BoolP
		} else if exp.Const.Real != nil {
			exp.Const.Type = RealP
		} else {
			panic("invalid const")
		}
	} else if exp.Block != nil {
		for _, a := range exp.Block.Assignments {
			initialise(a, ctx)
		}
		initialise(exp.Block.Value, ctx)
	} else if exp.Call != nil {
		for i := range exp.Call.Params {
			initialise(exp.Call.Params[i], ctx)
		}
		exp.Call.Type = MakeVariable()
	} else if exp.Def != nil {
		initialise(exp.Def.Body, ctx)
		for i := range exp.Def.Params {
			exp.Def.Params[i].Type = MakeVariable()
		}
	} else if exp.Id != nil {
		exp.Id.Type = MakeVariable()
	} else {
		panic("invalid expression")
	}
}

func inferExp(exp *ast.Exp, ctx *context, tgraph *graph.TypeGraph) error {
	if exp.Block != nil {
		blockCount++
		exp.Block.ID = blockCount
		if err := exp.Block.CheckValueCycles(); err != nil {
			return err
		}
		nctx := ctx.sub(exp.Block.ID)
		for k, a := range exp.Block.Assignments {
			typ := a.Type()
			nctx.setActiveType(k, &typ)
		}
		for _, a := range exp.Block.Assignments {
			if err := inferExp(a, nctx, tgraph); err != nil {
				return err
			}
		}
		// Apply substitutions to the assignments before inferring the expression
		// to avoid the expression affecting the assignment types
		sub, err := tgraph.Substitutions()
		if err != nil {
			return err
		}
		for k, a := range exp.Block.Assignments {
			nctx.stopInference(k)
			sub.Convert(a)
			nctx.ids[k] = &excon{v: a, c: nctx}
		}
		// infer and convert the block expression
		if err := inferExp(exp.Block.Value, nctx, tgraph); err != nil {
			return err
		}
		sub, err = tgraph.Substitutions()
		if err != nil {
			return err
		}
	} else if exp.Id != nil {
		var typ *Type
		if def := ctx.getActiveType(exp.Id.Name); def != nil {
			typ = def
		} else if at := ctx.getId(exp.Id.Name); at != nil {
			t := at.v.Type()
			typ = &t
		}
		if typ == nil {
			return errors.New("undefined id: " + exp.Id.Name)
		}
		if err := tgraph.Add(*typ, exp.Type()); err != nil {
			return err
		}
	} else if exp.Call != nil {
		name := exp.Call.Name
		var typ *Type

		if at := ctx.getActiveType(name); at != nil {
			typ = at
		} else if def := ctx.getId(name); def != nil {
			v := def.v.Type().Copy(NewTypeCopyCtx())
			typ = &v
		}
		if typ == nil {
			return errors.New("undefined function: " + name)
		}
		if !typ.IsFunction() && !typ.IsVariable() {
			return errors.New("not a function: " + name)
		}
		args := make([]Type, len(exp.Call.Params)+1)
		for i, a := range exp.Call.Params {
			if err := inferExp(a, ctx, tgraph); err != nil {
				return err
			}
			args[i] = a.Type()
		}
		args[len(exp.Call.Params)] = exp.Type()
		if err := tgraph.Add(MakeFunction(args...), *typ); err != nil {
			return err
		}
	} else if exp.Def != nil {
		blockCount++
		sc := ctx.sub(blockCount)
		for _, p := range exp.Def.Params {
			sc.setActiveType(p.Name, &p.Type)
		}
		if err := inferExp(exp.Def.Body, sc, tgraph); err != nil {
			return err
		}
		for _, p := range exp.Def.Params {
			sc.stopInference(p.Name)
		}
	}
	return nil
}
