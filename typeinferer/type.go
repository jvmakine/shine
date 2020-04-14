package typeinferer

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/typedef"
)

type TypeDef struct {
	Base *typedef.Primitive
	Fn   []*TypePtr
}

type TypePtr struct {
	Def *TypeDef
}

func base(t typedef.Primitive) *TypePtr {
	return &TypePtr{Def: &TypeDef{Base: &t}}
}

func function(ts ...*TypePtr) *TypePtr {
	return &TypePtr{&TypeDef{Fn: ts}}
}

func variable() *TypePtr {
	return &TypePtr{Def: &TypeDef{}}
}

func (t *TypePtr) Copy() *TypePtr {
	var params []*TypePtr = nil
	var def *TypeDef = nil
	if t.Def != nil {
		if t.Def.Fn != nil {
			params = make([]*TypePtr, len(t.Def.Fn))
			var seen map[*TypePtr]*TypePtr = map[*TypePtr]*TypePtr{}
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

	return &TypePtr{Def: def}
}

func (t *TypePtr) IsFunction() bool {
	return t.Def.Fn != nil
}

func (t *TypePtr) IsBase() bool {
	return t.Def.Base != nil
}

func (t *TypePtr) IsVariable() bool {
	return !t.IsBase() && !t.IsFunction()
}

func (t *TypePtr) ReturnType() *TypePtr {
	return t.Def.Fn[len(t.Def.Fn)-1]
}

func Unify(a *TypePtr, b *TypePtr) error {
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
