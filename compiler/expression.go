package compiler

import (
	"github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func compileExp(from *ast.Exp, ctx *context) value.Value {
	if from.Const != nil {
		return compileConst(from.Const, ctx)
	} else if from.Id != nil {
		return compileID(from.Id.Name, ctx)
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
		return constant.NewInt(IntType, *from.Int)
	} else if from.Bool != nil {
		return constant.NewBool(*from.Bool)
	} else if from.Real != nil {
		return constant.NewFloat(RealType, *from.Real)
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
	typ := getType(t.Type)
	resV := ctx.Block.NewAlloca(typ)

	cond := compileExp(c, ctx)
	ctx.Block.NewCondBr(cond, trueB, falseB)

	ctx.Block = trueB
	v := compileExp(t, ctx)
	trueB = ctx.Block
	trueB.NewStore(v, resV)

	ctx.Block = falseB
	v = compileExp(f, ctx)
	falseB = ctx.Block
	falseB.NewStore(v, resV)

	continueB := ctx.Func.NewBlock(ctx.newLabel())
	trueB.NewBr(continueB)
	falseB.NewBr(continueB)

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
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFMul(params[0], params[1])
			}
			return ctx.Block.NewMul(params[0], params[1])
		case "/":
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFDiv(params[0], params[1])
			}
			return ctx.Block.NewUDiv(params[0], params[1])
		case "%":
			return ctx.Block.NewURem(params[0], params[1])
		case "+":
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFAdd(params[0], params[1])
			}
			return ctx.Block.NewAdd(params[0], params[1])
		case "-":
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFSub(params[0], params[1])
			}
			return ctx.Block.NewSub(params[0], params[1])
		case ">":
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOGT, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSGT, params[0], params[1])
		case "<":
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOLT, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSLT, params[0], params[1])
		case ">=":
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOGE, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSGE, params[0], params[1])
		case "<=":
			if from.Params[0].Type.AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOLE, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSLE, params[0], params[1])
		case "==":
			return ctx.Block.NewICmp(enum.IPredEQ, params[0], params[1])
		case "||":
			return ctx.Block.NewOr(params[0], params[1])
		case "&&":
			return ctx.Block.NewAnd(params[0], params[1])
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
