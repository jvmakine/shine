package callresolver

import (
	"strconv"
	"strings"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
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
		case *Op:
			typ := e.Left.Type()
			contextual, ok := typ.(Contextual)
			if !ok {
				panic("no context available for type " + Signature(typ))
			}
			if contextual.GetContext() == nil {
				panic("nil context for type " + Signature(typ))
			}
			ic := contextual.GetContext().(*VisitContext)
			inter := ic.InterfaceWithType(e.Name, typ)
			return inter.Interf.ReplaceOp(e)
		case *FCall:
			resolveCall(e, ctx)
		case *Id:
			_, isFun := e.Type().(Function)
			if isFun && !strings.Contains(e.Name, "%%") {
				resolveIdFunct(e, ctx)
			}
		case *FDef:
			if ctx.NameOf(e) == "" {
				anonCount++
				typ := e.Type()
				fsig := MakeFSign("<anon"+strconv.Itoa(anonCount)+">", ctx.Definitions().ID, Signature(e.Type()))
				ctx.Definitions().Assignments[fsig] = e.CopyWithCtx(NewTypeCopyCtx())
				return &Id{Name: fsig, IdType: typ}
			}
		case *FieldAccessor:
			interfs := ctx.InterfacesWith(e.Field)
			for _, in := range interfs {
				if UnifiesWith(in.Interf.InterfaceType, e.Exp.Type(), ctx) {
					res := in.Interf
					newName := e.Field + "%interface%" + strconv.Itoa(res.Definitions.ID) + "%" + Signature(res.InterfaceType)
					id := NewId(newName)
					if res.InterfaceType == nil {
						panic("untyped interface at resolution for " + e.Field)
					}
					uni, err := Unifier(res.InterfaceType, e.Exp.Type(), ctx)
					if err != nil {
						panic(err)
					}
					ConvertTypes(e, uni)
					ts := append([]Type{e.Exp.Type()}, e.Type())
					typ := NewFunction(ts...)
					id.IdType = typ
					call := NewFCall(id, e.Exp)
					if fun, isFun := id.IdType.(Function); isFun {
						resolveIdFunct(id, ctx)
						call.CallType = fun.Fields[len(fun.Fields)-1]
					} else {
						panic("field accessor must be a function")
					}
					if HasFreeVars(id.IdType) {
						panic("free result type vars at resolver for " + e.Field)
					}
					return call
				}
			}
		}
		return v
	}, NewVisitCtx())
}

func resolveCall(v *FCall, ctx *VisitContext) {
	fun := v.MakeFunType()
	uni, err := Unifier(fun, v.Function.Type(), ctx)
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
				fsig := MakeFSign(v.Name, defin.ID, Signature(v.Type()))
				if defin.Assignments[fsig] == nil {
					cop := assig.CopyWithCtx(types.NewTypeCopyCtx())
					subs, err := Unifier(cop.Type(), v.Type(), ctx)
					if err != nil {
						panic(err)
					}
					ConvertTypes(cop, subs)
					if HasFreeVars(cop.Type()) {
						panic("could not unify " + fsig)
					}
					/*if varsAfter > 0 {
						panic("hidden variables after resolution for " + fsig)
					}*/
					defin.Assignments[fsig] = cop
				} else {
					f := defin.Assignments[v.Name]
					cop := f.CopyWithCtx(types.NewTypeCopyCtx())
					_, err := Unifier(cop.Type(), v.Type(), ctx)
					if err != nil {
						panic(err)
					}
				}
				v.Name = fsig
			}
		}
	}
}
