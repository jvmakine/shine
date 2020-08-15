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

func Unify(t Type, o Type, ctx UnificationCtx) (Type, error) {
	sub, err := Unifier(t, o, ctx)
	if err != nil {
		return nil, err
	}
	conv, _ := t.convert(sub, newSubstCtx())
	return conv, nil
}

func Unifier(t Type, o Type, ctx UnificationCtx) (Substitutions, error) {
	return unifier(t, o, ctx, newSubstCtx())
}

func unifier(t Type, o Type, ctx UnificationCtx, sctx substitutionCtx) (Substitutions, error) {
	sub, err := t.unifier(o, ctx, sctx)
	if err != nil {
		sub, err = o.unifier(t, ctx, sctx)
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
