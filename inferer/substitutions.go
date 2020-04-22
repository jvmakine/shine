package inferer

import (
	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

type Substitutions map[*TypeVar]Type

func (s Substitutions) Apply(t Type) Type {
	target := s[t.Variable]
	if t.IsVariable() && target.IsDefined() && !t.IsFunction() {
		return target
	} else if t.IsFunction() {
		ntyps := make([]Type, len(t.FunctTypes()))
		for i, v := range t.FunctTypes() {
			ntyps[i] = s.Apply(v)
		}
		return MakeFunction(ntyps...)
	}
	return t
}

func (s Substitutions) Convert(exp *ast.Exp) {
	exp.Type = s.Apply(exp.Type)
	if exp.Block != nil {
		s.Convert(exp.Block.Value)
	} else if exp.Call != nil {
		for _, p := range exp.Call.Params {
			s.Convert(p)
		}
	} else if exp.Def != nil {
		for _, p := range exp.Def.Params {
			p.Type = s.Apply(p.Type)
		}
		s.Convert(exp.Def.Body)
	}
}

func (s Substitutions) ConvertAssignment(ass *ast.Assign) {
	s.Convert(ass.Value)
}
