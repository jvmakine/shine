package compiler

import (
	"strconv"

	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func evalValue(block *ir.Block, val *grammar.Value, ctx *globalContext) value.Value {
	if val.Int != nil {
		return constant.NewInt(types.I32, int64(*val.Int))
	} else if val.Sub != nil {
		return evalExpression(block, val.Sub, ctx)
	} else if val.Call != nil {
		name := *val.Call.Name
		if ctx.Functions[name] == nil {
			panic("unknown function: " + name)
		}
		gotParms := len(val.Call.Params)
		expParms := len(ctx.Functions[name].From.Params)
		if gotParms != expParms {
			panic("invalid number of args for " + name + ". Got " + strconv.Itoa(gotParms) + ", expected " + strconv.Itoa(expParms))
		}

		var params []value.Value
		for _, p := range val.Call.Params {
			v := evalExpression(block, p, ctx)
			params = append(params, v)
		}
		return block.NewCall(ctx.Functions[name].Fun, params...)
	} else if val.Id != nil {
		return constant.NewInt(types.I32, 1)
	}
	panic("invalid value")
}

func evalOpFactor(block *ir.Block, opf *grammar.OpFactor, left value.Value, ctx *globalContext) value.Value {
	right := evalValue(block, opf.Right, ctx)
	switch *opf.Operation {
	case "*":
		return block.NewMul(left, right)
	case "/":
		return block.NewUDiv(left, right)
	default:
		panic("invalid opfactor: " + *opf.Operation)
	}
}

func evalTerm(block *ir.Block, term *grammar.Term, ctx *globalContext) value.Value {
	v := evalValue(block, term.Left, ctx)
	for _, r := range term.Right {
		v = evalOpFactor(block, r, v, ctx)
	}
	return v
}

func evalOpTerm(block *ir.Block, opt *grammar.OpTerm, left value.Value, ctx *globalContext) value.Value {
	right := evalTerm(block, opt.Right, ctx)
	switch *opt.Operation {
	case "+":
		return block.NewAdd(left, right)
	case "-":
		return block.NewSub(left, right)
	default:
		panic("invalid opterm: " + *opt.Operation)
	}
}

func evalExpression(block *ir.Block, prg *grammar.Expression, ctx *globalContext) value.Value {
	v := evalTerm(block, prg.Left, ctx)
	for _, r := range prg.Right {
		v = evalOpTerm(block, r, v, ctx)
	}
	return v
}
