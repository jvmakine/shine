package compiler

import (
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	IntType  = types.I64
	BoolType = types.I1
)

func getType(from interface{}) types.Type {
	var rtype types.Type = nil
	typ := from.(*t.TypePtr)
	if typ == nil || typ.Def.Base == nil {
		panic("trying to use undefined type at compilation")
	}
	switch *(typ.Def.Base) {
	case t.Int:
		rtype = IntType
	case t.Bool:
		rtype = BoolType
	default:
		panic("unsupported type at compilation")
	}
	return rtype
}
