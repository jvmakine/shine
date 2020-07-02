package optimisation

import (
	. "github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
)

// Optimise sequential function definitions into one when called with multiple arguments
func SequentialFunctionPass(exp Expression) {
	tctx := types.NewTypeCopyCtx()

	CrawlBefore(exp, func(v Ast, ctx *VisitContext) error {
		if c, ok := v.(*FCall); ok {
			root := c.RootFunc()
			var def *FDef
			var block *Block
			var id string
			changed := false

			if i, ok := root.(*Id); ok {
				id = i.Name
				if block = ctx.BlockOf(id); block != nil {
					assig := block.Def.Assignments[id]
					if d, ok := assig.(*FDef); ok {
						def = d.CopyWithCtx(tctx).(*FDef)
					} else {
						def = nil
					}
				}
			} else if _, ok := root.(*FDef); ok {
				def = root.CopyWithCtx(tctx).(*FDef)
			}

			params := c.Params
			nid := id
			fcall, isFCall := c.Function.(*FCall)
			isDefB := false
			if def != nil {
				_, isDefB = def.Body.(*FDef)
			}
			for isFCall && isDefB {
				changed = true
				params = append(fcall.Params, params...)
				def2 := def.Body.(*FDef)
				def.Params = append(def.Params, def2.Params...)
				def.Body = def2.Body

				if block != nil {
					nid = nid + "%c"
				}

				_, isDefB = def.Body.(*FDef)
				fcall, isFCall = fcall.Function.(*FCall)
			}

			if changed {
				if block != nil {
					block.Def.Assignments[nid] = def
					c.Function = &Id{Name: nid, IdType: def.Type()}
					c.Params = params
				} else {
					c.Params = params
					c.Function = def
				}
			}
		}
		return nil
	})
}
