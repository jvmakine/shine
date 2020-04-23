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

func Assign(name string, value *ast.Exp) *ast.Assign {
	return &ast.Assign{
		Name:  name,
		Value: value,
	}
}

func Block(e *ast.Exp, assigns ...*ast.Assign) *ast.Exp {
	as := assigns
	if as == nil {
		as = []*ast.Assign{}
	}
	return &ast.Exp{
		Block: &ast.Block{Value: e, Assignments: as},
	}
}

func Fcall(name string, args ...*ast.Exp) *ast.Exp {
	call := &ast.FCall{
		Name:   name,
		Params: args,
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
