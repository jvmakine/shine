package types

type Primitive = string

const (
	Int  Primitive = "int"
	Bool Primitive = "bool"
	Real Primitive = "real"
)

var (
	IntP  = MakePrimitive(Int)
	BoolP = MakePrimitive(Bool)
	RealP = MakePrimitive(Real)
)

type Function []Type

type TypeVar struct {
	Restrictions Restrictions
	Function     *Function
}

type Type struct {
	Function  *Function
	Variable  *TypeVar
	Primitive *Primitive
}

func MakeVariable() Type {
	return Type{Variable: &TypeVar{}}
}

func MakePrimitive(p string) Type {
	return Type{Primitive: &p}
}

func WithType(t Type, f func(t Type) Type) Type {
	return f(t)
}

func MakeFunction(ts ...Type) Type {
	isvar := false
	for _, t := range ts {
		if t.IsVariable() {
			isvar = true
			break
		}
	}
	var f Function = ts
	if isvar {
		return Type{Variable: &TypeVar{Function: &f}}
	}
	return Type{Function: &f}
}

func MakeRestricted(ps ...Primitive) Type {
	return Type{Variable: &TypeVar{Restrictions: ps}}
}

func (t Type) FreeVars() []*TypeVar {
	if t.Variable != nil {
		return []*TypeVar{t.Variable}
	}
	if t.Function != nil {
		res := []*TypeVar{}
		for _, p := range *t.Function {
			res = append(res, p.FreeVars()...)
		}
		return res
	}
	return []*TypeVar{}
}

func (t Type) IsFunction() bool {
	return t.Function != nil || (t.Variable != nil && t.Variable.Function != nil)
}

func (t Type) FunctTypes() []Type {
	if !t.IsFunction() {
		panic("can not get params from a non-function")
	}
	f := t.Function
	if f == nil {
		f = t.Variable.Function
	}
	return *f
}

func (t Type) FunctTypesPtr() []*Type {
	if !t.IsFunction() {
		panic("can not get params from a non-function")
	}
	f := t.Function
	if f == nil {
		f = t.Variable.Function
	}
	ptrs := make([]*Type, len(*f))
	for i, t := range *f {
		ptrs[i] = &t
	}
	return ptrs
}

func (t Type) FunctParams() []Type {
	typs := t.FunctTypes()
	return typs[:len(typs)-1]
}

func (t Type) FunctReturn() Type {
	typs := t.FunctTypes()
	return typs[len(typs)-1]
}

func (t Type) IsPrimitive() bool {
	return t.Primitive != nil
}

func (t Type) IsVariable() bool {
	return t.Variable != nil
}

func (t Type) IsRestrictedVariable() bool {
	return t.IsVariable() && len(t.Variable.Restrictions) > 0
}

func (t Type) IsDefined() bool {
	return t.Function != nil || t.Variable != nil || t.Primitive != nil
}

func (t Type) HasFreeVars() bool {
	return len(t.FreeVars()) > 0
}

func (t Type) AsPrimitive() Primitive {
	if !t.IsPrimitive() {
		panic("type not primitive: " + t.Signature())
	}
	return *t.Primitive
}

func (t *Type) AssignFrom(o Type) {
	t.Variable = o.Variable
	t.Function = o.Function
	t.Primitive = o.Primitive
}
