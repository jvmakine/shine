package inferer

import (
	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

type Substitutions map[*TypeVar]Type

func (s Substitutions) Apply(t *Type) {
	if t.IsVariable() && s[t.Variable].IsDefined() {
		t.AssignFrom(s[t.Variable])
	} else if t.IsFunction() {
		for i, v := range *t.Function {
			s.Apply(&v)
			(*t.Function)[i] = v
		}
	}
}

func (s Substitutions) Convert(exp *ast.Exp) {
	s.Apply(&exp.Type)
	if exp.Block != nil {
		s.Convert(exp.Block.Value)
	} else if exp.Call != nil {
		for _, p := range exp.Call.Params {
			s.Convert(p)
		}
	} else if exp.Def != nil {
		for _, p := range exp.Def.Params {
			s.Apply(&p.Type)
		}
		s.Convert(exp.Def.Body)
	}
}

func (s Substitutions) ConvertAssignment(ass *ast.Assign) {
	s.Convert(ass.Value)
}
