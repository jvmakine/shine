package types

type Substitutions struct {
	substitutions *map[VariableID]Type
	references    *map[VariableID]map[VariableID]bool
}

func MakeSubstitutions() Substitutions {
	return Substitutions{
		substitutions: &map[VariableID]Type{},
		references:    &map[VariableID]map[VariableID]bool{},
	}
}

// internal context used to deal with recursive structural variables
type substitutionCtx struct {
	visited map[VariableID]bool
}

func newSubstCtx() *substitutionCtx {
	return &substitutionCtx{map[VariableID]bool{}}
}

func (s Substitutions) Apply(t Type) Type {
	if t == nil {
		return nil
	}
	conv, _ := t.Convert(s)
	return conv
}

func (s Substitutions) Update(from VariableID, to Type, ctx UnificationCtx) error {
	return s.update(from, to, ctx, newSubstCtx())
}

func (s Substitutions) update(from VariableID, to Type, ctx UnificationCtx, sctx *substitutionCtx) error {
	if v, ok := to.(Variable); ok && v.ID == from {
		return nil
	}

	result, changed := to.Convert(s)
	if !sctx.visited[from] {
		if p := (*s.substitutions)[from]; p != nil {
			uni, err := Unifier(result, p, ctx)
			if err != nil {
				return err
			}
			sctx.visited[from] = true
			err = s.combine(uni, ctx, sctx)
			if err != nil {
				return err
			}
			result = uni.Apply(result)
		}
	}

	(*s.substitutions)[from] = result

	for _, fv := range result.freeVars() {
		if (*s.references)[fv.ID] == nil {
			(*s.references)[fv.ID] = map[VariableID]bool{}
		}
		(*s.references)[fv.ID][from] = true
	}

	if rs := (*s.references)[from]; rs != nil && changed {
		(*s.references)[from] = nil
		subs := MakeSubstitutions()
		subs.Update(from, result, ctx)
		for k := range rs {
			if k != from {
				substit := (*s.substitutions)[k]
				(*s.substitutions)[k] = subs.Apply(substit)
				for _, fv := range (*s.substitutions)[k].freeVars() {
					if (*s.references)[fv.ID] == nil {
						(*s.references)[fv.ID] = map[VariableID]bool{}
					}
					(*s.references)[fv.ID][from] = true
				}
			}
		}
	}
	return nil
}
func (s Substitutions) Combine(o Substitutions, ctx UnificationCtx) error {
	return s.combine(o, ctx, newSubstCtx())
}

func (s Substitutions) combine(o Substitutions, ctx UnificationCtx, sctx *substitutionCtx) error {
	// do not modify s if the combination fails
	attempt := s.Copy()
	for f, t := range *o.substitutions {
		err := attempt.update(f, t, ctx, sctx)
		if err != nil {
			return err
		}
	}
	*s.references = *attempt.references
	*s.substitutions = *attempt.substitutions
	return nil
}

func (s Substitutions) Copy() Substitutions {
	newRef := map[VariableID]map[VariableID]bool{}
	newSub := map[VariableID]Type{}
	for k, m := range *s.references {
		newRef[k] = map[VariableID]bool{}
		for k2, v2 := range m {
			newRef[k][k2] = v2
		}
	}

	for k, v := range *s.substitutions {
		newSub[k] = v
	}

	return Substitutions{
		references:    &newRef,
		substitutions: &newSub,
	}
}
