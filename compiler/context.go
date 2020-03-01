package compiler

import (
	"errors"

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

func (c *context) resolveId(name string) (value.Value, error) {
	if c.constants[name] != nil {
		return c.constants[name], nil
	} else if c.parent != nil {
		return c.parent.resolveId(name)
	}
	return nil, errors.New("undefined id " + name)
}

func (c *context) addId(name string, val value.Value) (*context, error) {
	if c.constants == nil {
		c.constants = map[string]value.Value{}
	}
	if c.constants[name] != nil {
		return nil, errors.New("redefinition of " + name)
	}
	c.constants[name] = val
	return c, nil
}

func (c *context) resolveFun(name string) (*compiledFun, error) {
	if c.functions[name] != nil {
		return c.functions[name], nil
	} else if c.parent != nil {
		return c.parent.resolveFun(name)
	}
	return nil, errors.New("undefined fun " + name)
}

func (c *context) addFun(name string, fun *compiledFun) *context {
	if c.functions == nil {
		c.functions = map[string]*compiledFun{}
	}
	c.functions[name] = fun
	return c
}
