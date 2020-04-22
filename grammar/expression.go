package grammar

type Expression struct {
	Fun  *FunDef         `@@`
	If   *IfExpression   `| @@`
	Comp *CompExpression `| @@`
}

type IfExpression struct {
	Cond  *Expression `"if" "(" @@ ")"`
	True  *Expression `@@`
	False *Expression `"else" @@`
}

type CompExpression struct {
	Left  *Comp     `@@`
	Right []*OpComp `@@*`
}

type OpComp struct {
	Operation *string `@( "||" | "&&" )`
	Right     *Comp   `@@*`
}

type Comp struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type TermExpression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type OpTerm struct {
	Operation *string `@("+" | "-" | ">" | "<" | "<=" | ">=" | "==" )`
	Right     *Term   `@@*`
}

type Term struct {
	Left  *Value      `@@`
	Right []*OpFactor `@@*`
}

type OpFactor struct {
	Operation *string `@("*" | "/" | "%")`
	Right     *Value  `@@`
}

type Block struct {
	Assignments []*Assignment `@@*`
	Value       *Expression   `@@`
}

type Value struct {
	Int   *int64      `@Int`
	Real  *float64    `| @Real`
	Bool  *string     `| @("true" | "false")`
	Call  *FunCall    `| @@`
	Id    *string     `| @Ident`
	Block *Block      `| "{" @@ "}"`
	Sub   *Expression `| "(" @@ ")"`
}

type Assignment struct {
	Name  *string     `@Ident "="`
	Value *Expression `@@`
}

type FunParam struct {
	Name *string `@Ident`
}

type FunDef struct {
	Params []*FunParam `"(" (@@ ("," @@)*)? ")" "=>"`
	Body   *Block      `"{" @@ "}"`
}

type FunCall struct {
	Name   *string       `@Ident`
	Params []*Expression `"(" (@@ ("," @@)*)? ")"`
}
