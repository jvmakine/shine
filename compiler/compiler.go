package compiler

import (
	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func evalFunDef(module *ir.Module, fun *grammar.FunDef, ctx *context) error {
	var params []*ir.Param
	subCtx := ctx.subContext()
	for _, p := range fun.Params {
		param := ir.NewParam(*p.Name, types.I32)
		subCtx.addId(*p.Name, param)
		params = append(params, param)
	}

	compiled := module.NewFunc(*fun.Name, types.I32, params...)

	body := compiled.NewBlock("")
	result, err := evalExpression(body, fun.Body, subCtx)
	if err != nil {
		return err
	}
	body.NewRet(result)

	ctx.addFun(*fun.Name, &compiledFun{fun, compiled})
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
		evalFunDef(module, f, &ctx)
	}
	v, err := evalExpression(entry, prg.Exp, &ctx)
	if err != nil {
		return nil, err
	}

	ptr := entry.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(printf, ptr, v)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module, nil
}
