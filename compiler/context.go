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
	functions *map[string]function
	utils     *utils
}

type context struct {
	Module *ir.Module
	Func   *ir.Func
	Block  *ir.Block
	parent *context
	global *globalc
	ids    map[string]value.Value
}

func (c *context) subContext() *context {
	return &context{
		parent: c,
		Module: c.Module,
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
		Module: c.Module,
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
	return (*c.global.functions)[name].From != nil
}

func (c *context) resolveFun(name string) function {
	i := (*c.global.functions)[name]
	if i.Fun == nil {
		panic(name + " is not a function")
	}
	return i
}

func (c *context) makeStructure(struc *Structure) value.Value {
	if struc == nil || len(struc.Fields) == 0 {
		return constant.NewNull(types.I8Ptr)
	}
	ctyp := structureType(struc)
	ctypp := types.NewPointer(ctyp)
	sp := c.Block.NewGetElementPtr(ctyp, constant.NewNull(ctypp), constant.NewInt(types.I32, 1))
	size := c.Block.NewPtrToInt(sp, types.I32)
	mem := c.Block.NewBitCast(c.malloc(size), ctypp)

	// reference count
	refcp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	c.Block.NewStore(constant.NewInt(types.I32, 1), refcp)

	// closure count
	closures := 0
	for _, clj := range struc.Fields {
		if clj.Type.IsFunction() {
			closures++
		}
	}
	clscp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	c.Block.NewStore(constant.NewInt(types.I16, int64(closures)), clscp)

	// structure count
	structures := 0
	for _, clj := range struc.Fields {
		if clj.Type.IsStructure() {
			structures++
		}
	}
	scp := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 2))
	c.Block.NewStore(constant.NewInt(types.I16, int64(structures)), scp)

	closures = 0
	for _, clj := range struc.Fields {
		if clj.Type.IsFunction() {
			res, err := c.resolveId(clj.Name)
			if err != nil {
				panic(err)
			}
			fptr := c.Block.NewExtractElement(res, constant.NewInt(types.I32, 0))
			cptr := c.Block.NewExtractElement(res, constant.NewInt(types.I32, 1))
			c.increfStructure(cptr)
			ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+3)))
			c.Block.NewStore(fptr, ptr)
			closures++
			ptr = c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+3)))
			c.Block.NewStore(cptr, ptr)
			closures++
		}
	}

	structures = 0
	for _, clj := range struc.Fields {
		if clj.Type.IsStructure() {
			ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+structures+3)))
			res, err := c.resolveId(clj.Name)
			if err != nil {
				panic(err)
			}
			c.increfStructure(res)
			c.Block.NewStore(res, ptr)
			structures++
		}
	}

	primitives := 0
	for _, clj := range struc.Fields {
		if !clj.Type.IsFunction() && !clj.Type.IsStructure() {
			ptr := c.Block.NewGetElementPtr(ctyp, mem, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(primitives+structures+3+closures)))
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

func (c *context) loadStructure(struc *Structure, ptr value.Value) {
	if struc == nil || len(struc.Fields) == 0 {
		return
	}
	ctyp := structureType(struc)
	ctypp := types.NewPointer(ctyp)
	cptr := c.Block.NewBitCast(ptr, ctypp)

	closures := 0
	for _, clj := range struc.Fields {
		if clj.Type.IsFunction() {
			fptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+3)))
			fun := c.Block.NewLoad(FunPType, fptr)
			closures++
			cptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+3)))
			cls := c.Block.NewLoad(ClosurePType, cptr)
			closures++
			vec := c.Block.NewInsertElement(constant.NewUndef(FunType), fun, constant.NewInt(types.I32, 0))
			vec = c.Block.NewInsertElement(vec, cls, constant.NewInt(types.I32, 1))

			c.addId(clj.Name, vec)
		}
	}

	structures := 0
	for _, clj := range struc.Fields {
		if clj.Type.IsStructure() {
			ptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(closures+3+structures)))
			r := c.Block.NewLoad(getType(clj.Type), ptr)
			c.addId(clj.Name, r)
			structures++
		}
	}

	primitives := 0
	for _, clj := range struc.Fields {
		if !clj.Type.IsFunction() && !clj.Type.IsStructure() {
			ptr := c.Block.NewGetElementPtr(ctyp, cptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(primitives+structures+closures+3)))
			r := c.Block.NewLoad(getType(clj.Type), ptr)
			c.addId(clj.Name, r)
			primitives++
		}
	}
}

func (c *context) freeStructure(fp value.Value) {
	c.Block.NewCall(c.global.utils.freeStructure, fp)
}

func (c *context) freeClosure(fp value.Value) {
	cptr := c.Block.NewExtractElement(fp, constant.NewInt(types.I32, 1))
	c.Block.NewCall(c.global.utils.freeStructure, cptr)
}

func (c *context) increfStructure(fp value.Value) {
	c.Block.NewCall(c.global.utils.incRef, fp)
}

func (c *context) increfClosure(fp value.Value) {
	cptr := c.Block.NewExtractElement(fp, constant.NewInt(types.I32, 1))
	c.Block.NewCall(c.global.utils.incRef, cptr)
}

func (c *context) call(f value.Value, typ t.Type, params []value.Value) value.Value {
	fptr := c.Block.NewExtractElement(f, constant.NewInt(types.I32, 0))
	cptr := c.Block.NewExtractElement(f, constant.NewInt(types.I32, 1))
	fun := c.Block.NewBitCast(fptr, getFunctPtr(typ))

	return c.Block.NewCall(fun, append(params, cptr)...)
}

func (c *context) ret(v cresult) {
	block := c.Block
	if v.ast.Type().IsFunction() && v.ast.Id != nil && (*c.global.functions)[v.ast.Id.Name].Fun == nil {
		c.increfClosure(v.value)
	} else if v.ast.Type().IsStructure() {
		c.increfStructure(v.value)
	}
	block.NewRet(v.value)
}

func (c *context) malloc(size value.Value) value.Value {
	return c.Block.NewCall(c.global.utils.malloc, size)
}

func (c *context) freeIfUnboundRef(res cresult) {
	if res.ast != nil {
		if res.ast.Type().IsFunction() && res.ast.Id == nil {
			c.freeClosure(res.value)
		} else if res.ast.Type().IsFunction() {
			if c.isFun(res.ast.Id.Name) {
				f := c.resolveFun(res.ast.Id.Name)
				if f.From.HasClosure() {
					c.freeClosure(res.value)
				}
			}
		} else if res.ast.Type().IsStructure() {
			c.freeStructure(res.value)
		}
	}
}
