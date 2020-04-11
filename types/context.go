package types

import "github.com/jvmakine/shine/ast"

type excon struct {
	v *ast.Exp
	c *inferContext
}

type inferContext struct {
	parent *inferContext
	ids    map[string]*excon
	active map[string]bool
}

func (ctx *inferContext) startInference(id string) {
	if ctx.active == nil {
		ctx.active = map[string]bool{}
	}
	ctx.active[id] = true
}

func (ctx *inferContext) stopInference(id string) {
	ctx.active[id] = false
}

func (ctx *inferContext) isInferring(id string) bool {
	if ctx.active[id] {
		return true
	}
	if ctx.parent != nil {
		return ctx.parent.isInferring(id)
	}
	return false
}
