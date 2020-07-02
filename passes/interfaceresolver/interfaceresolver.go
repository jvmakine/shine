package interfaceresolver

import (
	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

func mergeInterfaces(ins []*Interface) (*Interface, error) {
	// TODO errors
	methods := map[string]Expression{}
	subins := map[types.Type][]*Interface{}
	for _, in := range ins {
		for n, m := range in.Definitions.Assignments {
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
		},
	}, nil
}

func Resolve(exp Ast) error {
	err := VisitBefore(exp, func(a Ast, ctx *VisitContext) error {
		if def, ok := a.(*Definitions); ok {
			for typ, ins := range def.Interfaces {
				in, err := mergeInterfaces(ins)
				if err != nil {
					return err
				}
				delete(def.Interfaces, typ)
				def.Interfaces[typ] = []*Interface{in}
				in.InterfaceType = typ

				for name, method := range in.Definitions.Assignments {
					newName := name + "%interface%" + typ.Signature()
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
				method := res.Definitions.Assignments[fa.Field]
				newName := fa.Field + "%interface%" + res.InterfaceType.Signature()
				id := NewId(newName)
				id.IdType = method.Type()
				call := NewFCall(id, fa.Exp)
				call.CallType = method.Type().FunctReturn()
				return call
			}
		}
		return from
	}, NewVisitCtx())
}
