package types

import (
	"errors"
)

type unificationCtx struct {
	seenStructures map[*Structure]map[*Structure]bool
}

func UnificationError(a Type, b Type) error {
	sa := a.Signature()
	sb := b.Signature()
	if sa < sb {
		return errors.New("can not unify " + sa + " with " + sb)
	} else {
		return errors.New("can not unify " + sb + " with " + sa)
	}
}

func (t Type) Unify(o Type) (Type, error) {
	sub, err := t.Unifier(o)
	if err != nil {
		return t, err
	}
	res := sub.Apply(t)
	if res.IsStructure() && t.IsStructure() && o.IsStructure() {
		if t.Structure.Name != o.Structure.Name {
			res.Structure.Name = ""
		}
	}
	return res, nil
}

func (t Type) Unifier(o Type) (Substitutions, error) {
	ctx := &unificationCtx{seenStructures: map[*Structure]map[*Structure]bool{}}
	return unifier(t, o, ctx)
}

func unifier(t Type, o Type, ctx *unificationCtx) (Substitutions, error) {
	if o.IsPrimitive() && t.IsPrimitive() && *o.Primitive != *t.Primitive {
		return Substitutions{}, UnificationError(o, t)
	}
	if (o.IsPrimitive() && t.IsFunction()) || (o.IsFunction() && t.IsPrimitive()) {
		return Substitutions{}, UnificationError(o, t)
	}
	if (o.IsPrimitive() && t.IsStructure()) || (o.IsStructure() && t.IsPrimitive()) {
		return Substitutions{}, UnificationError(o, t)
	}
	if (o.IsFunction() && t.IsStructure()) || (o.IsStructure() && t.IsFunction()) {
		return Substitutions{}, UnificationError(o, t)
	}
	if t.IsVariable() && o.IsVariable() {
		return unifyVariables(t, o, ctx)
	}
	if o.IsVariable() && !t.IsVariable() {
		return unifier(o, t, ctx)
	}
	if t.IsVariable() && o.IsTypeClassRef() {
		subs := MakeSubstitutions()
		subs.Update(t.Variable, o)
		return subs, nil
	}
	if t.IsVariable() && o.IsHVariable() {
		subs := MakeSubstitutions()
		subs.Update(t.Variable, o)
		return subs, nil
	}
	if o.IsTypeClassRef() && !t.IsTypeClassRef() {
		return unifier(o, t, ctx)
	}
	if t.IsVariable() && o.IsFunction() {
		if t.IsUnionVar() || t.IsUnionVar() {
			return Substitutions{}, UnificationError(o, t)
		}
		subs := MakeSubstitutions()
		subs.Update(t.Variable, o)
		return subs, nil
	}
	if t.IsVariable() && o.IsStructure() {
		if t.IsUnionVar() {
			return Substitutions{}, UnificationError(o, t)
		} else if t.IsStructuralVar() {
			return unifyStructureWithStructuralVar(t, o)
		}
		subs := MakeSubstitutions()
		subs.Update(t.Variable, o)
		return subs, nil
	}
	if t.IsTypeClassRef() && o.IsTypeClassRef() {
		subs := MakeSubstitutions()
		for i, first := range t.TCRef.TypeClassVars {
			second := o.TCRef.TypeClassVars[i]
			s, err := first.Unifier(second)
			if err != nil {
				return Substitutions{}, err
			}
			err = subs.Combine(s)
			if err != nil {
				return Substitutions{}, err
			}
		}

		return subs, nil
	}
	if t.IsTypeClassRef() {
		for _, b := range t.TCRef.LocalBindings {
			f1 := b.Args[t.TCRef.Place]

			subs := MakeSubstitutions()

			s, err := o.Unifier(f1)
			if err != nil {
				continue
			}
			err = subs.Combine(s)
			if err != nil {
				continue
			}
			s, err = t.TCRef.TypeClassVars[t.TCRef.Place].Unifier(f1)
			if err != nil {
				continue
			}
			err = subs.Combine(s)
			if err != nil {
				continue
			}

			return subs, nil
		}
		return Substitutions{}, UnificationError(o, t)
	}
	if o.IsFunction() && t.IsFunction() {
		return unifyFunctions(t, o, ctx)
	}
	if o.IsStructure() && t.IsStructure() {
		return unifyStructures(t, o, ctx)
	}
	if t.IsStructure() && o.IsHVariable() {
		return unifier(o, t, ctx)
	}
	if t.IsHVariable() && o.IsStructure() {
		subs := MakeSubstitutions()
		if len(t.HVariable.Params) != len(o.Structure.TypeArguments) {
			return Substitutions{}, UnificationError(o, t)
		}
		for i, p := range t.HVariable.Params {
			s, err := unifier(p, o.Structure.TypeArguments[i], ctx)
			if err != nil {
				return Substitutions{}, err
			}
			err = subs.Combine(s)
			if err != nil {
				return Substitutions{}, err
			}
		}
		err := subs.Update(t.HVariable.Root, o)
		if err != nil {
			return Substitutions{}, err
		}
		return subs, nil
	}
	if t.IsHVariable() && o.IsHVariable() {
		subs := MakeSubstitutions()
		err := subs.Update(t.HVariable.Root, Type{Variable: o.HVariable.Root})
		if err != nil {
			return Substitutions{}, err
		}
		if len(t.HVariable.Params) != len(o.HVariable.Params) {
			return Substitutions{}, UnificationError(o, t)
		}
		for i, p := range t.HVariable.Params {
			s, err := unifier(p, o.HVariable.Params[i], ctx)
			if err != nil {
				return Substitutions{}, err
			}
			err = subs.Combine(s)
			if err != nil {
				return Substitutions{}, err
			}
		}
		return subs, nil
	}
	if o.IsPrimitive() {
		if t.IsHVariable() {
			return Substitutions{}, UnificationError(o, t)
		}
		if t.IsUnionVar() {
			err := t.Variable.Union.Unifies(*o.Primitive)
			subs := MakeSubstitutions()
			subs.Update(t.Variable, o)
			return subs, err
		} else if t.IsVariable() {
			subs := MakeSubstitutions()
			subs.Update(t.Variable, o)
			return subs, nil
		}
		return Substitutions{}, nil
	}
	if !o.IsDefined() {
		return Substitutions{}, nil
	}
	return Substitutions{}, UnificationError(o, t)
}

func unifyStructureWithStructuralVar(v Type, s Type) (Substitutions, error) {
	tres := MakeSubstitutions()
	smap := map[string]Type{}
	for _, f := range s.Structure.Fields {
		smap[f.Name] = f.Type
	}

	for n, t := range v.Variable.Structural {
		sv := smap[n]
		if !sv.IsDefined() {
			return Substitutions{}, UnificationError(v, s)
		}
		sub, err := t.Unifier(sv)
		if err != nil {
			return Substitutions{}, err
		}
		err = tres.Combine(sub)
		if err != nil {
			return Substitutions{}, err
		}
	}
	tres.Update(v.Variable, s)
	return tres, nil
}

func unifyVariables(t Type, o Type, ctx *unificationCtx) (Substitutions, error) {
	if o.IsStructuralVar() && !t.IsStructuralVar() {
		return unifier(o, t, ctx)
	} else if o.IsUnionVar() && t.IsStructuralVar() {
		return Substitutions{}, UnificationError(o, t)
	} else if o.IsUnionVar() && !t.IsUnionVar() {
		return unifier(o, t, ctx)
	} else if t.IsUnionVar() && o.IsUnionVar() {
		resolv, err := t.Variable.Union.Resolve(o.Variable.Union)
		if len(resolv) == 1 {
			prim := MakePrimitive(resolv[0])
			subs := MakeSubstitutions()
			subs.Update(t.Variable, prim)
			subs.Update(o.Variable, prim)
			return subs, err
		}
		rv := MakeUnionVar(resolv...)
		subs := MakeSubstitutions()
		subs.Update(t.Variable, rv)
		subs.Update(o.Variable, rv)
		return subs, err
	} else if t.IsStructuralVar() && o.IsStructuralVar() {
		ts := t.Variable.Structural
		os := o.Variable.Structural
		res := map[string]Type{}

		for k, v := range ts {
			res[k] = v
		}
		subs := MakeSubstitutions()
		for k, v := range os {
			if res[k].IsDefined() {
				s, err := res[k].Unifier(v)
				if err != nil {
					return Substitutions{}, err
				}
				err = subs.Combine(s)
				if err != nil {
					return Substitutions{}, err
				}
			}
			res[k] = v
		}
		rv := MakeStructuralVar(res)
		subs.Update(t.Variable, rv)
		subs.Update(o.Variable, rv)
		return subs, nil
	}
	subs := MakeSubstitutions()
	subs.Update(o.Variable, t)
	return subs, nil
}

func unifyFunctions(t Type, o Type, ctx *unificationCtx) (Substitutions, error) {
	op := o.FunctTypes()
	tp := t.FunctTypes()
	if len(op) != len(tp) {
		return MakeSubstitutions(), UnificationError(o, t)
	}
	result := MakeSubstitutions()
	for i, p := range op {
		s, err := unifier(p, tp[i], ctx)
		if err != nil {
			return MakeSubstitutions(), err
		}
		err = result.Combine(s)
		if err != nil {
			return Substitutions{}, err
		}
	}
	return result, nil
}

func unifyStructures(t Type, o Type, ctx *unificationCtx) (Substitutions, error) {
	if t.Structure.Name != o.Structure.Name {
		return MakeSubstitutions(), UnificationError(o, t)
	}
	// handle recursice structures
	if ctx.seenStructures[t.Structure] != nil {
		if ctx.seenStructures[t.Structure][o.Structure] {
			return MakeSubstitutions(), nil
		}
	} else {
		ctx.seenStructures[t.Structure] = map[*Structure]bool{}
	}
	ctx.seenStructures[t.Structure][o.Structure] = true
	ofs := map[string]Type{}
	for _, f := range o.Structure.Fields {
		ofs[f.Name] = f.Type
	}

	tfs := map[string]Type{}
	result := MakeSubstitutions()
	for _, f := range t.Structure.Fields {
		tfs[f.Name] = f.Type
		ot := ofs[f.Name]
		if !ot.IsDefined() {
			return MakeSubstitutions(), UnificationError(o, t)
		}
		s, err := unifier(f.Type, ot, ctx)
		if err != nil {
			return MakeSubstitutions(), err
		}
		err = result.Combine(s)
		if err != nil {
			return Substitutions{}, err
		}
	}
	return result, nil
}
