package compiler

import (
	"github.com/jvmakine/shine/ast"
	"github.com/llir/llvm/ir/value"
)

type cresult struct {
	value value.Value
	ids   map[string]value.Value
	exp   map[*ast.Exp]value.Value
}

func makeCR(v value.Value) cresult {
	return cresult{value: v}
}

func (c cresult) cmb(os ...cresult) cresult {
	for _, o := range os {
		for k, v := range o.ids {
			c.ids[k] = v
		}
		for k, v := range o.exp {
			c.exp[k] = v
		}
	}
	return c
}
