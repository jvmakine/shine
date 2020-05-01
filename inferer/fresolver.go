package inferer

import (
	"strconv"

	"github.com/jvmakine/shine/resolved"
	. "github.com/jvmakine/shine/resolved"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

type FSign = string

func MakeFSign(name string, blockId int, sign string) FSign {
	return name + "%%" + strconv.Itoa(blockId) + "%%" + sign
}

type FDef struct {
	block int
	def   *ast.Exp
}

type FCat = map[FSign]*ast.FDef

type gctx struct {
	blockCount int
	cat        FCat
}

type lctx struct {
	blockID   int
	anonCount int
	global    *gctx
	parent    *lctx
	defs      map[string]*FDef
}

func (l *lctx) sub(id int) *lctx {
	return &lctx{global: l.global, parent: l, defs: map[string]*FDef{}, blockID: id}
}

func (l *lctx) resolve(name string) *FDef {
	if r := l.defs[name]; r != nil {
		return r
	}
	if l.parent != nil {
		return l.parent.resolve(name)
	}
	return nil
}

func Resolve(exp *ast.Exp) *FCat {
	ctx := lctx{global: &gctx{cat: FCat{}, blockCount: 0}, defs: map[string]*FDef{}}
	resolveExp(exp, &ctx)
	return &(ctx.global.cat)
}

func resolveExp(exp *ast.Exp, ctx *lctx) Closure {
	if exp.Block != nil {
		return resolveBlock(exp, ctx)
	} else if exp.Call != nil {
		return resolveCall(exp, ctx)
	} else if exp.Def != nil {
		return resolveDef(exp, ctx)
	} else if exp.Id != nil {
		return resolveId(exp, ctx)
	}
	return Closure{}
}

func resolveAnonFuncParams(call *ast.FCall, ctx *lctx) Closure {
	cjs := Closure{}
	for _, p := range call.Params {
		if p.Def != nil { // anonymous function
			ctx.anonCount++
			anonc := strconv.Itoa(ctx.anonCount)
			fsig := MakeFSign("<anon"+anonc+">", ctx.blockID, p.Type().Signature())
			p.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = p.Def
				cjs = append(cjs, resolveExp(p, ctx)...)
			}
		}
	}
	return cjs
}

func resolveCall(exp *ast.Exp, ctx *lctx) Closure {
	call := exp.Call
	cjs := Closure{}
	for _, p := range call.Params {
		cjs = append(cjs, resolveExp(p, ctx)...)
	}
	es := ctx.resolve(call.Name)
	if es != nil {
		typ := es.def.Type()
		if !typ.HasFreeVars() {
			sig := typ.Signature()
			fsig := MakeFSign(call.Name, es.block, sig)
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = es.def.Def
				cjs = append(cjs, resolveExp(es.def, ctx)...)
			}
			cjs = append(cjs, resolveAnonFuncParams(call, ctx)...)
		} else {
			ptypes := make([]Type, len(call.Params)+1)
			for i, p := range call.Params {
				ptypes[i] = p.Type()
			}
			ptypes[len(call.Params)] = exp.Type()
			cop := es.def.Copy()
			fun := MakeFunction(ptypes...)
			u1 := exp.Type().Signature()
			u2 := cop.Type().Signature()

			s, err := Unify(fun, cop.Type())
			if err != nil {
				panic(err)
			}

			s.Convert(cop)
			s.Convert(exp)
			if cop.Type().HasFreeVars() || exp.Type().HasFreeVars() {
				panic("type inference failed: " + u1 + " u " + u2 + " => " + cop.Type().Signature())
			}

			fsig := MakeFSign(call.Name, es.block, cop.Type().Signature())
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = cop.Def
				cjs = append(cjs, resolveExp(cop, ctx)...)
			}
			cjs = append(cjs, resolveAnonFuncParams(call, ctx)...)
		}
	}
	return cjs
}

func resolveBlock(exp *ast.Exp, pctx *lctx) Closure {
	cjs := Closure{}
	ctx := pctx.sub(pctx.global.blockCount + 1)
	ctx.global.blockCount++
	block := exp.Block
	assigns := map[string]bool{}
	for _, a := range block.Assignments {
		assigns[a.Name] = true
		if a.Value.Def != nil {
			ctx.defs[a.Name] = &FDef{ctx.global.blockCount, a.Value}
		} else {
			cjs = append(cjs, resolveExp(a.Value, pctx)...)
		}
	}
	cjs = append(cjs, resolveExp(block.Value, ctx)...)
	result := Closure{}
	for _, i := range cjs {
		if !assigns[i.Name] {
			result = append(result, i)
		}
	}
	return result
}

func resolveDef(exp *ast.Exp, ctx *lctx) Closure {
	def := exp.Def
	cjs := resolveExp(def.Body, ctx)
	params := map[string]bool{}
	for _, p := range def.Params {
		params[p.Name] = true
	}
	ClosureIds := Closure{}
	for _, i := range cjs {
		if !params[i.Name] {
			ClosureIds = append(ClosureIds, i)
		}
	}
	exp.Def.Resolved = &resolved.ResolvedFnDef{Closure: ClosureIds}
	return ClosureIds
}

func resolveId(exp *ast.Exp, ctx *lctx) Closure {
	id := exp.Id
	typ := exp.Type()
	if typ.IsFunction() {
		f := ctx.resolve(id.Name)
		if f == nil {
			// function argument has been already resolved
			return Closure{{Name: id.Name, Type: exp.Type()}}
		}
		var fsig string
		if f.def.Type().HasFreeVars() {
			cop := f.def.Copy()
			subs, err := Unify(cop.Type(), typ)
			if err != nil {
				panic(err)
			}
			subs.Convert(cop)
			if cop.Type().HasFreeVars() {
				panic("could not unify")
			}
			sig := cop.Type().Signature()
			fsig = MakeFSign(id.Name, f.block, sig)
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = cop.Def
				resolveExp(cop, ctx)
			}
		} else {
			fsig := MakeFSign(id.Name, f.block, f.def.Type().Signature())
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = f.def.Def
				resolveExp(f.def, ctx)
			}
		}
	}
	return Closure{{Name: id.Name, Type: exp.Type()}}
}
