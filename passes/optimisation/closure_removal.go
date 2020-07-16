package optimisation

import (
	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

func ClosureRemoval(exp Expression) {
	VisitBefore(exp, func(v Ast, ctx *VisitContext) error {
		if b, ok := v.(*Block); ok {
			for k, a := range b.Def.Assignments {
				if d, ok := a.(*FDef); ok && d.HasClosure() {
					newname := k + "%flat"
					newparams := d.Params
					for _, c := range d.Closure.Fields {
						newparams = append(newparams, &FParam{Name: c.Name, ParamType: c.Type})
					}
					cls := types.NewStructure()
					b.Def.Assignments[newname] = &FDef{
						Body:    d.Body,
						Params:  newparams,
						Closure: &cls,
					}
				}
			}
		} else if c, ok := v.(*FCall); ok {
			if i, ok := c.Function.(*Id); ok {
				id := i.Name
				defin := ctx.DefinitionOf(id)
				if defin != nil {
					if bd, ok := defin.Assignments[id].(*FDef); ok {
						if bd.HasClosure() {
							newid := id + "%flat"
							i.Name = newid
							newargs := c.Params
							for _, c := range bd.Closure.Fields {
								a := &Id{Name: c.Name, IdType: c.Type}
								newargs = append(newargs, a)
							}
							c.Params = newargs
						}
					}
				}
			}
		}
		return nil
	})
}
