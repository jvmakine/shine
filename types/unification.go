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
	if o.IsPrimitive() && t.IsPrimitive() && *o.Primitive != *t.Primitive {
		return t, UnificationError(o, t)
	}
	if t.IsVariable() && o.IsVariable() {
		if o.IsRestrictedVariable() && !t.IsRestrictedVariable() {
			return o.Unify(t)
		} else if t.IsRestrictedVariable() && o.IsRestrictedVariable() {
			resolv, err := t.Variable.Restrictions.Resolve(o.Variable.Restrictions)
			if len(resolv) == 1 {
				return MakePrimitive(resolv[0]), err
			}
			return MakeRestricted(resolv...), err
		} else if t.IsFunction() && !o.IsFunction() {
			return o.Unify(t)
		} else if t.IsRestrictedVariable() && o.IsFunction() {
			return o, UnificationError(t, o)
		} else if o.IsFunction() && !t.IsFunction() {
			return o, nil
		} else if t.IsFunction() && o.IsFunction() {
			ot := o.FunctTypes()
			tt := t.FunctTypes()
			if len(ot) != len(tt) {
				return o, UnificationError(o, t)
			}
			unified := make([]Type, len(ot))
			for i, p := range tt {
				u, err := p.Unify(ot[i])
				if err != nil {
					return o, err
				}
				unified[i] = u
			}
			return MakeFunction(unified...), nil
		}
		return t, nil
	}
	if o.IsVariable() && !t.IsVariable() {
		return o.Unify(t)
	}
	if o.IsFunction() && t.IsFunction() {
		op := o.FunctTypes()
		tp := t.FunctTypes()
		if len(op) != len(tp) {
			return o, UnificationError(o, t)
		}
		for i, p := range op {
			_, err := p.Unify(tp[i])
			if err != nil {
				return o, err
			}
		}
		return o, nil
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
