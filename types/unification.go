package types

import (
	"errors"
)

type unificationCtx struct {
	seenIDs map[VariableID]bool
}

func NewUnificationCtx() *unificationCtx {
	return &unificationCtx{map[VariableID]bool{}}
}

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
	sub, err := unifier(t, o, &unificationCtx{seenIDs: map[VariableID]bool{}})
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
	return Convert(t, sub), nil
}

func Convert(t Type, s Substitutions) Type {
	return t.convert(s, &unificationCtx{seenIDs: map[VariableID]bool{}})
}

func unifier(t Type, o Type, ctx *unificationCtx) (Substitutions, error) {
	sub, err := t.unifier(o, ctx)
	if err != nil {
		sub, err = o.unifier(t, ctx)
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
