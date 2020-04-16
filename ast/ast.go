// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

import "github.com/jvmakine/shine/types"

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
	Int  *int
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
	if a == nil {
		return nil
	}
	return &Exp{
		Const: a.Const,
		Block: a.Block.Copy(),
		Id:    a.Id,
		Call:  a.Call.Copy(),
		Def:   a.Def.Copy(),
		Type:  a.Type.Copy(),
	}
}

func (a *Block) Copy() *Block {
	ac := make([]*Assign, len(a.Assignments))
	for i, as := range a.Assignments {
		ac[i] = as.Copy()
	}
	return &Block{
		Assignments: ac,
		Value:       a.Value.Copy(),
	}
}

func (a *Assign) Copy() *Assign {
	return &Assign{
		Name:  a.Name,
		Value: a.Value.Copy(),
	}
}

func (a *FCall) Copy() *FCall {
	pc := make([]*Exp, len(a.Params))
	for i, p := range a.Params {
		pc[i] = p.Copy()
	}
	return &FCall{
		Name:     a.Name,
		Params:   pc,
		Resolved: a.Resolved,
	}
}

func (a *FDef) Copy() *FDef {
	pc := make([]*FParam, len(a.Params))
	for i, p := range a.Params {
		pc[i] = p.Copy()
	}
	return &FDef{
		Params: pc,
		Body:   a.Body.Copy(),
	}
}

func (a *FParam) Copy() *FParam {
	return &FParam{
		Name: a.Name,
		Type: a.Type.Copy(),
	}
}
