// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

import (
	"errors"

	"github.com/jvmakine/shine/types"
)

// Expressions

type Exp struct {
	Const *Const
	Block *Block
	Id    *Id
	Op    *Op
	Call  *FCall
	Def   *FDef
}

type Op struct {
	Name string

	Type types.Type
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
	Function *Exp
	Params   []*Exp

	Type types.Type
}

type FParam struct {
	Name string

	Type types.Type
}

type FDef struct {
	Params []*FParam
	Body   *Exp

	Closure *types.Closure
}

// Blocks

type Block struct {
	Assignments map[string]*Exp
	Value       *Exp

	ID int
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
		Id:    a.Id.copy(ctx),
		Op:    a.Op.copy(ctx),
		Call:  a.Call.copy(ctx),
		Def:   a.Def.copy(ctx),
	}
}

func (a *FDef) ParamOf(name string) *FParam {
	for _, p := range a.Params {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (a *FDef) HasClosure() bool {
	return a.Closure != nil && len(*a.Closure) > 0
}

func (a *FCall) RootFunc() *Exp {
	if a.Function.Call != nil {
		return a.Function.Call.RootFunc()
	}
	return a.Function
}

func (a *Block) copy(ctx *types.TypeCopyCtx) *Block {
	if a == nil {
		return nil
	}
	ac := map[string]*Exp{}
	for k, v := range a.Assignments {
		ac[k] = v.copy(ctx)
	}
	return &Block{
		Assignments: ac,
		Value:       a.Value.copy(ctx),
		ID:          a.ID,
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
		Function: a.Function.copy(ctx),
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
		Body:    a.Body.copy(ctx),
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

func (b *Block) CheckValueCycles() error {
	names := map[string]*Exp{}

	type ToDo struct {
		id   string
		path []string
	}
	todo := []ToDo{}

	for k, a := range b.Assignments {
		names[k] = a
		todo = append(todo, ToDo{id: k, path: []string{}})
	}

	for len(todo) > 0 {
		i := todo[0]
		todo = todo[1:]
		for _, p := range i.path {
			if p == i.id {
				return errors.New("recursive value: " + cycleToStr(i.path, i.id))
			}
		}
		exp := b.Assignments[i.id]
		if exp.Def == nil {
			ids := exp.CollectIds()
			for _, id := range ids {
				if names[id] != nil {
					todo = append(todo, ToDo{id: id, path: append(i.path, i.id)})
				}
			}
		}
	}
	return nil
}

func (exp *Exp) CollectIds() []string {
	if exp.Id != nil {
		return []string{exp.Id.Name}
	} else if exp.Call != nil {
		resultM := map[string]bool{}
		for _, p := range exp.Call.Params {
			for _, id := range p.CollectIds() {
				resultM[id] = true
			}
		}
		result := []string{}
		for k := range resultM {
			result = append(result, k)
		}
		return result
	} else if exp.Def != nil {
		return exp.Def.Body.CollectIds()
	} else if exp.Block != nil {
		resultM := map[string]bool{}
		for _, id := range exp.Block.Value.CollectIds() {
			resultM[id] = true
		}
		for _, a := range exp.Block.Assignments {
			for _, id := range a.CollectIds() {
				resultM[id] = true
			}
		}
		result := []string{}
		for k := range resultM {
			result = append(result, k)
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

func (call *FCall) MakeFunType() types.Type {
	funps := make([]types.Type, len(call.Params)+1)
	for i, p := range call.Params {
		funps[i] = p.Type()
	}
	funps[len(call.Params)] = call.Type
	return types.MakeFunction(funps...)
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
	} else if exp.Op != nil {
		return exp.Op.Type
	}
	panic("invalid exp")
}

func (exp *Exp) Convert(s types.Substitutions) {
	if exp.Block != nil {
		exp.Block.Value.Convert(s)
		for _, a := range exp.Block.Assignments {
			a.Convert(s)
		}
	} else if exp.Call != nil {
		exp.Call.Function.Convert(s)
		for _, p := range exp.Call.Params {
			p.Convert(s)
		}
		exp.Call.Type = s.Apply(exp.Call.Type)
	} else if exp.Def != nil {
		for _, p := range exp.Def.Params {
			p.Type = s.Apply(p.Type)
		}
		exp.Def.Body.Convert(s)
	} else if exp.Id != nil {
		exp.Id.Type = s.Apply(exp.Id.Type)
	}
}
