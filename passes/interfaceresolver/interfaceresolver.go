package interfaceresolver

import (
	. "github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

func Resolve(exp Ast) {
	rewriter := func(f Ast, ctx *VisitContext) Ast {
		if op, ok := f.(*Op); ok {
			ic := op.Left.Type().(Contextual).GetContext()
			if ic == nil {
				panic("nil context")
			}
		}
		return f
	}
	exp.Visit(NullFun, NullFun, false, rewriter, NewVisitCtx())
}
