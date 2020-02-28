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
	functions map[string]*compiledFun
	constants map[string]value.Value
}

func (c *context) subContext() *context {
	return &context{parent: c}
}

func (c *context) resolveId(name string) value.Value {
	if c.constants[name] != nil {
		return c.constants[name]
	} else if c.parent != nil {
		c.parent.resolveId(name)
	}
	panic("unknown id: " + name)
}

func (c *context) addId(name string, val value.Value) *context {
	if c.constants == nil {
		c.constants = map[string]value.Value{}
	}
	c.constants[name] = val
	return c
}

func (c *context) resolveFun(name string) *compiledFun {
	if c.functions[name] != nil {
		return c.functions[name]
	} else if c.parent != nil {
		c.parent.resolveFun(name)
	}
	panic("unknown fun: " + name)
}

func (c *context) addFun(name string, fun *compiledFun) *context {
	if c.functions == nil {
		c.functions = map[string]*compiledFun{}
	}
	c.functions[name] = fun
	return c
}
