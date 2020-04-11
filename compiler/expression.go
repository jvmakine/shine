package compiler

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func compileExp(from *ast.Exp, ctx *context) (value.Value, error) {
	if from.Const != nil {
		return compileConst(from.Const, ctx)
	} else if from.Id != nil {
		return compileID(*from.Id, ctx)
	} else if from.Call != nil {
		return compileCall(from.Call, ctx)
	} else if from.Def != nil {
		return nil, errors.New("can not return function as a value yet")
	} else if from.Block != nil {
		return compileBlock(from.Block, ctx)
	}
	panic("invalid empty expression")
}

func compileConst(from *ast.Const, ctx *context) (value.Value, error) {
	return constant.NewInt(types.I32, int64(*from.Int)), nil
}

func compileID(name string, ctx *context) (value.Value, error) {
	id, err := ctx.resolveVal(name)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func compileCall(from *ast.FCall, ctx *context) (value.Value, error) {
	name := from.Name
	if name == "if" { // Need to evaluate if parameters lazily
		trueL := ctx.newLabel()
		falseL := ctx.newLabel()
		continueL := ctx.newLabel()
		trueB := ctx.Func.NewBlock(trueL)
		falseB := ctx.Func.NewBlock(falseL)
		continueB := ctx.Func.NewBlock(continueL)
		resV := ctx.Block.NewAlloca(types.I32)

		cond, err := compileExp(from.Params[0], ctx)
		if err != nil {
			return nil, err
		}
		ctx.Block.NewCondBr(cond, trueB, falseB)

		trueV, err := compileExp(from.Params[1], ctx.blockContext(trueB))
		if err != nil {
			return nil, err
		}
		trueB.NewStore(trueV, resV)
		trueB.NewBr(continueB)

		falseV, err := compileExp(from.Params[2], ctx.blockContext(falseB))
		if err != nil {
			return nil, err
		}
		falseB.NewStore(falseV, resV)
		falseB.NewBr(continueB)

		ctx.Block = continueB
		r := continueB.NewLoad(types.I32, resV)
		return r, nil
	} else {
		var params []value.Value
		for _, p := range from.Params {
			v, err := compileExp(p, ctx)
			if err != nil {
				return nil, err
			}
			params = append(params, v)
		}

		switch name {
		case "*":
			return ctx.Block.NewMul(params[0], params[1]), nil
		case "/":
			return ctx.Block.NewUDiv(params[0], params[1]), nil
		case "+":
			return ctx.Block.NewAdd(params[0], params[1]), nil
		case "-":
			return ctx.Block.NewSub(params[0], params[1]), nil
		case ">":
			return ctx.Block.NewICmp(enum.IPredSGT, params[0], params[1]), nil
		case "<":
			return ctx.Block.NewICmp(enum.IPredSLT, params[0], params[1]), nil
		case "==":
			return ctx.Block.NewICmp(enum.IPredEQ, params[0], params[1]), nil
		default:
			comp, err := ctx.resolveFun(name)
			if err != nil {
				return nil, err
			}

			gotParms := len(from.Params)
			expParms := len(comp.From.Params)
			if gotParms != expParms {
				return nil, errors.New("invalid number of args for " + name + ". Got " + strconv.Itoa(gotParms) + ", expected " + strconv.Itoa(expParms))
			}
			return ctx.Block.NewCall(comp.Fun, params...), nil
		}
	}
}

func makeFDef(name string, fun *ast.FDef, ctx *context) error {
	var params []*ir.Param
	for _, p := range fun.Params {
		param := ir.NewParam(p.Name, types.I32)
		params = append(params, param)
	}

	compiled := ctx.Module.NewFunc(name, types.I32, params...)
	compiled.Linkage = enum.LinkageInternal

	_, err := ctx.addId(name, compiledFun{fun, compiled})
	return err
}

func compileFDefs(ctx *context) error {
	for _, f := range ctx.functions() {
		body := f.Fun.NewBlock("")
		subCtx := ctx.funcContext(body, f.Fun)
		var params []*ir.Param
		for _, p := range f.From.Params {
			param := ir.NewParam(p.Name, types.I32)
			_, err := subCtx.addId(p.Name, compiledValue{param})
			if err != nil {
				return err
			}
			params = append(params, param)
		}
		result, err := compileExp(f.From.Body, subCtx)
		if err != nil {
			return err
		}
		subCtx.Block.NewRet(result)
	}
	return nil
}

func compileBlock(from *ast.Block, ctx *context) (value.Value, error) {
	sub := ctx.subContext()
	for _, c := range from.Assignments {
		if c.Value.Def != nil {
			makeFDef(c.Name, c.Value.Def, ctx)
		}
	}
	err := compileFDefs(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range from.Assignments {
		if c.Value.Def == nil {
			v, err := compileExp(c.Value, sub)
			if err != nil {
				return nil, err
			}
			_, err = sub.addId(c.Name, compiledValue{v})
			if err != nil {
				return nil, err
			}
		}
	}
	res, err := compileExp(from.Value, sub)
	// TODO: refactor block handling
	ctx.Block = sub.Block
	return res, err
}
