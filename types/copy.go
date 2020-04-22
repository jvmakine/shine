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
	return t
}

func (t *TypeVar) Copy(ctx *TypeCopyCtx) *TypeVar {
	var fun *Function = nil
	if t.Function != nil {
		f := make(Function, len(*t.Function))
		fun = &f
		for i, p := range *t.Function {
			(*fun)[i] = p.Copy(ctx)
		}
	}
	return &TypeVar{
		Restrictions: t.Restrictions,
		Function:     fun,
	}
}