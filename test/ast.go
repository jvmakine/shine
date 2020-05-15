package test

import "github.com/jvmakine/shine/ast"

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

func Fcall(function *ast.Exp, args ...*ast.Exp) *ast.Exp {
	call := &ast.FCall{
		Function: function,
		Params:   args,
	}
	return &ast.Exp{
		Call: call,
	}
}

func Fdef(body *ast.Exp, args ...string) *ast.Exp {
	params := make([]*ast.FParam, len(args))
	for i, p := range args {
		params[i] = &ast.FParam{Name: p}
	}
	fdef := &ast.FDef{
		Body:   body,
		Params: params,
	}
	return &ast.Exp{
		Def: fdef,
	}
}
