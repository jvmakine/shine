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
			for n, a := range v.Block.Assignments {
				if a.Def != nil {
					result[n] = FEntry{Def: a.Def}
					delete(v.Block.Assignments, n)
				}
			}
			for n, a := range v.Block.TypeDefs {
				if a.Struct != nil {
					result[n] = FEntry{Struct: a.Struct}
					delete(v.Block.Assignments, n)
				}
			}
		}
		return nil
	})
	return result
}

func ResolveFunctions(exp *ast.Exp) {
	anonCount := 0
	bmap := makeBlockMap(exp)
	exp.Crawl(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Call != nil {
			resolveCall(v.Call, ctx)
		}
		if v.Id != nil && v.Type().IsFunction() && !strings.Contains(v.Id.Name, "%%") && ctx.Binding() == nil {
			resolveIdFunct(v, bmap, ctx)
		} else if v.Def != nil && ctx.NameOf(v) == "" && ctx.Binding() == nil {
			anonCount++
			typ := v.Type()
			fsig := MakeFSign("<anon"+strconv.Itoa(anonCount)+">", ctx.Block().ID, v.Type().Signature())
			ctx.Block().Assignments[fsig] = v.Copy()
			v.Def = nil
			v.Id = &ast.Id{Name: fsig, Type: typ}
		}
		return nil
	})
}

func makeBlockMap(exp *ast.Exp) map[int]*ast.Block {
	res := map[int]*ast.Block{}
	exp.Visit(func(p *ast.Exp, ctx *ast.VisitContext) error {
		if b := p.Block; b != nil {
			res[b.ID] = b
		}
		return nil
	})
	return res
}

func resolveCall(v *ast.FCall, ctx *ast.VisitContext) {
	bindings := ctx.GetBindings()
	fun := v.MakeFunType()

	nt := v.Function.Type().ApplyBindings(bindings)

	uni, err := fun.Unifier(nt)

	if err != nil {
		panic(err)
	}
	v.Function.Convert(uni)
}

func resolveIdFunct(v *ast.Exp, bmap map[int]*ast.Block, ctx *ast.VisitContext) {
	name := v.Id.Name
	block := ctx.BlockOf(name)
	if block == nil {
		return
	}
	if block.Assignments[name] != nil && block.Assignments[name].Def != nil {
		fsig := MakeFSign(v.Id.Name, block.ID, v.Type().Signature())
		if block.Assignments[fsig] == nil {
			f := block.Assignments[v.Id.Name]
			cop := f.Copy()
			nt := cop.Type().ApplyBindings(ctx.GetBindings())
			subs, err := nt.Unifier(v.Type())
			if err != nil {
				panic(err)
			}
			cop.Convert(subs)
			if cop.Type().HasFreeVars() {
				panic("could not unify " + f.Type().Signature() + " u " + v.Type().Signature() + " => " + cop.Type().Signature())
			}
			block.Assignments[fsig] = cop
		} else {
			f := block.Assignments[v.Id.Name]
			cop := f.Copy()
			_, err := cop.Type().Unifier(v.Type())
			if err != nil {
				panic(err)
			}
		}
		v.Id.Name = fsig
	} else if block.TypeDefs[name] != nil && block.TypeDefs[name].Struct != nil {
		fsig := MakeFSign(v.Id.Name, block.ID, v.Type().Signature())
		if block.TypeDefs[fsig] == nil {
			f := block.TypeDefs[v.Id.Name]

			cop := f.Copy()
			subs, err := cop.Struct.Constructor().Unifier(v.Type())
			if err != nil {
				panic(err)
			}
			cop.Convert(subs)
			if cop.Type().HasFreeVars() {
				panic("could not unify " + cop.Struct.Constructor().Signature() + " u " + v.Type().Signature() + " => " + cop.Type().Signature())
			}
			block.TypeDefs[fsig] = cop
		} else {
			f := block.TypeDefs[v.Id.Name]
			cop := f.Copy()
			_, err := cop.Struct.Constructor().Unifier(v.Type())
			if err != nil {
				panic(err)
			}
		}
		v.Id.Name = fsig
	} else if fd := block.TCFunctions[name]; fd != nil {
		bindings := ctx.GetBindings()
		blockId := 0
		var fexp *ast.Exp
		tcname := ""
		for n, td := range block.TypeDefs {
			if td == fd {
				tcname = n
				break
			}
		}

		for _, b := range bindings {
			if b.Name == tcname {
				b2 := bmap[b.BlockID]
				for _, b := range b2.TypeBindings {
					if b.Name == tcname {
						def := b.Bindings[name]
						exp := &ast.Exp{Def: def}

						if exp.Type().Unifies(v.Type()) {
							blockId = b2.ID
							fexp = exp.Copy()
							break
						}
					}
				}
			}
			if blockId != 0 {
				break
			}
		}

		if fexp == nil {
			panic("no binding found")
		}

		fsig := MakeFSign(v.Id.Name, blockId, v.Type().Signature())
		if block.Assignments[fsig] == nil {
			cop := fexp.Copy()
			subs, err := cop.Type().Unifier(v.Type())
			if err != nil {
				panic(err)
			}
			cop.Convert(subs)
			if cop.Type().HasFreeVars() {
				panic("could not unify " + fexp.Type().Signature() + " u " + v.Type().Signature() + " => " + cop.Type().Signature())
			}
			block.Assignments[fsig] = cop
		} else {
			cop := block.Assignments[fsig].Copy()
			_, err := cop.Type().Unifier(v.Type())
			if err != nil {
				panic(err)
			}
		}
		v.Id.Name = fsig
	}
}
