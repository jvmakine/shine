package inferer

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

func fun(ts ...interface{}) *excon {
	result := make([]*types.TypePtr, len(ts))
	var variables map[string]*types.TypePtr = map[string]*types.TypePtr{}
	for _, t := range ts {
		switch v := t.(type) {
		case string:
			if variables[v] == nil {
				variables[v] = variable()
			}
		}
	}

	for i, t := range ts {
		switch v := t.(type) {
		case *types.TypePtr:
			result[i] = v
		case string:
			result[i] = variables[v]
		}
	}
	return &excon{
		&ast.Exp{Type: function(result...)},
		&context{},
	}
}

func base(t types.Primitive) *types.TypePtr {
	return &types.TypePtr{Def: &types.TypeDef{Base: &t}}
}

func function(ts ...*types.TypePtr) *types.TypePtr {
	return &types.TypePtr{&types.TypeDef{Fn: ts}}
}

func variable() *types.TypePtr {
	return &types.TypePtr{Def: &types.TypeDef{}}
}

var global map[string]*excon = map[string]*excon{
	"+":  fun(base(types.Int), base(types.Int), base(types.Int)),
	"-":  fun(base(types.Int), base(types.Int), base(types.Int)),
	"*":  fun(base(types.Int), base(types.Int), base(types.Int)),
	"%":  fun(base(types.Int), base(types.Int), base(types.Int)),
	"/":  fun(base(types.Int), base(types.Int), base(types.Int)),
	"<":  fun(base(types.Int), base(types.Int), base(types.Bool)),
	">":  fun(base(types.Int), base(types.Int), base(types.Bool)),
	">=": fun(base(types.Int), base(types.Int), base(types.Bool)),
	"<=": fun(base(types.Int), base(types.Int), base(types.Bool)),
	"==": fun("A", "A", base(types.Bool)),
	"if": fun(base(types.Bool), "A", "A", "A"),
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
	return inferExp(exp, &context{ids: map[string]*excon{}, parent: parent}, nil)
}

func inferExp(exp *ast.Exp, ctx *context, name *string) error {
	if exp.Type == nil {
		if exp.Const != nil {
			if exp.Const.Int != nil {
				exp.Type = base(types.Int)
			} else if exp.Const.Bool != nil {
				exp.Type = base(types.Bool)
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
			typ, err := inferDef(exp.Def, &context{parent: ctx, ids: map[string]*excon{}}, name)
			exp.Type = typ
			return err
		} else if exp.Block != nil {
			typ, err := inferBlock(exp.Block, &context{parent: ctx, ids: map[string]*excon{}}, name)
			exp.Type = typ
			return err
		}
		panic("unexpected expression")
	}
	return nil
}

func inferCall(call *ast.FCall, ctx *context) (*types.TypePtr, error) {
	var params []*types.TypePtr = make([]*types.TypePtr, len(call.Params)+1)
	for i, p := range call.Params {
		err := inferExp(p, ctx, nil)
		if err != nil {
			return nil, err
		}
		params[i] = p.Type
	}
	// Recursive type definition
	it := ctx.getActiveType(call.Name)
	var ft *types.TypePtr = nil
	if it != nil {
		ft = it
	} else {
		ec := ctx.getId(call.Name)
		if ec == nil {
			return nil, errors.New("undefined function: '" + call.Name + "'")
		}
		if ec.v.Type == nil {
			err := inferExp(ec.v, ec.c, &call.Name)
			if err != nil {
				return nil, err
			}
		}
		ft = ec.v.Type
	}
	if !ft.IsFunction() {
		return nil, errors.New("not a function: '" + call.Name + "'")
	}
	params[len(call.Params)] = ft.ReturnType().Copy()

	ft2 := function(params...)
	unifier, err := Unify(ft2, ft)
	if err != nil {
		return nil, err
	}
	unifier.ApplySource(ft2)
	if it != nil {
		unifier.ApplyDest(ft)
	}
	return ft2.ReturnType(), nil
}

func inferDef(def *ast.FDef, ctx *context, name *string) (*types.TypePtr, error) {
	var paramTypes []*types.TypePtr = make([]*types.TypePtr, len(def.Params)+1)
	for i, p := range def.Params {
		if ctx.getId(p.Name) != nil {
			return nil, errors.New("redefinition of '" + p.Name + "'")
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
		return nil, err
	}
	inferred := def.Body.Type
	unifier, err := Unify(paramTypes[len(def.Params)], inferred)
	if err != nil {
		return nil, err
	}
	unifier.ApplySource(paramTypes[len(def.Params)])
	if name != nil {
		ctx.stopInference(*name)
	}
	for _, p := range def.Params {
		ctx.stopInference(p.Name)
	}
	return ftyp, nil
}

func inferId(id string, ctx *context) (*types.TypePtr, error) {
	def := ctx.getId(id)
	if def == nil {
		act := ctx.getActiveType(id)
		if act != nil {
			return act, nil
		}
		return nil, errors.New("undefined id '" + id + "'")
	}
	if def.v.Type == nil {
		err := inferExp(def.v, def.c, &id)
		if err != nil {
			return nil, err
		}
	}
	return def.v.Type.Copy(), nil
}

func inferBlock(block *ast.Block, ctx *context, name *string) (*types.TypePtr, error) {
	for _, a := range block.Assignments {
		if ctx.getId(a.Name) != nil {
			return nil, errors.New("redefinition of '" + a.Name + "'")
		}
		ctx.ids[a.Name] = &excon{v: a.Value, c: ctx}
	}
	for _, a := range block.Assignments {
		if a.Value.Type == nil {
			err := inferExp(a.Value, ctx, &a.Name)
			if err != nil {
				return nil, err
			}
		}
	}

	err := inferExp(block.Value, ctx, name)
	if err != nil {
		return nil, err
	}
	return block.Value.Type, nil
}
