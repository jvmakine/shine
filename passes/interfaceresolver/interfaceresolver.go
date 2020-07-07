package interfaceresolver

import (
	"errors"
	"strconv"

	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

/*func mergeInterfaces(ins []*Interface, typ types.Type) ([]*Interface, error) {
	methods := map[string]Expression{}
	subins := []*Interface{}
	for _, in := range ins {
		for n, m := range in.Definitions.Assignments {
			if methods[n] != nil {
				return nil, errors.New(n + " declared twice for the same type: " + typ.Signature())
			}
			methods[n] = m
		}
		for t, i := range in.Definitions.Interfaces {
			subins = append(subins, i)
		}
	}
	return &Interface{
		Definitions: &Definitions{
			Assignments: methods,
			Interfaces:  subins,
			ID:          ins[0].Definitions.ID,
		},
	}, nil
}*/

func Resolve(exp Ast) error {
	err := VisitBefore(exp, func(a Ast, ctx *VisitContext) error {
		if def, ok := a.(*Definitions); ok {
			methods := map[string][]types.Type{}
			ins := def.Interfaces
			/*in, err := mergeInterfaces(ins, typ)
			if err != nil {
				return err
			}

			def.Interfaces = []*Interface{in}
			in.InterfaceType = typ*/
			for _, in := range ins {
				typ := in.InterfaceType
				for n := range in.Definitions.Assignments {
					if methods[n] == nil {
						methods[n] = []types.Type{}
					}
					for _, oit := range methods[n] {
						as := typ.Signature()
						bs := oit.Signature()
						if oit.UnifiesWith(typ) {
							if as < bs {
								return errors.New(n + " declared twice for unifiable types: " + as + ", " + bs)
							} else {
								return errors.New(n + " declared twice for unifiable types: " + bs + ", " + as)
							}
						}
					}
					methods[n] = append(methods[n], typ)
				}

				for name, method := range in.Definitions.Assignments {
					newName := name + "%interface%" + strconv.Itoa(in.Definitions.ID) + "%" + typ.Signature()
					newDef := NewFDef(method, &FParam{Name: "$", ParamType: typ})
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
			interfs := ctx.InterfacesWith(fa.Field)
			for _, in := range interfs {
				if in.Interf.InterfaceType.UnifiesWith(fa.Exp.Type()) {
					res := in.Interf
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
		}
		return from
	}, NewVisitCtx())
}
