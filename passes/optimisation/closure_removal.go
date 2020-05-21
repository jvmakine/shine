package optimisation

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

func ClosureRemoval(exp *ast.Exp) {
	exp.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Block != nil {
			for k, a := range v.Block.Assignments {
				if a.Def != nil && a.Def.Closure != nil && len(*a.Def.Closure) > 0 {
					newname := k + "%flat"
					newparams := a.Def.Params
					for _, c := range *a.Def.Closure {
						newparams = append(newparams, &ast.FParam{Name: c.Name, Type: c.Type})
					}
					v.Block.Assignments[newname] = &ast.Exp{Def: &ast.FDef{
						Body:    a.Def.Body,
						Params:  newparams,
						Closure: &types.Closure{},
					}}
				}
			}
		} else if v.Call != nil && v.Call.Function.Id != nil {
			id := v.Call.Function.Id.Name
			block := ctx.BlockOf(id)
			if block != nil && block.Assignments[id].Def != nil {
				f := block.Assignments[id]
				if f.Def.Closure != nil && len(*f.Def.Closure) > 0 {
					newid := id + "%flat"
					v.Call.Function.Id.Name = newid
					newargs := v.Call.Params
					for _, c := range *f.Def.Closure {
						a := &ast.Exp{Id: &ast.Id{Name: c.Name, Type: c.Type}}
						newargs = append(newargs, a)
					}
					v.Call.Params = newargs
				}
			}
		}
		return nil
	})
}
