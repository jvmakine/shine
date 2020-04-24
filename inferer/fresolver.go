package inferer

import (
	"strconv"

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
	global *gctx
	parent *lctx
	defs   map[string]*FDef
}

func (l *lctx) sub() *lctx {
	return &lctx{global: l.global, parent: l, defs: map[string]*FDef{}}
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

func resolveExp(exp *ast.Exp, ctx *lctx) {
	if exp.Block != nil {
		resolveBlock(exp, ctx)
	} else if exp.Call != nil {
		resolveCall(exp, ctx)
	} else if exp.Def != nil {
		resolveDef(exp, ctx)
	} else if exp.Id != nil {
		resolveId(exp, ctx)
	}
}

func resolveCall(exp *ast.Exp, ctx *lctx) {
	call := exp.Call
	for _, p := range call.Params {
		resolveExp(p, ctx)
	}
	es := ctx.resolve(call.Name)
	if es != nil {
		if !es.def.Type.HasFreeVars() {
			sig := es.def.Type.Signature()
			fsig := MakeFSign(call.Name, es.block, sig)
			call.Resolved = fsig
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = es.def.Def
				resolveExp(es.def, ctx)
			}
		} else {
			ptypes := make([]Type, len(call.Params)+1)
			for i, p := range call.Params {
				ptypes[i] = p.Type
			}
			ptypes[len(call.Params)] = exp.Type
			cop := es.def.Copy()
			fun := MakeFunction(ptypes...)
			u1 := exp.Type.Signature()
			u2 := cop.Type.Signature()

			s, err := Unify(fun, cop.Type)
			if err != nil {
				panic(err)
			}
			s.Convert(cop)
			s.Convert(exp)
			if cop.Type.HasFreeVars() {
				panic("type inference failed: " + u1 + " u " + u2 + " => " + cop.Type.Signature())
			}
			fsig := MakeFSign(call.Name, es.block, exp.Type.Signature())
			call.Resolved = fsig
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = cop.Def
				resolveExp(cop, ctx)
			}
		}
	}
}

func resolveBlock(exp *ast.Exp, pctx *lctx) {
	ctx := pctx.sub()
	ctx.global.blockCount++
	block := exp.Block
	for _, a := range block.Assignments {
		if a.Value.Def != nil {
			ctx.defs[a.Name] = &FDef{ctx.global.blockCount, a.Value}
		}
	}
	resolveExp(block.Value, ctx)
}

func resolveDef(exp *ast.Exp, ctx *lctx) {
	def := exp.Def
	resolveExp(def.Body, ctx)
}

func resolveId(exp *ast.Exp, ctx *lctx) {
	id := exp.Id
	if exp.Type.IsFunction() {
		f := ctx.resolve(id.Name)
		var fsig string
		if f.def.Type.HasFreeVars() {
			cop := f.def.Copy()
			subs, err := Unify(cop.Type, exp.Type)
			if err != nil {
				panic(err)
			}
			subs.Convert(cop)
			if cop.Type.HasFreeVars() {
				panic("could not unify")
			}
			sig := cop.Type.Signature()
			fsig = MakeFSign(id.Name, f.block, sig)
			id.Resolved = fsig

			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = cop.Def
				resolveExp(cop, ctx)
			}
		} else {
			fsig := MakeFSign(id.Name, f.block, f.def.Type.Signature())
			id.Resolved = MakeFSign(id.Name, f.block, fsig)
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = f.def.Def
				resolveExp(f.def, ctx)
			}
		}
	}
}
