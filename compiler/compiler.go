package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/callresolver"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type utils struct {
	malloc *ir.Func
	freeRc *ir.Func
	incRef *ir.Func

	printInt    *ir.Func
	printReal   *ir.Func
	printBool   *ir.Func
	printString *ir.Func

	PVEqual16   *ir.Func
	PVCombine16 *ir.Func
}

func makeUtils(m *ir.Module) *utils {
	return &utils{
		malloc:      m.NewFunc("heap_malloc", types.I8Ptr, ir.NewParam("size", types.I32)),
		freeRc:      m.NewFunc("free_rc", types.Void, ir.NewParam("ptr", types.I8Ptr)),
		incRef:      m.NewFunc("increase_refcount", types.Void, ir.NewParam("cls", types.I8Ptr)),
		printInt:    m.NewFunc("print_int", types.Void, ir.NewParam("p", IntType)),
		printReal:   m.NewFunc("print_real", types.Void, ir.NewParam("p", RealType)),
		printBool:   m.NewFunc("print_bool", types.Void, ir.NewParam("p", BoolType)),
		printString: m.NewFunc("print_string", types.Void, ir.NewParam("p", StringPType)),
		PVEqual16:   m.NewFunc("pv_uint16_equals", types.I8, ir.NewParam("s1", types.I8Ptr), ir.NewParam("s2", types.I8Ptr)),
		PVCombine16: m.NewFunc("pv_concatenate", types.I8Ptr, ir.NewParam("l", types.I8Ptr), ir.NewParam("r", types.I8Ptr)),
	}
}

func Compile(prg ast.Expression, fcat *callresolver.FCat) *ir.Module {
	module := ir.NewModule()
	utils := makeUtils(module)

	mainfun := module.NewFunc("main", types.I32)
	global := globalc{Module: module, utils: utils, strings: map[string]value.Value{}}
	ctx := context{Block: mainfun.NewBlock(""), Func: mainfun, global: &global}
	makeFDefs(fcat, &ctx)
	compileFDefs(fcat, &ctx)

	v := compileExp(prg, &ctx, false)
	typ := prg.Type()
	prim, isPrim := typ.(t.Primitive)

	if isPrim {
		if prim.ID == "int" {
			ctx.Block.NewCall(utils.printInt, v.value)
		} else if prim.ID == "real" {
			ctx.Block.NewCall(utils.printReal, v.value)
		} else if prim.ID == "bool" {
			ctx.Block.NewCall(utils.printBool, v.value)
		} else if prim.ID == "string" {
			ctx.Block.NewCall(utils.printString, v.value)
		}
	}
	ctx.Block.NewRet(constant.NewInt(types.I32, 0))
	return module
}
