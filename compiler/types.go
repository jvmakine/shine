package compiler

import (
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	ClosurePType = types.I8Ptr
	IntType      = types.I64
	BoolType     = types.I1
	RealType     = types.Double
)

func getFunctPtr(fun t.Type) types.Type {
	ret := getType(fun.FunctReturn())
	fparams := fun.FunctParams()
	params := make([]types.Type, len(fparams))
	for i, p := range fparams {
		params[i] = getType(p)
	}
	return types.NewPointer(types.NewFunc(ret, params...))
}

func getType(typ t.Type) types.Type {
	if !typ.IsDefined() && !typ.IsVariable() {
		panic("trying to use undefined type at compilation")
	}
	if typ.IsPrimitive() {
		var rtype types.Type = nil
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
	} else if typ.IsFunction() {
		return getFunctPtr(typ)
	}
	panic("invalid type: " + typ.Signature())
}
