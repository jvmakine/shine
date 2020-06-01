package grammar

import (
	"github.com/jvmakine/shine/types"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/ebnf"
	"github.com/jvmakine/shine/ast"
)

type Program struct {
	Body *Block `@@`
}

func Parse(str string) (*Program, error) {
	lexer, err := ebnf.New(`
		LineComment = ("//") { "\u0000"…"\uffff"-"\n" } .
		BlockComment = ("/*") { "\u0000"…"\uffff"-"*/" } ("*/") .
		Fun = "=>" .
		Newline = "\n" .
		Whitespace = " " | "\r" | "\t" .
		Reserved = "if" | "else" | "true" | "false" .
		Comma = "," .
		Dot = "." .
		Brackets = "(" | ")" | "{" | "}" .
		COp = ">=" | "<=" .
		Op = "+" | "-" | "*" | "/" | "%" |  ">" | "<" | "==" | "!=" | "||" | "&&" .
		Typ = ":" .
		PrimitiveType = "int" | "real" | "bool" .
		Eq = "=" .
		Ident = alpha { alpha | digit } .
		Real = "0"…"9" { digit } "." "0"…"9" { digit } .
		Int = "0" | "1"…"9" { digit } .
		alpha = "a"…"z" | "A"…"Z" | "_" .
		digit = "0"…"9" .
	`)
	parser, err := participle.Build(
		&Program{},
		participle.UseLookahead(2),
		participle.Lexer(lexer),
		participle.Elide("Whitespace", "LineComment", "BlockComment"),
	)
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

func (prg *Program) ToAst() *ast.Exp {
	return &ast.Exp{
		Block: convBlock(prg.Body),
	}
}

func convBlock(from *Block) *ast.Block {
	assigns := map[string]*ast.Exp{}
	for _, a := range from.Assignments {
		assigns[*a.Name] = convExp(a.Value)
	}
	return &ast.Block{
		Assignments: assigns,
		Value:       convExp(from.Value),
	}
}

func convExp(from *Expression) *ast.Exp {
	ut := convUTExp(from.Exp)
	if from.Type != nil {
		return &ast.Exp{TDecl: &ast.TypeDecl{Exp: ut, Type: convTypeDef(from.Type)}}
	}
	return ut
}

func convUTExp(from *UTExpression) *ast.Exp {
	if from.Def != nil {
		return convDef(from.Def)
	} else if from.If != nil {
		return convIf(from.If)
	} else if from.Comp != nil {
		return convOpComp(convComp(from.Comp.Left), from.Comp.Right)
	}
	panic("invalid expression")
}

func convIf(from *IfExpression) *ast.Exp {
	return &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: "if"}},
			Params:   []*ast.Exp{convExp(from.Cond), convExp(from.True), convExp(from.False)},
		},
	}
}

func convDef(from *Definition) *ast.Exp {
	if fd := from.Funct; fd != nil {
		params := make([]*ast.FParam, len(from.Params))
		for i, p := range from.Params {
			params[i] = convFParam(p)
		}
		body := convExp(fd.Body)
		if fd.ReturnType != nil {
			body = &ast.Exp{TDecl: &ast.TypeDecl{Exp: body, Type: convTypeDef(fd.ReturnType)}}
		}
		return &ast.Exp{
			Def: &ast.FDef{
				Params: params,
				Body:   body,
			},
		}
	} else {
		fields := make([]*ast.StructField, len(from.Params))
		for i, p := range from.Params {
			td := types.Type{}
			if p.Type != nil {
				td = convTypeDef(p.Type)
			}
			fields[i] = &ast.StructField{
				Name: *p.Name,
				Type: td,
			}
		}
		return &ast.Exp{
			Struct: &ast.Struct{
				Fields: fields,
			},
		}
	}
}

func convTypeDef(t *TypeDef) types.Type {
	if t.Primitive != "" {
		switch t.Primitive {
		case "int":
			return types.IntP
		case "real":
			return types.RealP
		case "bool":
			return types.BoolP
		default:
			panic("invalid type: " + t.Primitive)
		}
	} else if t.Function != nil {
		ps := make([]types.Type, len(t.Function.Params)+1)
		for i, p := range t.Function.Params {
			ps[i] = convTypeDef(p)
		}
		ps[len(t.Function.Params)] = convTypeDef(t.Function.Return)
		return types.MakeFunction(ps...)
	}
	panic("invalid type")
}

func convFParam(from *FunParam) *ast.FParam {
	typ := types.Type{}
	if from.Type != nil {
		typ = convTypeDef(from.Type)
	}
	return &ast.FParam{
		Name: *from.Name,
		Type: typ,
	}
}

func convOpComp(left *ast.Exp, right []*OpComp) *ast.Exp {
	if right == nil || len(right) == 0 {
		return left
	}
	res := &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: *right[0].Operation}},
			Params:   []*ast.Exp{left, convComp(right[0].Right)},
		},
	}
	return convOpComp(res, right[1:])
}

func convComp(from *Comp) *ast.Exp {
	return convOpTerm(convTerm(from.Left), from.Right)
}

func convOpTerm(left *ast.Exp, right []*OpTerm) *ast.Exp {
	if right == nil || len(right) == 0 {
		return left
	}
	res := &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: *right[0].Operation}},
			Params:   []*ast.Exp{left, convTerm(right[0].Right)},
		},
	}
	return convOpTerm(res, right[1:])
}

func convTerm(from *Term) *ast.Exp {
	return convOpFact(convFVal(from.Left), from.Right)
}

func convOpFact(left *ast.Exp, right []*OpFactor) *ast.Exp {
	if right == nil || len(right) == 0 {
		return left
	}
	res := &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: *right[0].Operation}},
			Params:   []*ast.Exp{left, convFVal(right[0].Right)},
		},
	}
	return convOpFact(res, right[1:])
}

func convFVal(from *FValue) *ast.Exp {
	pval := convPVal(from.Value)
	accessors := from.Access
	for len(accessors) > 0 {
		t := accessors[0]
		accessors = accessors[1:]
		pval = &ast.Exp{
			FAccess: &ast.FieldAccessor{
				Exp:   pval,
				Field: t,
			},
		}
	}
	if len(from.Calls) > 0 {
		call, calls := from.Calls[0], from.Calls[1:]
		params := make([]*ast.Exp, len(call.Params))
		for i, p := range call.Params {
			params[i] = convExp(p)
		}
		res := &ast.FCall{
			Function: pval,
			Params:   params,
		}
		for len(calls) > 0 {
			call, calls = calls[0], calls[1:]
			params := make([]*ast.Exp, len(call.Params))
			for i, p := range call.Params {
				params[i] = convExp(p)
			}
			res = &ast.FCall{
				Function: &ast.Exp{Call: res},
				Params:   params,
			}
		}
		return &ast.Exp{
			Call: res,
		}
	}
	return pval
}

func convPVal(from *PValue) *ast.Exp {
	if from.Block != nil {
		return &ast.Exp{
			Block: convBlock(from.Block),
		}
	} else if from.Int != nil {
		return &ast.Exp{
			Const: &ast.Const{Int: from.Int},
		}
	} else if from.Real != nil {
		return &ast.Exp{
			Const: &ast.Const{Real: from.Real},
		}
	} else if from.Bool != nil {
		value := false
		if *from.Bool == "true" {
			value = true
		}
		return &ast.Exp{
			Const: &ast.Const{Bool: &value},
		}
	} else if from.Sub != nil {
		return convExp(from.Sub)
	} else if from.Id != nil {
		return &ast.Exp{
			Id: &ast.Id{Name: *from.Id},
		}
	}
	return nil
}
