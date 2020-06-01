package grammar

type Expression struct {
	Exp  *UTExpression `@@`
	Type *TypeDef      `(":" @@)?`
}

type UTExpression struct {
	Def  *Definition     `@@`
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
	Left  *Accessor   `@@`
	Right []*OpFactor `@@*`
}

type OpFactor struct {
	Operation *string   `@("*" | "/" | "%")`
	Right     *Accessor `@@`
}

type NamedFValue struct {
	Id    string        `@Ident`
	Calls []*CallParams `@@*`
}

type Accessor struct {
	Left  *FValue        `@@`
	Right []*NamedFValue `("." @@)*`
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

type TypeFunc struct {
	Params []*TypeDef `"(" (@@ ("," @@)*)? ")" "=>"`
	Return *TypeDef   `@@`
}

type TypeDef struct {
	Primitive string    `@PrimitiveType`
	Function  *TypeFunc `| @@`
}

type FunParam struct {
	Name *string  `@Ident`
	Type *TypeDef `(":" @@)?`
}

type Definition struct {
	Params []*FunParam  `"(" Newline* (@@ Newline* ("," Newline* @@)*)? ")"`
	Funct  *FunctionDef `(Newline* @@)?`
}

type FunctionDef struct {
	ReturnType *TypeDef    `(":" @@)?`
	Body       *Expression `"=>" Newline* @@`
}
