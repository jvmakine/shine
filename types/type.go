package types

type Primitive = string

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

type Named struct {
	Name          string
	TypeArguments []Type
}

type SField struct {
	Name string
	Type Type
}

type Structure struct {
	Name          string
	TypeArguments []Type
	Fields        []SField

	OrginalVars    []*TypeVar
	OriginalFields []SField
}

type TypeVar struct {
	Union      Union
	Structural map[string]Type
}

type HierarchicalVar struct {
	Root   *TypeVar
	Params []Type
}

type TCBinding struct {
	Name string
	Args []Type
}

type TypeClassRef struct {
	TypeClass     string
	TypeClassVars []Type
	Place         int
	LocalBindings []TCBinding
}

type Type struct {
	Function  *Function
	Structure *Structure
	Variable  *TypeVar
	HVariable *HierarchicalVar
	Primitive *Primitive
	Named     *Named
	TCRef     *TypeClassRef
}

type Types []Type

func (ts Types) Copy(ctx *TypeCopyCtx) Types {
	res := make(Types, len(ts))
	for i, t := range ts {
		res[i] = t.Copy(ctx)
	}
	return res
}

func WithType(t Type, f func(t Type) Type) Type {
	return f(t)
}

func MakeVariable() Type {
	return Type{Variable: &TypeVar{}}
}

func MakeHierarchicalVar(root *TypeVar, inner ...Type) Type {
	return Type{HVariable: &HierarchicalVar{Root: root, Params: inner}}
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

func MakeNamed(name string, types ...Type) Type {
	t := types
	if t == nil {
		t = []Type{}
	}
	return Type{Named: &Named{Name: name, TypeArguments: t}}
}

func MakeFunction(ts ...Type) Type {
	var f Function = ts
	return Type{Function: &f}
}

func MakeStructure(name string, fields ...SField) Type {
	vars := map[*TypeVar]bool{}
	copfields := make([]SField, len(fields))
	for i, f := range fields {
		copfields[i] = SField{Name: f.Name, Type: f.Type}
		for _, v := range f.Type.FreeVars() {
			vars[v] = true
		}
	}

	targlist := make([]Type, 0, len(vars))
	varlist := make([]*TypeVar, 0, len(vars))
	for v := range vars {
		varlist = append(varlist, v)
		targlist = append(targlist, Type{Variable: v})
	}

	return Type{Structure: &Structure{
		Name:          name,
		Fields:        fields,
		TypeArguments: targlist,

		OriginalFields: copfields,
		OrginalVars:    varlist,
	}}
}

func MakeTypeClassRef(name string, place int, fields ...Type) Type {
	return Type{TCRef: &TypeClassRef{TypeClass: name, TypeClassVars: fields, Place: place}}
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
	if t.HVariable != nil {
		res := []*TypeVar{t.HVariable.Root}
		for _, p := range t.HVariable.Params {
			res = append(res, p.FreeVars()...)
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
	if tc := t.TCRef; tc != nil {
		res := []*TypeVar{}
		for _, p := range tc.TypeClassVars {
			res = append(res, p.FreeVars()...)
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

func (t Type) IsTypeClassRef() bool {
	return t.TCRef != nil
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

func (t Type) IsHVariable() bool {
	return t.HVariable != nil
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
	return t.Function != nil ||
		t.Variable != nil ||
		t.Primitive != nil ||
		t.Structure != nil ||
		t.Named != nil ||
		t.TCRef != nil ||
		t.HVariable != nil
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

func (t Type) NamedTypes() map[string]bool {
	if t.IsNamed() {
		res := map[string]bool{t.Named.Name: true}
		for _, f := range t.Named.TypeArguments {
			v := f.NamedTypes()
			for n := range v {
				res[n] = true
			}
		}
		return res
	}
	if t.IsFunction() {
		res := map[string]bool{}
		for _, f := range *t.Function {
			v := f.NamedTypes()
			for n := range v {
				res[n] = true
			}
		}
		return res
	}
	if t.IsStructure() {
		res := map[string]bool{}
		for _, f := range t.Structure.Fields {
			v := f.Type.NamedTypes()
			for n := range v {
				res[n] = true
			}
		}
		return res
	}
	if t.IsTypeClassRef() {
		res := map[string]bool{}
		for _, f := range t.TCRef.TypeClassVars {
			v := f.NamedTypes()
			for n := range v {
				res[n] = true
			}
		}
		return res
	}
	if t.IsHVariable() {
		res := map[string]bool{}
		for _, f := range *&t.HVariable.Params {
			v := f.NamedTypes()
			for n := range v {
				res[n] = true
			}
		}
		return res
	}
	return map[string]bool{}
}

func (t Type) Rewrite(f func(Type) (Type, error)) (Type, error) {
	if t.IsFunction() {
		fn := make([]Type, len(*t.Function))
		for i, a := range *t.Function {
			b, err := a.Rewrite(f)
			if err != nil {
				return Type{}, err
			}
			fn[i] = b
		}
		return f(MakeFunction(fn...))
	} else if t.IsNamed() {
		fn := make([]Type, len(t.Named.TypeArguments))
		for i, a := range t.Named.TypeArguments {
			b, err := a.Rewrite(f)
			if err != nil {
				return Type{}, err
			}
			fn[i] = b
		}
		return f(MakeNamed(t.Named.Name, fn...))
	} else if s := t.Structure; s != nil {
		nf := make([]SField, len(s.Fields))
		for i, a := range s.Fields {
			b, err := a.Type.Rewrite(f)
			if err != nil {
				return Type{}, err
			}
			nf[i] = SField{Name: a.Name, Type: b}
		}
		nt := make([]Type, len(s.TypeArguments))
		for i, a := range s.TypeArguments {
			b, err := a.Rewrite(f)
			if err != nil {
				return Type{}, err
			}
			nt[i] = b
		}
		return f(Type{Structure: &Structure{
			Name:           s.Name,
			TypeArguments:  nt,
			Fields:         nf,
			OrginalVars:    s.OrginalVars,
			OriginalFields: s.OriginalFields,
		}})
	} else if t.IsTypeClassRef() {
		nf := make([]Type, len(t.TCRef.TypeClassVars))
		for i, a := range t.TCRef.TypeClassVars {
			b, err := a.Rewrite(f)
			if err != nil {
				return Type{}, err
			}
			nf[i] = b
		}
		c := MakeTypeClassRef(t.TCRef.TypeClass, t.TCRef.Place, nf...)
		c.TCRef.LocalBindings = t.TCRef.LocalBindings
		return f(c)
	} else if t.IsHVariable() {
		fn := make([]Type, len(t.HVariable.Params))
		for i, a := range t.HVariable.Params {
			b, err := a.Rewrite(f)
			if err != nil {
				return Type{}, err
			}
			fn[i] = b
		}
		return f(MakeHierarchicalVar(t.HVariable.Root, fn...))
	}
	return f(t)
}

func (s *Structure) GetField(name string) *Type {
	for _, f := range s.Fields {
		if f.Name == name {
			return &f.Type
		}
	}
	return nil
}

func (s *Structure) FieldTypes() []Type {
	res := make([]Type, len(s.Fields))
	for i, f := range s.Fields {
		res[i] = f.Type
	}
	return res
}

func (t Type) Instantiate(types []Type) Type {
	if s := t.Structure; s != nil {
		return Type{Structure: s.Instantiate(types)}
	}
	return t
}

func (s *Structure) Instantiate(types []Type) *Structure {
	subs := MakeSubstitutions()
	for i, t := range types {
		err := subs.Update(s.OrginalVars[i], t)
		if err != nil {
			panic(err)
		}
	}

	args := make([]Type, len(s.TypeArguments))
	for i, v := range s.OrginalVars {
		args[i] = subs.Apply(Type{Variable: v})
	}

	fields := make([]SField, len(s.Fields))
	for i, f := range s.Fields {
		fields[i] = SField{Name: f.Name, Type: subs.Apply(f.Type)}
	}

	return &Structure{
		Name:           s.Name,
		TypeArguments:  args,
		Fields:         fields,
		OrginalVars:    s.OrginalVars,
		OriginalFields: s.OriginalFields,
	}
}
