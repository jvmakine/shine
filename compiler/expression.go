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

func compileExp(from *ast.Exp, ctx *context, funcRoot bool) cresult {
	if from.Const != nil {
		return compileConst(from, ctx)
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

func compileConst(from *ast.Exp, ctx *context) cresult {
	if from.Const.Int != nil {
		return makeCR(from, constant.NewInt(IntType, *from.Const.Int))
	} else if from.Const.Bool != nil {
		return makeCR(from, constant.NewBool(*from.Const.Bool))
	} else if from.Const.Real != nil {
		return makeCR(from, constant.NewFloat(RealType, *from.Const.Real))
	}
	panic("invalid constant at compilation")
}

func compileID(exp *ast.Exp, ctx *context) cresult {
	name := exp.Id.Name
	if (*ctx.functions)[name].Fun != nil {
		f := (*ctx.functions)[name]
		nv := ctx.Block.NewBitCast(f.Fun, types.I8Ptr)
		clj := ctx.makeClosure(f.From.Closure)
		vec := ctx.Block.NewInsertElement(constant.NewUndef(FunType), nv, constant.NewInt(types.I32, 0))
		vec = ctx.Block.NewInsertElement(vec, clj, constant.NewInt(types.I32, 1))
		return makeCR(exp, vec)
	}
	id, err := ctx.resolveId(name)
	if err != nil {
		panic(err)
	}
	return makeCR(exp, id)
}

func compileIf(c *ast.Exp, t *ast.Exp, f *ast.Exp, ctx *context, funcRoot bool) cresult {
	trueB := ctx.Func.NewBlock(ctx.newLabel())
	falseB := ctx.Func.NewBlock(ctx.newLabel())
	typ := getType(t.Type())

	cond := compileExp(c, ctx, funcRoot)
	ctx.Block.NewCondBr(cond.value, trueB, falseB)
	var resV *ir.InstAlloca
	if !funcRoot {
		resV = ctx.Block.NewAlloca(typ)
	}

	ctx.Block = trueB
	truev := compileExp(t, ctx, funcRoot)
	if t.Type().IsFunction() && t.Id == nil {
		ctx.freeClosure(truev.value)
	}
	if cond.ast.Type().IsFunction() && cond.ast.Id == nil {
		ctx.freeClosure(cond.value)
	}

	trueB = ctx.Block
	if funcRoot && truev.value != nil {
		ctx.ret(makeCR(c, truev.value).cmb(cond))
	}

	ctx.Block = falseB
	falsev := compileExp(f, ctx, funcRoot)
	falseB = ctx.Block
	if f.Type().IsFunction() && f.Id == nil {
		ctx.freeClosure(falsev.value)
	}
	if cond.ast.Type().IsFunction() && cond.ast.Id == nil {
		ctx.freeClosure(cond.value)
	}

	if funcRoot && falsev.value != nil {
		ctx.ret(makeCR(c, falsev.value).cmb(cond))
	}

	if !funcRoot {
		trueB.NewStore(truev.value, resV)
		falseB.NewStore(falsev.value, resV)

		continueB := ctx.Func.NewBlock(ctx.newLabel())
		trueB.NewBr(continueB)
		falseB.NewBr(continueB)

		ctx.Block = continueB
		return makeCR(c, continueB.NewLoad(typ, resV)).cmb(cond)
	} else { // optimise root ifs at functions for tail recursion elimination
		return cresult{}
	}
}

func compileCall(exp *ast.Exp, ctx *context, funcRoot bool) cresult {
	from := exp.Call
	if from.Function.Op != nil {
		var params []cresult
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
				return makeCR(exp, ctx.Block.NewFMul(params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewMul(params[0].value, params[1].value)).cmb(params...)
		case "/":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFDiv(params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewUDiv(params[0].value, params[1].value)).cmb(params...)
		case "%":
			return makeCR(exp, ctx.Block.NewURem(params[0].value, params[1].value)).cmb(params...)
		case "+":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFAdd(params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewAdd(params[0].value, params[1].value)).cmb(params...)
		case "-":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFSub(params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewSub(params[0].value, params[1].value)).cmb(params...)
		case ">":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOGT, params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSGT, params[0].value, params[1].value)).cmb(params...)
		case "<":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOLT, params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSLT, params[0].value, params[1].value)).cmb(params...)
		case ">=":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOGE, params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSGE, params[0].value, params[1].value)).cmb(params...)
		case "<=":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOLE, params[0].value, params[1].value)).cmb(params...)
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSLE, params[0].value, params[1].value)).cmb(params...)
		case "==":
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredEQ, params[0].value, params[1].value)).cmb(params...)
		case "!=":
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredNE, params[0].value, params[1].value)).cmb(params...)
		case "||":
			return makeCR(exp, ctx.Block.NewOr(params[0].value, params[1].value)).cmb(params...)
		case "&&":
			return makeCR(exp, ctx.Block.NewAnd(params[0].value, params[1].value)).cmb(params...)
		default:
			panic("unknown op " + name)
		}
	} else {
		params := []cresult{}
		for _, p := range from.Params {
			v := compileExp(p, ctx, false)
			params = append(params, v)
		}

		vparams := make([]value.Value, len(params))
		for i, p := range params {
			vparams[i] = p.value
		}

		if from.Function.Id != nil {
			name := from.Function.Id.Name
			if (*ctx.functions)[name].Fun != nil {
				f := (*ctx.functions)[name]
				vps := make([]value.Value, len(params))
				for i, p := range params {
					vps[i] = p.value
				}
				res := ctx.Block.NewCall(f.Call, append(vps, constant.NewNull(ClosurePType))...)
				for _, p := range params {
					if p.ast.Type().IsFunction() && p.ast.Id == nil {
						ctx.freeClosure(p.value)
					}
				}
				return makeCR(exp, res)
			}
			id, err := ctx.resolveId(name)
			if err != nil {
				panic(err)
			}
			res := ctx.call(id, from.Function.Type(), vparams)
			for _, p := range params {
				if p.ast.Type().IsFunction() && p.ast.Id == nil {
					ctx.freeClosure(p.value)
				}
			}
			return makeCR(exp, res)
		}
		fval := compileExp(from.Function, ctx, false)
		res := ctx.call(fval.value, from.Function.Type(), vparams)
		for _, p := range params {
			if p.ast.Type().IsFunction() && p.ast.Id == nil {
				ctx.freeClosure(p.value)
			}
		}
		ctx.freeClosure(fval.value)
		return makeCR(exp, res)
	}
}

func compileBlock(from *ast.Block, ctx *context, funcRoot bool) cresult {
	sub := ctx.subContext()

	assigns := map[string]*ast.Exp{}
	deps := map[string]map[string]bool{}
	for k, c := range from.Assignments {
		assigns[k] = c
		deps[k] = map[string]bool{}
		for _, i := range c.CollectIds() {
			deps[k][i] = true
		}
	}

	memids := map[string]value.Value{}

	for len(assigns) > 0 {
		for k, c := range assigns {
			dependencies := false
			for d, _ := range deps[k] {
				if assigns[d] != nil {
					dependencies = true
					break
				}
			}
			if !dependencies {
				v := compileExp(c, sub, false)
				_, err := sub.addId(k, v.value)
				if err != nil {
					panic(err)
				}
				if c.Type().IsFunction() {
					memids[k] = v.value

					if c.Id != nil { // TODO: Optimise renames away
						sub.increfClosure(v.value)
					}
				}
				delete(assigns, k)
			}
		}
	}

	res := compileExp(from.Value, sub, funcRoot)
	for id, v := range memids {
		if from.Value.Id == nil || from.Value.Id.Name != id {
			sub.freeClosure(v)
		}
	}
	ctx.Block = sub.Block
	return res
}
