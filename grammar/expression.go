package grammar

type Expression struct {
	Exp  *UTExpression    `@@`
	Type *TypeDeclaration `(":" @@)?`
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
	Elements []*BlockElement `Newline* @@*`
	Value    *Expression     `@@ Newline*`
}

type BlockElement struct {
	TypeDef     *TypeDefinition `@@`
	TypeBinding *TypeBinding    `| @@`
	Assignment  *Assignment     `| @@`
}

type TypeBinding struct {
	Name        *string            `@Ident`
	Arguments   []*TypeDeclaration `"[" @@ ("," @@)* "]" Binding Newline*`
	Assignments []*Assignment      `"{" Newline* @@+ "}" Newline*`
}

type TypeDefinition struct {
	Name      *string              `@Ident`
	FreeVars  []string             `("[" (@Ident ("," @Ident)*)? "]")? TypeDef`
	Struct    *StructDescription   `(@@ Newline+`
	TypeClass *TypeClassDefinition ` | @@ Newline+`
	Type      *TypeDeclaration     ` | @@ Newline+)`
}

type TypeClassDefinition struct {
	Functions []*AbstractFunDef `"{" Newline* @@+ "}"`
}

type AbstractFunDef struct {
	Name     *string   `@Ident`
	FreeVars []string  `("[" @Ident ("," @Ident)* "]")? TypeDef`
	Function *TypeFunc `@@ Newline*`
}

type StructDescription struct {
	Fields []*StructFiels `"(" Newline* (@@ Newline* ("," Newline* @@)*)? ")"`
}

type StructFiels struct {
	Name string           `@Ident`
	Type *TypeDeclaration `":" @@`
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

type Assignment struct {
	Name  *string          `@Ident`
	Type  *TypeDeclaration `(":" @@)?`
	Value *Expression      `"=" @@ Newline*`
}

type TypeFunc struct {
	Params []*TypeDeclaration `"(" (@@ ("," @@)*)? ")" "=>"`
	Return *TypeDeclaration   `@@`
}

type TypeDeclaration struct {
	Primitive string     `@PrimitiveType`
	Function  *TypeFunc  `| @@`
	Named     *TypeNamed `| @@`
}

type TypeNamed struct {
	Name string             `@Ident`
	Vars []*TypeDeclaration `( "[" @@ ("," @@)* "]" )?`
}

type FunParam struct {
	Name *string          `@Ident`
	Type *TypeDeclaration `(":" @@)?`
}

type Definition struct {
	Params []*FunParam  `"(" Newline* (@@ Newline* ("," Newline* @@)*)? ")"`
	Funct  *FunctionDef `Newline* @@`
}

type FunctionDef struct {
	ReturnType *TypeDeclaration `(":" @@)?`
	Body       *Expression      `"=>" Newline* @@`
}
