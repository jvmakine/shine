package types

type Primitive = string

type Named = string

const (
	Int    Primitive = "int"
	Bool   Primitive = "bool"
	Real   Primitive = "real"
	String Primitive = "string"
)

var (
	IntP    = MakePrimitive(Int)
	BoolP   = MakePrimitive(Bool)
	RealP   = MakePrimitive(Real)
	StringP = MakePrimitive(String)
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

func MakeUnionVar(ts ...Type) Type {
	return Type{Variable: &TypeVar{Union: ts}}
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
		stru := t.Variable.Structural
		res := []*TypeVar{t.Variable}
		if len(stru) > 0 {
			for _, v := range stru {
				res = append(res, v.FreeVars()...)
			}
		}
		return res
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

func (t Type) IsString() bool {
	return t.IsPrimitive() && *(t.Primitive) == String
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

func (t Type) IsFreeVar() bool {
	return t.Variable != nil && t.Variable.Structural == nil && t.Variable.Union == nil
}

func (t Type) IsNamed() bool {
	return t.Named != nil
}

func (t Type) IsUnionVar() bool {
	return t.IsVariable() && len(t.Variable.Union) > 0
}

func (t Type) IsStructuralVar() bool {
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

func (t Type) IsGeneralisationOf(o Type) bool {
	if t.IsFreeVar() {
		return true
	}
	if t.IsStructuralVar() && o.IsStructuralVar() {
		for k, v := range t.Variable.Structural {
			if ot := o.Variable.Structural[k]; ot.IsDefined() {
				if !v.IsGeneralisationOf(ot) {
					return false
				}
			} else {
				return false
			}
		}
		return true
	}
	if t.IsPrimitive() && o.IsPrimitive() && *t.Primitive == *o.Primitive {
		return true
	}
	if t.IsFunction() && o.IsFunction() && t.UnifiesWith(o) {
		for i := range *t.Function {
			if !(*t.Function)[i].IsGeneralisationOf((*o.Function)[i]) {
				return false
			}
		}
		return true
	}
	if t.IsStructure() && o.IsStructure() {
		return t.UnifiesWith(o)
	}
	return false
}

func (t Type) AddToUnion(o Type) Type {
	var union Union
	if t.IsUnionVar() {
		union = append(t.Variable.Union, o).deduplicate()
	} else {
		union = Union{t, o}.deduplicate()
	}
	if len(union) == 1 {
		return union[0]
	}
	return MakeUnionVar(union...)
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
