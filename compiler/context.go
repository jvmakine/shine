package compiler

import (
	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type compiledFun struct {
	From *grammar.FunDef
	Fun  *ir.Func
}

type context struct {
	parent    *context
	Functions map[string]*compiledFun
	constants map[string]value.Value
}

func (c *context) resolveId(name string) value.Value {
	if c.constants[name] != nil {
		return c.constants[name]
	} else if c.parent != nil {
		c.parent.resolveId(name)
	}
	panic("unknown id: " + name)
}

func (c *context) withConstant(name string, val value.Value) *context {
	if c.constants == nil {
		c.constants = map[string]value.Value{}
	}
	c.constants[name] = val
	return c
}

func (c *context) subContext() *context {
	return &context{parent: c}

}
