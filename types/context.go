package types

import "github.com/jvmakine/shine/ast"

type excon struct {
	v *ast.Exp
	c *context
}

type context struct {
	parent *context
	ids    map[string]*excon
	active map[string]bool
}

func (ctx *context) startInference(id string) {
	if ctx.active == nil {
		ctx.active = map[string]bool{}
	}
	ctx.active[id] = true
}

func (ctx *context) stopInference(id string) {
	ctx.active[id] = false
}

func (ctx *context) isInferring(id string) bool {
	if ctx.active[id] {
		return true
	}
	if ctx.parent != nil {
		return ctx.parent.isInferring(id)
	}
	return false
}
