package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/ebnf"
)

type Program struct {
	Body *Block `@@`
}

func Parse(str string) (*Program, error) {
	lexer, err := ebnf.New(`
		Fun = "=>" .
		Whitespace = " " | "\n" | "\r" | "\t" .
		Comma = "," .
		Eq = "=" .
		Brackets = "(" | ")" | "{" | "}" .
		Op = "+" | "-" | "*" | "/" .
		Ident = alpha { alpha | digit } .
		Int = "1"…"9" { digit } .
		alpha = "a"…"z" | "A"…"Z" | "_" .
		digit = "0"…"9" .
	`)
	parser, err := participle.Build(&Program{}, participle.Lexer(lexer), participle.Elide("Whitespace"))
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
