package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/callresolver"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type utils struct {
	malloc      *ir.Func
	freeClosure *ir.Func
	incRef      *ir.Func
}

func makeUtils(m *ir.Module) *utils {
	return &utils{
		malloc:      m.NewFunc("malloc", types.I8Ptr, ir.NewParam("size", types.I32)),
		freeClosure: m.NewFunc("free_closure", types.Void, ir.NewParam("ptr", types.I8Ptr)),
		incRef:      m.NewFunc("increase_refcount", types.Void, ir.NewParam("cls", types.I8Ptr)),
	}
}

func iPrintF(m *ir.Module, b *ir.Block) (*ir.Func, *ir.InstGetElementPtr) {
	msg := m.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%ld\n"))
	printf := m.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true
	ptr := b.NewGetElementPtr(types.NewArray(4, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	return printf, ptr
}

func fPrintF(m *ir.Module, b *ir.Block) (*ir.Func, *ir.InstGetElementPtr) {
	msg := m.NewGlobalDef("realFormat", constant.NewCharArrayFromString("%f\n"))
	printf := m.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true
	ptr := b.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	return printf, ptr
}

func Compile(prg *ast.Exp, fcat *callresolver.FCat) *ir.Module {
	module := ir.NewModule()
	utils := makeUtils(module)

	mainfun := module.NewFunc("main", types.I32)
	ctx := context{Module: module, Block: mainfun.NewBlock(""), Func: mainfun, utils: utils}
	makeFDefs(fcat, &ctx)
	compileFDefs(fcat, &ctx)

	v := compileExp(prg, &ctx, false)

	if prg.Type().AsPrimitive() == t.Int {
		printf, ptr := iPrintF(module, ctx.Block)
		ctx.Block.NewCall(printf, ptr, v.value)
	} else if prg.Type().AsPrimitive() == t.Real {
		printf, ptr := fPrintF(module, ctx.Block)
		ctx.Block.NewCall(printf, ptr, v.value)
	}
	ctx.Block.NewRet(constant.NewInt(types.I32, 0))
	return module
}
