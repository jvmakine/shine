package types

import "errors"

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
	return sub.Apply(t), nil
}

func (t Type) Unifier(o Type) (Substitutions, error) {
	if o.IsPrimitive() && t.IsPrimitive() && *o.Primitive != *t.Primitive {
		return Substitutions{}, UnificationError(o, t)
	}
	if (o.IsPrimitive() && t.IsFunction()) || (o.IsFunction() && t.IsPrimitive()) {
		return Substitutions{}, UnificationError(o, t)
	}
	if t.IsVariable() && o.IsVariable() {
		if o.IsRestrictedVariable() && !t.IsRestrictedVariable() {
			return o.Unifier(t)
		} else if t.IsRestrictedVariable() && o.IsRestrictedVariable() {
			resolv, err := t.Variable.Restrictions.Resolve(o.Variable.Restrictions)
			if len(resolv) == 1 {
				prim := MakePrimitive(resolv[0])
				return Substitutions{t.Variable: prim, o.Variable: prim}, err
			}
			rv := MakeRestricted(resolv...)
			return Substitutions{t.Variable: rv, o.Variable: rv}, err
		}
		return Substitutions{o.Variable: t}, nil
	}
	if o.IsVariable() && !t.IsVariable() {
		return o.Unifier(t)
	}
	if t.IsVariable() && o.IsFunction() {
		return Substitutions{t.Variable: o}, nil
	}
	if o.IsFunction() && t.IsFunction() {
		op := o.FunctTypes()
		tp := t.FunctTypes()
		if len(op) != len(tp) {
			return Substitutions{}, UnificationError(o, t)
		}
		result := Substitutions{}
		for i, p := range op {
			s, err := p.Unifier(tp[i])
			if err != nil {
				return Substitutions{}, err
			}
			result, err = result.Combine(s)
			if err != nil {
				return Substitutions{}, err
			}
		}
		return result, nil
	}
	if o.IsPrimitive() {
		if t.IsRestrictedVariable() {
			err := t.Variable.Restrictions.Unifies(*o.Primitive)
			return Substitutions{t.Variable: o}, err
		} else if t.IsVariable() {
			return Substitutions{t.Variable: o}, nil
		}
		return Substitutions{}, nil
	}
	return Substitutions{}, nil
}
