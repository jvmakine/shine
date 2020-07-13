package grammar

import (
	"github.com/jvmakine/shine/types"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/ebnf"
	"github.com/jvmakine/shine/ast"
)

type ConvCtx struct {
	DefCount int
}

type Program struct {
	Body *Block `@@`
}

func Parse(str string) (*Program, error) {
	lexer, err := ebnf.New(`
		LineComment = ("//") { "\u0000"…"\uffff"-"\n" } .
		BlockComment = ("/*") { "\u0000"…"\uffff"-"*/" } ("*/") .
		Fun = "=>" .
		Binder = "~>" .
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
		participle.UseLookahead(3),
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
	ctx := &ConvCtx{}
	return convBlock(prg.Body, ctx)
}

func convInterface(name *TypedName, from *Definitions, ctx *ConvCtx) *ast.Interface {
	defs := convDefinitions(from, ctx)
	defs.Visit(ast.NullFun, ast.NullFun, false, func(a ast.Ast, ctx *ast.VisitContext) ast.Ast {
		if id, ok := a.(*ast.Id); ok && id.Name == (*name.Name) {
			return ast.NewId("$")
		}
		return a
	}, ast.NewVisitCtx())
	return &ast.Interface{
		Definitions:   defs,
		InterfaceType: convTypeDef(name.Type),
	}
}

func convDefinitions(from *Definitions, ctx *ConvCtx) *ast.Definitions {
	res := ast.NewDefinitions(ctx.DefCount)
	ctx.DefCount++
	for _, d := range from.Defs {
		if d.Assignment != nil {
			a := d.Assignment
			raw := convAst(a.Value, ctx)
			if e, ok := raw.(ast.Expression); ok {
				res.Assignments[*d.Name.Name] = e
			} else {
				panic("invalid assignment")
			}
			if d.Name.Type != nil {
				t := convTypeDef(d.Name.Type)
				res.Assignments[*d.Name.Name] = &ast.TypeDecl{Exp: res.Assignments[*d.Name.Name], DeclType: t}
			}
		} else if d.Assignment.Interface != nil {
			newI := convInterface(d.Name, d.Assignment.Interface, ctx)
			res.Interfaces = append(res.Interfaces, newI)
		}
	}
	return res
}

func convBlock(from *Block, ctx *ConvCtx) *ast.Block {
	assigns := map[string]ast.Expression{}
	interfs := []*ast.Interface{}
	for _, d := range from.Def.Defs {
		if d.Assignment.Value != nil {
			a := d.Assignment
			raw := convAst(a.Value, ctx)
			if e, ok := raw.(ast.Expression); ok {
				assigns[*d.Name.Name] = e
			} else {
				panic("invalid assignment")
			}
			if d.Name.Type != nil {
				t := convTypeDef(d.Name.Type)
				assigns[*d.Name.Name] = &ast.TypeDecl{Exp: assigns[*d.Name.Name], DeclType: t}
			}
		} else if d.Assignment.Interface != nil {
			b := d.Assignment
			newI := convInterface(d.Name, b.Interface, ctx)
			interfs = append(interfs, newI)
		} else {
			panic("invalid definition")
		}
	}
	id := ctx.DefCount
	ctx.DefCount++
	return &ast.Block{
		Def:   &ast.Definitions{Assignments: assigns, Interfaces: interfs, ID: id},
		Value: convExp(from.Value, ctx),
	}
}

func convExp(from *Expression, ctx *ConvCtx) ast.Expression {
	ut := convAst(from, ctx)
	if e, ok := ut.(ast.Expression); ok {
		return e
	}
	panic("non expression AST")
}

func convAst(from *Expression, ctx *ConvCtx) ast.Ast {
	ut := convUTE(from.Exp, ctx)
	if e, ok := ut.(ast.Expression); ok && from.Type != nil {
		return &ast.TypeDecl{Exp: e, DeclType: convTypeDef(from.Type)}
	}
	return ut
}

func convUTE(from *UTExpression, ctx *ConvCtx) ast.Ast {
	if from.Def != nil {
		return convFDef(from.Def, ctx)
	} else if from.If != nil {
		return convIf(from.If, ctx)
	} else if from.Comp != nil {
		return convOpComp(convComp(from.Comp.Left, ctx), from.Comp.Right, ctx)
	}
	panic("invalid expression")
}

func convIf(from *IfExpression, ctx *ConvCtx) ast.Expression {
	return ast.NewBranch(
		convExp(from.Cond, ctx),
		convExp(from.True, ctx),
		convExp(from.False, ctx),
	)
}

func convFDef(from *FDefinition, ctx *ConvCtx) ast.Ast {
	if fd := from.Funct; fd != nil {
		params := make([]*ast.FParam, len(from.Params))
		for i, p := range from.Params {
			params[i] = convFParam(p)
		}
		body := convExp(fd.Body, ctx)
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
			var td types.Type
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
	if t == nil {
		return nil
	}
	if t.Primitive != "" {
		switch t.Primitive {
		case "int":
			return types.Int
		case "real":
			return types.Real
		case "bool":
			return types.Bool
		case "string":
			return types.String
		default:
			panic("invalid type: " + t.Primitive)
		}
	} else if t.Function != nil {
		ps := make([]types.Type, len(t.Function.Params)+1)
		for i, p := range t.Function.Params {
			ps[i] = convTypeDef(p)
		}
		ps[len(t.Function.Params)] = convTypeDef(t.Function.Return)
		return types.NewFunction(ps...)
	} else if t.Named != "" {
		return types.NewNamed(t.Named, nil)
	}
	panic("invalid type")
}

func convFParam(from *TypedName) *ast.FParam {
	var typ types.Type
	if from.Type != nil {
		typ = convTypeDef(from.Type)
	}
	return &ast.FParam{
		Name:      *from.Name,
		ParamType: typ,
	}
}

func convOpComp(left ast.Expression, right []*OpComp, ctx *ConvCtx) ast.Expression {
	if right == nil || len(right) == 0 {
		return left
	}
	res := ast.NewOp(*right[0].Operation, left, convComp(right[0].Right, ctx))
	return convOpComp(res, right[1:], ctx)
}

func convComp(from *Comp, ctx *ConvCtx) ast.Expression {
	return convOpTerm(convTerm(from.Left, ctx), from.Right, ctx)
}

func convOpTerm(left ast.Expression, right []*OpTerm, ctx *ConvCtx) ast.Expression {
	if right == nil || len(right) == 0 {
		return left
	}
	res := ast.NewOp(*right[0].Operation, left, convTerm(right[0].Right, ctx))
	return convOpTerm(res, right[1:], ctx)
}

func convTerm(from *Term, ctx *ConvCtx) ast.Expression {
	return convOpFact(convAccessor(from.Left, ctx), from.Right, ctx)
}

func convAccessor(from *Accessor, ctx *ConvCtx) ast.Expression {
	acc := from.Right
	res := convFVal(from.Left, ctx)
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
				params[i] = convExp(p, ctx)
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

func convOpFact(left ast.Expression, right []*OpFactor, ctx *ConvCtx) ast.Expression {
	if right == nil || len(right) == 0 {
		return left
	}
	res := ast.NewOp(*right[0].Operation, left, convAccessor(right[0].Right, ctx))
	return convOpFact(res, right[1:], ctx)
}

func convFVal(from *FValue, ctx *ConvCtx) ast.Expression {
	pval := convPVal(from.Value, ctx)
	if len(from.Calls) > 0 {
		call, calls := from.Calls[0], from.Calls[1:]
		params := make([]ast.Expression, len(call.Params))
		for i, p := range call.Params {
			params[i] = convExp(p, ctx)
		}
		res := &ast.FCall{
			Function: pval,
			Params:   params,
		}
		for len(calls) > 0 {
			call, calls = calls[0], calls[1:]
			params := make([]ast.Expression, len(call.Params))
			for i, p := range call.Params {
				params[i] = convExp(p, ctx)
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

func convPVal(from *PValue, ctx *ConvCtx) ast.Expression {
	if from.Block != nil {
		return convBlock(from.Block, ctx)
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
		return convExp(from.Sub, ctx)
	} else if from.Id != nil {
		return &ast.Id{Name: *from.Id}
	}
	return nil
}
