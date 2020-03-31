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
		subCtx := ctx.blockContext(body)
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
		body.NewRet(result)
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
	return compileExp(from.Value, sub)
}

/*func compileBlock(block *ir.Block, from *grammar.Block, ctx *context) (value.Value, error) {
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
}*/
