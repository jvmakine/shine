package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func compileExp(from *ast.Exp, ctx *context) value.Value {
	if from.Const != nil {
		return compileConst(from.Const, ctx)
	} else if from.Id != nil {
		return compileID(*from.Id, ctx)
	} else if from.Call != nil {
		return compileCall(from.Call, ctx)
	} else if from.Def != nil {
		panic("can not return function as a value yet")
	} else if from.Block != nil {
		return compileBlock(from.Block, ctx)
	}
	panic("invalid empty expression")
}

func compileConst(from *ast.Const, ctx *context) value.Value {
	return constant.NewInt(types.I32, int64(*from.Int))
}

func compileID(name string, ctx *context) value.Value {
	id, err := ctx.resolveVal(name)
	if err != nil {
		panic(err)
	}
	return id
}

func compileCall(from *ast.FCall, ctx *context) value.Value {
	name := from.Name
	if name == "if" { // Need to evaluate if parameters lazily
		trueL := ctx.newLabel()
		falseL := ctx.newLabel()
		continueL := ctx.newLabel()
		trueB := ctx.Func.NewBlock(trueL)
		falseB := ctx.Func.NewBlock(falseL)
		continueB := ctx.Func.NewBlock(continueL)
		resV := ctx.Block.NewAlloca(types.I32)

		cond := compileExp(from.Params[0], ctx)
		ctx.Block.NewCondBr(cond, trueB, falseB)

		trueV := compileExp(from.Params[1], ctx.blockContext(trueB))
		trueB.NewStore(trueV, resV)
		trueB.NewBr(continueB)

		falseV := compileExp(from.Params[2], ctx.blockContext(falseB))
		falseB.NewStore(falseV, resV)
		falseB.NewBr(continueB)

		ctx.Block = continueB
		r := continueB.NewLoad(types.I32, resV)
		return r
	} else {
		var params []value.Value
		for _, p := range from.Params {
			v := compileExp(p, ctx)
			params = append(params, v)
		}

		switch name {
		case "*":
			return ctx.Block.NewMul(params[0], params[1])
		case "/":
			return ctx.Block.NewUDiv(params[0], params[1])
		case "+":
			return ctx.Block.NewAdd(params[0], params[1])
		case "-":
			return ctx.Block.NewSub(params[0], params[1])
		case ">":
			return ctx.Block.NewICmp(enum.IPredSGT, params[0], params[1])
		case "<":
			return ctx.Block.NewICmp(enum.IPredSLT, params[0], params[1])
		case ">=":
			return ctx.Block.NewICmp(enum.IPredSGE, params[0], params[1])
		case "<=":
			return ctx.Block.NewICmp(enum.IPredSLE, params[0], params[1])
		case "==":
			return ctx.Block.NewICmp(enum.IPredEQ, params[0], params[1])
		default:
			comp, err := ctx.resolveFun(name)
			if err != nil {
				panic(err)
			}
			return ctx.Block.NewCall(comp.Fun, params...)
		}
	}
}

func makeFDef(name string, fun *ast.FDef, ctx *context) {
	var params []*ir.Param
	for _, p := range fun.Params {
		param := ir.NewParam(p.Name, types.I32)
		params = append(params, param)
	}

	compiled := ctx.Module.NewFunc(name, types.I32, params...)
	compiled.Linkage = enum.LinkageInternal

	ctx.addId(name, compiledFun{fun, compiled})
}

func compileFDefs(ctx *context) {
	for _, f := range ctx.functions() {
		body := f.Fun.NewBlock("")
		subCtx := ctx.funcContext(body, f.Fun)
		var params []*ir.Param
		for _, p := range f.From.Params {
			param := ir.NewParam(p.Name, types.I32)
			_, err := subCtx.addId(p.Name, compiledValue{param})
			if err != nil {
				panic(err)
			}
			params = append(params, param)
		}
		result := compileExp(f.From.Body, subCtx)
		subCtx.Block.NewRet(result)
	}
}

func compileBlock(from *ast.Block, ctx *context) value.Value {
	sub := ctx.subContext()
	for _, c := range from.Assignments {
		if c.Value.Def != nil {
			makeFDef(c.Name, c.Value.Def, ctx)
		}
	}
	compileFDefs(ctx)
	for _, c := range from.Assignments {
		if c.Value.Def == nil {
			v := compileExp(c.Value, sub)
			_, err := sub.addId(c.Name, compiledValue{v})
			if err != nil {
				panic(err)
			}
		}
	}
	res := compileExp(from.Value, sub)
	// TODO: refactor block handling
	ctx.Block = sub.Block
	return res
}
