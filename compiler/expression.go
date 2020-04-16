package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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
	if from.Int != nil {
		return constant.NewInt(IntType, int64(*from.Int))
	} else if from.Bool != nil {
		return constant.NewBool(*from.Bool)
	}
	panic("invalid constant at compilation")
}

func compileID(name string, ctx *context) value.Value {
	id, err := ctx.resolveVal(name)
	if err != nil {
		panic(err)
	}
	return id
}

func compileIf(c *ast.Exp, t *ast.Exp, f *ast.Exp, ctx *context) value.Value {
	trueB := ctx.Func.NewBlock(ctx.newLabel())
	falseB := ctx.Func.NewBlock(ctx.newLabel())
	continueB := ctx.Func.NewBlock(ctx.newLabel())
	typ := getType(t.Type)
	resV := ctx.Block.NewAlloca(typ)

	cond := compileExp(c, ctx)
	ctx.Block.NewCondBr(cond, trueB, falseB)

	ctx.Block = trueB
	ctx.Block.NewStore(compileExp(t, ctx), resV)
	ctx.Block.NewBr(continueB)

	ctx.Block = falseB
	ctx.Block.NewStore(compileExp(f, ctx), resV)
	ctx.Block.NewBr(continueB)

	ctx.Block = continueB
	return continueB.NewLoad(typ, resV)
}

func compileCall(from *ast.FCall, ctx *context) value.Value {
	name := from.Name
	if name == "if" { // Need to evaluate if parameters lazily
		return compileIf(from.Params[0], from.Params[1], from.Params[2], ctx)
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
		case "%":
			return ctx.Block.NewURem(params[0], params[1])
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
			comp := ctx.resolveFun(from.Resolved)
			return ctx.Block.NewCall(comp.Fun, params...)
		}
	}
}

func compileBlock(from *ast.Block, ctx *context) value.Value {
	sub := ctx.subContext()
	for _, c := range from.Assignments {
		if c.Value.Def == nil {
			v := compileExp(c.Value, sub)
			_, err := sub.addId(c.Name, val{v})
			if err != nil {
				panic(err)
			}
		}
	}
	res := compileExp(from.Value, sub)
	ctx.Block = sub.Block
	return res
}
