package inferer

import (
	"github.com/jvmakine/shine/ast"
	. "github.com/jvmakine/shine/types"
)

type Substitutions map[*TypeVar]Type

func (s Substitutions) Apply(t Type) Type {
	target := s[t.Variable]
	if !target.IsDefined() {
		target = t
	}
	if target.IsFunction() {
		ntyps := make([]Type, len(target.FunctTypes()))
		for i, v := range target.FunctTypes() {
			ntyps[i] = s.Apply(v)
		}
		return MakeFunction(ntyps...)
	}
	return target
}

func (s Substitutions) Convert(exp *ast.Exp) {
	if exp.Block != nil {
		s.Convert(exp.Block.Value)
		for _, a := range exp.Block.Assignments {
			s.Convert(a.Value)
		}
	} else if exp.Call != nil {
		for _, p := range exp.Call.Params {
			s.Convert(p)
		}
		exp.Call.Type = s.Apply(exp.Call.Type)
	} else if exp.Def != nil {
		for _, p := range exp.Def.Params {
			p.Type = s.Apply(p.Type)
		}
		s.Convert(exp.Def.Body)
	} else if exp.Id != nil {
		exp.Id.Type = s.Apply(exp.Id.Type)
	}
}

func (s Substitutions) ConvertAssignment(ass *ast.Assign) {
	s.Convert(ass.Value)
}
