// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

// Expressions

type Exp struct {
	Const *Const
	Block *Block
	Id    *string
	Call  *FCall
	Def   *FDef
	// Type may contain data related to the type inference
	// The actual structure of the type is up to the inference algorithm
	Type interface{}
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
	Type interface{}
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
