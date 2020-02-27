package grammar

import (
	"github.com/alecthomas/participle"
)

type Program struct {
	Functions []*Function `@@*`
	Exp       *Expression `@@`
}

type Function struct {
	Name   *string     `"fun" @Ident`
	Params []*string   `"(" (@Ident ("," @Ident)*)? ")"`
	Body   *Expression `"{" @@ "}"`
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
