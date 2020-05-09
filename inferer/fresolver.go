package inferer

import (
	"strconv"

	"github.com/jvmakine/shine/resolved"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

type FSign = string

func MakeFSign(name string, blockId int, sign string) FSign {
	return name + "%%" + strconv.Itoa(blockId) + "%%" + sign
}

type Source int

const (
	Assignment Source = iota
	Parameter         = iota
)

type FDef struct {
	def    *ast.Exp
	source Source
}

type FCat = map[FSign]*ast.FDef

type gctx struct {
	cat FCat
}

type lctx struct {
	anonCount int
	global    *gctx
	parent    *lctx
	defs      map[string]*FDef
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
	ctx := lctx{global: &gctx{cat: FCat{}}, defs: map[string]*FDef{}, parent: nil}
	resolveExp(exp, &ctx, "")
	return &(ctx.global.cat)
}

func resolveExp(exp *ast.Exp, ctx *lctx, name string) {
	// TODO: remove from AST
	if exp.HasBeenResolved {
		return
	}
	exp.HasBeenResolved = true
	if exp.Block != nil {
		resolveBlock(exp, ctx)
	} else if exp.Call != nil {
		resolveCall(exp, ctx)
	} else if exp.Def != nil {
		resolveDef(exp, ctx, name)
	} else if exp.Id != nil {
		resolveId(exp, ctx)
	}
}

func resolveCall(exp *ast.Exp, ctx *lctx) {
	call := exp.Call
	for _, p := range call.Params {
		resolveExp(p, ctx, "")
	}
	es := ctx.resolve(call.Name)
	if es != nil {
		typ := es.def.Type()
		if !typ.HasFreeVars() {
			sig := typ.Signature()
			fsig := MakeFSign(call.Name, es.def.BlockID, sig)
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			resolveExp(es.def, ctx, call.Name)
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

			fsig := MakeFSign(call.Name, es.def.BlockID, cop.Type().Signature())
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			resolveExp(cop, ctx, call.Name)
		}
	}
}

func resolveBlock(exp *ast.Exp, pctx *lctx) {
	ctx := pctx.sub()
	block := exp.Block
	assigns := map[string]bool{}
	for _, a := range block.Assignments {
		assigns[a.Name] = true
		if a.Value.Def != nil {
			ctx.defs[a.Name] = &FDef{a.Value, Assignment}
		}
	}
	resolveExp(block.Value, ctx, "")
}

func resolveDef(exp *ast.Exp, ctx *lctx, name string) {
	def := exp.Def
	resolveExp(def.Body, ctx, "")
	if name != "" {
		fsig := MakeFSign(name, exp.BlockID, exp.Type().Signature())
		exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
		if ctx.global.cat[fsig] == nil {
			ctx.global.cat[fsig] = exp.Def
		}
	} else {
		ctx.anonCount++
		anonc := strconv.Itoa(ctx.anonCount)
		if exp.Resolved == nil {
			fsig := MakeFSign("<anon"+anonc+">", exp.BlockID, exp.Type().Signature())
			ctx.global.cat[fsig] = exp.Def
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
		}
	}
}

func resolveId(exp *ast.Exp, ctx *lctx) {
	id := exp.Id
	typ := exp.Type()
	if typ.IsFunction() {
		f := ctx.resolve(id.Name)
		if f != nil {
			sig := typ.Signature()
			fsig := MakeFSign(id.Name, f.def.BlockID, sig)
			exp.Resolved = &resolved.ResolvedFnCall{ID: fsig}
			if f.def.Type().HasFreeVars() {
				cop := f.def.Copy()
				subs, err := Unify(cop.Type(), typ)
				if err != nil {
					panic(err)
				}
				subs.Convert(cop)
				if cop.Type().HasFreeVars() {
					panic("could not unify " + f.def.Type().Signature() + " u " + typ.Signature() + " => " + cop.Type().Signature())
				}
				resolveExp(cop, ctx, id.Name)
			} else {
				resolveExp(f.def, ctx, id.Name)
			}
		}
	}
}
