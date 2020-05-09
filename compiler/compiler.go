package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/callresolver"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func makeFDefs(fcat *callresolver.FCat, ctx *context) {
	for name, fun := range *fcat {
		rtype := getType(fun.Body.Type())
		var params []*ir.Param
		for _, p := range fun.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			params = append(params, param)
		}

		/*if len(fun.Resolved.Closure) > 0 {
			param := ir.NewParam("%%closure", types.I8Ptr)
			param.Attrs = append(param.Attrs, ir.AttrString("nest"))
			params = append(params, param)
		}*/

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
		var params []*ir.Param
		for _, p := range f.From.Params {
			param := ir.NewParam(p.Name, getType(p.Type))
			_, err := subCtx.addId(p.Name, val{param})
			if err != nil {
				panic(err)
			}
			params = append(params, param)
		}
		/*if len(f.From.Resolved.Closure) > 0 {
			cparam := ir.NewParam("%%closure", ClosurePType)
			for _, c := range f.From.Resolved.Closure {
				t := getType(c.Type)
				v := subCtx.Block.NewLoad(t, cparam)
				subCtx.addId(c.Name, val{v})
			}
		}*/
		result := compileExp(f.From.Body, subCtx, true)
		if result != nil { // result can be nil if it has already been returned from the function
			subCtx.Block.NewRet(result)
		}
	}
}

func IPrintF(m *ir.Module, b *ir.Block) (*ir.Func, *ir.InstGetElementPtr) {
	msg := m.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%ld\n"))
	printf := m.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true
	ptr := b.NewGetElementPtr(types.NewArray(4, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	return printf, ptr
}

func FPrintF(m *ir.Module, b *ir.Block) (*ir.Func, *ir.InstGetElementPtr) {
	msg := m.NewGlobalDef("realFormat", constant.NewCharArrayFromString("%f\n"))
	printf := m.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true
	ptr := b.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	return printf, ptr
}

func Compile(prg *ast.Exp, fcat *callresolver.FCat) *ir.Module {
	module := ir.NewModule()

	mainfun := module.NewFunc("main", types.I32)
	ctx := context{Module: module, Block: mainfun.NewBlock(""), Func: mainfun}
	makeFDefs(fcat, &ctx)
	compileFDefs(fcat, &ctx)

	v := compileExp(prg, &ctx, false)

	if prg.Type().AsPrimitive() == t.Int {
		printf, ptr := IPrintF(module, ctx.Block)
		ctx.Block.NewCall(printf, ptr, v)
	} else if prg.Type().AsPrimitive() == t.Real {
		printf, ptr := FPrintF(module, ctx.Block)
		ctx.Block.NewCall(printf, ptr, v)
	}
	ctx.Block.NewRet(constant.NewInt(types.I32, 0))
	return module
}
