package types

type Substitutions struct {
	substitutions map[VariableID]Type
	references    map[VariableID]map[VariableID]bool
}

func MakeSubstitutions() Substitutions {
	return Substitutions{
		substitutions: map[VariableID]Type{},
		references:    map[VariableID]map[VariableID]bool{},
	}
}

type substCtx struct {
	visited map[Named]bool
}

func (s Substitutions) Apply(t Type) Type {
	return Convert(t, s)
}

func (s Substitutions) Update(from VariableID, to Type) error {
	if v, ok := to.(Variable); ok && v.ID == from {
		return nil
	}

	result := s.Apply(to)

	if p := s.substitutions[from]; p != nil {
		uni, err := Unifier(result, p)
		if err != nil {
			return err
		}
		err = s.Combine(uni)
		if err != nil {
			return err
		}
		result = uni.Apply(result)
	}

	s.substitutions[from] = result

	for _, fv := range result.freeVars(&unificationCtx{seenIDs: map[VariableID]bool{}}) {
		if s.references[fv.ID] == nil {
			s.references[fv.ID] = map[VariableID]bool{}
		}
		s.references[fv.ID][from] = true
	}

	if rs := s.references[from]; rs != nil {
		s.references[from] = nil
		subs := MakeSubstitutions()
		subs.Update(from, result)
		for k := range rs {
			substit := s.substitutions[k]
			s.substitutions[k] = subs.Apply(substit)
			for _, fv := range s.substitutions[k].freeVars(&unificationCtx{seenIDs: map[VariableID]bool{}}) {
				if s.references[fv.ID] == nil {
					s.references[fv.ID] = map[VariableID]bool{}
				}
				s.references[fv.ID][from] = true
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

func (s Substitutions) Copy() Substitutions {
	newRef := map[VariableID]map[VariableID]bool{}
	newSub := map[VariableID]Type{}
	for k := range s.references {
		newRef[k] = map[VariableID]bool{}
		for k2, v2 := range newRef[k] {
			newRef[k][k2] = v2
		}
	}

	for k, v := range s.substitutions {
		newSub[k] = v
	}

	return Substitutions{
		references:    newRef,
		substitutions: newSub,
	}
}
