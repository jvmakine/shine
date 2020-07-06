package types

import "errors"

type Union []Type

func (u Union) deduplicate() Union {
	res := Union{}
	for _, tt := range u {
		if len(res) == 0 {
			res = Union{tt}
		} else {
			shouldAppend := false
			for i, rt := range res {
				if tt.IsGeneralisationOf(rt) {
					res = append(res[:i], res[i+1:]...)
					shouldAppend = true
				} else if !rt.IsGeneralisationOf(tt) {
					shouldAppend = true
				}
			}
			if shouldAppend {
				res = append(res, tt)
			}
		}
	}
	return res
}

func (r Union) Unify(o Union) (Union, error) {
	union := Union{}
	for _, rp := range r {
		for _, op := range o {
			un, err := rp.Unify(op)
			if err == nil {
				union = append(union, un)
			}
		}
	}
	res := union.deduplicate()
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
