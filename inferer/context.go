package inferer

import (
	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

type excon struct {
	v *ast.Exp
	c *context
}

type context struct {
	parent     *context
	blockID    int
	ids        map[string]*excon
	active     map[string]*Type
	activeVals *[]string
}

func (c *context) sub(bid int) *context {
	return &context{parent: c, ids: map[string]*excon{}, active: map[string]*Type{}, blockID: bid}
}

func (ctx *context) setActiveType(id string, typ *Type) {
	if ctx.active == nil {
		ctx.active = map[string]*Type{}
	}
	ctx.active[id] = typ
}

func (ctx *context) stopInference(id string) {
	ctx.active[id] = nil
}

func (ctx *context) getActiveType(id string) *Type {
	if ctx.active[id] != nil {
		return ctx.active[id]
	}
	if ctx.parent != nil {
		return ctx.parent.getActiveType(id)
	}
	return nil
}
