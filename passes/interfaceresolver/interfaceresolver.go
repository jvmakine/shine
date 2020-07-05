package interfaceresolver

import (
	"errors"
	"strconv"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

func mergeInterfaces(ins []*Interface, typ types.Type) (*Interface, error) {
	methods := map[string]Expression{}
	subins := map[types.Type][]*Interface{}
	for _, in := range ins {
		for n, m := range in.Definitions.Assignments {
			if methods[n] != nil {
				return nil, errors.New(n + " declared twice for the same type: " + typ.Signature())
			}
			methods[n] = m
		}
		for t, i := range in.Definitions.Interfaces {
			subins[t] = i
		}
	}
	return &Interface{
		Definitions: &Definitions{
			Assignments: methods,
			Interfaces:  subins,
			ID:          ins[0].Definitions.ID,
		},
	}, nil
}

func Resolve(exp Ast) error {
	err := VisitBefore(exp, func(a Ast, ctx *VisitContext) error {
		if def, ok := a.(*Definitions); ok {
			methods := map[string][]types.Type{}
			for typ, ins := range def.Interfaces {
				in, err := mergeInterfaces(ins, typ)
				if err != nil {
					return err
				}
				delete(def.Interfaces, typ)

				def.Interfaces[typ] = []*Interface{in}
				in.InterfaceType = typ

				for n := range in.Definitions.Assignments {
					if methods[n] == nil {
						methods[n] = []types.Type{}
					}
					for _, oit := range methods[n] {
						if oit.UnifiesWith(typ) {
							return errors.New(n + " declared twice for unifiable types: " + typ.Signature() + ", " + oit.Signature())
						}
					}
					methods[n] = append(methods[n], typ)
				}

				for name, method := range in.Definitions.Assignments {
					newName := name + "%interface%" + strconv.Itoa(in.Definitions.ID) + "%" + typ.Signature()
					newDef := NewFDef(method, "$")
					newDef.Params[0].ParamType = typ
					def.Assignments[newName] = newDef
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return exp.Visit(NullFun, NullFun, false, func(from Ast, ctx *VisitContext) Ast {
		if fa, ok := from.(*FieldAccessor); ok {
			res, _ := ctx.InterfaceWith(fa.Field)
			if res != nil {
				newName := fa.Field + "%interface%" + strconv.Itoa(res.Definitions.ID) + "%" + res.InterfaceType.Signature()
				id := NewId(newName)
				ts := append([]types.Type{fa.Exp.Type()}, fa.Type())
				typ := types.MakeFunction(ts...)
				id.IdType = typ
				call := NewFCall(id, fa.Exp)
				call.CallType = typ.FunctReturn()
				return call
			}
		}
		return from
	}, NewVisitCtx())
}
