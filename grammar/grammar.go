package grammar

import (
	"github.com/alecthomas/participle"
)

type Program struct {
	Functions []*FunDef `@@*`
	Body      *Block    `@@`
}

type Assignment struct {
	Name  *string     `@Ident "="`
	Value *Expression `@@`
}

type Block struct {
	Assignments []*Assignment `@@*`
	Value       *Expression   `@@`
}

type FunParam struct {
	Name *string `@Ident`
}

type FunDef struct {
	Name   *string     `"fun" @Ident`
	Params []*FunParam `"(" (@@ ("," @@)*)? ")"`
	Body   *Block      `"{" @@ "}"`
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
