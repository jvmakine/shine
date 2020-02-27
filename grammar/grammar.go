package grammar

import (
	"github.com/alecthomas/participle"
)

type Program struct {
	Exp *Expression `@@`
}

type Value struct {
	Int *int        `@Int`
	Sub *Expression `| "(" @@ ")"`
}

type Expression struct {
	Value *Value     `@@`
	Op    *Operation `@@?`
}

type Operation struct {
	Add *Expression `"+" @@`
	Sub *Expression `| "-" @@`
	Mul *Expression `| "*" @@`
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
