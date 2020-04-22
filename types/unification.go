package types

import "errors"

func (t Type) Unify(o Type) (Type, error) {
	if o.IsPrimitive() && t.IsPrimitive() && *o.Primitive != *t.Primitive {
		op := *o.Primitive
		tp := *t.Primitive
		if op < tp {
			return t, errors.New("can not unify " + op + " with " + tp)
		}
		return t, errors.New("can not unify " + tp + " with " + op)
	}
	if t.IsVariable() && o.IsVariable() {
		if t.IsRestrictedVariable() && !o.IsRestrictedVariable() {
			return o.Unify(t)
		} else if t.IsRestrictedVariable() && o.IsRestrictedVariable() {
			resolv, err := t.Variable.Restrictions.Resolve(o.Variable.Restrictions)
			if len(resolv) == 1 {
				return MakePrimitive(resolv[0]), err
			}
			return Type{Variable: &TypeVar{resolv}}, err
		}
		return t, nil
	}
	if o.IsVariable() && !t.IsVariable() {
		return o.Unify(t)
	}
	if o.IsFunction() && t.IsFunction() {
		panic("functions as parameters not upported yet")
	}
	if o.IsPrimitive() {
		if t.IsRestrictedVariable() {
			err := t.Variable.Restrictions.Unifies(*o.Primitive)
			return o, err
		}
		return o, nil
	}
	if o.IsFunction() {
		return o, nil
	}
	return t, nil
}
