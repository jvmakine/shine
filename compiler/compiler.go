package compiler

import (
	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func makeFunDefinition(module *ir.Module, fun *grammar.FunDef, ctx *context) error {
	var params []*ir.Param
	for _, p := range fun.Params {
		param := ir.NewParam(*p.Name, types.I32)
		params = append(params, param)
	}

	compiled := module.NewFunc(*fun.Name, types.I32, params...)
	compiled.Linkage = enum.LinkageInternal

	ctx.addFun(*fun.Name, &compiledFun{fun, compiled})
	return nil
}

func compileBlock(block *ir.Block, from *grammar.Block, ctx *context) (value.Value, error) {
	for _, c := range from.Assignments {
		v, err := evalExpression(block, c.Value, ctx)
		if err != nil {
			return nil, err
		}
		ctx.addId(*c.Name, v)
	}
	result, err := evalExpression(block, from.Value, ctx)
	return result, err
}

func compileFunBodies(ctx *context) error {
	for _, f := range ctx.functions {
		body := f.Fun.NewBlock("")
		subCtx := ctx.subContext()
		var params []*ir.Param
		for _, p := range f.From.Params {
			param := ir.NewParam(*p.Name, types.I32)
			subCtx.addId(*p.Name, param)
			params = append(params, param)
		}
		result, err := compileBlock(body, f.From.Body, subCtx)
		if err != nil {
			return err
		}
		body.NewRet(result)
	}
	return nil
}

func Compile(prg *grammar.Program) (*ir.Module, error) {
	module := ir.NewModule()

	msg := module.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%d\n"))
	printf := module.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true

	mainfun := module.NewFunc("main", types.I32)
	entry := mainfun.NewBlock("")

	ctx := context{}
	for _, f := range prg.Functions {
		err := makeFunDefinition(module, f, &ctx)
		if err != nil {
			return nil, err
		}
	}
	compileFunBodies(&ctx)
	v, err := compileBlock(entry, prg.Body, &ctx)
	if err != nil {
		return nil, err
	}

	ptr := entry.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(printf, ptr, v)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module, nil
}
