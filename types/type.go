package types

import (
	"errors"
	"strconv"
)

type TypeBase = string

type TypeDef struct {
	Base *TypeBase
	Fn   []*Type
}

type Type struct {
	Def *TypeDef
}

var (
	Int  *Type = base("int")
	Bool *Type = base("bool")
)

func base(t string) *Type {
	return &Type{Def: &TypeDef{Base: &t}}
}

func function(ts ...*Type) *Type {
	return &Type{&TypeDef{Fn: ts}}
}

func variable() *Type {
	return &Type{Def: &TypeDef{}}
}

func (t *Type) copy() *Type {
	var params []*Type = nil
	var def *TypeDef = nil
	if t.Def != nil {
		if t.Def.Fn != nil {
			params = make([]*Type, len(t.Def.Fn))
			var seen map[*Type]*Type = map[*Type]*Type{}
			for i, p := range t.Def.Fn {
				if seen[p] == nil {
					seen[p] = p.copy()
				}
				params[i] = seen[p]
			}
		}
		def = &TypeDef{
			Fn:   params,
			Base: t.Def.Base,
		}
	}

	return &Type{Def: def}
}

func (t *Type) isFunction() bool {
	return t.Def.Fn != nil
}

func (t *Type) isBase() bool {
	return t.Def.Base != nil
}

func (t *Type) isVariable() bool {
	return !t.isBase() && !t.isFunction()
}

func (t *Type) returnType() *Type {
	return t.Def.Fn[len(t.Def.Fn)-1]
}

func unify(a *Type, b *Type) error {
	if a.isVariable() && b.isVariable() {
		a.Def = b.Def
		return nil
	}
	if a.isVariable() {
		a.Def.Fn = b.Def.Fn
		a.Def.Base = b.Def.Base
		return nil
	}
	if b.isVariable() {
		b.Def.Fn = a.Def.Fn
		b.Def.Base = a.Def.Base
		return nil
	}
	if a.isBase() && b.isBase() {
		if *(a.Def.Base) != *(b.Def.Base) {
			return errors.New("can not unify " + *(a.Def.Base) + " with " + *(b.Def.Base))
		}
		return nil
	}
	if a.isFunction() && b.isFunction() {
		if len(a.Def.Fn) != len(b.Def.Fn) {
			return errors.New("wrong number of function arguments " + strconv.Itoa(len(a.Def.Fn)) + "given " + strconv.Itoa(len(b.Def.Fn)) + "required")
		}
		for i := range a.Def.Fn {
			err := unify(a.Def.Fn[i], b.Def.Fn[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
	if a.isFunction() {
		return errors.New("not a function")
	}
	if b.isFunction() {
		return errors.New("a function required")
	}
	panic("missing unification rule")
}
