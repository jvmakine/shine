package compiler

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir/value"

	"github.com/llir/llvm/ir"
)

type function struct {
	From *ast.FDef
	Fun  *ir.Func
}

type val struct {
	Value value.Value
}

type context struct {
	Module *ir.Module
	Func   *ir.Func
	Block  *ir.Block
	parent *context
	ids    map[string]interface{}
	labels int
}

func (c *context) subContext() *context {
	return &context{parent: c, Module: c.Module, Block: c.Block, Func: c.Func}
}

func (c *context) newLabel() string {
	count := c.labels
	c.labels = c.labels + 1
	return "label_" + strconv.Itoa(count)
}

func (c *context) funcContext(block *ir.Block, fun *ir.Func) *context {
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

func (c *context) resolveFun(name string) (function, error) {
	i, err := c.resolveId(name)
	if err != nil {
		return function{}, err
	}
	switch i.(type) {
	case function:
		return i.(function), nil
	}
	return function{}, errors.New(name + " is not a function")
}

func (c *context) resolveVal(name string) (value.Value, error) {
	i, err := c.resolveId(name)
	if err != nil {
		return nil, err
	}
	switch i.(type) {
	case val:
		return i.(val).Value, nil
	}
	return nil, errors.New(name + " is not a value")
}

func (c *context) functions() []function {
	var res []function
	for _, i := range c.ids {
		switch i.(type) {
		case function:
			res = append(res, i.(function))
		}
	}
	return res
}
