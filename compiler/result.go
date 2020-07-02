package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir/value"
)

type cresult struct {
	value value.Value
	ast   ast.Expression
}

func makeCR(e ast.Expression, v value.Value) cresult {
	return cresult{value: v, ast: e}
}
