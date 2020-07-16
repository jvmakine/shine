package optimisation

import (
	. "github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

// Optimise sequential function definitions into one when called with multiple arguments
func SequentialFunctionPass(exp Expression) {
	tctx := NewTypeCopyCtx()

	CrawlBefore(exp, func(v Ast, ctx *VisitContext) error {
		if c, ok := v.(*FCall); ok {
			root := c.RootFunc()
			var def *FDef
			var defin *Definitions
			var id string
			var typ Type
			changed := false

			if i, ok := root.(*Id); ok {
				id = i.Name
				typ = i.IdType
				if defin = ctx.DefinitionOf(id); defin != nil {
					assig := defin.Assignments[id]
					if d, ok := assig.(*FDef); ok {
						def = d.CopyWithCtx(tctx).(*FDef)
					} else {
						def = nil
					}
				}
			} else if _, ok := root.(*FDef); ok {
				def = root.CopyWithCtx(tctx).(*FDef)
				typ = def.Type()
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

				if defin != nil {
					nid = nid + "%c"
				}

				_, isDefB = def.Body.(*FDef)
				fcall, isFCall = fcall.Function.(*FCall)
				ftyp := typ.(Function)
				ts := append(ftyp.Params(), (ftyp.Return().(Function).Fields)...)
				typ = NewFunction(ts...)
			}

			if changed {
				if defin != nil {
					defin.Assignments[nid] = def
					c.Function = &Id{Name: nid, IdType: typ}
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
