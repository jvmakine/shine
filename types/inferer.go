package types

import (
	"errors"

	"github.com/jvmakine/shine/ast"
)

func fun(ts ...interface{}) *excon {
	result := make([]*Type, len(ts))
	var variables map[string]*Type = map[string]*Type{}
	for _, t := range ts {
		switch v := t.(type) {
		case string:
			if variables[v] == nil {
				variables[v] = &Type{}
			}
		}
	}

	for i, t := range ts {
		switch v := t.(type) {
		case *Type:
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

var global map[string]*excon = map[string]*excon{
	"+":  fun(Int, Int, Int),
	"-":  fun(Int, Int, Int),
	"*":  fun(Int, Int, Int),
	"%":  fun(Int, Int, Int),
	"/":  fun(Int, Int, Int),
	"<":  fun(Int, Int, Bool),
	">":  fun(Int, Int, Bool),
	">=": fun(Int, Int, Bool),
	"<=": fun(Int, Int, Bool),
	"==": fun(Int, Int, Bool),
	"if": fun(Bool, "A", "A", "A"),
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
	return inferExp(exp, &context{ids: map[string]*excon{}, parent: parent})
}

func inferExp(exp *ast.Exp, ctx *context) error {
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
			typ, err := inferDef(exp.Def, &context{parent: ctx, ids: map[string]*excon{}})
			exp.Type = typ
			return err
		} else if exp.Block != nil {
			typ, err := inferBlock(exp.Block, &context{parent: ctx, ids: map[string]*excon{}})
			exp.Type = typ
			return err
		}
		panic("unexpected expression")
	}
	return nil
}

func inferCall(call *ast.FCall, ctx *context) (*Type, error) {
	var params []*Type = make([]*Type, len(call.Params)+1)
	for i, p := range call.Params {
		err := inferExp(p, ctx)
		if err != nil {
			return nil, err
		}
		params[i] = p.Type.(*Type)
	}
	// Recursive type definition
	if ctx.isInferring(call.Name) {
		return &Type{}, nil
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
	ft := ec.v.Type.(*Type)
	if !ft.isFunction() {
		return nil, errors.New("not a function: '" + call.Name + "'")
	}
	params[len(call.Params)] = ft.returnType()

	ft2 := function(params...)
	err := unify(&ft2, &ft)
	if err != nil {
		return nil, err
	}

	return ft2.returnType(), nil
}

func inferDef(def *ast.FDef, ctx *context) (*Type, error) {
	var paramTypes []*Type = make([]*Type, len(def.Params)+1)
	for i, p := range def.Params {
		if ctx.getId(p.Name) != nil {
			return nil, errors.New("redefinition of '" + p.Name + "'")
		}
		ctx.ids[p.Name] = &excon{
			v: &ast.Exp{
				Id:   &p.Name,
				Type: Int,
			},
			c: ctx,
		}
		paramTypes[i] = Int
	}
	paramTypes[len(def.Params)] = Int
	err := inferExp(def.Body, ctx)
	if err != nil {
		return nil, err
	}
	return function(paramTypes...), nil
}

func inferId(id string, ctx *context) (*Type, error) {
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
	return def.v.Type.(*Type), nil
}

func inferBlock(block *ast.Block, ctx *context) (*Type, error) {
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
	return block.Value.Type.(*Type), nil
}
