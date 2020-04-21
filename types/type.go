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

type TVarID = int64

type TypeVar struct {
	Restrictions Restrictions
}

type Type struct {
	Function  *[]Type
	Variable  *TypeVar
	Primitive *Primitive
}

func MakeVariable() Type {
	return Type{Variable: &TypeVar{}}
}

func MakePrimitive(p string) Type {
	return Type{Primitive: &p}
}

func MakeFunction(ts ...Type) Type {
	return Type{Function: &ts}
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
	return t.Function != nil
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
	return len(t.FreeVars()) == 0
}

func (t Type) AsPrimitive() Primitive {
	if !t.IsPrimitive() {
		panic("type not primitive type")
	}
	return *t.Primitive
}

func (t Type) ReturnType() Type {
	return (*t.Function)[len(*t.Function)-1]
}

func (t *Type) AssignFrom(o Type) {
	t.Variable = o.Variable
	t.Function = o.Function
	t.Primitive = o.Primitive
}
