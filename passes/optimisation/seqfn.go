package optimisation

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

// Optimise sequential function definitions into one when called with multiple arguments
func SequentialFunctionPass(exp *ast.Exp) {
	tctx := types.NewTypeCopyCtx()
	exp.Crawl(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Call != nil && v.Call.Function.Call != nil {
			root := v.Call.RootFunc()

			var def *ast.FDef
			var block *ast.Block
			var id string
			changed := false

			if root.Id != nil {
				id = root.Id.Name
				if block = ctx.BlockOf(id); block != nil {
					def = block.Def.Assignments[id].CopyWithCtx(tctx).Def
				}
			} else if root.Def != nil {
				def = root.CopyWithCtx(tctx).Def
			}

			params := v.Call.Params
			ptr := v
			nid := id
			for ptr.Call.Function.Call != nil && def != nil && def.Body.Def != nil {
				changed = true
				params = append(ptr.Call.Function.Call.Params, params...)
				ptr = ptr.Call.Function

				def2 := def.Body.Def
				def.Params = append(def.Params, def2.Params...)
				def.Body = def2.Body

				if block != nil {
					nid = nid + "%c"
				}
			}

			if changed {
				if block != nil {
					block.Def.Assignments[nid] = &ast.Exp{Def: def}
					v.Call.Function = &ast.Exp{Id: &ast.Id{Name: nid, Type: block.Def.Assignments[nid].Type()}}
					v.Call.Params = params
				} else {
					v.Call.Params = params
					v.Call.Function = &ast.Exp{Def: def}
				}
			}
		}
		return nil
	})
}
