package types

import (
	"errors"

	"github.com/chewxy/hm"
	"github.com/jvmakine/shine/ast"
)

func globalFun(ts ...hm.Type) *excon {
	return &excon{
		&ast.Exp{Type: hm.NewFnType(ts...)},
		&inferContext{},
	}
}

var globalConsts map[string]*excon = map[string]*excon{
	"+":  globalFun(Int, Int, Int),
	"-":  globalFun(Int, Int, Int),
	"*":  globalFun(Int, Int, Int),
	"/":  globalFun(Int, Int, Int),
	"<":  globalFun(Int, Int, Bool),
	">":  globalFun(Int, Int, Bool),
	">=": globalFun(Int, Int, Bool),
	"<=": globalFun(Int, Int, Bool),
	"==": globalFun(Int, Int, Bool),
	"if": globalFun(Bool, hm.TypeVariable('a'), hm.TypeVariable('a'), hm.TypeVariable('a')),
}

func (ctx *inferContext) getId(id string) *excon {
	if ctx.ids[id] != nil {
		return ctx.ids[id]
	} else if ctx.parent != nil {
		return ctx.parent.getId(id)
	}
	return nil
}

func Infer(exp *ast.Exp) error {
	parent := &inferContext{ids: globalConsts}
	return inferExp(exp, &inferContext{ids: map[string]*excon{}, parent: parent})
}

func inferExp(exp *ast.Exp, ctx *inferContext) error {
	if exp.Type == nil {
		if exp.Const != nil {
			if exp.Const.Int != nil {
				exp.Type = Int
			} else if exp.Const.Bool != nil {
				exp.Type = Bool
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
			typ, err := inferDef(exp.Def, &inferContext{parent: ctx, ids: map[string]*excon{}})
			exp.Type = typ
			return err
		} else if exp.Block != nil {
			typ, err := inferBlock(exp.Block, &inferContext{parent: ctx, ids: map[string]*excon{}})
			exp.Type = typ
			return err
		}
		panic("unexpected expression")
	}
	return nil
}

func inferCall(call *ast.FCall, ctx *inferContext) (hm.Type, error) {
	var params []hm.Type = make([]hm.Type, len(call.Params)+1)
	for i, p := range call.Params {
		err := inferExp(p, ctx)
		if err != nil {
			return nil, err
		}
		params[i] = p.Type
	}
	// Recursive type definition
	if ctx.isInferring(call.Name) {
		return hm.TypeVariable('r'), nil
	}
	ec := ctx.getId(call.Name)
	if ec == nil {
		return nil, errors.New("undefined function: '" + call.Name + "'")
	}
	if ec.v.Type == nil {
		err := inferExp(ec.v, ec.c)
		if err != nil {
			return nil, err
		}
	}
	ft, ok := ec.v.Type.(*hm.FunctionType)
	if !ok {
		return nil, errors.New("not a function: '" + call.Name + "'")
	}
	params[len(call.Params)] = ft.Ret(true)

	ft2 := hm.NewFnType(params...)
	_, err := hm.Unify(ft, ft2)
	if err != nil {
		return nil, err
	}

	return ft.Ret(true), nil
}

func inferDef(def *ast.FDef, ctx *inferContext) (hm.Type, error) {
	var paramTypes []hm.Type = make([]hm.Type, len(def.Params)+1)
	for i, p := range def.Params {
		if ctx.getId(p.Name) != nil {
			return nil, errors.New("redefinition of '" + p.Name + "'")
		}
		ctx.ids[p.Name] = &excon{
			v: &ast.Exp{
				Id:   &p.Name,
				Type: hm.TypeVariable('a' + i),
			},
			c: ctx,
		}
		paramTypes[i] = hm.TypeVariable('a' + i)
	}
	paramTypes[len(def.Params)] = Int
	err := inferExp(def.Body, ctx)
	if err != nil {
		return nil, err
	}
	return hm.NewFnType(paramTypes...), nil
}

func inferId(id string, ctx *inferContext) (hm.Type, error) {
	def := ctx.getId(id)
	if def == nil {
		return nil, errors.New("undefined id '" + id + "'")
	}
	if def.v.Type == nil {
		err := inferExp(def.v, def.c)
		if err != nil {
			return nil, err
		}
	}
	return def.v.Type, nil
}

func inferBlock(block *ast.Block, ctx *inferContext) (hm.Type, error) {
	for _, a := range block.Assignments {
		if ctx.getId(a.Name) != nil {
			return nil, errors.New("redefinition of '" + a.Name + "'")
		}
		ctx.ids[a.Name] = &excon{v: a.Value, c: ctx}
	}
	for _, a := range block.Assignments {
		ctx.startInference(a.Name)
		if a.Value.Type == nil {
			err := inferExp(a.Value, ctx)
			if err != nil {
				return nil, err
			}
		}
		ctx.stopInference(a.Name)
	}

	err := inferExp(block.Value, ctx)
	if err != nil {
		return nil, err
	}
	return block.Value.Type, nil
}
