package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/inferer"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func makeFDefs(fcat *inferer.FCat, ctx *context) {
	for name, fun := range *fcat {
		var params []*ir.Param
		for _, p := range fun.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			params = append(params, param)
		}

		rtype := getType(fun.Body.Type)
		compiled := ctx.Module.NewFunc(name, rtype, params...)
		compiled.Linkage = enum.LinkageInternal

		ctx.addId(name, function{fun, compiled})
	}
}

func compileFDefs(fcat *inferer.FCat, ctx *context) {
	for name, _ := range *fcat {
		f := ctx.resolveFun(name)
		body := f.Fun.NewBlock("")
		subCtx := ctx.funcContext(body, f.Fun)
		var params []*ir.Param
		for _, p := range f.From.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			_, err := subCtx.addId(p.Name, val{param})
			if err != nil {
				panic(err)
			}
			params = append(params, param)
		}
		result := compileExp(f.From.Body, subCtx)
		subCtx.Block.NewRet(result)
	}
}

func Compile(prg *ast.Exp, fcat *inferer.FCat) *ir.Module {
	module := ir.NewModule()

	msg := module.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%d\n"))
	printf := module.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true

	mainfun := module.NewFunc("main", types.I32)

	ctx := context{Module: module, Block: mainfun.NewBlock(""), Func: mainfun}
	makeFDefs(fcat, &ctx)
	compileFDefs(fcat, &ctx)

	v := compileExp(prg, &ctx)

	ptr := ctx.Block.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	ctx.Block.NewCall(printf, ptr, v)
	ctx.Block.NewRet(constant.NewInt(types.I32, 0))
	return module
}
