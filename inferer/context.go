package inferer

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

type excon struct {
	v *ast.Exp
	c *context
}

type context struct {
	parent *context
	ids    map[string]*excon
	active map[string]*types.TypePtr
}

func (ctx *context) setActiveType(id string, typ *types.TypePtr) {
	if ctx.active == nil {
		ctx.active = map[string]*types.TypePtr{}
	}
	ctx.active[id] = typ
}

func (ctx *context) stopInference(id string) {
	ctx.active[id] = nil
}

func (ctx *context) getActiveType(id string) *types.TypePtr {
	if ctx.active[id] != nil {
		return ctx.active[id]
	}
	if ctx.parent != nil {
		return ctx.parent.getActiveType(id)
	}
	return nil
}
