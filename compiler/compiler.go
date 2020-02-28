package compiler

import (
	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type compiledFun struct {
	From *grammar.FunDef
	Fun  *ir.Func
}

type globalContext struct {
	Functions map[string]*compiledFun
}

func evalFunDef(module *ir.Module, fun *grammar.FunDef, ctx *globalContext) {
	var params []*ir.Param
	for _, p := range fun.Params {
		params = append(params, ir.NewParam(*p.Name, types.I32))
	}

	compiled := module.NewFunc(*fun.Name, types.I32, params...)

	body := compiled.NewBlock("")
	result := evalExpression(body, fun.Body, ctx)
	body.NewRet(result)

	ctx.Functions[*fun.Name] = &compiledFun{fun, compiled}
}

func Compile(prg *grammar.Program) *ir.Module {
	module := ir.NewModule()

	msg := module.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%d\n"))
	printf := module.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true

	mainfun := module.NewFunc("main", types.I32)
	entry := mainfun.NewBlock("")

	ctx := globalContext{map[string]*compiledFun{}}
	for _, f := range prg.Functions {
		evalFunDef(module, f, &ctx)
	}
	v := evalExpression(entry, prg.Exp, &ctx)

	ptr := entry.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(printf, ptr, v)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module
}
