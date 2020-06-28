package callresolver

import (
	"strconv"
	"strings"

	"github.com/jvmakine/shine/ast"
)

type FSign = string

func MakeFSign(name string, blockId int, sign string) FSign {
	return name + "%%" + strconv.Itoa(blockId) + "%%" + sign
}

type FEntry struct {
	Def    *ast.FDef
	Struct *ast.Struct
}

type FCat = map[FSign]FEntry

func Collect(exp *ast.Exp) FCat {
	result := FCat{}
	exp.VisitAfter(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Block != nil {
			for n, a := range v.Block.Defin.Assignments {
				if a.Def != nil {
					result[n] = FEntry{Def: a.Def}
					delete(v.Block.Defin.Assignments, n)
				}
				if a.Struct != nil {
					result[n] = FEntry{Struct: a.Struct}
					delete(v.Block.Defin.Assignments, n)
				}
			}
		}
		return nil
	})
	return result
}

func ResolveFunctions(exp *ast.Exp) {
	anonCount := 0
	exp.Crawl(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Call != nil {
			resolveCall(v.Call)
		}
		if v.Id != nil && v.Type().IsFunction() && !strings.Contains(v.Id.Name, "%%") {
			resolveIdFunct(v, ctx)
		} else if v.Def != nil && ctx.NameOf(v) == "" {
			anonCount++
			typ := v.Type()
			fsig := MakeFSign("<anon"+strconv.Itoa(anonCount)+">", ctx.Block().ID, v.Type().TSignature())
			ctx.Block().Defin.Assignments[fsig] = v.Copy()
			v.Def = nil
			v.Id = &ast.Id{Name: fsig, Type: typ}
		}
		return nil
	})
}

func resolveCall(v *ast.FCall) {
	fun := v.MakeFunType()
	uni, err := fun.Unifier(v.Function.Type())
	if err != nil {
		panic(err)
	}
	v.Function.Convert(uni)
}

func resolveIdFunct(v *ast.Exp, ctx *ast.VisitContext) {
	name := v.Id.Name
	if block := ctx.BlockOf(name); block != nil && (block.Defin.Assignments[name].Def != nil || block.Defin.Assignments[name].Struct != nil) {
		fsig := MakeFSign(v.Id.Name, block.ID, v.Type().TSignature())
		if block.Defin.Assignments[fsig] == nil {
			f := block.Defin.Assignments[v.Id.Name]
			cop := f.Copy()
			subs, err := cop.Type().Unifier(v.Type())
			if err != nil {
				panic(err)
			}
			cop.Convert(subs)
			if cop.Type().HasFreeVars() {
				panic("could not unify " + f.Type().Signature() + " u " + v.Type().Signature() + " => " + cop.Type().Signature())
			}
			block.Defin.Assignments[fsig] = cop
		} else {
			f := block.Defin.Assignments[v.Id.Name]
			cop := f.Copy()
			_, err := cop.Type().Unifier(v.Type())
			if err != nil {
				panic(err)
			}
		}
		v.Id.Name = fsig
	}
}
