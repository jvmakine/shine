package types

import (
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type VariableID string

func NewVariableID() VariableID {
	uid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return VariableID(uid.String())
}

type TypeCopyCtx struct {
	vars map[VariableID]Type
}

func NewTypeCopyCtx() *TypeCopyCtx {
	return &TypeCopyCtx{
		vars: map[VariableID]Type{},
	}
}

type UnificationCtx interface {
	StructuralTypeFor(name string, typ Type) Type
}

// Primitive types
var (
	Int    = NewPrimitive("int")
	Bool   = NewPrimitive("bool")
	Real   = NewPrimitive("real")
	String = NewPrimitive("string")
)

type Type interface {
	Copy(ctx *TypeCopyCtx) Type
	Convert(s Substitutions) (Type, bool)

	freeVars() []Variable
	unifier(o Type, ctx UnificationCtx) (Substitutions, error)
	signature(ctx *signatureContext) string
}

type Replacable interface {
	GetVariableID() VariableID
}

type Primitive struct {
	ID string
}

func NewPrimitive(id string) Primitive {
	return Primitive{ID: id}
}

func (t Primitive) Copy(ctx *TypeCopyCtx) Type {
	return t
}

func (t Primitive) unifier(o Type, ctx UnificationCtx) (Substitutions, error) {
	if p, ok := o.(Primitive); ok && p.ID == t.ID {
		return MakeSubstitutions(), nil
	} else {
		return MakeSubstitutions(), UnificationError(t, o)
	}
}

func (t Primitive) Convert(s Substitutions) (Type, bool) {
	return t, false
}

func (t Primitive) freeVars() []Variable {
	return []Variable{}
}

func (t Primitive) signature(ctx *signatureContext) string {
	return t.ID
}

type Function struct {
	Fields []Type
}

func NewFunction(ts ...Type) Function {
	return Function{Fields: ts}
}

func (t Function) Copy(ctx *TypeCopyCtx) Type {
	ts := make([]Type, len(t.Fields))
	for i, f := range t.Fields {
		ts[i] = f.Copy(ctx)
	}
	return NewFunction(ts...)
}

func (t Function) unifier(o Type, ctx UnificationCtx) (Substitutions, error) {
	if fun, ok := o.(Function); ok && len(fun.Fields) == len(t.Fields) {
		result := MakeSubstitutions()
		for i, f := range t.Fields {
			s, err := unifier(f, fun.Fields[i], ctx)
			if err != nil {
				return MakeSubstitutions(), err
			}
			err = result.Combine(s, ctx)
			if err != nil {
				return MakeSubstitutions(), err
			}
		}
		return result, nil
	}
	return MakeSubstitutions(), UnificationError(t, o)
}

func (t Function) Convert(s Substitutions) (Type, bool) {
	changed := false
	ts := make([]Type, len(t.Fields))
	for i, f := range t.Fields {
		nt, c := f.Convert(s)
		ts[i] = nt
		changed = changed || c
	}
	return NewFunction(ts...), changed
}

func (t Function) freeVars() []Variable {
	res := []Variable{}
	for _, f := range t.Fields {
		res = append(res, f.freeVars()...)
	}
	return res
}

func (t Function) signature(ctx *signatureContext) string {
	var sb strings.Builder
	sb.WriteString("(")
	if len(t.Fields) > 1 {
		for i, p := range t.Fields {
			sb.WriteString(p.signature(ctx))
			if i < len(t.Fields)-2 {
				sb.WriteString(",")
			} else if i < len(t.Fields)-1 {
				sb.WriteString(")=>")
			}
		}
	} else {
		sb.WriteString(")=>")
		sb.WriteString(t.Fields[0].signature(ctx))
	}
	return sb.String()
}

type Named struct {
	Name string
	Type Type
}

func NewNamed(name string, typ Type) Named {
	return Named{
		Name: name,
		Type: typ,
	}
}

func (t Named) Copy(ctx *TypeCopyCtx) Type {
	return Named{
		Name: t.Name,
		Type: t.Type.Copy(ctx),
	}
}

func (t Named) unifier(o Type, ctx UnificationCtx) (Substitutions, error) {
	if n, ok := o.(Named); ok && n.Name != t.Name {
		return MakeSubstitutions(), UnificationError(t, o)
	}
	return unifier(t.Type, o, ctx)
}

func (t Named) Convert(s Substitutions) (Type, bool) {
	nt, changed := t.Type.Convert(s)
	return Named{
		Name: t.Name,
		Type: nt,
	}, changed
}

func (t Named) freeVars() []Variable {
	return t.Type.freeVars()
}

func (t Named) signature(ctx *signatureContext) string {
	return t.Name
}

type Structure struct {
	Fields []Named
}

func NewStructure(Fields ...Named) Structure {
	return Structure{
		Fields: Fields,
	}
}

func (t Structure) Copy(ctx *TypeCopyCtx) Type {
	ts := make([]Named, len(t.Fields))
	for i, f := range t.Fields {
		ts[i] = NewNamed(f.Name, f.Type.Copy(ctx))
	}
	return NewStructure(ts...)
}

func (t Structure) unifier(o Type, ctx UnificationCtx) (Substitutions, error) {
	if s, ok := o.(Structure); ok && len(s.Fields) == len(t.Fields) {
		resmap := map[string]Type{}
		for _, f := range t.Fields {
			resmap[f.Name] = f.Type
		}
		for _, f := range s.Fields {
			p := resmap[f.Name]
			if p != nil {
				_, err := unifier(f.Type, p, ctx)
				if err != nil {
					return MakeSubstitutions(), err
				}
			} else {
				return MakeSubstitutions(), UnificationError(t, o)
			}
		}
		return MakeSubstitutions(), nil
	}
	return MakeSubstitutions(), UnificationError(t, o)
}

func (t Structure) Convert(s Substitutions) (Type, bool) {
	ts := make([]Named, len(t.Fields))
	changed := false
	for i, f := range t.Fields {
		n, c := f.Type.Convert(s)
		ts[i] = NewNamed(f.Name, n)
		changed = changed || c
	}
	return NewStructure(ts...), changed
}

func (t Structure) freeVars() []Variable {
	res := []Variable{}
	for _, f := range t.Fields {
		res = append(res, f.Type.freeVars()...)
	}
	return res
}

func (s Structure) signature(ctx *signatureContext) string {
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for _, f := range s.Fields {
		p := f.Type
		sb.WriteString(f.Name)
		sb.WriteString(":")
		sb.WriteString(p.signature(ctx))
		if i < len(s.Fields)-1 {
			sb.WriteString(",")
		}
		i++
	}
	sb.WriteString("}")
	return sb.String()
}

type Variable struct {
	Replacable

	ID     VariableID
	Fields map[string]Type
}

func NewVariable(Fields ...Named) Variable {
	fs := map[string]Type{}
	for _, f := range Fields {
		fs[f.Name] = f.Type
	}
	return Variable{
		ID:     NewVariableID(),
		Fields: fs,
	}
}

func (t Variable) GetVariableID() VariableID {
	return t.ID
}

func (t Variable) Copy(ctx *TypeCopyCtx) Type {
	c := ctx.vars[t.ID]
	if c != nil {
		return c
	}
	fs := map[string]Type{}
	for n, f := range t.Fields {
		fs[n] = f.Copy(ctx)
	}
	copy := Variable{
		ID:     NewVariableID(),
		Fields: fs,
	}
	ctx.vars[t.ID] = copy
	return copy
}

func (t Variable) unifier(o Type, ctx UnificationCtx) (Substitutions, error) {
	if len(t.Fields) == 0 {
		result := MakeSubstitutions()
		if err := result.Update(t.ID, o, ctx); err != nil {
			return MakeSubstitutions(), err
		}
		if v, ok := o.(Variable); ok {
			if err := result.Update(v.ID, t, ctx); err != nil {
				return MakeSubstitutions(), err
			}
		}
		return result, nil
	}
	if v, ok := o.(Variable); ok {
		result := MakeSubstitutions()
		sum := map[string]Type{}
		for n, f := range t.Fields {
			if of := v.Fields[n]; of != nil {
				sub, err := Unifier(of, f, ctx)
				if err != nil {
					return MakeSubstitutions(), err
				}
				err = result.Combine(sub, ctx)
				if err != nil {
					return MakeSubstitutions(), err
				}
				sum[n], _ = f.Convert(sub)
			} else {
				sum[n] = f
			}
		}
		for n, f := range v.Fields {
			if sum[n] == nil {
				sum[n] = f
			}
		}
		res := Variable{
			Fields: sum,
			ID:     NewVariableID(),
		}
		result.Update(v.ID, res, ctx)
		result.Update(t.ID, res, ctx)
		return result, nil
	}
	if v, ok := o.(Structure); ok {
		sFields := map[string]Type{}
		for _, f := range v.Fields {
			sFields[f.Name] = f.Type
		}
		for n, f := range t.Fields {
			if p := sFields[n]; p != nil {
				_, err := Unifier(p, f, ctx)
				if err != nil {
					return MakeSubstitutions(), err
				}
			} else {
				return MakeSubstitutions(), UnificationError(t, o)
			}
		}
		result := MakeSubstitutions()
		result.Update(t.ID, v, ctx)
		return result, nil
	}
	if f, ok := o.(Function); ok && len(f.freeVars()) > 0 {
		stru := NewVariable(NewNamed("%call", f))
		return unifier(stru, t, ctx)
	}
	result := MakeSubstitutions()
	for name, typ := range t.Fields {
		in := ctx.StructuralTypeFor(name, o)
		if in == nil {
			return MakeSubstitutions(), UnificationError(t, o)
		}
		sub, err := unifier(in, typ, ctx)
		if err != nil {
			return MakeSubstitutions(), err
		}
		err = result.Combine(sub, ctx)
		if err != nil {
			return MakeSubstitutions(), err
		}
	}
	result.Update(t.ID, o, ctx)
	return result, nil
}

func (t Variable) Convert(s Substitutions) (Type, bool) {
	if r := (*s.substitutions)[t.ID]; r != nil {
		return r, true
	}
	if len(t.Fields) == 0 {
		return t, false
	}
	changed := false
	res := map[string]Type{}
	for n, f := range t.Fields {
		nv, c := f.Convert(s)
		res[n] = nv
		changed = changed || c
	}
	if !changed {
		return t, false
	}
	return Variable{
		Fields: res,
		ID:     NewVariableID(),
	}, true
}

func (t Variable) freeVars() []Variable {
	res := []Variable{t}
	for _, f := range t.Fields {
		res = append(res, f.freeVars()...)
	}
	return res
}

func (t Variable) signature(ctx *signatureContext) string {
	keys := []string{}
	for n := range t.Fields {
		keys = append(keys, n)
	}
	sort.Strings(keys)

	var sb strings.Builder
	if ctx.variables[t.ID] != "" {
		sb.WriteString(ctx.variables[t.ID])
	} else {
		ctx.variableCount++
		str := "V" + strconv.Itoa(ctx.variableCount)
		sb.WriteString(str)
		ctx.variables[t.ID] = str
	}
	if len(t.Fields) > 0 {
		sb.WriteString("{")
		i := 0
		for _, n := range keys {
			f := t.Fields[n]
			sb.WriteString(n)
			sb.WriteString(":")
			sb.WriteString(f.signature(ctx))
			if i < len(t.Fields)-1 {
				sb.WriteString(",")
			}
			i++
		}
		sb.WriteString("}")
	}
	return sb.String()
}
