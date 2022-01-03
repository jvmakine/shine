package types

type TypeCopyCtx struct {
	vars map[*TypeVar]*TypeVar
}

func NewTypeCopyCtx() *TypeCopyCtx {
	return &TypeCopyCtx{
		vars: map[*TypeVar]*TypeVar{},
	}

}

func (t Type) Copy(ctx *TypeCopyCtx) Type {
	if t.Variable != nil {
		if ctx.vars[t.Variable] == nil {
			ctx.vars[t.Variable] = t.Variable.Copy(ctx)
		}
		return Type{Variable: ctx.vars[t.Variable]}
	}
	if t.Function != nil {
		ps := make([]Type, len(*t.Function))
		for i, p := range *t.Function {
			ps[i] = p.Copy(ctx)
		}
		return MakeFunction(ps...)
	}
	if t.Structure != nil {
		return Type{Structure: t.Structure.Copy(ctx)}
	}
	if t.TCRef != nil {
		ps := make([]Type, len(*&t.TCRef.TypeClassVars))
		for i, p := range t.TCRef.TypeClassVars {
			ps[i] = p.Copy(ctx)
		}
		c := MakeTypeClassRef(t.TCRef.TypeClass, t.TCRef.Place, ps...)
		c.TCRef.LocalBindings = t.TCRef.LocalBindings
		return c
	}
	if t.HVariable != nil {
		ps := make([]Type, len(t.HVariable.Params))
		for i, p := range t.HVariable.Params {
			ps[i] = p.Copy(ctx)
		}
		r := t.HVariable.Root
		if ctx.vars[r] == nil {
			ctx.vars[r] = r.Copy(ctx)
		}
		return MakeHierarchicalVar(ctx.vars[r], ps...)
	}
	return t
}

func (s *Structure) Copy(ctx *TypeCopyCtx) *Structure {
	if s == nil {
		return nil
	}
	ps := make([]SField, len(s.Fields))
	for i, p := range s.Fields {
		ps[i] = SField{
			Name: p.Name,
			Type: p.Type.Copy(ctx),
		}
	}
	nt := make([]Type, len(s.TypeArguments))
	for i, a := range s.TypeArguments {
		nt[i] = a.Copy(ctx)
	}
	return &Structure{
		Name:           s.Name,
		TypeArguments:  nt,
		Fields:         ps,
		OrginalVars:    s.OrginalVars,
		OriginalFields: s.OriginalFields,
	}
}

func (t *TypeVar) Copy(ctx *TypeCopyCtx) *TypeVar {
	structural := map[string]Type{}
	for k, v := range t.Structural {
		structural[k] = v.Copy(ctx)
	}
	return &TypeVar{
		Union:      t.Union,
		Structural: structural,
	}
}
