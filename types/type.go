package types

type Primitive = string

type Named = string

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

type SField struct {
	Name string
	Type Type
}

type Structure struct {
	Name   string
	Fields []SField
}

type TypeVar struct {
	Union      Union
	Structural map[string]Type
}

type Type struct {
	Function  *Function
	Structure *Structure
	Variable  *TypeVar
	Primitive *Primitive
	Named     *Named
}

func WithType(t Type, f func(t Type) Type) Type {
	return f(t)
}

func MakeVariable() Type {
	return Type{Variable: &TypeVar{}}
}

func MakeUnionVar(ps ...Primitive) Type {
	return Type{Variable: &TypeVar{Union: ps}}
}

func MakeStructuralVar(s map[string]Type) Type {
	return Type{Variable: &TypeVar{Structural: s}}
}

func MakePrimitive(p string) Type {
	return Type{Primitive: &p}
}

func MakeNamed(name string) Type {
	return Type{Named: &name}
}

func MakeFunction(ts ...Type) Type {
	var f Function = ts
	return Type{Function: &f}
}

func MakeStructure(name string, fields ...SField) Type {
	return Type{Structure: &Structure{Name: name, Fields: fields}}
}

func (t Type) FreeVars() []*TypeVar {
	if t.Variable != nil {
		return []*TypeVar{t.Variable}
	}
	if fun := t.Function; fun != nil {
		res := []*TypeVar{}
		for _, p := range *fun {
			res = append(res, p.FreeVars()...)
		}
		return res
	}
	if stru := t.Structure; stru != nil {
		res := []*TypeVar{}
		for _, p := range stru.Fields {
			res = append(res, p.Type.FreeVars()...)
		}
		return res
	}
	return []*TypeVar{}
}

func (t Type) IsFunction() bool {
	return t.Function != nil
}

func (t Type) IsStructure() bool {
	return t.Structure != nil
}

func (t Type) FunctTypes() []Type {
	if !t.IsFunction() {
		panic("can not get params from a non-function")
	}
	f := t.Function
	return *f
}

func (t Type) FunctTypesPtr() []*Type {
	if !t.IsFunction() {
		panic("can not get params from a non-function")
	}
	f := t.Function
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

func (t Type) IsNamed() bool {
	return t.Named != nil
}

func (t Type) IsUnionVar() bool {
	return t.IsVariable() && len(t.Variable.Union) > 0
}

func (t Type) IsStructurealVar() bool {
	return t.IsVariable() && len(t.Variable.Structural) > 0
}

func (t Type) IsDefined() bool {
	return t.Function != nil || t.Variable != nil || t.Primitive != nil || t.Structure != nil || t.Named != nil
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
	t.Structure = o.Structure
}

func (s *Structure) GetField(name string) *Type {
	for _, f := range s.Fields {
		if f.Name == name {
			return &f.Type
		}
	}
	return nil
}
