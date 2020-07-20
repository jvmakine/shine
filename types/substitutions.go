package types

import (
	"sort"
	"strings"
)

type Substitutions struct {
	substitutions *map[VariableID]Type
	references    *map[VariableID]map[VariableID]bool
	contexts      *map[VariableID]UnificationCtx
}

func MakeSubstitutions() Substitutions {
	return Substitutions{
		substitutions: &map[VariableID]Type{},
		references:    &map[VariableID]map[VariableID]bool{},
		contexts:      &map[VariableID]UnificationCtx{},
	}
}

// internal context used to deal with recursive structural variables
type substitutionCtx struct {
	visited map[VariableID]map[VariableID]bool
}

func newSubstCtx() *substitutionCtx {
	return &substitutionCtx{map[VariableID]map[VariableID]bool{}}
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

func (s Substitutions) update(from VariableID, to Type, ctx UnificationCtx, sctx *substitutionCtx) error {
	if v, ok := to.(Variable); ok && v.ID == from {
		return nil
	}

	result, _ := to.convert(s, sctx)
	// deal with recursive variables
	if v, ok := result.(Variable); ok {
		if sctx.visited[from] == nil {
			sctx.visited[from] = map[VariableID]bool{}
		}
		if sctx.visited[from][v.ID] {
			return nil
		}
		sctx.visited[from][v.ID] = true
	}
	if p := (*s.substitutions)[from]; p != nil {
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

	(*s.substitutions)[from] = result
	if v, ok := result.(Variable); ok {
		if (*s.contexts)[v.ID] != nil {
			s.AddContext(from, (*s.contexts)[v.ID])
		}
	}

	for _, fv := range result.freeVars() {
		if (*s.references)[fv.ID] == nil {
			(*s.references)[fv.ID] = map[VariableID]bool{}
		}
		(*s.references)[fv.ID][from] = true
	}

	if rs := (*s.references)[from]; rs != nil {
		(*s.references)[from] = nil
		subs := MakeSubstitutions()
		subs.update(from, result, ctx, sctx)
		for k := range rs {
			if k != from {
				substit := (*s.substitutions)[k]
				c, _ := substit.convert(subs, sctx)
				(*s.substitutions)[k] = c
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

func (s Substitutions) Add(from Type, to Type, ctx UnificationCtx) error {
	sub, err := Unifier(from, to, ctx)
	if err != nil {
		return err
	}
	return s.combine(sub, ctx, newSubstCtx())
}

func (s Substitutions) combine(o Substitutions, ctx UnificationCtx, sctx *substitutionCtx) error {
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
	return nil
}
func (s Substitutions) AddContext(v VariableID, ctx UnificationCtx) {
	if (*s.contexts)[v] != nil {
		panic("context already set")
	}
	(*s.contexts)[v] = ctx
}

func (s Substitutions) Copy() Substitutions {
	newRef := map[VariableID]map[VariableID]bool{}
	newSub := map[VariableID]Type{}
	newCon := map[VariableID]UnificationCtx{}

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

	return Substitutions{
		references:    &newRef,
		substitutions: &newSub,
		contexts:      &newCon,
	}
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
