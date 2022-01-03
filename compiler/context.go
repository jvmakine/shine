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

type globalc struct {
	functions map[string]function
	strings   map[string]value.Value
	utils     *utils
	Module    *ir.Module
}

type context struct {
	Func   *ir.Func
	Block  *ir.Block
	parent *context
	global *globalc
	ids    map[string]value.Value
}

func (c *context) subContext() *context {
	return &context{
		parent: c,
		Block:  c.Block,
		Func:   c.Func,
		global: c.global,
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
		parent: c,
		Block:  block,
		Func:   fun,
		global: c.global,
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

func (c *context) isFun(name string) bool {
	return c.global.functions[name].Fun != nil
}

func (c *context) resolveFun(name string) function {
	i := c.global.functions[name]
	if i.Fun == nil {
		panic(name + " is not a function")
	}
	return i
}

func (c *context) makeStructure(struc *Structure, fun value.Value) value.Value {
	if struc == nil {
		return constant.NewNull(types.I8Ptr)
	}
	ctyp := structureType(struc, fun != nil)
	ctypp := types.NewPointer(ctyp)
	sp := c.Block.NewGetElementPtr(ctyp, constant.NewNull(ctypp), constant.NewInt(types.I32, 1))
	size := c.Block.NewPtrToInt(sp, types.I32)
	mem := c.Block.NewBitCast(c.malloc(size), ctypp)

	extra := 3
	if fun != nil {
		extra = 4
	}

	// reference type
	reftp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	rtype := 2
	if fun != nil {
		rtype = 1
	}
	c.Block.NewStore(constant.NewInt(types.I8, int64(rtype)), reftp)

	// reference count
	refcp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	c.Block.NewStore(constant.NewInt(types.I32, 1), refcp)

	if fun != nil {
		fptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 2))
		f := c.Block.NewBitCast(fun, types.I8Ptr)
		c.Block.NewStore(f, fptr)
	}

	// structure count
	structures := 0
	for _, clj := range struc.Fields {
		if clj.Type.IsStructure() || clj.Type.IsString() || clj.Type.IsFunction() {
			structures++
		}
	}
	scp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(extra-1)))
	c.Block.NewStore(constant.NewInt(types.I16, int64(structures)), scp)

	structures = 0
	for _, clj := range struc.Fields {
		if clj.Type.IsStructure() || clj.Type.IsString() || clj.Type.IsFunction() {
			ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(structures+extra)))
			res, err := c.resolveId(clj.Name)
			if err != nil {
				panic(err)
			}
			c.incRef(res)
			c.Block.NewStore(res, ptr)
			structures++
		}
	}

	primitives := 0
	for _, clj := range struc.Fields {
		if !clj.Type.IsFunction() && !clj.Type.IsStructure() && !clj.Type.IsString() {
			ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(primitives+structures+extra)))
			res, err := c.resolveId(clj.Name)
			if err != nil {
				panic(err)
			}
			c.Block.NewStore(res, ptr)
			primitives++
		}
	}
	return c.Block.NewBitCast(mem, types.I8Ptr)
}

func (c *context) loadClosure(struc *Structure, ptr value.Value) {
	if struc == nil || len(struc.Fields) == 0 {
		return
	}
	ctyp := structureType(struc, true)
	ctypp := types.NewPointer(ctyp)
	cptr := c.Block.NewBitCast(ptr, ctypp)

	closures := 0
	for _, clj := range struc.Fields {
		if clj.Type.IsFunction() {
			fptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+4)))
			fun := c.Block.NewLoad(FunType, fptr)
			closures++
			c.addId(clj.Name, fun)
		}
	}

	structures := 0
	for _, clj := range struc.Fields {
		if clj.Type.IsStructure() || clj.Type.IsString() {
			ptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+4+structures)))
			r := c.Block.NewLoad(getType(clj.Type), ptr)
			c.addId(clj.Name, r)
			structures++
		}
	}

	primitives := 0
	for _, clj := range struc.Fields {
		if !clj.Type.IsFunction() && !clj.Type.IsStructure() && !clj.Type.IsString() {
			ptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(primitives+structures+closures+4)))
			r := c.Block.NewLoad(getType(clj.Type), ptr)
			c.addId(clj.Name, r)
			primitives++
		}
	}
}

func (c *context) freeRef(fp value.Value) {
	c.Block.NewCall(c.global.utils.freeRc, fp)
}

func (c *context) incRef(fp value.Value) {
	c.Block.NewCall(c.global.utils.incRef, fp)
}

func (c *context) call(f value.Value, typ t.Type, params []value.Value) value.Value {
	bc := c.Block.NewBitCast(f, types.NewPointer(ClosureCallType))
	ptr := c.Block.NewGetElementPtr(ClosureCallType, bc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 2))
	fptr := c.Block.NewLoad(FunType, ptr)
	fun := c.Block.NewBitCast(fptr, getFunctPtr(typ))

	return c.Block.NewCall(fun, append(params, f)...)
}

func (c *context) ret(v cresult) {
	block := c.Block
	if v.ast.Type().IsFunction() && v.ast.Id != nil && c.global.functions[v.ast.Id.Name].Fun == nil {
		c.incRef(v.value)
	} else if v.ast.Type().IsStructure() || v.ast.Type().IsString() {
		c.incRef(v.value)
	}
	block.NewRet(v.value)
}

func (c *context) malloc(size value.Value) value.Value {
	return c.Block.NewCall(c.global.utils.malloc, size)
}

func (c *context) freeIfUnboundRef(res cresult) {
	if res.ast != nil {
		if res.ast.Type().IsFunction() && res.ast.Id == nil {
			c.freeRef(res.value)
		} else if res.ast.Type().IsFunction() {
			if c.isFun(res.ast.Id.Name) {
				f := c.resolveFun(res.ast.Id.Name)
				if f.From.HasClosure() {
					c.freeRef(res.value)
				}
			}
		} else if (res.ast.Type().IsStructure() || res.ast.Type().IsString()) && res.ast.Id == nil {
			c.freeRef(res.value)
		}
	}
}
