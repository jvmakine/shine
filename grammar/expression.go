package grammar

type Expression struct {
	Fun  *FunDef         `@@`
	If   *IfExpression   `| @@`
	Comp *CompExpression `| @@`
}

type IfExpression struct {
	Cond  *Expression `"if" "(" @@ ")"`
	True  *Expression `@@ Newline*`
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
	Operation *string `@("+" | "-" | ">" | "<" | "<=" | ">=" | "==" | "!=" )`
	Right     *Term   `@@*`
}

type Term struct {
	Left  *FValue     `@@`
	Right []*OpFactor `@@*`
}

type OpFactor struct {
	Operation *string `@("*" | "/" | "%")`
	Right     *FValue `@@`
}

type Block struct {
	Assignments []*Assignment `@@*`
	Value       *Expression   `@@ Newline*`
}

type CallParams struct {
	Params []*Expression `"(" Newline* (@@ Newline* ("," Newline* @@)*)? Newline* ")"`
}

type FValue struct {
	Value *PValue       `@@`
	Calls []*CallParams `@@*`
}

type PValue struct {
	Int   *int64      `@Int`
	Real  *float64    `| @Real`
	Bool  *string     `| @("true" | "false")`
	Id    *string     `| @Ident`
	Block *Block      `| "{" Newline* @@ Newline* "}"`
	Sub   *Expression `| "(" @@ ")"`
}

type Assignment struct {
	Name  *string     `Newline* @Ident "="`
	Value *Expression `@@ Newline+`
}

type TypeDef struct {
	Primitive string `@PrimitiveType`
}

type FunParam struct {
	Name *string  `@Ident`
	Type *TypeDef `(":" @@)?`
}

type FunDef struct {
	Params     []*FunParam `"(" Newline* (@@ Newline* ("," Newline* @@)*)? ")" Newline*`
	ReturnType *TypeDef    `(":" @@)?`
	Body       *Expression `"=>" Newline* @@`
}
