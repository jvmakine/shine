package types

type TypeCopyCtx struct {
	vars map[VariableID]Type
}

func NewTypeCopyCtx() *TypeCopyCtx {
	return &TypeCopyCtx{
		vars: map[VariableID]Type{},
	}

}
