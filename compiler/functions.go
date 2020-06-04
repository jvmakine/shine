package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/passes/callresolver"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
)

func makeFDefs(fcat *callresolver.FCat, ctx *context) {
	ctx.functions = &map[string]function{}
	for name, fun := range *fcat {
		if fun.Def != nil {
			def := fun.Def
			rtype := getType(def.Body.Type())
			params := []*ir.Param{}
			for _, p := range def.Params {
				param := ir.NewParam(p.Name, getType(p.Type))
				params = append(params, param)
			}
			params = append(params, ir.NewParam("+cls", ClosurePType))
			compiled := ctx.Module.NewFunc(name, rtype, params...)
			compiled.Linkage = enum.LinkageInternal

			(*ctx.functions)[name] = function{def, compiled, compiled}
		} else if fun.Struct != nil {
			stru := fun.Struct
			params := []*ir.Param{}
			for _, p := range stru.Fields {
				param := ir.NewParam(p.Name, getType(p.Type))
				params = append(params, param)
			}
			rtype := getType(stru.Type)
			compiled := ctx.Module.NewFunc(name, rtype, params...)
			compiled.Linkage = enum.LinkageInternal

			(*ctx.functions)[name] = function{nil, compiled, compiled}
		}
	}
}

func compileFDefs(fcat *callresolver.FCat, ctx *context) {
	for name, v := range *fcat {
		if v.Def != nil {
			f := ctx.resolveFun(name)
			body := f.Fun.NewBlock("")
			subCtx := ctx.funcContext(body, f.Fun)
			for _, p := range v.Def.Params {
				param := ir.NewParam(p.Name, getType(p.Type))
				_, err := subCtx.addId(p.Name, param)
				if err != nil {
					panic(err)
				}
			}
			subCtx.loadStructure(v.Def.Closure, ir.NewParam("+cls", ClosurePType))
			result := compileExp(v.Def.Body, subCtx, true)
			if result.value != nil { // result can be nil if it has already been returned from the function
				subCtx.ret(result)
			}
		} else if v.Struct != nil {
			f := ctx.resolveFun(name)
			body := f.Fun.NewBlock("")
			subCtx := ctx.funcContext(body, f.Fun)
			for _, p := range v.Struct.Fields {
				param := ir.NewParam(p.Name, getType(p.Type))
				_, err := subCtx.addId(p.Name, param)
				if err != nil {
					panic(err)
				}
			}
			s := subCtx.makeStructure(v.Struct.Type.Structure)
			subCtx.ret(makeCR(&ast.Exp{Struct: v.Struct}, s))
		}
	}
}
