package types

import (
	"errors"
	"strconv"
	"strings"
)

type TypeSignature = string

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

func (t *Type) Copy() *Type {
	var params []*Type = nil
	var def *TypeDef = nil
	if t.Def != nil {
		if t.Def.Fn != nil {
			params = make([]*Type, len(t.Def.Fn))
			var seen map[*Type]*Type = map[*Type]*Type{}
			for i, p := range t.Def.Fn {
				if seen[p] == nil {
					seen[p] = p.Copy()
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

func (t *Type) IsFunction() bool {
	return t.Def.Fn != nil
}

func (t *Type) IsBase() bool {
	return t.Def.Base != nil
}

func (t *Type) IsVariable() bool {
	return !t.IsBase() && !t.IsFunction()
}

func (t *Type) ReturnType() *Type {
	return t.Def.Fn[len(t.Def.Fn)-1]
}

func (t *Type) Signature() (TypeSignature, error) {
	if t.IsBase() {
		return *t.Def.Base, nil
	}
	if t.IsFunction() {
		var sb strings.Builder
		sb.WriteString("(")
		for i, p := range t.Def.Fn {
			s, err := p.Signature()
			if err != nil {
				return "", err
			}
			sb.WriteString(s)
			if i < len(t.Def.Fn)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		return sb.String(), nil
	}
	return "", errors.New("can not get signature for a type variable")
}

func Unify(a *Type, b *Type) error {
	if a.IsVariable() && b.IsVariable() {
		a.Def = b.Def
		return nil
	}
	if a.IsVariable() {
		a.Def.Fn = b.Def.Fn
		a.Def.Base = b.Def.Base
		return nil
	}
	if b.IsVariable() {
		b.Def.Fn = a.Def.Fn
		b.Def.Base = a.Def.Base
		return nil
	}
	if a.IsBase() && b.IsBase() {
		if *(a.Def.Base) != *(b.Def.Base) {
			return errors.New("can not unify " + *(a.Def.Base) + " with " + *(b.Def.Base))
		}
		return nil
	}
	if a.IsFunction() && b.IsFunction() {
		if len(a.Def.Fn) != len(b.Def.Fn) {
			return errors.New("wrong number of function arguments " + strconv.Itoa(len(a.Def.Fn)) + "given " + strconv.Itoa(len(b.Def.Fn)) + "required")
		}
		for i := range a.Def.Fn {
			err := Unify(a.Def.Fn[i], b.Def.Fn[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
	if a.IsFunction() {
		return errors.New("not a function")
	}
	if b.IsFunction() {
		return errors.New("a function required")
	}
	panic("missing unification rule")
}
