package ast

import "github.com/chewxy/hm"

// Expressions

type Exp struct {
	Const *Const
	Block *Block
	Id    *string
	Call  *FCall
	Def   *FDef
	Type  hm.Type
}

type Const struct {
	Int  *int
	Bool *bool
}

// Functions

type FCall struct {
	Name   string
	Params []*Exp
}

type FParam struct {
	Name string
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
