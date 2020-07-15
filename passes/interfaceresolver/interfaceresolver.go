package interfaceresolver

import (
	"errors"
	"strconv"

	. "github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func Resolve(exp Ast) error {
	return VisitBefore(exp, func(a Ast, ctx *VisitContext) error {
		if def, ok := a.(*Definitions); ok {
			methods := map[string][]Type{}
			ins := def.Interfaces
			for _, in := range ins {
				typ := in.InterfaceType
				for n := range in.Definitions.Assignments {
					if methods[n] == nil {
						methods[n] = []Type{}
					}
					for _, oit := range methods[n] {
						as := Signature(typ)
						bs := Signature(oit)
						if UnifiesWith(oit, typ, ctx) {
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
					newName := name + "%interface%" + strconv.Itoa(in.Definitions.ID) + "%" + Signature(typ)
					newDef := NewFDef(method, &FParam{Name: "$", ParamType: typ})
					def.Assignments[newName] = newDef
				}
			}
		}
		return nil
	})
}
