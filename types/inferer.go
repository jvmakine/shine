package types

import (
	"errors"

	"github.com/chewxy/hm"
	"github.com/jvmakine/shine/ast"
)

type excon struct {
	v *ast.Exp
	c *inferContext
}

type inferContext struct {
	parent *inferContext
	ids    map[string]*excon
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
	return inferExp(exp, &inferContext{ids: map[string]*excon{}})
}

func inferExp(exp *ast.Exp, ctx *inferContext) error {
	if exp.Const != nil {
		exp.Type = Int
		return nil
	} else if exp.Id != nil {
		typ, err := inferId(*exp.Id, ctx)
		exp.Type = typ
		return err
	} else if exp.Call != nil {
		return nil
	} else if exp.Def != nil {
		return nil
	} else if exp.Block != nil {
		typ, err := inferBlock(exp.Block, &inferContext{parent: ctx, ids: map[string]*excon{}})
		exp.Type = typ
		return err
	}
	panic("unexpected expression")
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
		if a.Value.Type == nil {
			err := inferExp(a.Value, ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	err := inferExp(block.Value, ctx)
	if err != nil {
		return nil, err
	}
	return block.Value.Type, nil
}
