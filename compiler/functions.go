package compiler

import (
	"github.com/jvmakine/shine/passes/callresolver"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func makeFDefs(fcat *callresolver.FCat, ctx *context) {
	for name, fun := range *fcat {
		rtype := getType(fun.Body.Type())
		params := []*ir.Param{}
		for _, p := range fun.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			params = append(params, param)
		}
		params = append(params, ir.NewParam("+cls", ClosurePType))
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
		for _, p := range f.From.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			_, err := subCtx.addId(p.Name, val{param})
			if err != nil {
				panic(err)
			}
		}
		subCtx.loadClosure(f.From.Closure, ir.NewParam("+cls", ClosurePType))
		result := compileExp(f.From.Body, subCtx, true)
		if result != nil { // result can be nil if it has already been returned from the function
			compileRet(result, f.From.Body.Type(), subCtx.Block)
		}
	}
}

func compileRet(v value.Value, typ t.Type, block *ir.Block) {
	_, isvect := v.Type().(*types.VectorType)
	if typ.IsFunction() && !isvect {
		nv := block.NewBitCast(v, types.I8Ptr)
		vec := block.NewInsertElement(constant.NewUndef(FunType), nv, constant.NewInt(types.I32, 0))
		block.NewRet(vec)
	} else {
		block.NewRet(v)
	}
}
