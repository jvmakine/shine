package types

import (
	"sort"
	"strings"
)

type Substitutions struct {
	substitutions *map[VariableID]Type
	references    *map[VariableID]map[VariableID]bool
	contexts      *map[VariableID]UnificationCtx
	equalities    *map[VariableID]map[VariableID]bool
}

func MakeSubstitutions() Substitutions {
	return Substitutions{
		substitutions: &map[VariableID]Type{},
		references:    &map[VariableID]map[VariableID]bool{},
		contexts:      &map[VariableID]UnificationCtx{},
		equalities:    &map[VariableID]map[VariableID]bool{},
	}
}

// internal context used to deal with recursive structural variables
type substitutionCtx struct {
	updated   map[VariableID]bool
	converted map[VariableID]bool
	unifying  map[VariableID]bool
	resolved  map[VariableID]map[string]bool
}

func newSubstCtx() substitutionCtx {
	return substitutionCtx{
		map[VariableID]bool{},
		map[VariableID]bool{},
		map[VariableID]bool{},
		map[VariableID]map[string]bool{},
	}
}

type freeVarsCtx = map[VariableID]bool

func newFreeVarsCtx() freeVarsCtx {
	return freeVarsCtx{}
}

func (s Substitutions) Apply(t Type) Type {
	if t == nil {
		return nil
	}
	conv, _ := t.convert(s, newSubstCtx())
	return conv
}

func (s Substitutions) Update(from VariableID, to Type, ctx UnificationCtx) error {
	return s.update(from, to, ctx, newSubstCtx())
}

func (s Substitutions) update(from VariableID, to Type, ctx UnificationCtx, sctx substitutionCtx) error {
	if v, ok := to.(Variable); ok {
		if v.ID == from {
			return nil
		}
		s.addEquality(from, v.ID)
	}

	result, changed := to.convert(s, sctx)

	// deal with recursive variables
	if _, ok := result.(Variable); ok {
		if sctx.updated[from] {
			return nil
		}
		sctx.updated[from] = true
	}

	// If we have changed the target variable because of existing subsitutions,
	// Add that substitution so that future references to the variable
	// unify to the same result
	if v, ok := to.(Variable); ok && changed {
		err := s.update(v.ID, result, ctx, sctx)
		if err != nil {
			return err
		}
	}

	if p := (*s.substitutions)[from]; p != nil {
		pv, pok := p.(Variable)
		rv, rok := result.(Variable)

		if pok && rok && pv.ID == rv.ID {
			return nil
		}

		uni, err := Unifier(result, p, ctx)
		if err != nil {
			return err
		}
		err = s.combine(uni, ctx, sctx)
		if err != nil {
			return err
		}
		result = s.Apply(result)
	}

	if v, ok := result.(Variable); ok && v.ID == from {
		return nil
	}

	(*s.substitutions)[from] = result
	if v, ok := result.(Variable); ok {
		if (*s.contexts)[v.ID] != nil {
			s.AddContext(from, (*s.contexts)[v.ID])
		}
	}

	if rs := (*s.references)[from]; rs != nil {
		delete(*s.references, from)
		subs := MakeSubstitutions()
		subs.update(from, result, ctx, newSubstCtx())
		for k := range rs {
			if k != from {
				substit := (*s.substitutions)[k]
				c, _ := substit.convert(subs, sctx)
				(*s.substitutions)[k] = c
				for _, fv := range (*s.substitutions)[k].freeVars(newFreeVarsCtx()) {
					if (*s.references)[fv.ID] == nil {
						(*s.references)[fv.ID] = map[VariableID]bool{}
					}
					(*s.references)[fv.ID][k] = true
				}
			}
		}
	}

	for _, fv := range result.freeVars(newFreeVarsCtx()) {
		if (*s.references)[fv.ID] == nil {
			(*s.references)[fv.ID] = map[VariableID]bool{}
		}
		(*s.references)[fv.ID][from] = true
	}

	return nil
}

func (s Substitutions) Combine(o Substitutions, ctx UnificationCtx) error {
	return s.combine(o, ctx, newSubstCtx())
}

func (s Substitutions) Add(from Type, to Type, ctx UnificationCtx) error {
	return s.add(from, to, ctx, newSubstCtx())
}

func (s Substitutions) add(from Type, to Type, ctx UnificationCtx, sctx substitutionCtx) error {
	sub, err := unifier(from, to, ctx, sctx)
	if err != nil {
		return err
	}
	return s.combine(sub, ctx, sctx)
}

func (s Substitutions) combine(o Substitutions, ctx UnificationCtx, sctx substitutionCtx) error {
	// do not modify s if the combination fails
	attempt := s.Copy()
	for k, v := range *o.contexts {
		attempt.AddContext(k, v)
	}

	for f, t := range *o.substitutions {
		err := attempt.update(f, t, ctx, sctx)
		if err != nil {
			return err
		}
	}

	*s.references = *attempt.references
	*s.substitutions = *attempt.substitutions
	*s.contexts = *attempt.contexts
	*s.equalities = *attempt.equalities
	return nil
}

func (s Substitutions) AddContext(v VariableID, ctx UnificationCtx) {
	if (*s.contexts)[v] != nil {
		return
	}
	(*s.contexts)[v] = ctx
}

func (s Substitutions) Copy() Substitutions {
	newRef := map[VariableID]map[VariableID]bool{}
	newSub := map[VariableID]Type{}
	newCon := map[VariableID]UnificationCtx{}
	newEq := map[VariableID]map[VariableID]bool{}

	for k, m := range *s.references {
		newRef[k] = map[VariableID]bool{}
		for k2, v2 := range m {
			newRef[k][k2] = v2
		}
	}

	for k, v := range *s.substitutions {
		newSub[k] = v
	}

	for k, v := range *s.contexts {
		newCon[k] = v
	}

	for k, v := range *s.equalities {
		eqar := map[VariableID]bool{}
		for i, x := range v {
			eqar[i] = x
		}
		newEq[k] = eqar
	}

	return Substitutions{
		references:    &newRef,
		substitutions: &newSub,
		contexts:      &newCon,
		equalities:    &newEq,
	}
}

func (s Substitutions) addEquality(from VariableID, to VariableID) {
	if (*s.equalities)[from] == nil {
		(*s.equalities)[from] = map[VariableID]bool{}
	}
	(*s.equalities)[from][to] = true

	if (*s.equalities)[to] == nil {
		(*s.equalities)[to] = map[VariableID]bool{}
	}
	(*s.equalities)[to][from] = true
}

func (s Substitutions) GetContext(from VariableID) UnificationCtx {
	visited := map[VariableID]bool{}
	todo := []VariableID{from}
	for len(todo) > 0 {
		v, nt := todo[0], todo[1:]
		todo = nt
		if !visited[v] {
			visited[v] = true
			if c := (*s.contexts)[v]; c != nil {
				return c
			}
			for n := range (*s.equalities)[v] {
				todo = append(todo, n)
			}
		}
	}
	return nil
}

func (s Substitutions) String() string {
	var sb strings.Builder
	keys := make([]string, 0, len(*s.substitutions))
	for k := range *s.substitutions {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, k := range keys {
		t := (*s.substitutions)[VariableID(k)]
		sb.WriteString(k)
		sb.WriteString(" => ")
		sb.WriteString(Signature(t))
		sb.WriteString("\n")
	}
	return sb.String()
}
