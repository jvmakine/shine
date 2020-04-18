// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

import (
	"github.com/jvmakine/shine/types"
)

// Expressions

type Exp struct {
	Const *Const
	Block *Block
	Id    *string
	Call  *FCall
	Def   *FDef
	// Type may contain data related to the type inference
	Type *types.TypePtr
}

type Const struct {
	Int  *int64
	Real *float64
	Bool *bool
}

// Functions

type FCall struct {
	Name   string
	Params []*Exp

	Resolved string
}

type FParam struct {
	Name string
	Type *types.TypePtr
}

type FDef struct {
	Params []*FParam
	Body   *Exp

	Resolved string
}

// Blocks

type Assign struct {
	Name  string
	Value *Exp
}

type Block struct {
	Assignments []*Assign
	Value       *Exp
}

func (a *Exp) Copy() *Exp {
	return a.copy(types.NewTypeCopyCtx())
}

func (a *Exp) copy(ctx *types.TypeCopyCtx) *Exp {
	if a == nil {
		return nil
	}
	return &Exp{
		Const: a.Const,
		Block: a.Block.copy(ctx),
		Id:    a.Id,
		Call:  a.Call.copy(ctx),
		Def:   a.Def.copy(ctx),
		Type:  a.Type.Copy(ctx),
	}
}

func (a *Block) copy(ctx *types.TypeCopyCtx) *Block {
	if a == nil {
		return nil
	}
	ac := make([]*Assign, len(a.Assignments))
	for i, as := range a.Assignments {
		ac[i] = as.copy(ctx)
	}
	return &Block{
		Assignments: ac,
		Value:       a.Value.copy(ctx),
	}
}

func (a *Assign) copy(ctx *types.TypeCopyCtx) *Assign {
	if a == nil {
		return nil
	}
	return &Assign{
		Name:  a.Name,
		Value: a.Value.copy(ctx),
	}
}

func (a *FCall) copy(ctx *types.TypeCopyCtx) *FCall {
	if a == nil {
		return nil
	}
	pc := make([]*Exp, len(a.Params))
	for i, p := range a.Params {
		pc[i] = p.copy(ctx)
	}
	return &FCall{
		Name:     a.Name,
		Params:   pc,
		Resolved: a.Resolved,
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
		Params: pc,
		Body:   a.Body.copy(ctx),
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
