// Package resolved contains an AST for program where
// all the types have been resolved, functions converted to global
// names and different type instantations of functions converted to
// different functions
package resolved

import (
	"github.com/jvmakine/shine/typedef"
)

type FunctionId = string

type FunctionRef struct {
	Name      string
	BlockId   string
	Signature string
}

func (f *FunctionRef) Id() FunctionId {
	return f.Name + "%%" + f.BlockId + "%%" + f.Signature
}

type Exp struct {
	Type *typedef.Type

	Call  *Call
	Const *Const
	Block *Block
}

type Const struct {
	Bool *bool
	Int  *int
}

type Call struct {
	Id     FunctionId
	Params []*Exp
}

type Assignment struct {
	Name  string
	Value *Exp
}

type Block struct {
	AList []*Assignment
	Value *Exp
}

type FParam struct {
	Type *typedef.Type
	Name string
}

type FDef struct {
	Args []*typedef.Type
	Body *Exp
}

type Resolved struct {
	Functions map[FunctionId]*FDef
	Body      *Exp
}
