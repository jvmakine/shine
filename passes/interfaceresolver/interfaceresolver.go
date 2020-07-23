package interfaceresolver

import (
	. "github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func Resolve(exp Ast) {
	rewriter := func(f Ast, ctx *VisitContext) Ast {
		if op, ok := f.(*Op); ok {
			typ := op.Left.Type()
			ic := typ.(Contextual).GetContext().(*VisitContext)
			if ic == nil {
				panic("nil context")
			}
			inter := ic.InterfaceWithType(op.Name, typ)
			return inter.Interf.ReplaceOp(op)
		}
		return f
	}
	exp.Visit(NullFun, NullFun, false, rewriter, NewVisitCtx())
}
