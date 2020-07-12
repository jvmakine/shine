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

// Primitive types
var (
	Int    = NewPrimitive("int")
	Bool   = NewPrimitive("bool")
	Real   = NewPrimitive("real")
	String = NewPrimitive("string")
)

type Type interface {
	Copy(ctx *TypeCopyCtx) Type
	convert(s Substitutions, ctx *unificationCtx) Type

	freeVars(ctx *unificationCtx) []Variable
	unifier(o Type, ctx *unificationCtx) (Substitutions, error)
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

func (t Primitive) unifier(o Type, ctx *unificationCtx) (Substitutions, error) {
	if p, ok := o.(Primitive); ok && p.ID == t.ID {
		return MakeSubstitutions(), nil
	} else {
		return MakeSubstitutions(), UnificationError(t, o)
	}
}

func (t Primitive) convert(s Substitutions, ctx *unificationCtx) Type {
	return t
}

func (t Primitive) freeVars(ctx *unificationCtx) []Variable {
	return []Variable{}
}

func (t Primitive) signature(ctx *signatureContext) string {
	return t.ID
}

type Function struct {
	Replacable

	ID     VariableID
	fields []Type
}

func NewFunction(ts ...Type) Function {
	return Function{fields: ts, ID: NewVariableID()}
}

func (t Function) GetVariableID() VariableID {
	return t.ID
}

func (t Function) Copy(ctx *TypeCopyCtx) Type {
	ts := make([]Type, len(t.fields))
	for i, f := range t.fields {
		ts[i] = f.Copy(ctx)
	}
	return NewFunction(ts...)
}

func (t Function) unifier(o Type, ctx *unificationCtx) (Substitutions, error) {
	if fun, ok := o.(Function); ok && len(fun.fields) == len(t.fields) {
		result := MakeSubstitutions()
		for i, f := range t.fields {
			s, err := unifier(f, fun.fields[i], ctx)
			if err != nil {
				return MakeSubstitutions(), err
			}
			err = result.Combine(s)
			if err != nil {
				return MakeSubstitutions(), err
			}
		}
		return result, nil
	}
	return MakeSubstitutions(), UnificationError(t, o)
}

func (t Function) convert(s Substitutions, ctx *unificationCtx) Type {
	if s.substitutions[t.ID] != nil {
		return s.substitutions[t.ID]
	}
	ts := make([]Type, len(t.fields))
	for i, f := range t.fields {
		ts[i] = f.convert(s, ctx)
	}
	return NewFunction(ts...)
}

func (t Function) freeVars(ctx *unificationCtx) []Variable {
	res := []Variable{}
	for _, f := range t.fields {
		res = append(res, f.freeVars(ctx)...)
	}
	return res
}

func (t Function) signature(ctx *signatureContext) string {
	var sb strings.Builder
	sb.WriteString("(")
	if len(t.fields) > 1 {
		for i, p := range t.fields {
			sb.WriteString(p.signature(ctx))
			if i < len(t.fields)-2 {
				sb.WriteString(",")
			} else if i < len(t.fields)-1 {
				sb.WriteString(")=>")
			}
		}
	} else {
		sb.WriteString(")=>")
		sb.WriteString(t.fields[0].signature(ctx))
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

func (t Named) unifier(o Type, ctx *unificationCtx) (Substitutions, error) {
	if n, ok := o.(Named); ok && n.Name != t.Name {
		return MakeSubstitutions(), UnificationError(t, o)
	}
	return unifier(t.Type, o, ctx)
}

func (t Named) convert(s Substitutions, ctx *unificationCtx) Type {
	return Named{
		Name: t.Name,
		Type: t.Type.convert(s, ctx),
	}
}

func (t Named) freeVars(ctx *unificationCtx) []Variable {
	return t.Type.freeVars(ctx)
}

func (t Named) signature(ctx *signatureContext) string {
	if ctx.definingNamed[t.Name] {
		return t.Name
	}
	ctx.definingNamed[t.Name] = true
	var sb strings.Builder
	sb.WriteString(t.Name)
	sb.WriteString("[")
	sb.WriteString(t.Type.signature(ctx))
	sb.WriteString("]")
	return sb.String()
}

type Structure struct {
	fields []Named
}

func NewStructure(fields ...Named) Structure {
	for _, n := range fields {
		if len(n.freeVars(NewUnificationCtx())) > 0 {
			panic("free variables in a structure")
		}
	}
	return Structure{
		fields: fields,
	}
}

func (t Structure) Copy(ctx *TypeCopyCtx) Type {
	return t
}

func (t Structure) unifier(o Type, ctx *unificationCtx) (Substitutions, error) {
	if s, ok := o.(Structure); ok && len(s.fields) == len(t.fields) {
		resmap := map[string]Type{}
		for _, f := range t.fields {
			resmap[f.Name] = f.Type
		}
		for _, f := range s.fields {
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

func (t Structure) convert(s Substitutions, ctx *unificationCtx) Type {
	return t
}

func (t Structure) freeVars(ctx *unificationCtx) []Variable {
	return []Variable{}
}

func (s Structure) signature(ctx *signatureContext) string {
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for _, f := range s.fields {
		p := f.Type
		sb.WriteString(f.Name)
		sb.WriteString(":")
		sb.WriteString(p.signature(ctx))
		if i < len(s.fields)-1 {
			sb.WriteString(",")
		}
		i++
	}
	sb.WriteString("}")
	return sb.String()
}

type Variable struct {
	Replacable

	ID VariableID
}

func NewVariable() Variable {
	return Variable{
		ID: NewVariableID(),
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
	copy := Variable{
		ID: NewVariableID(),
	}
	ctx.vars[t.ID] = copy
	return copy
}

func (t Variable) unifier(o Type, ctx *unificationCtx) (Substitutions, error) {
	result := MakeSubstitutions()
	if v, ok := o.(Variable); ok {
		result.Update(v.ID, t)
	}
	result.Update(t.ID, o)
	return result, nil
}

func (t Variable) convert(s Substitutions, ctx *unificationCtx) Type {
	if r := s.substitutions[t.ID]; r != nil {
		return r
	}
	return t
}

func (t Variable) freeVars(ctx *unificationCtx) []Variable {
	return []Variable{t}
}

func (t Variable) signature(ctx *signatureContext) string {
	c := ctx.variables[t.ID]
	if c != "" {
		return c
	}
	ctx.variableCount++
	str := "V" + strconv.Itoa(ctx.variableCount)
	ctx.variables[t.ID] = str
	return str
}

type UnionVariable struct {
	Replacable

	ID      VariableID
	Options []Type
}

func NewUnionVariable(types ...Type) UnionVariable {
	return UnionVariable{
		ID:      NewVariableID(),
		Options: types,
	}
}

func (t UnionVariable) GetVariableID() VariableID {
	return t.ID
}

func (t UnionVariable) Copy(ctx *TypeCopyCtx) Type {
	c := ctx.vars[t.ID]
	if c != nil {
		return c
	}
	ts := make([]Type, len(t.Options))
	for i, t := range t.Options {
		ts[i] = t.Copy(ctx)
	}
	copy := UnionVariable{
		ID:      NewVariableID(),
		Options: ts,
	}
	ctx.vars[t.ID] = copy
	return copy
}

func deduplicate(types []Type) ([]Type, Substitutions) {
	sub := MakeSubstitutions()
	res := []Type{}
	handled := map[int]bool{}

	for i, t1 := range types {
		if handled[i] {
			continue
		}
		handled[i] = true
		gen := t1

		for j, t2 := range types {
			if handled[j] {
				continue
			}
			u1, e1 := gen.unifier(t2, NewUnificationCtx())
			u2, e2 := t2.unifier(gen, NewUnificationCtx())
			if e1 == nil && e2 != nil {
				handled[j] = true
				if err := sub.Combine(u1); err != nil {
					panic(err)
				}
			} else if e1 != nil && e2 == nil {
				handled[j] = true
				if err := sub.Combine(u2); err != nil {
					panic(err)
				}
				gen = t2
			} else if e1 == nil && e2 == nil {
				handled[j] = true
				gen = t1.convert(u1, NewUnificationCtx())
				if err := sub.Combine(u1); err != nil {
					panic(err)
				}
				if err := sub.Combine(u2); err != nil {
					panic(err)
				}
			}
		}
		res = append(res, gen)
	}

	return res, sub
}

func (t UnionVariable) unifier(o Type, ctx *unificationCtx) (Substitutions, error) {
	result := MakeSubstitutions()
	other := []Type{o}
	if u, ok := o.(UnionVariable); ok {
		other = u.Options
	}
	conversions := []Type{}
	for _, a := range t.Options {
		for _, b := range other {
			subs, e := Unifier(a, b)
			if e == nil {
				conversions = append(conversions, Convert(a, subs))
			}
		}
	}

	dd, s := deduplicate(conversions)
	err := result.Combine(s)
	if err != nil {
		return MakeSubstitutions(), err
	}
	conversions = dd

	if len(conversions) == 0 {
		return MakeSubstitutions(), UnificationError(t, o)
	}
	to := conversions[0]
	if len(conversions) > 1 {
		to = NewUnionVariable(conversions...)
	}
	result.Update(t.ID, to)
	if v, ok := o.(Replacable); ok {
		result.Update(v.GetVariableID(), to)
	}
	return result, nil
}

func (t UnionVariable) convert(s Substitutions, ctx *unificationCtx) Type {
	if sub := s.substitutions[t.ID]; sub != nil {
		return sub
	}
	return t
}

func (t UnionVariable) freeVars(ctx *unificationCtx) []Variable {
	res := []Variable{}
	for _, t := range t.Options {
		res = append(res, t.freeVars(ctx)...)
	}
	return res
}

func (t UnionVariable) signature(ctx *signatureContext) string {
	c := ctx.variables[t.ID]
	if c != "" {
		return c
	}
	ctx.variableCount++
	var sb strings.Builder
	sb.WriteString("V" + strconv.Itoa(ctx.variableCount))
	sb.WriteString("[")
	for i, o := range t.Options {
		sb.WriteString(o.signature(ctx))
		if i < len(t.Options)-1 {
			sb.WriteString("|")
		}
	}
	sb.WriteString("]")
	str := sb.String()
	ctx.variables[t.ID] = str
	return str
}

type StructuralVar struct {
	Replacable

	ID     VariableID
	fields map[string]Type
}

func NewStructuralVar(fields ...Named) StructuralVar {
	fs := map[string]Type{}
	for _, f := range fields {
		fs[f.Name] = f.Type
	}
	return StructuralVar{
		ID:     NewVariableID(),
		fields: fs,
	}
}

func (t StructuralVar) GetVariableID() VariableID {
	return t.ID
}

func (t StructuralVar) Copy(ctx *TypeCopyCtx) Type {
	c := ctx.vars[t.ID]
	if c != nil {
		return c
	}
	fs := map[string]Type{}
	for n, f := range t.fields {
		fs[n] = f.Copy(ctx)
	}
	copy := StructuralVar{
		ID:     NewVariableID(),
		fields: fs,
	}
	ctx.vars[t.ID] = copy
	return copy
}

func (t StructuralVar) unifier(o Type, ctx *unificationCtx) (Substitutions, error) {
	if v, ok := o.(StructuralVar); ok {
		result := MakeSubstitutions()
		sum := map[string]Type{}
		for n, f := range t.fields {
			if of := v.fields[n]; of != nil {
				sub, err := Unifier(of, f)
				if err != nil {
					return MakeSubstitutions(), err
				}
				err = result.Combine(sub)
				if err != nil {
					return MakeSubstitutions(), err
				}
				sum[n] = Convert(f, sub)
			} else {
				sum[n] = f
			}
		}
		for n, f := range v.fields {
			if sum[n] == nil {
				sum[n] = f
			}
		}
		res := StructuralVar{
			fields: sum,
			ID:     NewVariableID(),
		}
		result.Update(v.ID, res)
		result.Update(t.ID, res)
		return result, nil
	}
	if v, ok := o.(Structure); ok {
		sfields := map[string]Type{}
		for _, f := range v.fields {
			sfields[f.Name] = f.Type
		}
		for n, f := range t.fields {
			if p := sfields[n]; p != nil {
				_, err := Unifier(p, f)
				if err != nil {
					return MakeSubstitutions(), err
				}
			} else {
				return MakeSubstitutions(), UnificationError(t, o)
			}
		}
		result := MakeSubstitutions()
		result.Update(t.ID, v)
		return result, nil
	}
	return MakeSubstitutions(), UnificationError(t, o)
}

func (t StructuralVar) convert(s Substitutions, ctx *unificationCtx) Type {
	if r := s.substitutions[t.ID]; r != nil {
		return r
	}
	res := map[string]Type{}
	for n, f := range t.fields {
		res[n] = f.convert(s, ctx)
	}
	return StructuralVar{
		fields: res,
		ID:     NewVariableID(),
	}
}

func (t StructuralVar) freeVars(ctx *unificationCtx) []Variable {
	res := []Variable{}
	for _, f := range t.fields {
		res = append(res, f.freeVars(ctx)...)
	}
	return res
}

func (t StructuralVar) signature(ctx *signatureContext) string {
	keys := []string{}
	for n := range t.fields {
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
	sb.WriteString("{")
	i := 0
	for _, n := range keys {
		f := t.fields[n]
		sb.WriteString(n)
		sb.WriteString(":")
		sb.WriteString(f.signature(ctx))
		if i < len(t.fields)-1 {
			sb.WriteString(",")
		}
		i++
	}
	sb.WriteString("}")
	return sb.String()
}
