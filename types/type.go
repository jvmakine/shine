package types

import (
	"errors"
	"strconv"
)

type TypeBase = string

type Type struct {
	Base *TypeBase
	Fn   []*Type
}

var (
	Int  *Type = base("int")
	Bool *Type = base("bool")
)

func base(t string) *Type {
	return &Type{Base: &t}
}

func function(ts ...*Type) *Type {
	return &Type{Fn: ts}
}

func variable() *Type {
	return &Type{}
}

func (t *Type) isFunction() bool {
	return t.Fn != nil
}

func (t *Type) isBase() bool {
	return t.Base != nil
}

func (t *Type) isVariable() bool {
	return !t.isBase() && !t.isFunction()
}

func (t *Type) returnType() *Type {
	return t.Fn[len(t.Fn)-1]
}

func unify(a **Type, b **Type) error {
	if (*a).isVariable() && (*b).isVariable() {
		*a = *b
		return nil
	}
	if (*a).isVariable() {
		(*a).Fn = (*b).Fn
		(*a).Base = (*b).Base
		return nil
	}
	if (*b).isVariable() {
		(*b).Fn = (*a).Fn
		(*b).Base = (*a).Base
		return nil
	}
	if (*a).isBase() && (*b).isBase() {
		if *((*a).Base) != *((*b).Base) {
			return errors.New("can not unify " + *((*a).Base) + " with " + *((*b).Base))
		}
		return nil
	}
	if (*a).isFunction() && (*b).isFunction() {
		if len((*a).Fn) != len((*b).Fn) {
			return errors.New("wrong number of function arguments " + strconv.Itoa(len((*a).Fn)) + "given " + strconv.Itoa(len((*b).Fn)) + "required")
		}
		for i := range (*a).Fn {
			err := unify(&(*a).Fn[i], &(*b).Fn[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
	if (*a).isFunction() {
		return errors.New("not a function")
	}
	if (*b).isFunction() {
		return errors.New("a function required")
	}
	panic("missing unification rule")
}
