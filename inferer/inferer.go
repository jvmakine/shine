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
	ptrs := make([]Type, len(un))
	for i, u := range un {
		ptrs[i] = base(u)
	}
	return Type{Variable: &TypeVar{Restrictions: ptrs}}
}

func function(ts ...Type) Type {
	return Type{Function: &ts}
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
	parent := &context{ids: global}
	return inferExp(exp, &context{ids: map[string]*excon{}, parent: parent, activeVals: &[]string{}}, nil)
}

func inferExp(exp *ast.Exp, ctx *context, name *string) error {
	if !exp.Type.IsDefined() {
		if exp.Const != nil {
			if exp.Const.Int != nil {
				exp.Type = base(Int)
			} else if exp.Const.Bool != nil {
				exp.Type = base(Bool)
			} else if exp.Const.Real != nil {
				exp.Type = base(Real)
			} else {
				panic("unknown constant")
			}
			return nil
		} else if exp.Id != nil {
			typ, err := inferId(*exp.Id, ctx)
			exp.Type = typ
			return err
		} else if exp.Call != nil {
			typ, err := inferCall(exp.Call, ctx)
			exp.Type = typ
			return err
		} else if exp.Def != nil {
			typ, err := inferDef(exp.Def, &context{parent: ctx, ids: map[string]*excon{}, activeVals: ctx.activeVals}, name)
			exp.Type = typ
			return err
		} else if exp.Block != nil {
			typ, err := inferBlock(exp.Block, &context{parent: ctx, ids: map[string]*excon{}, activeVals: ctx.activeVals}, name)
			exp.Type = typ
			return err
		}
		panic("unexpected expression")
	}
	return nil
}

func inferCall(call *ast.FCall, ctx *context) (Type, error) {
	var params []Type = make([]Type, len(call.Params)+1)
	for i, p := range call.Params {
		err := inferExp(p, ctx, nil)
		if err != nil {
			return Type{}, err
		}
		params[i] = p.Type
	}
	// Recursive type definition
	it := ctx.getActiveType(call.Name)
	var ft Type
	if it.IsDefined() {
		ft = it
	} else {
		ec := ctx.getId(call.Name)
		if ec == nil {
			return Type{}, errors.New("undefined function: '" + call.Name + "'")
		}
		if !ec.v.Type.IsDefined() {
			err := inferExp(ec.v, ec.c, &call.Name)
			if err != nil {
				return Type{}, err
			}
		}
		ft = ec.v.Type
	}
	if !ft.IsFunction() {
		return Type{}, errors.New("not a function: '" + call.Name + "'")
	}
	params[len(call.Params)] = ft.ReturnType().Copy(NewTypeCopyCtx())

	ft2 := function(params...)
	unifier, err := Unify(ft2, ft)
	if err != nil {
		return Type{}, err
	}
	unifier.Apply(&ft2)
	if it.IsDefined() {
		unifier.Apply(&ft)
	}
	return ft2.ReturnType(), nil
}

func inferDef(def *ast.FDef, ctx *context, name *string) (Type, error) {
	var paramTypes []Type = make([]Type, len(def.Params)+1)
	for i, p := range def.Params {
		if ctx.getId(p.Name) != nil {
			return Type{}, errors.New("redefinition of '" + p.Name + "'")
		}
		typ := variable()
		ctx.setActiveType(p.Name, typ)
		paramTypes[i] = typ
		p.Type = typ
	}
	paramTypes[len(def.Params)] = variable()
	ftyp := function(paramTypes...)
	if name != nil {
		ctx.setActiveType(*name, ftyp)
	}
	err := inferExp(def.Body, ctx, nil)
	if err != nil {
		return Type{}, err
	}
	paramTypes[len(def.Params)] = def.Body.Type
	if name != nil {
		ctx.stopInference(*name)
	}
	for _, p := range def.Params {
		ctx.stopInference(p.Name)
	}
	return ftyp, nil
}

func inferId(id string, ctx *context) (Type, error) {
	def := ctx.getId(id)
	if def == nil {
		act := ctx.getActiveType(id)
		if act.IsDefined() {
			return act, nil
		}
		return Type{}, errors.New("undefined id '" + id + "'")
	}
	if !def.v.Type.IsDefined() {
		if contains((*ctx.activeVals), id) {
			return Type{}, errors.New("recursive value: " + cycleToStr((*ctx.activeVals), id))
		}
		(*ctx.activeVals) = append((*ctx.activeVals), id)
		err := inferExp(def.v, def.c, &id)
		if err != nil {
			return Type{}, err
		}
		(*ctx.activeVals) = (*ctx.activeVals)[:1]
	}
	return def.v.Type.Copy(NewTypeCopyCtx()), nil
}

func contains(arr []string, v string) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func cycleToStr(arr []string, v string) string {
	res := ""
	for _, a := range arr {
		res = res + a + " -> "
	}
	return res + v
}

func inferBlock(block *ast.Block, ctx *context, name *string) (Type, error) {
	for _, a := range block.Assignments {
		if ctx.getId(a.Name) != nil {
			return Type{}, errors.New("redefinition of '" + a.Name + "'")
		}
		ctx.ids[a.Name] = &excon{v: a.Value, c: ctx}
	}
	for _, a := range block.Assignments {
		if !a.Value.Type.IsDefined() {
			err := inferExp(a.Value, ctx, &a.Name)
			if err != nil {
				return Type{}, err
			}
		}
	}

	err := inferExp(block.Value, ctx, name)
	if err != nil {
		return Type{}, err
	}
	return block.Value.Type, nil
}
