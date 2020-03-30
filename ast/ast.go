package ast

// Expressions

type Exp struct {
	Value *Val
	Call  *FCall
}

type Val struct {
	Int *int
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
	Functions   []*FDef
	Value       *Exp
}
