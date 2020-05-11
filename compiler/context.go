package compiler

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir/value"

	"github.com/llir/llvm/ir"
)

type function struct {
	// Source AST
	From *ast.FDef
	// Compiled IR function
	Fun *ir.Func
	// Value used to call the function
	Call value.Value
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
}

func (c *context) subContext() *context {
	return &context{parent: c, Module: c.Module, Block: c.Block, Func: c.Func}
}

var labels = 0

func (c *context) newLabel() string {
	count := labels
	labels++
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

func (c *context) resolveFun(name string) function {
	i, err := c.resolveId(name)
	if err != nil {
		panic(name + " is not a function")
	}
	switch i.(type) {
	case function:
		return i.(function)
	}
	panic(name + " is not a function")
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
