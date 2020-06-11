package compiler

import (
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	FunType      = types.I8Ptr
	ClosurePType = types.I8Ptr
	StruType     = types.I8Ptr

	StringPType     = types.I8Ptr
	StringType      = types.NewStruct(types.I32, types.I16, types.I16, types.I8Ptr, types.I8Ptr)
	ClosureCallType = types.NewStruct(types.I8, types.I32, types.I8Ptr)

	IntType  = types.I64
	BoolType = types.I1
	RealType = types.Double
)

func structureType(s *t.Structure, closure bool) types.Type {
	extra := 3
	if closure {
		extra = 4
	}
	ps := make([]types.Type, len(s.Fields)+extra)
	ps[0] = types.I8  // reference type
	ps[1] = types.I32 // reference count
	if closure {
		ps[2] = FunType
	}
	ps[extra-1] = types.I16 // number of structures

	structures := 0
	for _, p := range s.Fields {
		if p.Type.IsStructure() {
			ps[extra+structures] = getType(p.Type)
			structures++
		} else if p.Type.IsFunction() {
			ps[extra+structures] = FunType
			structures++
		}
	}
	primitives := 0
	for _, p := range s.Fields {
		if !p.Type.IsFunction() && !p.Type.IsStructure() {
			ps[extra+structures+primitives] = getType(p.Type)
			primitives++
		}
	}
	return types.NewStruct(ps...)
}

func getFunctPtr(fun t.Type) types.Type {
	ret := getType(fun.FunctReturn())
	fparams := fun.FunctParams()
	params := make([]types.Type, len(fparams)+1)
	for i, p := range fparams {
		params[i] = getType(p)
	}
	params[len(fparams)] = ClosurePType
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
		case t.String:
			rtype = StringPType
		default:
			panic("unsupported type at compilation")
		}
		return rtype
	} else if typ.IsFunction() {
		return FunType
	} else if typ.IsStructure() {
		return StruType
	} else if typ.IsNamed() {
		panic("trying to use named type at compilation")
	}
	panic("invalid type: " + typ.Signature())
}
