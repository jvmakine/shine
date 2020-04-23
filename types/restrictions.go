package types

import "errors"

type Restrictions []Primitive

func (r Restrictions) Resolve(o Restrictions) (Restrictions, error) {
	res := Restrictions{}
	found := map[Primitive]bool{}
	for _, p := range o {
		found[p] = true
	}
	for _, p := range r {
		if found[p] {
			res = append(res, p)
		}
	}
	if len(res) == 0 {
		return nil, UnificationError(MakeRestricted(r...), MakeRestricted(o...))
	}
	return res, nil
}

func (r Restrictions) Unifies(o Primitive) error {
	for _, r := range r {
		if r == o {
			return nil
		}
	}
	sig := MakeRestricted(r...).Signature()
	return errors.New("can not unify " + o + " with " + sig)
}
