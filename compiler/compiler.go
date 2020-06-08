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
	malloc        *ir.Func
	freeStructure *ir.Func
	incRef        *ir.Func
	printInt      *ir.Func
	printReal     *ir.Func
	printBool     *ir.Func
	printString   *ir.Func
}

func makeUtils(m *ir.Module) *utils {
	return &utils{
		malloc:        m.NewFunc("malloc", types.I8Ptr, ir.NewParam("size", types.I32)),
		freeStructure: m.NewFunc("free_structure", types.Void, ir.NewParam("ptr", types.I8Ptr)),
		incRef:        m.NewFunc("increase_refcount", types.Void, ir.NewParam("cls", types.I8Ptr)),
		printInt:      m.NewFunc("print_int", types.Void, ir.NewParam("p", IntType)),
		printReal:     m.NewFunc("print_real", types.Void, ir.NewParam("p", RealType)),
		printBool:     m.NewFunc("print_bool", types.Void, ir.NewParam("p", BoolType)),
		printString:   m.NewFunc("print_string", types.Void, ir.NewParam("p", StringType)),
	}
}

func Compile(prg *ast.Exp, fcat *callresolver.FCat) *ir.Module {
	module := ir.NewModule()
	utils := makeUtils(module)

	mainfun := module.NewFunc("main", types.I32)
	global := globalc{Module: module, utils: utils, strings: map[string]value.Value{}}
	ctx := context{Block: mainfun.NewBlock(""), Func: mainfun, global: &global}
	makeFDefs(fcat, &ctx)
	compileFDefs(fcat, &ctx)

	v := compileExp(prg, &ctx, false)

	if prg.Type().AsPrimitive() == t.Int {
		ctx.Block.NewCall(utils.printInt, v.value)
	} else if prg.Type().AsPrimitive() == t.Real {
		ctx.Block.NewCall(utils.printReal, v.value)
	} else if prg.Type().AsPrimitive() == t.Bool {
		ctx.Block.NewCall(utils.printBool, v.value)
	} else if prg.Type().AsPrimitive() == t.String {
		ctx.Block.NewCall(utils.printString, v.value)
	}
	ctx.Block.NewRet(constant.NewInt(types.I32, 0))
	return module
}
