package ast

import "github.com/jvmakine/shine/types"

func (a *Exp) Copy() *Exp {
	return a.CopyWithCtx(types.NewTypeCopyCtx())
}

func (a *Exp) CopyWithCtx(ctx *types.TypeCopyCtx) *Exp {
	if a == nil {
		return nil
	}
	return &Exp{
		Const:   a.Const,
		Block:   a.Block.copy(ctx),
		Id:      a.Id.copy(ctx),
		Op:      a.Op.copy(ctx),
		Call:    a.Call.copy(ctx),
		Def:     a.Def.copy(ctx),
		TDecl:   a.TDecl.copy(ctx),
		FAccess: a.FAccess.copy(ctx),
		Struct:  a.Struct.copy(ctx),
	}
}

func (a *Block) copy(ctx *types.TypeCopyCtx) *Block {
	if a == nil {
		return nil
	}
	ac := map[string]*Exp{}
	for k, v := range a.Assignments {
		ac[k] = v.CopyWithCtx(ctx)
	}
	return &Block{
		Assignments: ac,
		Value:       a.Value.CopyWithCtx(ctx),
		ID:          a.ID,
	}
}

func (a *FieldAccessor) copy(ctx *types.TypeCopyCtx) *FieldAccessor {
	if a == nil {
		return nil
	}
	return &FieldAccessor{
		Exp:   a.Exp.CopyWithCtx(ctx),
		Field: a.Field,
		Type:  a.Type.Copy(ctx),
	}
}

func (a *Struct) copy(ctx *types.TypeCopyCtx) *Struct {
	if a == nil {
		return nil
	}
	fs := make([]*StructField, len(a.Fields))
	for i, f := range a.Fields {
		fs[i] = &StructField{
			Name: f.Name,
			Type: f.Type.Copy(ctx),
		}
	}
	return &Struct{
		Fields: fs,
		Type:   a.Type.Copy(ctx),
	}
}

func (a *FCall) copy(ctx *types.TypeCopyCtx) *FCall {
	if a == nil {
		return nil
	}
	pc := make([]*Exp, len(a.Params))
	for i, p := range a.Params {
		pc[i] = p.CopyWithCtx(ctx)
	}
	return &FCall{
		Function: a.Function.CopyWithCtx(ctx),
		Params:   pc,
		Type:     a.Type.Copy(ctx),
	}
}

func (a *FDef) copy(ctx *types.TypeCopyCtx) *FDef {
	if a == nil {
		return nil
	}
	pc := make([]*FParam, len(a.Params))
	for i, p := range a.Params {
		pc[i] = p.copy(ctx)
	}
	return &FDef{
		Params:  pc,
		Body:    a.Body.CopyWithCtx(ctx),
		Closure: a.Closure.Copy(ctx),
	}
}

func (a *FParam) copy(ctx *types.TypeCopyCtx) *FParam {
	if a == nil {
		return nil
	}
	return &FParam{
		Name: a.Name,
		Type: a.Type.Copy(ctx),
	}
}

func (a *Id) copy(ctx *types.TypeCopyCtx) *Id {
	if a == nil {
		return nil
	}
	return &Id{
		Name: a.Name,
		Type: a.Type.Copy(ctx),
	}
}

func (a *Op) copy(ctx *types.TypeCopyCtx) *Op {
	if a == nil {
		return nil
	}
	return &Op{
		Name: a.Name,
		Type: a.Type.Copy(ctx),
	}
}

func (a *TypeDecl) copy(ctx *types.TypeCopyCtx) *TypeDecl {
	if a == nil {
		return nil
	}
	return &TypeDecl{
		Exp:  a.Exp.CopyWithCtx(ctx),
		Type: a.Type.Copy(ctx),
	}
}
