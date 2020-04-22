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
		s1 := MakeRestricted(r...).Signature()
		s2 := MakeRestricted(o...).Signature()
		return nil, errors.New("can not unify " + s1 + " with " + s2)
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
