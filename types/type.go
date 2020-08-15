package types

import (
	"sort"
	"strconv"
	"strings"
)

type VariableID string

var variables = 0

func NewVariableID() VariableID {
	variables++
	return VariableID(strconv.Itoa(variables))
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

	convert(s Substitutions, sctx substitutionCtx) (Type, bool)
	freeVars(ctx freeVarsCtx) []Variable
	unifier(o Type, ctx UnificationCtx, sctx substitutionCtx) (Substitutions, error)
	signature(ctx *signatureContext) string
}

type Contextual interface {
	WithContext(ctx UnificationCtx) Contextual
	GetContext() UnificationCtx
}

type Primitive struct {
	ID string

	ctx UnificationCtx
}

func NewPrimitive(id string) Primitive {
	return Primitive{ID: id}
}

func (t Primitive) Copy(ctx *TypeCopyCtx) Type {
	return t
}

func (t Primitive) unifier(o Type, ctx UnificationCtx, sctx substitutionCtx) (Substitutions, error) {
	if p, ok := o.(Primitive); ok && p.ID == t.ID {
		return MakeSubstitutions(), nil
	} else {
		return MakeSubstitutions(), UnificationError(t, o)
	}
}

func (t Primitive) convert(s Substitutions, sctx substitutionCtx) (Type, bool) {
	return t, false
}

func (t Primitive) freeVars(ctx freeVarsCtx) []Variable {
	return []Variable{}
}

func (t Primitive) signature(ctx *signatureContext) string {
	return t.ID
}

func (t Primitive) WithContext(ctx UnificationCtx) Contextual {
	c := t.Copy(NewTypeCopyCtx()).(Primitive)
	c.ctx = ctx
	return c
}

func (t Primitive) GetContext() UnificationCtx {
	return t.ctx
}

type Function struct {
	Fields []Type

	ctx UnificationCtx
}

func NewFunction(ts ...Type) Function {
	return Function{Fields: ts}
}

func (t Function) Copy(ctx *TypeCopyCtx) Type {
	ts := make([]Type, len(t.Fields))
	for i, f := range t.Fields {
		ts[i] = f.Copy(ctx)
	}
	nf := NewFunction(ts...)
	nf.ctx = t.ctx
	return nf
}

func (t Function) unifier(o Type, ctx UnificationCtx, sctx substitutionCtx) (Substitutions, error) {
	if fun, ok := o.(Function); ok && len(fun.Fields) == len(t.Fields) {
		result := MakeSubstitutions()
		for i, f := range t.Fields {
			s, err := unifier(f, fun.Fields[i], ctx, sctx)
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

func (t Function) convert(s Substitutions, sctx substitutionCtx) (Type, bool) {
	changed := false
	ts := make([]Type, len(t.Fields))
	for i, f := range t.Fields {
		nt, c := f.convert(s, sctx)
		ts[i] = nt
		changed = changed || c
	}
	return NewFunction(ts...), changed
}

func (t Function) freeVars(ctx freeVarsCtx) []Variable {
	res := []Variable{}
	for _, f := range t.Fields {
		res = append(res, f.freeVars(ctx)...)
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

func (t Function) Params() []Type {
	return t.Fields[:((len(t.Fields)) - 1)]
}

func (t Function) Return() Type {
	return t.Fields[(len(t.Fields))-1]
}

func (t Function) WithContext(ctx UnificationCtx) Contextual {
	c := t.Copy(NewTypeCopyCtx()).(Function)
	c.ctx = ctx
	return c
}

func (t Function) GetContext() UnificationCtx {
	return t.ctx
}

type Named struct {
	Name string
	Type Type

	ctx UnificationCtx
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
		ctx:  t.ctx,
	}
}

func (t Named) unifier(o Type, ctx UnificationCtx, sctx substitutionCtx) (Substitutions, error) {
	if n, ok := o.(Named); ok && n.Name != t.Name {
		return MakeSubstitutions(), UnificationError(t, o)
	}
	return unifier(t.Type, o, ctx, sctx)
}

func (t Named) convert(s Substitutions, sctx substitutionCtx) (Type, bool) {
	nt, changed := t.Type.convert(s, sctx)
	return Named{
		Name: t.Name,
		Type: nt,
	}, changed
}

func (t Named) freeVars(ctx freeVarsCtx) []Variable {
	return t.Type.freeVars(ctx)
}

func (t Named) signature(ctx *signatureContext) string {
	return t.Name
}

func (t Named) WithContext(ctx UnificationCtx) Contextual {
	c := t.Copy(NewTypeCopyCtx()).(Named)
	c.ctx = ctx
	return c
}

func (t Named) GetContext() UnificationCtx {
	return t.ctx
}

type Structure struct {
	Fields []Named

	ctx UnificationCtx
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
	ns := NewStructure(ts...)
	ns.ctx = t.ctx
	return ns
}

func (t Structure) unifier(o Type, ctx UnificationCtx, sctx substitutionCtx) (Substitutions, error) {
	if s, ok := o.(Structure); ok && len(s.Fields) == len(t.Fields) {
		resmap := map[string]Type{}
		for _, f := range t.Fields {
			resmap[f.Name] = f.Type
		}
		for _, f := range s.Fields {
			p := resmap[f.Name]
			if p != nil {
				_, err := unifier(f.Type, p, ctx, sctx)
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

func (t Structure) convert(s Substitutions, sctx substitutionCtx) (Type, bool) {
	ts := make([]Named, len(t.Fields))
	changed := false
	for i, f := range t.Fields {
		n, c := f.Type.convert(s, sctx)
		ts[i] = NewNamed(f.Name, n)
		changed = changed || c
	}
	return NewStructure(ts...), changed
}

func (t Structure) freeVars(ctx freeVarsCtx) []Variable {
	res := []Variable{}
	for _, f := range t.Fields {
		res = append(res, f.Type.freeVars(ctx)...)
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

func (t Structure) WithContext(ctx UnificationCtx) Contextual {
	c := t.Copy(NewTypeCopyCtx()).(Structure)
	c.ctx = ctx
	return c
}

func (t Structure) GetContext() UnificationCtx {
	return t.ctx
}

type Variable struct {
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

func (t Variable) Copy(ctx *TypeCopyCtx) Type {
	c := ctx.vars[t.ID]
	if c != nil {
		if len(t.Fields) != len(c.(Variable).Fields) {
			panic("variable copy mismatch")
		}
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

func (t Variable) unifier(o Type, ctx UnificationCtx, sctx substitutionCtx) (Substitutions, error) {
	if len(t.Fields) == 0 {
		result := MakeSubstitutions()
		if err := result.Update(t.ID, o, ctx); err != nil {
			return MakeSubstitutions(), err
		}
		return result, nil
	}
	if v, ok := o.(Variable); ok {
		if len(v.Fields) == 0 {
			return MakeSubstitutions(), UnificationError(t, o)
		}
		if sctx.unifying[t.ID] {
			return MakeSubstitutions(), nil
		}
		sctx.unifying[t.ID] = true
		result := MakeSubstitutions()
		sum := map[string]Type{}
		for n, f := range t.Fields {
			if of := v.Fields[n]; of != nil {
				u, err := unifier(f, of, ctx, sctx)
				if err != nil {
					return MakeSubstitutions(), err
				}
				err = result.Combine(u, ctx)
				if err != nil {
					return MakeSubstitutions(), err
				}
				sum[n] = u.Apply(f)
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
		if err := result.Update(v.ID, res, ctx); err != nil {
			return MakeSubstitutions(), err
		}
		if err := result.Update(t.ID, res, ctx); err != nil {
			return MakeSubstitutions(), err
		}
		return result, nil
	}
	if v, ok := o.(Structure); ok {
		sFields := map[string]Type{}
		for _, f := range v.Fields {
			sFields[f.Name] = f.Type
		}
		for n, f := range t.Fields {
			if p := sFields[n]; p != nil {
				_, err := unifier(p, f, ctx, sctx)
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
	if f, ok := o.(Function); ok && len(f.freeVars(newFreeVarsCtx())) > 0 {
		stru := NewVariable(NewNamed("%call", f))
		return unifier(stru, t, ctx, sctx)
	}
	result := MakeSubstitutions()
	for name, typ := range t.Fields {
		if !sctx.resolved[t.ID][name] {
			if sctx.resolved[t.ID] == nil {
				sctx.resolved[t.ID] = map[string]bool{}
			}
			sctx.resolved[t.ID][name] = true
			in := ctx.StructuralTypeFor(name, o)
			if in == nil {
				return MakeSubstitutions(), UnificationError(t, o)
			}
			in = in.Copy(NewTypeCopyCtx())
			ftyp := NewFunction(o, typ)
			if err := result.add(in, ftyp, ctx, sctx); err != nil {
				return MakeSubstitutions(), err
			}
		}
	}
	result.AddContext(t.ID, ctx)
	result.Update(t.ID, o, ctx)
	return result, nil
}

func (t Variable) convert(s Substitutions, sctx substitutionCtx) (Type, bool) {
	if r := (*s.substitutions)[t.ID]; r != nil {
		if v, ok := r.(Variable); ok {
			if v.ID == t.ID {
				return v, false
			}
		}
		if sctx.converted[t.ID] {
			return r, true
		}
		sctx.converted[t.ID] = true
		// Deals with recursive variables
		c, _ := r.convert(s, sctx)
		if con, ok := c.(Contextual); ok {
			context := s.GetContext(t.ID)
			if context != nil && con.GetContext() == nil {
				c = con.WithContext(context).(Type)
			}
		}
		return c, true
	}
	if len(t.Fields) == 0 {
		return t, false
	}
	changed := false
	res := map[string]Type{}
	if !sctx.converted[t.ID] {
		sctx.converted[t.ID] = true
		for n, f := range t.Fields {
			nv, c := f.convert(s, sctx)
			res[n] = nv
			changed = changed || c
		}
	}
	if !changed {
		return t, false
	}
	rest := Variable{
		Fields: res,
		ID:     t.ID,
	}
	return rest, true
}

func (t Variable) freeVars(ctx freeVarsCtx) []Variable {
	if ctx[t.ID] {
		return []Variable{}
	}
	ctx[t.ID] = true
	res := []Variable{t}
	for _, f := range t.Fields {
		res = append(res, f.freeVars(ctx)...)
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
		// Print out the fields only on the first occurrence of the variable
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
	}
	return sb.String()
}

func HasFreeVars(t Type) bool {
	return len(t.freeVars(newFreeVarsCtx())) > 0
}

func IsFunction(t Type) bool {
	_, isFun := t.(Function)
	return isFun
}

func IsStructure(t Type) bool {
	_, isStructure := t.(Structure)
	return isStructure
}

func IsString(t Type) bool {
	p, isPrim := t.(Primitive)
	return isPrim && p.ID == "string"
}
