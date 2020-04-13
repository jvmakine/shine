package compiler

import (
	"github.com/jvmakine/shine/ast"
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	IntType  = types.I64
	BoolType = types.I1
)

func getType(exp *ast.Exp) types.Type {
	var rtype types.Type = nil
	switch *(exp.Type.(*t.Type)).Base {
	case *(t.Int.Base):
		rtype = IntType
	case *(t.Bool.Base):
		rtype = BoolType
	default:
		panic("unsupported type at compilation")
	}
	return rtype
}
