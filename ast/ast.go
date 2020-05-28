// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

import (
	"errors"

	"github.com/jvmakine/shine/types"
)

// Expressions

type Exp struct {
	Const *Const    // Constant value
	Block *Block    // Block with assignments and a body
	Id    *Id       // Id referring to a value or parameter defined elsewhere
	Op    *Op       // Operator from a set of predefined operations like +, *, etc
	Call  *FCall    // Call of a function
	Def   *FDef     // Definition of a function
	TDecl *TypeDecl // Manually defined type for an expression
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

type TypeDecl struct {
	Exp  *Exp
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
	return a.CopyWithCtx(types.NewTypeCopyCtx())
}

func (a *Exp) CopyWithCtx(ctx *types.TypeCopyCtx) *Exp {
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
		TDecl: a.TDecl.copy(ctx),
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
		ac[k] = v.CopyWithCtx(ctx)
	}
	return &Block{
		Assignments: ac,
		Value:       a.Value.CopyWithCtx(ctx),
		ID:          a.ID,
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
	ids := map[string]bool{}
	exp.Visit(func(v *Exp, _ *VisitContext) error {
		if v.Id != nil {
			ids[v.Id.Name] = true
		}
		return nil
	})
	result := []string{}
	for k := range ids {
		result = append(result, k)
	}
	return result
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
	} else if exp.TDecl != nil {
		return exp.TDecl.Type
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
	} else if exp.TDecl != nil {
		exp.TDecl.Type = s.Apply(exp.TDecl.Type)
		exp.TDecl.Exp.Convert(s)
	}
}
