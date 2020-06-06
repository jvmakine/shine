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
	return MakeStructure(s.Name, ps...).Structure
}

func (t *TypeVar) Copy(ctx *TypeCopyCtx) *TypeVar {
	return &TypeVar{
		Union: t.Union,
	}
}
