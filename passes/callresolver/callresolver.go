package callresolver

import (
	"strconv"
	"strings"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

type FSign = string

func MakeFSign(name string, blockId int, sign string) FSign {
	return name + "%%" + strconv.Itoa(blockId) + "%%" + sign
}

type FEntry struct {
	Def    *FDef
	Struct *Struct
}

type FCat = map[FSign]FEntry

func Collect(exp Expression) FCat {
	result := FCat{}
	VisitAfter(exp, func(v Ast, _ *VisitContext) error {
		if b, ok := v.(*Block); ok {
			for n, a := range b.Def.Assignments {
				if d, ok := a.(*FDef); ok {
					result[n] = FEntry{Def: d}
					delete(b.Def.Assignments, n)
				}
				if s, ok := a.(*Struct); ok {
					result[n] = FEntry{Struct: s}
					delete(b.Def.Assignments, n)
				}
			}
		}
		return nil
	})
	return result
}

func ResolveFunctions(exp Expression) {
	anonCount := 0
	exp.Visit(NullFun, NullFun, true, func(v Ast, ctx *VisitContext) Ast {
		switch e := v.(type) {
		case *FCall:
			resolveCall(e)
		case *Id:
			if e.Type().IsFunction() && !strings.Contains(e.Name, "%%") {
				resolveIdFunct(e, ctx)
			}
		case *FDef:
			if ctx.NameOf(e) == "" {
				anonCount++
				typ := e.Type()
				fsig := MakeFSign("<anon"+strconv.Itoa(anonCount)+">", ctx.Definitions().ID, e.Type().TSignature())
				ctx.Definitions().Assignments[fsig] = e.CopyWithCtx(types.NewTypeCopyCtx())
				return &Id{Name: fsig, IdType: typ}
			}
		case *FieldAccessor:
			interfs := ctx.InterfacesWith(e.Field)
			for _, in := range interfs {
				if in.Interf.InterfaceType.UnifiesWith(e.Exp.Type()) {
					res := in.Interf
					newName := e.Field + "%interface%" + strconv.Itoa(res.Definitions.ID) + "%" + res.InterfaceType.Signature()
					id := NewId(newName)
					if !res.InterfaceType.IsDefined() {
						panic("untyped interface at resolution for " + e.Field)
					}
					uni, err := res.InterfaceType.Unifier(e.Exp.Type())
					if err != nil {
						panic(err)
					}
					ConvertTypes(e, uni)
					ts := append([]types.Type{e.Exp.Type()}, e.Type())
					typ := types.MakeFunction(ts...)
					id.IdType = typ
					if id.IdType.IsFunction() {
						resolveIdFunct(id, ctx)
					}
					call := NewFCall(id, e.Exp)
					call.CallType = id.IdType.FunctReturn()
					if id.IdType.HasFreeVars() {
						panic("free result type vars at resolver for " + e.Field)
					}
					return call
				}
			}
		}
		return v
	}, NewVisitCtx())
}

func resolveCall(v *FCall) {
	fun := v.MakeFunType()
	uni, err := fun.Unifier(v.Function.Type())
	if err != nil {
		panic(err)
	}
	ConvertTypes(v.Function, uni)
}

func resolveIdFunct(v *Id, ctx *VisitContext) {
	name := v.Name
	if defin := ctx.DefinitionOf(name); defin != nil {
		assig := defin.Assignments[name]
		if assig != nil {
			_, isDef := assig.(*FDef)
			_, isStruct := assig.(*Struct)
			_, isBlock := assig.(*Block)
			if isDef || isStruct || isBlock {
				fsig := MakeFSign(v.Name, defin.ID, v.Type().TSignature())
				if defin.Assignments[fsig] == nil {
					cop := assig.CopyWithCtx(types.NewTypeCopyCtx())
					subs, err := cop.Type().Unifier(v.Type())
					if err != nil {
						panic(err)
					}
					ConvertTypes(cop, subs)
					if cop.Type().HasFreeVars() {
						panic("could not unify " + assig.Type().Signature() + " u " + v.Type().Signature() + " => " + cop.Type().Signature())
					}
					defin.Assignments[fsig] = cop
				} else {
					f := defin.Assignments[v.Name]
					cop := f.CopyWithCtx(types.NewTypeCopyCtx())
					_, err := cop.Type().Unifier(v.Type())
					if err != nil {
						panic(err)
					}
				}
				v.Name = fsig
			}
		}
	}
}
