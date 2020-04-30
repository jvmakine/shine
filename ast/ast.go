// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

import (
	"errors"

	"github.com/jvmakine/shine/resolved"
	"github.com/jvmakine/shine/types"
)

// Expressions

type Exp struct {
	Const *Const
	Block *Block
	Id    *Id
	Call  *FCall
	Def   *FDef

	// for function value expressions. reference to the compiled function
	Resolved *resolved.ResolvedFn
}

type Id struct {
	Name string

	Type types.Type
}

type Const struct {
	Int  *int64
	Real *float64
	Bool *bool

	Type types.Type
}

// Functions

type FCall struct {
	Name   string
	Params []*Exp

	Type types.Type
}

type FParam struct {
	Name string

	Type types.Type
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
	return a.copy(types.NewTypeCopyCtx())
}

func (a *Exp) copy(ctx *types.TypeCopyCtx) *Exp {
	if a == nil {
		return nil
	}
	return &Exp{
		Const:    a.Const,
		Block:    a.Block.copy(ctx),
		Id:       a.Id.copy(ctx),
		Call:     a.Call.copy(ctx),
		Def:      a.Def.copy(ctx),
		Resolved: a.Resolved,
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
		Name:   a.Name,
		Params: pc,
		Type:   a.Type.Copy(ctx),
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

func (a *Id) copy(ctx *types.TypeCopyCtx) *Id {
	if a == nil {
		return nil
	}
	return &Id{
		Name: a.Name,
		Type: a.Type.Copy(ctx),
	}
}

func (b *Block) CheckValueCycles() error {
	names := map[string]*Exp{}
	verified := map[string]bool{}
	for _, a := range b.Assignments {
		names[a.Name] = a.Value
	}
	for _, a := range b.Assignments {
		if !verified[a.Name] {
			todo := a.Value.collectIds()
			visited := []string{a.Name}
			visitedb := map[string]bool{a.Name: true}
			verified[a.Name] = true
			for len(todo) > 0 {
				i := todo[0]
				todo = todo[1:]
				if names[i] != nil && names[i].Def == nil {
					if visitedb[i] {
						return errors.New("recursive value: " + cycleToStr(visited, i))
					}
					verified[i] = true
					visitedb[i] = true
					visited = append(visited, i)
					todo = append(todo, names[i].collectIds()...)
				}
			}
		}
	}
	return nil
}

func (exp *Exp) collectIds() []string {
	if exp.Id != nil {
		return []string{exp.Id.Name}
	} else if exp.Call != nil {
		result := []string{}
		for _, p := range exp.Call.Params {
			result = append(result, p.collectIds()...)
		}
		return result
	} else if exp.Def != nil {
		return exp.Def.Body.collectIds()
	} else if exp.Block != nil {
		result := exp.Block.Value.collectIds()
		for _, a := range exp.Block.Assignments {
			result = append(result, a.Value.collectIds()...)
		}
		return result
	}
	return []string{}
}

func cycleToStr(arr []string, v string) string {
	res := ""
	for _, a := range arr {
		res = res + a + " -> "
	}
	return res + v
}

func (exp *Exp) Type() types.Type {
	if exp.Block != nil {
		return exp.Block.Value.Type()
	} else if exp.Call != nil {
		return exp.Call.Type
	} else if exp.Const != nil {
		return exp.Const.Type
	} else if exp.Def != nil {
		ts := make([]types.Type, len(exp.Def.Params)+1)
		for i, p := range exp.Def.Params {
			ts[i] = p.Type
		}
		ts[len(exp.Def.Params)] = exp.Def.Body.Type()
		return types.MakeFunction(ts...)
	} else if exp.Id != nil {
		return exp.Id.Type
	}
	panic("invalid exp")
}
