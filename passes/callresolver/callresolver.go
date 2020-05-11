package callresolver

import (
	"strconv"

	"github.com/jvmakine/shine/passes/typeinference"

	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

type FSign = string

func MakeFSign(name string, blockId int, sign string) FSign {
	return name + "%%" + strconv.Itoa(blockId) + "%%" + sign
}

type FCat = map[FSign]*ast.FDef

func Collect(exp *ast.Exp) FCat {
	result := FCat{}
	exp.Visit(func(v *ast.Exp) {
		if v.Block != nil {
			for n, a := range v.Block.Assignments {
				if a.Def != nil {
					result[n] = a.Def
				}
			}
		}
	})
	return result
}

func ResolveFunctions(exp *ast.Exp) {
	anonCount := 0
	exp.Crawl(func(v *ast.Exp, ctx *ast.CrawlContext) {
		if v.Id != nil && v.Type().IsFunction() {
			if block := ctx.BlockOf(v.Id.Name); block != nil {
				fsig := MakeFSign(v.Id.Name, block.ID, v.Type().Signature())
				if block.Assignments[fsig] == nil {
					f := block.Assignments[v.Id.Name]
					cop := f.Copy()
					subs, err := typeinference.Unify(cop.Type(), v.Type())
					if err != nil {
						panic(err)
					}
					subs.Convert(cop)
					if cop.Type().HasFreeVars() {
						panic("could not unify " + f.Type().Signature() + " u " + v.Type().Signature() + " => " + cop.Type().Signature())
					}
					block.Assignments[fsig] = cop
				}
				v.Id.Name = fsig
			}
		} else if v.Call != nil {
			params := make([]Type, len(v.Call.Params)+1)
			for i, p := range v.Call.Params {
				params[i] = p.Type()
			}
			params[len(v.Call.Params)] = v.Type()
			fun := MakeFunction(params...)
			if block := ctx.BlockOf(v.Call.Name); block != nil {
				fsig := MakeFSign(v.Call.Name, block.ID, fun.Signature())
				if block.Assignments[fsig] == nil {
					f := block.Assignments[v.Call.Name]
					cop := f.Copy()

					subs, err := typeinference.Unify(cop.Type(), fun)
					if err != nil {
						panic(err)
					}
					subs.Convert(cop)
					if cop.Type().HasFreeVars() {
						panic("could not unify " + f.Type().Signature() + " u " + fun.Signature() + " => " + cop.Type().Signature())
					}
					block.Assignments[fsig] = cop
				}
				v.Call.Name = fsig
			}
		} else if v.Def != nil && ctx.NameOf(v) == "" {
			anonCount++
			typ := v.Type()
			fsig := MakeFSign("<anon"+strconv.Itoa(anonCount)+">", ctx.Block().ID, v.Type().Signature())
			ctx.Block().Assignments[fsig] = v.Copy()
			v.Def = nil
			v.Id = &ast.Id{Name: fsig, Type: typ}
		}
	})
}
