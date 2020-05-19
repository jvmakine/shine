package compiler

import (
	"github.com/jvmakine/shine/passes/callresolver"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
)

func makeFDefs(fcat *callresolver.FCat, ctx *context) {
	for name, fun := range *fcat {
		rtype := getType(fun.Body.Type())
		params := []*ir.Param{ir.NewParam("+cls", ClosurePType)}
		if fun.Body.Type().IsFunction() {
			params = append(params, ir.NewParam("+rc", ClosureRType))
		}
		for _, p := range fun.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			params = append(params, param)
		}
		// TODO: handle closure in non simple cases
		for _, p := range *fun.Closure {
			param := ir.NewParam(p.Name, getType(p.Type))
			params = append(params, param)
		}

		compiled := ctx.Module.NewFunc(name, rtype, params...)
		compiled.Linkage = enum.LinkageInternal

		ctx.addId(name, function{fun, compiled, compiled})
	}
}

func compileFDefs(fcat *callresolver.FCat, ctx *context) {
	for name, _ := range *fcat {
		f := ctx.resolveFun(name)
		body := f.Fun.NewBlock("")
		subCtx := ctx.funcContext(body, f.Fun)
		closureParam := ir.NewParam("+cls", ClosurePType)
		params := []*ir.Param{closureParam}
		if f.From.Body.Type().IsFunction() {
			params = append(params, ir.NewParam("+rc", ClosureRType))
		}
		for _, p := range f.From.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			_, err := subCtx.addId(p.Name, val{param})
			if err != nil {
				panic(err)
			}
			params = append(params, param)
		}
		// TODO: handle closure in non simple cases
		for _, p := range *f.From.Closure {
			param := ir.NewParam(p.Name, getType(p.Type))
			_, err := subCtx.addId(p.Name, val{param})
			if err != nil {
				panic(err)
			}
			params = append(params, param)
		}
		result := compileExp(f.From.Body, subCtx, true)
		if result != nil { // result can be nil if it has already been returned from the function
			subCtx.Block.NewRet(result)
		}
	}
}
