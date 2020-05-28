package test

import (
	"github.com/jvmakine/shine/types"

	"github.com/jvmakine/shine/ast"
)

func IConst(v int64) *ast.Exp {
	return &ast.Exp{
		Const: &ast.Const{Int: &v},
	}
}

func RConst(v float64) *ast.Exp {
	return &ast.Exp{
		Const: &ast.Const{Real: &v},
	}
}

func BConst(v bool) *ast.Exp {
	return &ast.Exp{
		Const: &ast.Const{Bool: &v},
	}
}

func Id(name string) *ast.Exp {
	return &ast.Exp{
		Id: &ast.Id{Name: name},
	}
}

func Op(name string) *ast.Exp {
	return &ast.Exp{
		Op: &ast.Op{Name: name},
	}
}

type Assgs = map[string]*ast.Exp

func Block(assign Assgs, e *ast.Exp) *ast.Exp {
	return &ast.Exp{
		Block: &ast.Block{Value: e, Assignments: assign},
	}
}

func TDecl(e *ast.Exp, t types.Type) *ast.Exp {
	return &ast.Exp{
		TDecl: &ast.TypeDecl{Exp: e, Type: t},
	}
}

func Fcall(function *ast.Exp, args ...*ast.Exp) *ast.Exp {
	call := &ast.FCall{
		Function: function,
		Params:   args,
	}
	return &ast.Exp{
		Call: call,
	}
}

func Param(name string, typ types.Type) *ast.FParam {
	return &ast.FParam{
		Name: name,
		Type: typ,
	}
}

func Fdef(body *ast.Exp, args ...interface{}) *ast.Exp {
	params := make([]*ast.FParam, len(args))
	for i, p := range args {
		if s, ok := p.(string); ok {
			params[i] = &ast.FParam{Name: s}
		} else {
			params[i] = p.(*ast.FParam)
		}
	}
	fdef := &ast.FDef{
		Body:   body,
		Params: params,
	}
	return &ast.Exp{
		Def: fdef,
	}
}
