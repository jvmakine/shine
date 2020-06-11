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
	} else if from.TDecl != nil {
		return compileExp(from.TDecl.Exp, ctx, funcRoot)
	} else if from.Struct != nil {
		panic("non resolved struct at compilation")
	} else if from.FAccess != nil {
		return compileFAccess(from, ctx)
	}
	panic("invalid empty expression")
}

func getStructFieldIndex(s *t.Structure, name string) int {
	index := 0
	found := false
	for _, v := range s.Fields {
		if v.Type.IsFunction() {
			if v.Name == name {
				found = true
				break
			}
			index++
		}
	}
	if !found {
		for _, v := range s.Fields {
			if v.Type.IsStructure() || v.Type.IsString() {
				if v.Name == name {
					found = true
					break
				}
				index++
			}
		}
	}
	if !found {
		for _, v := range s.Fields {
			if !v.Type.IsStructure() && !v.Type.IsFunction() {
				if v.Name == name {
					found = true
					break
				}
				index++
			}
		}
	}
	if !found {
		panic("field not found: " + name)
	}
	return index
}

func compileFAccess(from *ast.Exp, ctx *context) cresult {
	fa := from.FAccess
	cstru := compileExp(fa.Exp, ctx, false)
	tstru := fa.Exp.Type()
	ctyp := structureType(tstru.Structure, false)
	typ := types.NewPointer(ctyp)
	bc := ctx.Block.NewBitCast(cstru.value, typ)
	index := getStructFieldIndex(tstru.Structure, fa.Field)
	ptr := ctx.Block.NewGetElementPtr(ctyp, bc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(index+3)))
	res := ctx.Block.NewLoad(getType(from.Type()), ptr)
	return makeCR(from, res)
}

func compileConst(from *ast.Exp, ctx *context) cresult {
	if from.Const.Int != nil {
		return makeCR(from, constant.NewInt(IntType, *from.Const.Int))
	} else if from.Const.Bool != nil {
		return makeCR(from, constant.NewBool(*from.Const.Bool))
	} else if from.Const.Real != nil {
		return makeCR(from, constant.NewFloat(RealType, *from.Const.Real))
	} else if from.Const.String != nil {
		sref := ctx.makeStringRefRoot(*from.Const.String)
		bc := ctx.Block.NewBitCast(sref, StringPType)
		return makeCR(from, bc)
	}
	panic("invalid constant at compilation")
}

func compileID(exp *ast.Exp, ctx *context) cresult {
	name := exp.Id.Name
	if ctx.isFun(name) {
		f := ctx.global.functions[name]
		clj := ctx.makeStructure(f.From.Closure, f.Fun)
		return makeCR(exp, clj)
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
	ctx.freeIfUnboundRef(truev)
	ctx.freeIfUnboundRef(cond)

	trueB = ctx.Block
	if funcRoot && truev.value != nil {
		ctx.ret(makeCR(c, truev.value))
	}

	ctx.Block = falseB
	falsev := compileExp(f, ctx, funcRoot)
	falseB = ctx.Block
	ctx.freeIfUnboundRef(falsev)
	ctx.freeIfUnboundRef(cond)

	if funcRoot && falsev.value != nil {
		ctx.ret(makeCR(c, falsev.value))
	}

	if !funcRoot {
		trueB.NewStore(truev.value, resV)
		falseB.NewStore(falsev.value, resV)

		continueB := ctx.Func.NewBlock(ctx.newLabel())
		trueB.NewBr(continueB)
		falseB.NewBr(continueB)

		ctx.Block = continueB
		return makeCR(c, continueB.NewLoad(typ, resV))
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
				return makeCR(exp, ctx.Block.NewFMul(params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewMul(params[0].value, params[1].value))
		case "/":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFDiv(params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewUDiv(params[0].value, params[1].value))
		case "%":
			return makeCR(exp, ctx.Block.NewURem(params[0].value, params[1].value))
		case "+":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFAdd(params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewAdd(params[0].value, params[1].value))
		case "-":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFSub(params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewSub(params[0].value, params[1].value))
		case ">":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOGT, params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSGT, params[0].value, params[1].value))
		case "<":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOLT, params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSLT, params[0].value, params[1].value))
		case ">=":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOGE, params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSGE, params[0].value, params[1].value))
		case "<=":
			if from.Params[0].Type().AsPrimitive() == t.Real {
				return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOLE, params[0].value, params[1].value))
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredSLE, params[0].value, params[1].value))
		case "==":
			if from.Params[0].Type().IsString() {
				v := ctx.Block.NewCall(ctx.global.utils.stringsEqual, params[0].value, params[1].value)
				r := ctx.Block.NewICmp(enum.IPredEQ, v, constant.NewInt(types.I8, int64(1)))
				return makeCR(exp, r)
			}
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredEQ, params[0].value, params[1].value))
		case "!=":
			return makeCR(exp, ctx.Block.NewICmp(enum.IPredNE, params[0].value, params[1].value))
		case "||":
			return makeCR(exp, ctx.Block.NewOr(params[0].value, params[1].value))
		case "&&":
			return makeCR(exp, ctx.Block.NewAnd(params[0].value, params[1].value))
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
			if ctx.global.functions[name].Fun != nil {
				f := ctx.global.functions[name]
				vps := make([]value.Value, len(params))
				for i, p := range params {
					vps[i] = p.value
				}
				res := ctx.Block.NewCall(f.Call, append(vps, constant.NewNull(ClosurePType))...)
				for _, p := range params {
					ctx.freeIfUnboundRef(p)
				}
				return makeCR(exp, res)
			}
			id, err := ctx.resolveId(name)
			if err != nil {
				panic(err)
			}
			res := ctx.call(id, from.Function.Type(), vparams)
			for _, p := range params {
				ctx.freeIfUnboundRef(p)
			}
			return makeCR(exp, res)
		}
		fval := compileExp(from.Function, ctx, false)
		res := ctx.call(fval.value, from.Function.Type(), vparams)
		for _, p := range params {
			ctx.freeIfUnboundRef(p)
		}
		ctx.freeIfUnboundRef(fval)
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
		for _, i := range collectDeps(c, ctx) {
			deps[k][i] = true
		}
	}

	closureids := map[string]value.Value{}
	structids := map[string]value.Value{}

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
					closureids[k] = v.value
					if c.Id != nil { // TODO: Optimise renames away
						sub.incRef(v.value)
					}
				} else if c.Type().IsStructure() || c.Type().IsString() {
					structids[k] = v.value
					if c.Id != nil { // TODO: Optimise renames away
						sub.incRef(v.value)
					}
				}
				delete(assigns, k)
			}
		}
	}

	res := compileExp(from.Value, sub, funcRoot)
	for id, v := range closureids {
		if from.Value.Id == nil || from.Value.Id.Name != id {
			sub.freeRef(v)
		}
	}
	for id, v := range structids {
		if from.Value.Id == nil || from.Value.Id.Name != id {
			sub.freeRef(v)
		}
	}
	ctx.Block = sub.Block
	return res
}

func collectDeps(exp *ast.Exp, c *context) []string {
	ids := map[string]bool{}
	exp.Visit(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Id != nil {
			name := v.Id.Name
			ids[name] = true
			if v.Type().IsFunction() && c.isFun(name) {
				f := c.resolveFun(name)
				if f.From.HasClosure() {
					for _, c := range f.From.Closure.Fields {
						ids[c.Name] = true
					}
				}
			}
		}
		return nil
	})
	result := []string{}
	for k := range ids {
		result = append(result, k)
	}
	return result
}
