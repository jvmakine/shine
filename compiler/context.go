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

type context struct {
	Module    *ir.Module
	Func      *ir.Func
	Block     *ir.Block
	parent    *context
	functions *map[string]function
	utils     *utils
	ids       map[string]value.Value
}

func (c *context) subContext() *context {
	return &context{
		parent:    c,
		Module:    c.Module,
		Block:     c.Block,
		Func:      c.Func,
		utils:     c.utils,
		functions: c.functions,
	}
}

var labels = 0

func (c *context) newLabel() string {
	count := labels
	labels++
	return "label_" + strconv.Itoa(count)
}

func (c *context) funcContext(block *ir.Block, fun *ir.Func) *context {
	return &context{
		parent:    c,
		Module:    c.Module,
		Block:     block,
		Func:      fun,
		utils:     c.utils,
		functions: c.functions,
	}
}

func (c *context) resolveId(name string) (value.Value, error) {
	if c.ids[name] != nil {
		return c.ids[name], nil
	} else if c.parent != nil {
		return c.parent.resolveId(name)
	}
	return nil, errors.New("undefined id " + name)
}

func (c *context) addId(name string, val value.Value) (*context, error) {
	if c.ids == nil {
		c.ids = map[string]value.Value{}
	}
	if c.ids[name] != nil {
		return nil, errors.New("redefinition of " + name)
	}
	c.ids[name] = val
	return c, nil
}

func (c *context) resolveFun(name string) function {
	i := (*c.functions)[name]
	if i.Fun == nil {
		panic(name + " is not a function")
	}
	return i
}

func (c *context) makeClosure(closure *Closure) value.Value {
	if closure == nil || len(*closure) == 0 {
		return constant.NewNull(types.I8Ptr)
	}
	ctyp := closureType(closure)
	ctypp := types.NewPointer(ctyp)
	sp := c.Block.NewGetElementPtr(ctyp, constant.NewNull(ctypp), constant.NewInt(types.I32, 1))
	size := c.Block.NewPtrToInt(sp, types.I32)
	mem := c.Block.NewBitCast(c.malloc(size), ctypp)

	// reference count
	refcp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	c.Block.NewStore(constant.NewInt(types.I32, 1), refcp)

	// closure count
	closures := 0
	for _, clj := range *closure {
		if clj.Type.IsFunction() {
			closures++
		}
	}
	clscp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	c.Block.NewStore(constant.NewInt(types.I16, int64(closures)), clscp)

	closures = 0
	for _, clj := range *closure {
		if clj.Type.IsFunction() {
			ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+2)))
			fptr := c.Block.NewBitCast(ptr, types.NewPointer(FunType))
			f := c.Block.NewLoad(FunType, fptr)
			res, err := c.resolveId(clj.Name)
			if err != nil {
				panic(err)
			}
			cptr := c.Block.NewExtractElement(f, constant.NewInt(types.I32, 1))
			c.Block.NewCall(c.utils.incRef, cptr)
			c.Block.NewStore(res, ptr)
			closures++
		}
	}

	nonclosures := 0
	for _, clj := range *closure {
		if !clj.Type.IsFunction() {
			ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(nonclosures+2+closures)))
			res, err := c.resolveId(clj.Name)
			if err != nil {
				panic(err)
			}
			c.Block.NewStore(res, ptr)
			nonclosures++
		}
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

	closures := 0
	for _, clj := range *closure {
		if clj.Type.IsFunction() {
			ptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+2)))
			r := c.Block.NewLoad(getType(clj.Type), ptr)
			c.addId(clj.Name, r)
			closures++
		}
	}

	nonclosures := 0
	for _, clj := range *closure {
		if !clj.Type.IsFunction() {
			ptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(nonclosures+closures+2)))
			r := c.Block.NewLoad(getType(clj.Type), ptr)
			c.addId(clj.Name, r)
			nonclosures++
		}
	}
}

func (c *context) freeClosure(fp value.Value) {
	cptr := c.Block.NewExtractElement(fp, constant.NewInt(types.I32, 1))
	// TODO: Reference counting
	c.freeClosure(cptr)
}

func (c *context) call(f value.Value, typ t.Type, params []value.Value) value.Value {
	fptr := c.Block.NewExtractElement(f, constant.NewInt(types.I32, 0))
	cptr := c.Block.NewExtractElement(f, constant.NewInt(types.I32, 1))
	fun := c.Block.NewBitCast(fptr, getFunctPtr(typ))
	return c.Block.NewCall(fun, append(params, cptr)...)
}

func (c *context) ret(v value.Value) {
	_, isfunc := v.Type().(*types.FuncType)
	block := c.Block
	if isfunc {
		nv := block.NewBitCast(v, types.I8Ptr)
		vec := block.NewInsertElement(constant.NewUndef(FunType), nv, constant.NewInt(types.I32, 0))
		block.NewRet(vec)
	} else {
		block.NewRet(v)
	}
}

func (c *context) malloc(size value.Value) value.Value {
	return c.Block.NewCall(c.utils.malloc, size)
}
