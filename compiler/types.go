package compiler

import (
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	ClosurePType = types.I8Ptr
	FunPType     = types.I8Ptr
	FunType      = types.NewVector(2, types.I8Ptr)
	StruType     = types.I8Ptr

	IntType    = types.I64
	BoolType   = types.I1
	RealType   = types.Double
	StringType = types.I8Ptr
)

func structureType(s *t.Structure) types.Type {
	cc := 0
	for _, p := range s.Fields {
		if p.Type.IsFunction() {
			cc++
		}
	}

	ps := make([]types.Type, len(s.Fields)+3+cc)
	ps[0] = types.I32 // reference count
	ps[1] = types.I16 // number of closures
	ps[2] = types.I16 // number of structures

	closures := 0
	for _, p := range s.Fields {
		if p.Type.IsFunction() {
			ps[closures+3] = FunPType
			closures++
			ps[closures+3] = ClosurePType
			closures++
		}
	}
	structures := 0
	for _, p := range s.Fields {
		if p.Type.IsStructure() {
			ps[closures+3+structures] = getType(p.Type)
			structures++
		}
	}
	primitives := 0
	for _, p := range s.Fields {
		if !p.Type.IsFunction() && !p.Type.IsStructure() {
			ps[3+closures+structures+primitives] = getType(p.Type)
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
			rtype = StringType
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
