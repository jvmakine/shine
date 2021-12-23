// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

import (
	"errors"

	"github.com/jvmakine/shine/types"
)

// Expressions

type Exp struct {
	Const   *Const         // Constant value
	Block   *Block         // Block with assignments and a body
	Id      *Id            // Id referring to a value or parameter defined elsewhere
	Op      *Op            // Operator from a set of predefined operations like +, *, etc
	Call    *FCall         // Call of a function
	Def     *FDef          // Definition of a function
	TDecl   *TypeDecl      // Manually defined type for an expression
	FAccess *FieldAccessor // Accessing a field / method of a structure
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
	Int    *int64
	Real   *float64
	Bool   *bool
	String *string

	Type types.Type
}

type TypeDecl struct {
	Exp  *Exp
	Type types.Type
}

type FieldAccessor struct {
	Exp   *Exp
	Field string
	Type  types.Type
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

	Closure *types.Structure
}

// Blocks

type Block struct {
	Assignments map[string]*Exp
	TypeDefs    map[string]*TypeDefinition
	Value       *Exp

	ID int
}

// Types

type TypeDefinition struct {
	FreeVariables []string
	Struct        *Struct
}

func (t *TypeDefinition) WithFreeVars(vars ...string) *TypeDefinition {
	t.FreeVariables = vars
	return t
}

type StructField struct {
	Name string
	Type types.Type
}

type Struct struct {
	Fields []*StructField
	Type   types.Type
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
	return a.Closure != nil && len(a.Closure.Fields) > 0
}

func (a *FCall) RootFunc() *Exp {
	if a.Function.Call != nil {
		return a.Function.Call.RootFunc()
	}
	return a.Function
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
	exp.Visit(func(v *Exp, c *VisitContext) error {
		if v.Id != nil {
			name := v.Id.Name
			ids[name] = true
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
	} else if exp.FAccess != nil {
		return exp.FAccess.Type
	}
	panic("invalid exp")
}

func (t *TypeDefinition) Type() types.Type {
	if t.Struct != nil {
		return t.Struct.Type
	}
	panic("invalid type def")
}

func (exp *Exp) Convert(s types.Substitutions) {
	exp.RewriteTypes(func(t types.Type, ctx *VisitContext) (types.Type, error) {
		return s.Apply(t), nil
	})
}

func (t *TypeDefinition) Convert(s types.Substitutions) {
	if t.Struct != nil {
		t.Struct.Type = s.Apply(t.Struct.Type)
		for _, f := range t.Struct.Fields {
			f.Type = s.Apply(f.Type)
		}
	}
}
