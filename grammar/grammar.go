package grammar

import (
	"errors"

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
		Binding = "->" .
		Newline = "\n" .
		Whitespace = " " | "\r" | "\t" .
		Reserved = "if" | "else" | "true" | "false" .
		Comma = "," .
		Dot = "." .
		Brackets = "(" | ")" | "{" | "}" | "[" | "]" .
		COp = ">=" | "<=" .
		Op = "+" | "-" | "*" | "/" | "%" |  ">" | "<" | "==" | "!=" | "||" | "&&" .
		TypeDef = "::" .
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

func (prg *Program) ToAst() (*ast.Exp, error) {
	b, err := convBlock(prg.Body)
	if err != nil {
		return nil, err
	}
	return &ast.Exp{
		Block: b,
	}, nil
}

func convBlock(from *Block) (*ast.Block, error) {
	assigns := map[string]*ast.Exp{}
	typedefs := map[string]*ast.TypeDefinition{}
	bindings := []*ast.TypeBinding{}
	for _, e := range from.Elements {
		if e.Assignment != nil {
			a := e.Assignment
			name := *a.Name
			if assigns[name] != nil || typedefs[name] != nil {
				return nil, errors.New("redefinition of " + name)
			}
			exp, err := convExp(a.Value)
			if err != nil {
				return nil, err
			}
			assigns[name] = exp
			if a.Type != nil {
				t := convTypeDecl(a.Type)
				assigns[name] = &ast.Exp{TDecl: &ast.TypeDecl{Exp: assigns[name], Type: t}}
			}
		} else if e.TypeDef != nil {
			from := e.TypeDef
			name := *from.Name
			if assigns[name] != nil || typedefs[name] != nil {
				return nil, errors.New("redefinition of " + name)
			}
			if from.Struct != nil {
				fields := make([]*ast.StructField, len(from.Struct.Fields))
				for i, p := range from.Struct.Fields {
					td := convTypeDecl(p.Type)
					fields[i] = &ast.StructField{
						Name: p.Name,
						Type: td,
					}
				}
				typedefs[name] = &ast.TypeDefinition{
					FreeVariables: from.FreeVars,
					Struct:        &ast.Struct{Fields: fields},
				}
			} else if from.TypeClass != nil {
				funs := map[string]*ast.TypeDefinition{}
				for _, f := range from.TypeClass.Functions {
					funs[*f.Name] = &ast.TypeDefinition{
						FreeVariables: f.FreeVars,
						TypeDecl:      convTypeFunction(f.Function),
					}
				}
				typedefs[name] = &ast.TypeDefinition{
					FreeVariables: from.FreeVars,
					TypeClass:     &ast.TypeClass{Functions: funs},
				}
			} else {
				t := convTypeDecl(from.Type)
				typedefs[name] = &ast.TypeDefinition{
					FreeVariables: from.FreeVars,
					TypeDecl:      t,
				}
			}
		} else if e.TypeBinding != nil {
			from := e.TypeBinding
			params := make([]types.Type, len(from.Arguments))
			for i, a := range from.Arguments {
				params[i] = convTypeDecl(a)
			}
			fbinds := map[string]*ast.FDef{}
			for _, b := range from.Assignments {
				e, err := convExp(b.Value)
				if err != nil {
					return nil, err
				}
				if e.Def == nil {
					return nil, errors.New("bound value must be a function")
				}
				fbinds[*b.Name] = e.Def
			}
			bindings = append(bindings, &ast.TypeBinding{
				Name:       *from.Name,
				Parameters: params,
				Bindings:   fbinds,
			})
		}
	}
	exp, err := convExp(from.Value)
	if err != nil {
		return nil, err
	}
	return &ast.Block{
		Assignments:  assigns,
		TypeDefs:     typedefs,
		Value:        exp,
		TypeBindings: bindings,
		TCFunctions:  map[string]*ast.TypeDefinition{},
	}, nil
}

func convExp(from *Expression) (*ast.Exp, error) {
	ut, err := convUTExp(from.Exp)
	if err != nil {
		return nil, err
	}
	if from.Type != nil {
		return &ast.Exp{TDecl: &ast.TypeDecl{Exp: ut, Type: convTypeDecl(from.Type)}}, nil
	}
	return ut, nil
}

func convUTExp(from *UTExpression) (*ast.Exp, error) {
	if from.Def != nil {
		return convDef(from.Def)
	} else if from.If != nil {
		return convIf(from.If)
	} else if from.Comp != nil {
		comp, err := convComp(from.Comp.Left)
		if err != nil {
			return nil, err
		}
		return convOpComp(comp, from.Comp.Right)
	}
	panic("invalid expression")
}

func convIf(from *IfExpression) (*ast.Exp, error) {
	c, err := convExp(from.Cond)
	if err != nil {
		return nil, err
	}
	t, err := convExp(from.True)
	if err != nil {
		return nil, err
	}
	f, err := convExp(from.False)
	if err != nil {
		return nil, err
	}
	return &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: "if"}},
			Params:   []*ast.Exp{c, t, f},
		},
	}, nil
}

func convDef(from *Definition) (*ast.Exp, error) {
	fd := from.Funct
	params := make([]*ast.FParam, len(from.Params))
	for i, p := range from.Params {
		params[i] = convFParam(p)
	}
	body, err := convExp(fd.Body)
	if err != nil {
		return nil, err
	}
	if fd.ReturnType != nil {
		body = &ast.Exp{TDecl: &ast.TypeDecl{Exp: body, Type: convTypeDecl(fd.ReturnType)}}
	}
	return &ast.Exp{
		Def: &ast.FDef{
			Params: params,
			Body:   body,
		},
	}, nil
}

func convTypeFunction(t *TypeFunc) types.Type {
	ps := make([]types.Type, len(t.Params)+1)
	for i, p := range t.Params {
		ps[i] = convTypeDecl(p)
	}
	ps[len(t.Params)] = convTypeDecl(t.Return)
	return types.MakeFunction(ps...)
}

func convTypeDecl(t *TypeDeclaration) types.Type {
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
		return convTypeFunction(t.Function)
	} else if t.Named != nil {
		vars := make([]types.Type, len(t.Named.Vars))
		for i, v := range t.Named.Vars {
			vars[i] = convTypeDecl(v)
		}
		return types.MakeNamed(t.Named.Name, vars...)
	}
	panic("invalid type")
}

func convFParam(from *FunParam) *ast.FParam {
	typ := types.Type{}
	if from.Type != nil {
		typ = convTypeDecl(from.Type)
	}
	return &ast.FParam{
		Name: *from.Name,
		Type: typ,
	}
}

func convOpComp(left *ast.Exp, right []*OpComp) (*ast.Exp, error) {
	if right == nil || len(right) == 0 {
		return left, nil
	}
	c, err := convComp(right[0].Right)
	if err != nil {
		return nil, err
	}
	res := &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: *right[0].Operation}},
			Params:   []*ast.Exp{left, c},
		},
	}
	return convOpComp(res, right[1:])
}

func convComp(from *Comp) (*ast.Exp, error) {
	term, err := convTerm(from.Left)
	if err != nil {
		return nil, err
	}
	return convOpTerm(term, from.Right)
}

func convOpTerm(left *ast.Exp, right []*OpTerm) (*ast.Exp, error) {
	if right == nil || len(right) == 0 {
		return left, nil
	}
	term, err := convTerm(right[0].Right)
	if err != nil {
		return nil, err
	}
	res := &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: *right[0].Operation}},
			Params:   []*ast.Exp{left, term},
		},
	}
	return convOpTerm(res, right[1:])
}

func convTerm(from *Term) (*ast.Exp, error) {
	acc, err := convAccessor(from.Left)
	if err != nil {
		return nil, err
	}
	return convOpFact(acc, from.Right)
}

func convAccessor(from *Accessor) (*ast.Exp, error) {
	acc := from.Right
	res, err := convFVal(from.Left)
	if err != nil {
		return nil, err
	}
	for len(acc) > 0 {
		res = &ast.Exp{
			FAccess: &ast.FieldAccessor{
				Exp:   res,
				Field: acc[0].Id,
			},
		}
		calls := acc[0].Calls
		for len(calls) > 0 {
			call := calls[0]
			calls = calls[1:]
			params := make([]*ast.Exp, len(call.Params))
			for i, p := range call.Params {
				e, err := convExp(p)
				if err != nil {
					return nil, err
				}
				params[i] = e
			}
			res = &ast.Exp{
				Call: &ast.FCall{
					Function: res,
					Params:   params,
				}}
		}
		acc = acc[1:]
	}
	return res, nil
}

func convOpFact(left *ast.Exp, right []*OpFactor) (*ast.Exp, error) {
	if right == nil || len(right) == 0 {
		return left, nil
	}
	acc, err := convAccessor(right[0].Right)
	if err != nil {
		return nil, err
	}
	res := &ast.Exp{
		Call: &ast.FCall{
			Function: &ast.Exp{Op: &ast.Op{Name: *right[0].Operation}},
			Params:   []*ast.Exp{left, acc},
		},
	}
	return convOpFact(res, right[1:])
}

func convFVal(from *FValue) (*ast.Exp, error) {
	pval, err := convPVal(from.Value)
	if err != nil {
		return nil, err
	}
	if len(from.Calls) > 0 {
		call, calls := from.Calls[0], from.Calls[1:]
		params := make([]*ast.Exp, len(call.Params))
		for i, p := range call.Params {
			e, err := convExp(p)
			if err != nil {
				return nil, err
			}
			params[i] = e
		}
		res := &ast.FCall{
			Function: pval,
			Params:   params,
		}
		for len(calls) > 0 {
			call, calls = calls[0], calls[1:]
			params := make([]*ast.Exp, len(call.Params))
			for i, p := range call.Params {
				e, err := convExp(p)
				if err != nil {
					return nil, err
				}
				params[i] = e
			}
			res = &ast.FCall{
				Function: &ast.Exp{Call: res},
				Params:   params,
			}
		}
		return &ast.Exp{
			Call: res,
		}, nil
	}
	return pval, nil
}

func convPVal(from *PValue) (*ast.Exp, error) {
	if from.Block != nil {
		b, err := convBlock(from.Block)
		if err != nil {
			return nil, err
		}
		return &ast.Exp{
			Block: b,
		}, nil
	} else if from.Int != nil {
		return &ast.Exp{
			Const: &ast.Const{Int: from.Int},
		}, nil
	} else if from.Real != nil {
		return &ast.Exp{
			Const: &ast.Const{Real: from.Real},
		}, nil
	} else if from.Bool != nil {
		value := false
		if *from.Bool == "true" {
			value = true
		}
		return &ast.Exp{
			Const: &ast.Const{Bool: &value},
		}, nil
	} else if from.String != nil {
		str := *from.String
		str = str[1:(len(str) - 1)]
		return &ast.Exp{
			Const: &ast.Const{String: &str},
		}, nil
	} else if from.Sub != nil {
		return convExp(from.Sub)
	} else if from.Id != nil {
		return &ast.Exp{
			Id: &ast.Id{Name: *from.Id},
		}, nil
	}
	return nil, nil
}
