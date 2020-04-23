package inferer

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

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
		&ast.Exp{Type: function(result...)},
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
	root := &context{ids: global}
	initialise(exp, root)
	graph := MakeTypeGraph()
	if err := inferExp(exp, root.sub(), &graph); err != nil {
		return err
	}
	sub, err := graph.Substitutions()
	if err != nil {
		return err
	}
	sub.Convert(exp)
	return nil
}

func UnifyCall(def *ast.Exp, call *ast.Exp) error {
	graph := MakeTypeGraph()
	for i, p := range call.Call.Params {
		if err := graph.Add(p.Type, def.Def.Params[i].Type); err != nil {
			return err
		}
	}
	if err := graph.Add(def.Type, call.Type); err != nil {
		return err
	}
	sub, err := graph.Substitutions()
	if err != nil {
		return err
	}
	sub.Convert(def)
	sub.Convert(call)
	return nil
}

func Unify(a Type, b Type) (Substitutions, error) {
	graph := MakeTypeGraph()
	if err := graph.Add(a, b); err != nil {
		return nil, err
	}
	return graph.Substitutions()
}

func initialise(exp *ast.Exp, ctx *context) {
	if exp.Const != nil {
		if exp.Const.Int != nil {
			exp.Type = IntP
		} else if exp.Const.Bool != nil {
			exp.Type = BoolP
		} else if exp.Const.Real != nil {
			exp.Type = RealP
		} else {
			panic("invalid const")
		}
	} else if exp.Block != nil {
		for _, a := range exp.Block.Assignments {
			initialise(a.Value, ctx)
		}
		initialise(exp.Block.Value, ctx)
		exp.Type = exp.Block.Value.Type
	} else if exp.Call != nil {
		for i := range exp.Call.Params {
			initialise(exp.Call.Params[i], ctx)
		}
		exp.Type = MakeVariable()
	} else if exp.Def != nil {
		initialise(exp.Def.Body, ctx)
		ftps := make([]Type, len(exp.Def.Params)+1)
		for i := range exp.Def.Params {
			v := MakeVariable()
			exp.Def.Params[i].Type = v
			ftps[i] = v
		}
		ftps[len(exp.Def.Params)] = exp.Def.Body.Type
		exp.Type = MakeFunction(ftps...)
	} else if exp.Id != nil {
		exp.Type = MakeVariable()
	} else {
		panic("invalid expression")
	}
}

func inferExp(exp *ast.Exp, ctx *context, graph *TypeGraph) error {
	if exp.Block != nil {
		if err := exp.Block.CheckValueCycles(); err != nil {
			return err
		}
		nctx := ctx.sub()
		for _, a := range exp.Block.Assignments {
			nctx.setActiveType(a.Name, &a.Value.Type)
		}
		for _, a := range exp.Block.Assignments {
			if err := inferExp(a.Value, nctx, graph); err != nil {
				return err
			}
		}
		// Apply substitutions to the assignments before inferring the expression
		// to avoid the expression affecting the assignment types
		sub, err := graph.Substitutions()
		if err != nil {
			return err
		}
		for _, a := range exp.Block.Assignments {
			nctx.stopInference(a.Name)
			sub.ConvertAssignment(a)
			nctx.ids[a.Name] = &excon{v: a.Value, c: nctx}
		}
		// infer and convert the block expression
		if err := inferExp(exp.Block.Value, nctx, graph); err != nil {
			return err
		}
		sub, err = graph.Substitutions()
		if err != nil {
			return err
		}
	} else if exp.Id != nil {
		var typ *Type
		if def := ctx.getActiveType(exp.Id.Name); def != nil {
			typ = def
		} else if at := ctx.getId(exp.Id.Name); at != nil {
			typ = &at.v.Type
		}
		if typ == nil {
			return errors.New("undefined id: " + exp.Id.Name)
		}
		if err := graph.Add(*typ, exp.Type); err != nil {
			return err
		}
	} else if exp.Call != nil {
		name := exp.Call.Name
		var typ *Type

		if at := ctx.getActiveType(name); at != nil {
			typ = at
		} else if def := ctx.getId(name); def != nil {
			v := def.v.Type.Copy(NewTypeCopyCtx())
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
			if err := inferExp(a, ctx, graph); err != nil {
				return err
			}
			args[i] = a.Type
		}
		args[len(exp.Call.Params)] = exp.Type
		if err := graph.Add(MakeFunction(args...), *typ); err != nil {
			return err
		}
	} else if exp.Def != nil {
		sc := ctx.sub()
		for _, p := range exp.Def.Params {
			sc.setActiveType(p.Name, &p.Type)
		}
		if err := inferExp(exp.Def.Body, sc, graph); err != nil {
			return err
		}
		for _, p := range exp.Def.Params {
			sc.stopInference(p.Name)
		}
	}
	return nil
}
