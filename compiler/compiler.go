package compiler

import (
	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func makeFunDefinition(module *ir.Module, name string, fun *grammar.FunDef, ctx *context) error {
	var params []*ir.Param
	for _, p := range fun.Params {
		param := ir.NewParam(*p.Name, types.I32)
		params = append(params, param)
	}

	compiled := module.NewFunc(name, types.I32, params...)
	compiled.Linkage = enum.LinkageInternal

	_, err := ctx.addId(name, compiledFun{fun, compiled})
	return err
}

func compileFunBodies(ctx *context) error {
	for _, f := range ctx.functions() {
		body := f.Fun.NewBlock("")
		subCtx := ctx.subContext()
		var params []*ir.Param
		for _, p := range f.From.Params {
			param := ir.NewParam(*p.Name, types.I32)
			_, err := subCtx.addId(*p.Name, compiledValue{param})
			if err != nil {
				return err
			}
			params = append(params, param)
		}
		result, err := evalExpression(body, f.From.Body, subCtx)
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
	for _, f := range prg.Body.Assignments {
		if f.Value.Fun != nil {
			err := makeFunDefinition(module, *f.Name, f.Value.Fun, &ctx)
			if err != nil {
				return nil, err
			}
		}
	}
	err := compileFunBodies(&ctx)
	if err != nil {
		return nil, err
	}
	v, err := compileBlock(entry, prg.Body, &ctx)
	if err != nil {
		return nil, err
	}

	ptr := entry.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(printf, ptr, v)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module, nil
}
