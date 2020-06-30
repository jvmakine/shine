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
		PrimitiveType = "int" | "real" | "bool" | "string" .
		Eq = "=" .
		Ident = alpha { alpha | digit } .
		Real = "0"…"9" { digit } "." "0"…"9" { digit } .
		Int = "0" | "1"…"9" { digit } .
		String = "\"" { "\u0000"…"\uffff"-"\n"-"\""} "\"" .
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

func (prg *Program) ToAst() ast.Expression {
	return convBlock(prg.Body)
}

func convInterface(from *Definitions) *ast.Interface {
	return &ast.Interface{}
}

func convBlock(from *Block) *ast.Block {
	assigns := map[string]ast.Expression{}
	interfs := map[string]*ast.Interface{}
	typedecl := map[string]*ast.Struct{}
	for _, d := range from.Def.Defs {
		if d.Assignment != nil {
			a := d.Assignment
			raw := convAst(a.Value)
			if e, ok := raw.(ast.Expression); ok {
				assigns[*a.Name] = e
			} else if s, ok := raw.(*ast.Struct); ok {
				typedecl[*a.Name] = s
			} else {
				panic("invalid assignment")
			}
			if a.Type != nil {
				t := convTypeDef(a.Type)
				assigns[*a.Name] = &ast.TypeDecl{Exp: assigns[*a.Name], DeclType: t}
			}
		} else if d.Binding != nil {
			b := d.Binding
			name := *b.Name
			interfs[name] = convInterface(b.Interface)
		} else {
			panic("invalid definition")
		}
	}
	return &ast.Block{
		Def:   &ast.Definitions{Assignments: assigns, Interfaces: interfs, TypeDefs: typedecl},
		Value: convExp(from.Value),
	}
}

func convExp(from *Expression) ast.Expression {
	ut := convAst(from)
	if e, ok := ut.(ast.Expression); ok {
		return e
	}
	panic("non expression AST")
}

func convAst(from *Expression) ast.Ast {
	ut := convUTE(from.Exp)
	if e, ok := ut.(ast.Expression); ok && from.Type != nil {
		return &ast.TypeDecl{Exp: e, DeclType: convTypeDef(from.Type)}
	}
	return ut
}

func convUTE(from *UTExpression) ast.Ast {
	if from.Def != nil {
		return convFDef(from.Def)
	} else if from.If != nil {
		return convIf(from.If)
	} else if from.Comp != nil {
		return convOpComp(convComp(from.Comp.Left), from.Comp.Right)
	}
	panic("invalid expression")
}

func convIf(from *IfExpression) ast.Expression {
	return &ast.FCall{
		Function: &ast.Op{Name: "if"},
		Params:   []ast.Expression{convExp(from.Cond), convExp(from.True), convExp(from.False)},
	}
}

func convFDef(from *FDefinition) ast.Ast {
	if fd := from.Funct; fd != nil {
		params := make([]*ast.FParam, len(from.Params))
		for i, p := range from.Params {
			params[i] = convFParam(p)
		}
		body := convExp(fd.Body)
		if fd.ReturnType != nil {
			body = &ast.TypeDecl{Exp: body, DeclType: convTypeDef(fd.ReturnType)}
		}
		return &ast.FDef{
			Params: params,
			Body:   body,
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
		return &ast.Struct{
			Fields: fields,
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
		case "string":
			return types.StringP
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
	} else if t.Named != "" {
		return types.MakeNamed(t.Named)
	}
	panic("invalid type")
}

func convFParam(from *FunParam) *ast.FParam {
	typ := types.Type{}
	if from.Type != nil {
		typ = convTypeDef(from.Type)
	}
	return &ast.FParam{
		Name:      *from.Name,
		ParamType: typ,
	}
}

func convOpComp(left ast.Expression, right []*OpComp) ast.Expression {
	if right == nil || len(right) == 0 {
		return left
	}
	res := &ast.FCall{
		Function: &ast.Op{Name: *right[0].Operation},
		Params:   []ast.Expression{left, convComp(right[0].Right)},
	}
	return convOpComp(res, right[1:])
}

func convComp(from *Comp) ast.Expression {
	return convOpTerm(convTerm(from.Left), from.Right)
}

func convOpTerm(left ast.Expression, right []*OpTerm) ast.Expression {
	if right == nil || len(right) == 0 {
		return left
	}
	res := &ast.FCall{
		Function: &ast.Op{Name: *right[0].Operation},
		Params:   []ast.Expression{left, convTerm(right[0].Right)},
	}
	return convOpTerm(res, right[1:])
}

func convTerm(from *Term) ast.Expression {
	return convOpFact(convAccessor(from.Left), from.Right)
}

func convAccessor(from *Accessor) ast.Expression {
	acc := from.Right
	res := convFVal(from.Left)
	for len(acc) > 0 {
		res = &ast.FieldAccessor{
			Exp:   res,
			Field: acc[0].Id,
		}
		calls := acc[0].Calls
		for len(calls) > 0 {
			call := calls[0]
			calls = calls[1:]
			params := make([]ast.Expression, len(call.Params))
			for i, p := range call.Params {
				params[i] = convExp(p)
			}
			res = &ast.FCall{
				Function: res,
				Params:   params,
			}
		}
		acc = acc[1:]
	}
	return res
}

func convOpFact(left ast.Expression, right []*OpFactor) ast.Expression {
	if right == nil || len(right) == 0 {
		return left
	}
	res := &ast.FCall{
		Function: &ast.Op{Name: *right[0].Operation},
		Params:   []ast.Expression{left, convAccessor(right[0].Right)},
	}
	return convOpFact(res, right[1:])
}

func convFVal(from *FValue) ast.Expression {
	pval := convPVal(from.Value)
	if len(from.Calls) > 0 {
		call, calls := from.Calls[0], from.Calls[1:]
		params := make([]ast.Expression, len(call.Params))
		for i, p := range call.Params {
			params[i] = convExp(p)
		}
		res := &ast.FCall{
			Function: pval,
			Params:   params,
		}
		for len(calls) > 0 {
			call, calls = calls[0], calls[1:]
			params := make([]ast.Expression, len(call.Params))
			for i, p := range call.Params {
				params[i] = convExp(p)
			}
			res = &ast.FCall{
				Function: res,
				Params:   params,
			}
		}
		return res
	}
	return pval
}

func convPVal(from *PValue) ast.Expression {
	if from.Block != nil {
		return convBlock(from.Block)
	} else if from.Int != nil {
		return &ast.Const{Int: from.Int}
	} else if from.Real != nil {
		return &ast.Const{Real: from.Real}
	} else if from.Bool != nil {
		value := false
		if *from.Bool == "true" {
			value = true
		}
		return &ast.Const{Bool: &value}
	} else if from.String != nil {
		str := *from.String
		str = str[1:(len(str) - 1)]
		return &ast.Const{String: &str}
	} else if from.Sub != nil {
		return convExp(from.Sub)
	} else if from.Id != nil {
		return &ast.Id{Name: *from.Id}
	}
	return nil
}
