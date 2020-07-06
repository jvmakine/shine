package types

import "errors"

type Union []Type

func (r Union) Unify(o Union) (Union, error) {
	res := Union{}
	for _, rp := range r {
		for _, op := range o {
			un, err := rp.Unify(op)
			if err == nil {
				res = append(res, un)
			}
		}
	}
	if len(res) == 0 {
		return nil, UnificationError(MakeUnionVar(r...), MakeUnionVar(o...))
	}
	return res, nil
}

func (r Union) Unifies(o Primitive) error {
	for _, r := range r {
		if r.IsPrimitive() && *r.Primitive == o {
			return nil
		}
	}
	sig := MakeUnionVar(r...).Signature()
	return errors.New("can not unify " + o + " with " + sig)
}
