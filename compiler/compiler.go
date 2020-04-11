package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func Compile(prg *ast.Exp) (*ir.Module, error) {
	module := ir.NewModule()

	msg := module.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%d\n"))
	printf := module.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true

	mainfun := module.NewFunc("main", types.I32)

	ctx := context{Module: module, Block: mainfun.NewBlock(""), Func: mainfun}

	v, err := compileExp(prg, &ctx)
	if err != nil {
		return nil, err
	}

	ptr := ctx.Block.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	ctx.Block.NewCall(printf, ptr, v)
	ctx.Block.NewRet(constant.NewInt(types.I32, 0))
	return module, nil
}
