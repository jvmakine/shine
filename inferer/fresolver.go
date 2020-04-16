package inferer

import (
	"strconv"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
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
	}
}

func resolveCall(exp *ast.Exp, ctx *lctx) {
	call := exp.Call
	for _, p := range call.Params {
		resolveExp(p, ctx)
	}
	es := ctx.resolve(call.Name)
	if es != nil {
		if es.def.Type.IsDefined() {
			sig := es.def.Type.Signature()
			fsig := MakeFSign(call.Name, es.block, sig)
			call.Resolved = fsig
			if ctx.global.cat[fsig] == nil {
				ctx.global.cat[fsig] = es.def.Def
				resolveExp(es.def, ctx)
			}
		} else {
			ptypes := make([]*types.TypePtr, len(call.Params)+1)
			for i, p := range call.Params {
				ptypes[i] = p.Type
			}
			ptypes[len(call.Params)] = exp.Type
			ftype := types.MakeFun(ptypes...)
			cop := es.def.Copy()
			u1 := ftype.Signature()
			u2 := cop.Type.Signature()

			uni, err := Unify(ftype, cop.Type)
			if err != nil {
				panic(err)
			}
			uni.ApplySource(ftype)
			uni.ApplyDest(cop.Type)
			if !ftype.IsDefined() {
				panic("type inference failed: " + u1 + " u " + u2 + " => " + ftype.Signature())
			}
			fsig := MakeFSign(call.Name, es.block, ftype.Signature())
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
