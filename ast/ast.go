package ast

// Expressions

type Exp struct {
	Const *Const
	Block *Block
	Id    *string
	Call  *FCall
	Def   *FDef
}

type Const struct {
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
	Value       *Exp
}
