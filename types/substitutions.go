package types

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

func (s Substitutions) Combine(o Substitutions) (Substitutions, error) {
	result := Substitutions{}
	for k, v := range o {
		if s[k].IsDefined() {
			s, err := s[k].Unifier(v)
			if err != nil {
				return Substitutions{}, err
			}
			result[k] = s.Apply(v)
		} else {
			result[k] = v
		}
	}
	return result, nil
}
