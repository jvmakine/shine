package compiler

import (
	"github.com/jvmakine/shine/typedef"
	t "github.com/jvmakine/shine/typeinferer"
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
	case typedef.Int:
		rtype = IntType
	case typedef.Bool:
		rtype = BoolType
	default:
		panic("unsupported type at compilation")
	}
	return rtype
}
