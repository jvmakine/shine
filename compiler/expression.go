package compiler

import (
	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func compileExp(from ast.Expression, ctx *context, funcRoot bool) cresult {
	if c, ok := from.(*ast.Const); ok {
		return compileConst(c, ctx)
	} else if i, ok := from.(*ast.Id); ok {
		return compileID(i, ctx)
	} else if c, ok := from.(*ast.FCall); ok {
		return compileCall(c, ctx, funcRoot)
	} else if d, ok := from.(*ast.FDef); ok {
		panic("non resolved anonymous function: " + Signature(d.Type()))
	} else if b, ok := from.(*ast.Block); ok {
		return compileBlock(b, ctx, funcRoot)
	} else if t, ok := from.(*ast.TypeDecl); ok {
		return compileExp(t.Exp, ctx, funcRoot)
	} else if _, ok := from.(*ast.Struct); ok {
		panic("non resolved struct at compilation")
	} else if a, ok := from.(*ast.FieldAccessor); ok {
		return compileFAccess(a, ctx)
	} else if p, ok := from.(*ast.PrimitiveOp); ok {
		return compilePrimitiveOp(p, ctx)
	}
	panic("invalid empty expression")
}

func compilePrimitiveOp(from *ast.PrimitiveOp, ctx *context) cresult {
	left := compileExp(from.Left, ctx, false)
	right := compileExp(from.Right, ctx, false)

	var res cresult
	switch from.ID {
	case "int_+":
		res = makeCR(from, ctx.Block.NewAdd(left.value, right.value))
	case "int_*":
		res = makeCR(from, ctx.Block.NewMul(left.value, right.value))
	case "int_-":
		res = makeCR(from, ctx.Block.NewSub(left.value, right.value))
	case "int_%":
		res = makeCR(from, ctx.Block.NewURem(left.value, right.value))
	case "int_/":
		res = makeCR(from, ctx.Block.NewUDiv(left.value, right.value))
	case "int_>":
		res = makeCR(from, ctx.Block.NewICmp(enum.IPredSGT, left.value, right.value))
	case "int_<":
		res = makeCR(from, ctx.Block.NewICmp(enum.IPredSLT, left.value, right.value))
	case "int_>=":
		res = makeCR(from, ctx.Block.NewICmp(enum.IPredSGE, left.value, right.value))
	case "int_<=":
		res = makeCR(from, ctx.Block.NewICmp(enum.IPredSLE, left.value, right.value))
	case "int_==":
		res = makeCR(from, ctx.Block.NewICmp(enum.IPredEQ, left.value, right.value))
	case "int_!=":
		res = makeCR(from, ctx.Block.NewICmp(enum.IPredNE, left.value, right.value))
	case "real_+":
		res = makeCR(from, ctx.Block.NewFAdd(left.value, right.value))
	case "real_*":
		res = makeCR(from, ctx.Block.NewFMul(left.value, right.value))
	case "real_-":
		res = makeCR(from, ctx.Block.NewFSub(left.value, right.value))
	case "real_/":
		res = makeCR(from, ctx.Block.NewFDiv(left.value, right.value))
	case "string_+":
		res = makeCR(from, ctx.Block.NewCall(ctx.global.utils.PVCombine16, left.value, right.value))
	default:
		panic("unknown primary op " + from.ID)
	}
	ctx.freeIfUnboundRef(left)
	ctx.freeIfUnboundRef(right)
	return res
}

func compileBinOp(from *ast.FCall, exp ast.Expression, op string, params []cresult, ctx *context) cresult {
	switch op {
	case ">":
		if from.Params[0].Type() == Real {
			return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOGT, params[0].value, params[1].value))
		}
		return makeCR(exp, ctx.Block.NewICmp(enum.IPredSGT, params[0].value, params[1].value))
	case "<":
		if from.Params[0].Type() == Real {
			return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOLT, params[0].value, params[1].value))
		}
		return makeCR(exp, ctx.Block.NewICmp(enum.IPredSLT, params[0].value, params[1].value))
	case ">=":
		if from.Params[0].Type() == Real {
			return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOGE, params[0].value, params[1].value))
		}
		return makeCR(exp, ctx.Block.NewICmp(enum.IPredSGE, params[0].value, params[1].value))
	case "<=":
		if from.Params[0].Type() == Real {
			return makeCR(exp, ctx.Block.NewFCmp(enum.FPredOLE, params[0].value, params[1].value))
		}
		return makeCR(exp, ctx.Block.NewICmp(enum.IPredSLE, params[0].value, params[1].value))
	case "==":
		if from.Params[0].Type() == String {
			v := ctx.Block.NewCall(ctx.global.utils.PVEqual16, params[0].value, params[1].value)
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
		panic("unknown op " + op)
	}
}

func getStructFieldIndex(s Structure, name string) int {
	index := 0
	found := false
	for _, v := range s.Fields {
		if IsFunction(v.Type) {
			if v.Name == name {
				found = true
				break
			}
			index++
		}
	}
	if !found {
		for _, v := range s.Fields {
			if IsStructure(v.Type) || IsString(v.Type) {
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
			if !IsStructure(v.Type) && !IsFunction(v.Type) {
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

func compileFAccess(fa *ast.FieldAccessor, ctx *context) cresult {
	cstru := compileExp(fa.Exp, ctx, false)
	tstru := fa.Exp.Type()
	ctyp := structureType(tstru.(Structure), false)
	typ := types.NewPointer(ctyp)
	bc := ctx.Block.NewBitCast(cstru.value, typ)
	index := getStructFieldIndex(tstru.(Structure), fa.Field)
	ptr := ctx.Block.NewGetElementPtr(ctyp, bc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(index+3)))
	res := ctx.Block.NewLoad(getType(fa.Type()), ptr)
	return makeCR(fa, res)
}

func compileConst(from *ast.Const, ctx *context) cresult {
	if from.Int != nil {
		return makeCR(from, constant.NewInt(IntType, *from.Int))
	} else if from.Bool != nil {
		return makeCR(from, constant.NewBool(*from.Bool))
	} else if from.Real != nil {
		return makeCR(from, constant.NewFloat(RealType, *from.Real))
	} else if from.String != nil {
		sref := ctx.makeStringRefRoot(*from.String)
		bc := ctx.Block.NewBitCast(sref, StringPType)
		return makeCR(from, bc)
	}
	panic("invalid constant at compilation")
}

func compileID(id *ast.Id, ctx *context) cresult {
	name := id.Name
	if ctx.isFun(name) {
		f := ctx.global.functions[name]
		clj := ctx.makeStructure(*f.From.Closure, f.Fun)
		return makeCR(id, clj)
	}
	r, err := ctx.resolveId(name)
	if err != nil {
		panic(err)
	}
	return makeCR(id, r)
}

func compileIf(c ast.Expression, t ast.Expression, f ast.Expression, ctx *context, funcRoot bool) cresult {
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
	ctx.freeIfUnboundRef(cond)

	trueB = ctx.Block
	if funcRoot && truev.value != nil {
		ctx.ret(makeCR(c, truev.value))
	}

	ctx.Block = falseB
	falsev := compileExp(f, ctx, funcRoot)
	falseB = ctx.Block
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

func compileCall(from *ast.FCall, ctx *context, funcRoot bool) cresult {
	if op, ok := from.Function.(*ast.Op); ok {
		var params []cresult
		name := op.Name
		if name == "if" { // Need to evaluate if parameters lazily
			return compileIf(from.Params[0], from.Params[1], from.Params[2], ctx, funcRoot)
		}
		for _, p := range from.Params {
			v := compileExp(p, ctx, false)
			params = append(params, v)
		}
		result := compileBinOp(from, from, name, params, ctx)
		for _, p := range params {
			ctx.freeIfUnboundRef(p)
		}
		return result
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

		if id, ok := from.Function.(*ast.Id); ok {
			name := id.Name
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
				return makeCR(from, res)
			}
			id, err := ctx.resolveId(name)
			if err != nil {
				panic(err)
			}
			res := ctx.call(id, from.Function.Type(), vparams)
			for _, p := range params {
				ctx.freeIfUnboundRef(p)
			}
			return makeCR(from, res)
		}
		fval := compileExp(from.Function, ctx, false)
		res := ctx.call(fval.value, from.Function.Type(), vparams)
		for _, p := range params {
			ctx.freeIfUnboundRef(p)
		}
		ctx.freeIfUnboundRef(fval)
		return makeCR(from, res)
	}
}

func compileBlock(from *ast.Block, ctx *context, funcRoot bool) cresult {
	sub := ctx.subContext()

	assigns := map[string]ast.Expression{}
	deps := map[string]map[string]bool{}
	for k, c := range from.Def.Assignments {
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
				if IsFunction(c.Type()) {
					closureids[k] = v.value
					if _, isId := c.(*ast.Id); isId { // TODO: Optimise renames away
						sub.incRef(v.value)
					}
				} else if IsStructure(c.Type()) || IsString(c.Type()) {
					structids[k] = v.value
					if _, isId := c.(*ast.Id); isId { // TODO: Optimise renames away
						sub.incRef(v.value)
					}
				}
				delete(assigns, k)
			}
		}
	}

	res := compileExp(from.Value, sub, funcRoot)
	for id, v := range closureids {
		if i, isId := from.Value.(*ast.Id); !isId || i.Name != id {
			sub.freeRef(v)
		}
	}
	for id, v := range structids {
		if i, isId := from.Value.(*ast.Id); !isId || i.Name != id {
			sub.freeRef(v)
		}
	}
	ctx.Block = sub.Block
	return res
}

func collectDeps(exp ast.Expression, c *context) []string {
	ids := map[string]bool{}
	ast.VisitBefore(exp, func(v ast.Ast, _ *ast.VisitContext) error {
		if id, ok := v.(*ast.Id); ok {
			name := id.Name
			ids[name] = true
			if IsFunction(id.Type()) && c.isFun(name) {
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
