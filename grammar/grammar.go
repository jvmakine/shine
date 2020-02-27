package grammar

import (
	"github.com/alecthomas/participle"
)

type Program struct {
	Exp *Expression `@@`
}

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

func Parse(str string) (*Program, error) {
	parser, err := participle.Build(&Program{})
	if err != nil {
		panic(err)
	}

	ast := &Program{}
	err = parser.ParseString(str, ast)
	if err != nil {
		return nil, err
	}
	return ast, nil
}
