package compiler

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir/value"

	"github.com/llir/llvm/ir"
)

type compiledFun struct {
	From *ast.FDef
	Fun  *ir.Func
}

type compiledValue struct {
	Value value.Value
}

type context struct {
	Module *ir.Module
	Func   *ir.Func
	Block  *ir.Block
	parent *context
	ids    map[string]interface{}
}

func (c *context) subContext() *context {
	return &context{parent: c, Module: c.Module, Block: c.Block, Func: c.Func}
}

func (c *context) blockContext(block *ir.Block, fun *ir.Func) *context {
	return &context{parent: c, Module: c.Module, Block: block, Func: fun}
}

func (c *context) resolveId(name string) (interface{}, error) {
	if c.ids[name] != nil {
		return c.ids[name], nil
	} else if c.parent != nil {
		return c.parent.resolveId(name)
	}
	return nil, errors.New("undefined identifier " + name)
}

func (c *context) addId(name string, val interface{}) (*context, error) {
	if c.ids == nil {
		c.ids = map[string]interface{}{}
	}
	if c.ids[name] != nil {
		return nil, errors.New("redefinition of " + name)
	}
	c.ids[name] = val
	return c, nil
}

func (c *context) resolveFun(name string) (compiledFun, error) {
	i, err := c.resolveId(name)
	if err != nil {
		return compiledFun{}, err
	}
	switch i.(type) {
	case compiledFun:
		return i.(compiledFun), nil
	}
	return compiledFun{}, errors.New(name + " is not a function")
}

func (c *context) resolveVal(name string) (value.Value, error) {
	i, err := c.resolveId(name)
	if err != nil {
		return nil, err
	}
	switch i.(type) {
	case compiledValue:
		return i.(compiledValue).Value, nil
	}
	return nil, errors.New(name + " is not a value")
}

func (c *context) functions() []compiledFun {
	var res []compiledFun
	for _, i := range c.ids {
		switch i.(type) {
		case compiledFun:
			res = append(res, i.(compiledFun))
		}
	}
	return res
}
