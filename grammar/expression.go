package grammar

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type Value struct {
	Int  *int        `@Int`
	Call *FunCall    `| @@`
	Id   *string     `| @Ident`
	Sub  *Expression `| "(" @@ ")"`
}

type FunCall struct {
	Name   *string       `@Ident`
	Params []*Expression `"(" (@@ ("," @@)*)? ")"`
}

type OpFactor struct {
	Operation *string `@("*" | "/")`
	Right     *Value  `@@`
}

type Term struct {
	Left  *Value      `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Operation *string `@("+" | "-")`
	Right     *Term   `@@*`
}
