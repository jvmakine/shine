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
	exp.Type = s.Apply(exp.Type)
	if exp.Block != nil {
		s.Convert(exp.Block.Value)
		for _, a := range exp.Block.Assignments {
			s.Convert(a.Value)
		}
	} else if exp.Call != nil {
		for _, p := range exp.Call.Params {
			s.Convert(p)
		}
	} else if exp.Def != nil {
		f := make([]Type, len(exp.Def.Params)+1)
		for i, p := range exp.Def.Params {
			p.Type = s.Apply(p.Type)
			f[i] = p.Type
		}
		s.Convert(exp.Def.Body)
		f[len(exp.Def.Params)] = exp.Def.Body.Type
		// TODO: Why is this needed?
		exp.Type = MakeFunction(f...)
	}
}

func (s Substitutions) ConvertAssignment(ass *ast.Assign) {
	s.Convert(ass.Value)
}
