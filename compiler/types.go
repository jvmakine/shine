package compiler

import (
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	IntType  = types.I64
	BoolType = types.I1
	RealType = types.FP128
)

func getType(from interface{}) types.Type {
	var rtype types.Type = nil
	typ := from.(*t.TypePtr)
	if typ == nil || typ.Def.Bases == nil || len(typ.Def.Bases) > 1 {
		panic("trying to use undefined type at compilation")
	}
	switch *(typ.Def.Bases[0]) {
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
