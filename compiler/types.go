package compiler

import (
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	IntType  = types.I64
	BoolType = types.I1
	RealType = types.Double
)

func getType(typ t.Type) types.Type {
	var rtype types.Type = nil
	if !typ.IsDefined() {
		panic("trying to use undefined type at compilation")
	}
	switch typ.AsPrimitive() {
	case t.Int:
		rtype = IntType
	case t.Bool:
		rtype = BoolType
	case t.Real:
		rtype = RealType
	default:
		panic("unsupported type at compilation")
	}
	return rtype
}
