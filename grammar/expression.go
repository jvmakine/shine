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
	Value       *Expression   `@@ Newline*`
}

type EValue struct {
	Id  *string     `@Ident`
	Sub *Expression `| "(" @@ ")"`
}

type Value struct {
	Call  *FunCall `@@`
	Int   *int64   `| @Int`
	Real  *float64 `| @Real`
	Bool  *string  `| @("true" | "false")`
	Block *Block   `| "{" Newline* @@ Newline* "}"`
	Eval  *EValue  `| @@`
}

type Assignment struct {
	Name  *string     `Newline* @Ident "="`
	Value *Expression `@@ Newline+`
}

type FunParam struct {
	Name *string `@Ident`
}

type FunDef struct {
	Params []*FunParam `"(" (@@ ("," @@)*)? ")" "=>"`
	Body   *Block      `"{" Newline* @@ Newline* "}"`
}

type CallParams struct {
	Params []*Expression `"(" ( Newline* @@ ("," Newline* @@)*)? Newline* ")"`
}

type FunCall struct {
	Function *EValue       `@@`
	Calls    []*CallParams `@@+`
}
