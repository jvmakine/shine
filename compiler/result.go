package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir/value"
)

type cresult struct {
	value value.Value
	ast   *ast.Exp
	exp   map[*ast.Exp]value.Value
}

func makeCR(e *ast.Exp, v value.Value) cresult {
	if e.Type().IsFunction() && e.Id == nil {
		return cresult{value: v, exp: map[*ast.Exp]value.Value{e: v}, ast: e}
	}
	return cresult{value: v, ast: e}
}

func (c cresult) cmb(os ...cresult) cresult {
	for _, o := range os {
		if c.exp == nil {
			c.exp = map[*ast.Exp]value.Value{}
		}
		for k, v := range o.exp {
			c.exp[k] = v
		}
	}
	return c
}
