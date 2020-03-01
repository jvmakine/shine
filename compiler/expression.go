package compiler

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func compileBlock(block *ir.Block, from *grammar.Block, ctx *context) (value.Value, error) {
	sub := ctx.subContext()
	for _, c := range from.Assignments {
		if c.Value.Term != nil {
			v, err := evalExpression(block, c.Value, sub)
			if err != nil {
				return nil, err
			}
			_, err = sub.addId(*c.Name, compiledValue{v})
			if err != nil {
				return nil, err
			}
		}
	}
	result, err := evalExpression(block, from.Value, sub)
	return result, err
}

func evalValue(block *ir.Block, val *grammar.Value, ctx *context) (value.Value, error) {
	if val.Int != nil {
		return constant.NewInt(types.I32, int64(*val.Int)), nil
	} else if val.Sub != nil {
		return evalExpression(block, val.Sub, ctx)
	} else if val.Call != nil {
		name := *val.Call.Name
		comp, err := ctx.resolveFun(name)
		if err != nil {
			return nil, err
		}

		gotParms := len(val.Call.Params)
		expParms := len(comp.From.Params)
		if gotParms != expParms {
			return nil, errors.New("invalid number of args for " + name + ". Got " + strconv.Itoa(gotParms) + ", expected " + strconv.Itoa(expParms))
		}

		var params []value.Value
		for _, p := range val.Call.Params {
			v, err := evalExpression(block, p, ctx)
			if err != nil {
				return nil, err
			}
			params = append(params, v)
		}
		return block.NewCall(comp.Fun, params...), nil
	} else if val.Id != nil {
		id, err := ctx.resolveVal(*val.Id)
		if err != nil {
			return nil, err
		}
		return id, nil
	} else if val.Block != nil {
		return compileBlock(block, val.Block, ctx)
	}
	panic("invalid value")
}

func evalOpFactor(block *ir.Block, opf *grammar.OpFactor, left value.Value, ctx *context) (value.Value, error) {
	right, err := evalValue(block, opf.Right, ctx)
	if err != nil {
		return nil, err
	}
	switch *opf.Operation {
	case "*":
		return block.NewMul(left, right), nil
	case "/":
		return block.NewUDiv(left, right), nil
	default:
		panic("invalid opfactor: " + *opf.Operation)
	}
}

func evalTerm(block *ir.Block, term *grammar.Term, ctx *context) (value.Value, error) {
	v, err := evalValue(block, term.Left, ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range term.Right {
		v, err = evalOpFactor(block, r, v, ctx)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func evalOpTerm(block *ir.Block, opt *grammar.OpTerm, left value.Value, ctx *context) (value.Value, error) {
	right, err := evalTerm(block, opt.Right, ctx)
	if err != nil {
		return nil, err
	}
	switch *opt.Operation {
	case "+":
		return block.NewAdd(left, right), nil
	case "-":
		return block.NewSub(left, right), nil
	default:
		panic("invalid opterm: " + *opt.Operation)
	}
}

func evalExpression(block *ir.Block, prg *grammar.Expression, ctx *context) (value.Value, error) {
	if prg.Term != nil {
		t := prg.Term
		v, err := evalTerm(block, t.Left, ctx)
		if err != nil {
			return nil, err
		}
		for _, r := range t.Right {
			v, err = evalOpTerm(block, r, v, ctx)
			if err != nil {
				return nil, err
			}
		}
		return v, nil
	}
	return nil, errors.New("function can not be used as return value")
}
