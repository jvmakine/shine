package grammar

type Expression struct {
	Exp  *UTExpression `@@`
	Type *TypeDef      `(":" @@)?`
}

type UTExpression struct {
	Def  *FDefinition    `@@`
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
	Def   *Definitions `@@`
	Value *Expression  `@@ Newline*`
}

type Definition struct {
	Assignment *Assignment `@@`
	Binding    *Binding    `| @@`
}

type Definitions struct {
	Defs []*Definition `@@*`
}

type Binding struct {
	Name      *TypedName   `Newline* @@`
	Interface *Definitions `"~>" "{" Newline* @@ Newline* "}" Newline*`
}

type Assignment struct {
	Name  *TypedName  `Newline* @@`
	Value *Expression `"=" @@ Newline+`
}

type CallParams struct {
	Params []*Expression `"(" Newline* (@@ Newline* ("," Newline* @@)*)? Newline* ")"`
}

type FValue struct {
	Value *PValue       `@@`
	Calls []*CallParams `@@*`
}

type PValue struct {
	Int    *int64      `@Int`
	Real   *float64    `| @Real`
	Bool   *string     `| @("true" | "false")`
	String *string     `| @String`
	Id     *string     `| @Ident`
	Block  *Block      `| "{" Newline* @@ Newline* "}"`
	Sub    *Expression `| "(" @@ ")"`
}

type TypeFunc struct {
	Params []*TypeDef `"(" (@@ ("," @@)*)? ")" "=>"`
	Return *TypeDef   `@@`
}

type TypeDef struct {
	Primitive string    `@PrimitiveType`
	Function  *TypeFunc `| @@`
	Named     string    `| @Ident`
}

type TypedName struct {
	Name *string  `@Ident`
	Type *TypeDef `(":" @@)?`
}

type FDefinition struct {
	Params []*TypedName `"(" Newline* (@@ Newline* ("," Newline* @@)*)? ")"`
	Funct  *FunctionDef `(Newline* @@)?`
}

type FunctionDef struct {
	ReturnType *TypeDef    `(":" @@)?`
	Body       *Expression `"=>" Newline* @@`
}
