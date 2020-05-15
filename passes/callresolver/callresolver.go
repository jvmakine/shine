package callresolver

import (
	"strconv"

	"github.com/jvmakine/shine/ast"
)

type FSign = string

func MakeFSign(name string, blockId int, sign string) FSign {
	return name + "%%" + strconv.Itoa(blockId) + "%%" + sign
}

type FCat = map[FSign]*ast.FDef

func Collect(exp *ast.Exp) FCat {
	result := FCat{}
	exp.Visit(func(v *ast.Exp, _ *ast.VisitContext) error {
		if v.Block != nil {
			for n, a := range v.Block.Assignments {
				if a.Def != nil {
					result[n] = a.Def
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
			fun := v.Call.MakeFunType()
			uni, err := fun.Unifier(v.Call.Function.Type())
			if err != nil {
				panic(err)
			}
			v.Call.Function.Convert(uni)
		}
		if v.Id != nil && v.Type().IsFunction() {
			if block := ctx.BlockOf(v.Id.Name); block != nil {
				fsig := MakeFSign(v.Id.Name, block.ID, v.Type().Signature())
				if block.Assignments[fsig] == nil {
					f := block.Assignments[v.Id.Name]
					cop := f.Copy()
					subs, err := cop.Type().Unifier(v.Type())
					if err != nil {
						panic(err)
					}
					cop.Convert(subs)
					if cop.Type().HasFreeVars() {
						panic("could not unify " + f.Type().Signature() + " u " + v.Type().Signature() + " => " + cop.Type().Signature())
					}
					block.Assignments[fsig] = cop
				}
				v.Id.Name = fsig
			}
		} else if v.Def != nil && ctx.NameOf(v) == "" {
			anonCount++
			typ := v.Type()
			fsig := MakeFSign("<anon"+strconv.Itoa(anonCount)+">", ctx.Block().ID, v.Type().Signature())
			ctx.Block().Assignments[fsig] = v.Copy()
			v.Def = nil
			v.Id = &ast.Id{Name: fsig, Type: typ}
		} else if v.Call != nil && v.Call.Function.Call != nil && v.Call.Function.Op == nil {
			anonCount++
			typ := v.Call.Function.Type()
			fsig := MakeFSign("<anon"+strconv.Itoa(anonCount)+">", ctx.Block().ID, typ.Signature())
			ctx.Block().Assignments[fsig] = v.Call.Function.Copy()
			v.Call.Function.Call = nil
			v.Call.Function.Id = &ast.Id{Name: fsig, Type: typ}
		}
		return nil
	})
}
