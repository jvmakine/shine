package compiler

import (
	t "github.com/jvmakine/shine/types"
	"github.com/llir/llvm/ir/types"
)

var (
	ClosurePType = types.I8Ptr
	FunType      = types.NewVector(2, types.I8Ptr)
	StruType     = types.I8Ptr

	IntType  = types.I64
	BoolType = types.I1
	RealType = types.Double
)

func structureType(s *t.Structure) types.Type {
	ps := make([]types.Type, len(s.Fields)+2)
	ps[0] = types.I32 // reference count
	ps[1] = types.I16 // number of closures
	// TODO: handle structures
	closures := 0
	for _, p := range s.Fields {
		if p.Type.IsFunction() {
			ps[closures+2] = getType(p.Type)
			closures++
		}
	}
	nonclosures := 0
	for _, p := range s.Fields {
		if !p.Type.IsFunction() {
			ps[2+closures+nonclosures] = getType(p.Type)
			nonclosures++
		}
	}
	return types.NewStruct(ps...)
}

func closureType(c *t.Closure) types.Type {
	ps := make([]types.Type, len(*c)+2)
	ps[0] = types.I32 // reference count
	ps[1] = types.I16 // number of closures
	// TODO: handle structures
	closures := 0
	for _, p := range *c {
		if p.Type.IsFunction() {
			ps[closures+2] = getType(p.Type)
			closures++
		}
	}
	nonclosures := 0
	for _, p := range *c {
		if !p.Type.IsFunction() {
			ps[2+closures+nonclosures] = getType(p.Type)
			nonclosures++
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
