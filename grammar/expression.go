package grammar

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type Value struct {
	Int *int        `@Int`
	Sub *Expression `| "(" @@ ")"`
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
