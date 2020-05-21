package compiler

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
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
	utils  *utils
	ids    map[string]interface{}
}

func (c *context) subContext() *context {
	return &context{parent: c, Module: c.Module, Block: c.Block, Func: c.Func, utils: c.utils}
}

var labels = 0

func (c *context) newLabel() string {
	count := labels
	labels++
	return "label_" + strconv.Itoa(count)
}

func (c *context) funcContext(block *ir.Block, fun *ir.Func) *context {
	return &context{parent: c, Module: c.Module, Block: block, Func: fun, utils: c.utils}
}

func (c *context) resolveId(name string) (interface{}, error) {
	if c.ids[name] != nil {
		return c.ids[name], nil
	} else if c.parent != nil {
		return c.parent.resolveId(name)
	}
	return nil, errors.New("undefined id " + name)
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

func (c *context) makeClosure(closure *Closure) value.Value {
	if closure == nil || len(*closure) == 0 {
		return constant.NewNull(types.I8Ptr)
	}
	ctyp := closureType(closure)
	ctypp := types.NewPointer(ctyp)
	sp := c.Block.NewGetElementPtr(ctyp, constant.NewNull(ctypp), constant.NewInt(types.I32, 1))
	size := c.Block.NewPtrToInt(sp, types.I32)
	mem := c.Block.NewBitCast(c.Block.NewCall(c.utils.malloc, size), ctypp)
	for i, clj := range *closure {
		ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(i)))
		res, err := c.resolveId(clj.Name)
		if err != nil {
			panic(err)
		}
		c.Block.NewStore(res.(val).Value, ptr)
	}
	return c.Block.NewBitCast(mem, types.I8Ptr)
}

func (c *context) loadClosure(closure *Closure, ptr value.Value) {
	if closure == nil || len(*closure) == 0 {
		return
	}
	ctyp := closureType(closure)
	ctypp := types.NewPointer(ctyp)
	cptr := c.Block.NewBitCast(ptr, ctypp)
	for i, clj := range *closure {
		ptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(i)))
		r := c.Block.NewLoad(getType(clj.Type), ptr)
		c.addId(clj.Name, val{r})
	}
}

func (c *context) callClosureFunction(f value.Value, typ t.Type, params []value.Value) value.Value {
	fptr := c.Block.NewExtractElement(f, constant.NewInt(types.I32, 0))
	cptr := c.Block.NewExtractElement(f, constant.NewInt(types.I32, 1))
	fun := c.Block.NewBitCast(fptr, getFunctPtr(typ))
	return c.Block.NewCall(fun, append(params, cptr)...)
}
