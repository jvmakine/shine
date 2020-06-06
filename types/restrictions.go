package types

import "errors"

type Union []Primitive

func (r Union) Resolve(o Union) (Union, error) {
	res := Union{}
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
		return nil, UnificationError(MakeUnionVar(r...), MakeUnionVar(o...))
	}
	return res, nil
}

func (r Union) Unifies(o Primitive) error {
	for _, r := range r {
		if r == o {
			return nil
		}
	}
	sig := MakeUnionVar(r...).Signature()
	return errors.New("can not unify " + o + " with " + sig)
}
