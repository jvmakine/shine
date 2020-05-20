package compiler

import (
	"github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func compileExp(from *ast.Exp, ctx *context, funcRoot bool) value.Value {
	if from.Const != nil {
		return compileConst(from.Const, ctx)
	} else if from.Id != nil {
		return compileID(from, ctx)
	} else if from.Call != nil {
		return compileCall(from, ctx, funcRoot)
	} else if from.Def != nil {
		panic("non resolved anonymous function: " + from.Type().Signature())
	} else if from.Block != nil {
		return compileBlock(from.Block, ctx, funcRoot)
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

func compileID(exp *ast.Exp, ctx *context) value.Value {
	id, err := ctx.resolveId(exp.Id.Name)
	if err != nil {
		panic(err)
	}
	if f, ok := id.(function); ok {
		nv := ctx.Block.NewBitCast(f.Fun, types.I8Ptr)
		vec := ctx.Block.NewInsertElement(constant.NewUndef(FunType), nv, constant.NewInt(types.I32, 0))
		return vec
	}
	return id.(val).Value
}

func compileIf(c *ast.Exp, t *ast.Exp, f *ast.Exp, ctx *context, funcRoot bool) value.Value {
	trueB := ctx.Func.NewBlock(ctx.newLabel())
	falseB := ctx.Func.NewBlock(ctx.newLabel())
	typ := getType(t.Type())

	cond := compileExp(c, ctx, funcRoot)
	ctx.Block.NewCondBr(cond, trueB, falseB)
	var resV *ir.InstAlloca
	if !funcRoot {
		resV = ctx.Block.NewAlloca(typ)
	}

	ctx.Block = trueB
	truev := compileExp(t, ctx, funcRoot)
	trueB = ctx.Block

	ctx.Block = falseB
	falsev := compileExp(f, ctx, funcRoot)
	falseB = ctx.Block
	if !funcRoot {
		trueB.NewStore(truev, resV)
		falseB.NewStore(falsev, resV)

		continueB := ctx.Func.NewBlock(ctx.newLabel())
		trueB.NewBr(continueB)
		falseB.NewBr(continueB)

		ctx.Block = continueB
		return continueB.NewLoad(typ, resV)
	} else { // optimise root ifs at functions for tail recursion elimination
		if truev != nil {
			compileRet(truev, t.Type(), trueB)
		}
		if falsev != nil {
			compileRet(falsev, f.Type(), falseB)
		}
		return nil
	}
}

func compileCall(exp *ast.Exp, ctx *context, funcRoot bool) value.Value {
	from := exp.Call
	if from.Function.Op != nil {
		var params []value.Value
		name := from.Function.Op.Name
		if name == "if" { // Need to evaluate if parameters lazily
			return compileIf(from.Params[0], from.Params[1], from.Params[2], ctx, funcRoot)
		}
		for _, p := range from.Params {
			v := compileExp(p, ctx, false)
			params = append(params, v)
		}
		switch name {
		case "*":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFMul(params[0], params[1])
			}
			return ctx.Block.NewMul(params[0], params[1])
		case "/":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFDiv(params[0], params[1])
			}
			return ctx.Block.NewUDiv(params[0], params[1])
		case "%":
			return ctx.Block.NewURem(params[0], params[1])
		case "+":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFAdd(params[0], params[1])
			}
			return ctx.Block.NewAdd(params[0], params[1])
		case "-":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFSub(params[0], params[1])
			}
			return ctx.Block.NewSub(params[0], params[1])
		case ">":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOGT, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSGT, params[0], params[1])
		case "<":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOLT, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSLT, params[0], params[1])
		case ">=":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOGE, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSGE, params[0], params[1])
		case "<=":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return ctx.Block.NewFCmp(enum.FPredOLE, params[0], params[1])
			}
			return ctx.Block.NewICmp(enum.IPredSLE, params[0], params[1])
		case "==":
			return ctx.Block.NewICmp(enum.IPredEQ, params[0], params[1])
		case "!=":
			return ctx.Block.NewICmp(enum.IPredNE, params[0], params[1])
		case "||":
			return ctx.Block.NewOr(params[0], params[1])
		case "&&":
			return ctx.Block.NewAnd(params[0], params[1])
		default:
			panic("unknown op " + name)
		}
	} else {
		params := []value.Value{constant.NewIntToPtr(constant.NewInt(types.I64, 0), ClosurePType)}
		name := from.Function.Id.Name
		for _, p := range from.Params {
			v := compileExp(p, ctx, false)
			params = append(params, v)
		}

		id, err := ctx.resolveId(name)
		if err != nil {
			panic(err)
		}
		if f, ok := id.(function); ok {
			for _, p := range *f.From.Closure {
				r, err := ctx.resolveId(p.Name)
				if err != nil {
					panic(err)
				}
				params = append(params, r.(val).Value)
			}
			return ctx.Block.NewCall(f.Call, params...)
		}
		fptr := ctx.Block.NewExtractElement(id.(val).Value, constant.NewInt(types.I32, 0))
		f := ctx.Block.NewBitCast(fptr, getFunctPtr(from.Function.Type()))
		return ctx.Block.NewCall(f, params...)
	}
}

func compileBlock(from *ast.Block, ctx *context, funcRoot bool) value.Value {
	sub := ctx.subContext()
	for k, c := range from.Assignments {
		if c.Def == nil {
			v := compileExp(c, sub, false)
			_, err := sub.addId(k, val{v})
			if err != nil {
				panic(err)
			}
		}
	}
	res := compileExp(from.Value, sub, funcRoot)
	ctx.Block = sub.Block
	return res
}
