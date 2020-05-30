package types

type Substitutions struct {
	substitutions map[*TypeVar]Type
	references    map[*TypeVar]map[*TypeVar]bool
}

func MakeSubstitutions() Substitutions {
	return Substitutions{
		substitutions: map[*TypeVar]Type{},
		references:    map[*TypeVar]map[*TypeVar]bool{},
	}
}

func (s Substitutions) Apply(t Type) Type {
	target := s.substitutions[t.Variable]
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
	if target.IsStructure() {
		ntyps := make([]SField, len(target.Structure.Fields))
		for i, v := range target.Structure.Fields {
			ntyps[i] = SField{
				Name: v.Name,
				Type: s.Apply(v.Type),
			}
		}
		return MakeStructure(target.Structure.Name, ntyps...)
	}
	return target
}

func (s Substitutions) Update(from *TypeVar, to Type) error {
	if from == to.Variable {
		return nil
	}

	result := s.Apply(to)

	if p := s.substitutions[from]; p.IsDefined() {
		if result != p {
			uni, err := result.Unifier(p)
			if err != nil {
				return err
			}
			s.Combine(uni)
			result = uni.Apply(result)
		}
	}

	s.substitutions[from] = result

	for _, fv := range result.FreeVars() {
		if s.references[fv] == nil {
			s.references[fv] = map[*TypeVar]bool{}
		}
		s.references[fv][from] = true
	}

	if rs := s.references[from]; rs != nil {
		s.references[from] = nil
		subs := MakeSubstitutions()
		subs.Update(from, result)
		for k := range rs {
			substit := s.substitutions[k]
			s.substitutions[k] = subs.Apply(substit)
			for _, fv := range s.substitutions[k].FreeVars() {
				if s.references[fv] == nil {
					s.references[fv] = map[*TypeVar]bool{}
				}
				s.references[fv][from] = true
			}
		}
	}
	return nil
}

func (s Substitutions) Combine(o Substitutions) error {
	for f, t := range o.substitutions {
		err := s.Update(f, t)
		if err != nil {
			return err
		}
	}
	return nil
}
