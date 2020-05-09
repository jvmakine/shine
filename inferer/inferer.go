package inferer

import (
	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/inferer/typeinference"
)

func Infer(e *ast.Exp) error {
	return typeinference.Infer(e)
}
