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

func Unifier(t Type, o Type, ctx UnificationCtx) (Substitutions, error) {
	sub, err := unifier(t, o, ctx)
	if err != nil {
		return MakeSubstitutions(), err
	}
	return sub, nil
}

func Unify(t Type, o Type, ctx UnificationCtx) (Type, error) {
	sub, err := Unifier(t, o, ctx)
	if err != nil {
		return nil, err
	}
	conv, _ := t.Convert(sub)
	return conv, nil
}

func unifier(t Type, o Type, ctx UnificationCtx) (Substitutions, error) {
	sub, err := t.unifier(o, ctx)
	if err != nil {
		sub, err = o.unifier(t, ctx)
		if err != nil {
			return MakeSubstitutions(), err
		}
	}
	return sub, nil
}

func UnifiesWith(t Type, o Type, ctx UnificationCtx) bool {
	_, e := Unifier(t, o, ctx)
	return e == nil
}
