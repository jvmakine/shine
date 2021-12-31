package optimisation

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

func ClosureRemoval(exp *ast.Exp) {
	exp.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Block != nil {
			for k, a := range v.Block.Assignments {
				if a.Def != nil && a.Def.HasClosure() {
					newname := k + "%flat"
					newparams := a.Def.Params
					for _, c := range a.Def.Closure.Fields {
						newparams = append(newparams, &ast.FParam{Name: c.Name, Type: c.Type})
					}
					v.Block.Assignments[newname] = &ast.Exp{Def: &ast.FDef{
						Body:    a.Def.Body,
						Params:  newparams,
						Closure: types.MakeStructure("", []types.Type{}).Structure,
					}}
				}
			}
		} else if v.Call != nil && v.Call.Function.Id != nil {
			id := v.Call.Function.Id.Name
			block := ctx.BlockOf(id)
			if block != nil && block.Assignments[id] != nil && block.Assignments[id].Def != nil {
				f := block.Assignments[id]
				if f.Def.HasClosure() {
					newid := id + "%flat"
					v.Call.Function.Id.Name = newid
					newargs := v.Call.Params
					for _, c := range f.Def.Closure.Fields {
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
