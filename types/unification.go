package types

import (
	"errors"
)

func UnificationError(a Type, b Type) error {
	sa := Signature(a)
	sb := Signature(b)
	if sa < sb {
		return errors.New("can not unify " + sa + " with " + sb)
	} else {
		return errors.New("can not unify " + sb + " with " + sa)
	}
}

func Unifier(t Type, o Type) (Substitutions, error) {
	sub, err := unifier(t, o)
	if err != nil {
		return MakeSubstitutions(), err
	}
	return sub, nil
}

func Unify(t Type, o Type) (Type, error) {
	sub, err := Unifier(t, o)
	if err != nil {
		return nil, err
	}
	return t.Convert(sub), nil
}

func unifier(t Type, o Type) (Substitutions, error) {
	sub, err := t.unifier(o)
	if err != nil {
		sub, err = o.unifier(t)
		if err != nil {
			return MakeSubstitutions(), err
		}
	}
	return sub, nil
}

func UnifiesWith(t Type, o Type) bool {
	_, e := Unifier(t, o)
	return e == nil
}
