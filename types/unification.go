package types

import (
	"errors"
)

type unificationCtx struct {
	seenStructures map[*Structure]map[*Structure]bool
	sctx           *signctx
	resolver       BindingResolver
}

type BindingResolver interface {
	FindBindings(name string, types []Type) []*TCBinding
}

type NullResolver struct{}

func (r *NullResolver) FindBindings(name string, types []Type) []*TCBinding {
	return nil
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

func (t Type) Unify(o Type, resolver BindingResolver) (Type, error) {
	sub, err := t.Unifier(o, resolver)
	if err != nil {
		return t, err
	}
	res := sub.Apply(t, resolver)
	if res.IsStructure() && t.IsStructure() && o.IsStructure() {
		if t.Structure.Name != o.Structure.Name {
			res.Structure.Name = ""
		}
	}
	return res, nil
}

func (t Type) Unifier(o Type, resolver BindingResolver) (Substitutions, error) {
	ctx := &unificationCtx{
		seenStructures: map[*Structure]map[*Structure]bool{},
		sctx:           newSignCtx(),
		resolver:       resolver,
	}
	return unifier(t, o, ctx)
}

func (t Type) Unifies(o Type, resolver BindingResolver) bool {
	_, err := t.Unifier(o, resolver)
	return err == nil
}

func unifier(t Type, o Type, ctx *unificationCtx) (Substitutions, error) {
	if o.IsPrimitive() && t.IsPrimitive() && *o.Primitive == *t.Primitive {
		return Substitutions{}, nil
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
		subs.Update(t.Variable, o, ctx.resolver)
		return subs, nil
	}
	if t.IsVariable() && o.IsHVariable() {
		subs := MakeSubstitutions()
		subs.Update(t.Variable, o, ctx.resolver)
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
		subs.Update(t.Variable, o, ctx.resolver)
		return subs, nil
	}
	if t.IsVariable() && o.IsStructure() {
		if t.IsUnionVar() {
			return Substitutions{}, UnificationError(o, t)
		} else if t.IsStructuralVar() {
			return unifyStructureWithStructuralVar(t, o, ctx)
		}
		subs := MakeSubstitutions()
		subs.Update(t.Variable, o, ctx.resolver)
		return subs, nil
	}
	if t.IsTypeClassRef() && o.IsTypeClassRef() {
		f := t.TCRef.TypeClassVars[t.TCRef.Place]
		s := o.TCRef.TypeClassVars[o.TCRef.Place]
		return f.Unifier(s, ctx.resolver)
	}
	if t.IsTypeClassRef() {
		bs := ctx.resolver.FindBindings(t.TCRef.TypeClass, t.TCRef.TypeClassVars)
		if len(bs) == 0 {
			return Substitutions{}, UnificationError(o, t)
		}
		for _, b := range bs {
			f1 := b.Args[t.TCRef.Place].Copy(NewTypeCopyCtx())

			subs := MakeSubstitutions()

			s, err := o.Unifier(f1, ctx.resolver)
			if err != nil {
				continue
			}
			err = subs.Combine(s, ctx.resolver)
			if err != nil {
				continue
			}
			s, err = t.TCRef.TypeClassVars[t.TCRef.Place].Unifier(f1, ctx.resolver)
			if err != nil {
				continue
			}
			err = subs.Combine(s, ctx.resolver)
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
		err := subs.Update(t.HVariable.Root, Type{Constructor: &TypeConstructor{Name: o.Structure.Name, Arguments: len(o.Structure.TypeArguments), Underlying: o}}, ctx.resolver)
		if err != nil {
			return Substitutions{}, err
		}
		for i, p1 := range t.HVariable.Params {
			p2 := o.Structure.TypeArguments[i]
			s, err := p1.Unifier(p2, ctx.resolver)
			if err != nil {
				return Substitutions{}, err
			}
			err = subs.Combine(s, ctx.resolver)
			if err != nil {
				return Substitutions{}, err
			}
		}
		return subs, nil
	}
	if t.IsHVariable() && o.IsHVariable() {
		subs := MakeSubstitutions()
		err := subs.Update(t.HVariable.Root, Type{Variable: o.HVariable.Root}, ctx.resolver)
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
			err = subs.Combine(s, ctx.resolver)
			if err != nil {
				return Substitutions{}, err
			}
		}
		return subs, nil
	}
	if o.IsPrimitive() {
		if t.IsUnionVar() {
			err := t.Variable.Union.Unifies(*o.Primitive)
			if err != nil {
				return Substitutions{}, err
			}
			subs := MakeSubstitutions()
			err = subs.Update(t.Variable, o, ctx.resolver)
			return subs, err
		} else if t.IsVariable() {
			subs := MakeSubstitutions()
			err := subs.Update(t.Variable, o, ctx.resolver)
			return subs, err
		}
		return Substitutions{}, UnificationError(o, t)
	}
	if t.IsVariable() {
		subs := MakeSubstitutions()
		subs.Update(t.Variable, o, ctx.resolver)
		return subs, nil
	}
	if o.IsConstructor() && t.IsConstructor() {
		if o.Constructor.Name != t.Constructor.Name {
			return Substitutions{}, UnificationError(o, t)
		}
		return Substitutions{}, nil
	}
	if o.IsConstructor() {
		return unifier(o, t, ctx)
	}
	if t.IsConstructor() && o.IsHVariable() {
		inst := t.Constructor.Underlying.Instantiate(o.HVariable.Params, ctx.resolver)
		subs := MakeSubstitutions()
		subs.Update(o.HVariable.Root, inst, ctx.resolver)
		return subs, nil
	}
	if !o.IsDefined() {
		return Substitutions{}, nil
	}
	return Substitutions{}, UnificationError(o, t)
}

func unifyStructureWithStructuralVar(v Type, s Type, ctx *unificationCtx) (Substitutions, error) {
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
		sub, err := unifier(t, sv, ctx)
		if err != nil {
			return Substitutions{}, err
		}
		err = tres.Combine(sub, ctx.resolver)
		if err != nil {
			return Substitutions{}, err
		}
	}
	tres.Update(v.Variable, s, ctx.resolver)
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
			subs.Update(t.Variable, prim, ctx.resolver)
			subs.Update(o.Variable, prim, ctx.resolver)
			return subs, err
		}
		rv := MakeUnionVar(resolv...)
		subs := MakeSubstitutions()
		subs.Update(t.Variable, rv, ctx.resolver)
		subs.Update(o.Variable, rv, ctx.resolver)
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
				s, err := unifier(res[k], v, ctx)
				if err != nil {
					return Substitutions{}, err
				}
				err = subs.Combine(s, ctx.resolver)
				if err != nil {
					return Substitutions{}, err
				}
			}
			res[k] = v
		}
		rv := MakeStructuralVar(res)
		subs.Update(t.Variable, rv, ctx.resolver)
		subs.Update(o.Variable, rv, ctx.resolver)
		return subs, nil
	}
	subs := MakeSubstitutions()
	subs.Update(o.Variable, t, ctx.resolver)
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
		err = result.Combine(s, ctx.resolver)
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
		err = result.Combine(s, ctx.resolver)
		if err != nil {
			return Substitutions{}, err
		}
	}
	return result, nil
}
