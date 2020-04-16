package inferer

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/types"
)

type Subs struct {
	Variables map[*types.TypeDef]*types.TypeDef
}

func apply(t *types.TypePtr, s *Subs, d *Subs) {
	if fv := s.Variables[t.Def]; fv != nil {
		if dv := d.Variables[fv]; dv != nil {
			t.Def = dv
		} else {
			if fv.Base == nil && fv.Fn == nil {
				panic("invalid substitution")
			}
			t.Def.Fn = fv.Fn
			t.Def.Base = fv.Base
		}

	} else if t.IsFunction() {
		for _, a := range t.Def.Fn {
			apply(a, s, d)
		}
	}
}

type Unifier struct {
	source Subs
	dest   Subs
}

func (u *Unifier) ApplySource(t *types.TypePtr) {
	apply(t, &u.source, &u.dest)
}

func (u *Unifier) ApplyDest(t *types.TypePtr) {
	apply(t, &u.dest, &u.source)
}

func NewUnifier() *Unifier {
	return &Unifier{
		source: Subs{Variables: map[*types.TypeDef]*types.TypeDef{}},
		dest:   Subs{Variables: map[*types.TypeDef]*types.TypeDef{}},
	}
}

func doesConflict(x *types.TypeDef, y *types.TypeDef) error {
	if x == nil || y == nil || x.Base == nil || y.Base == nil {
		return nil
	}
	a := *(x.Base)
	b := *(y.Base)
	if a != b {
		return errors.New("can not unify " + b + " with " + a)
	}
	return nil
}

func (u *Unifier) combine(o *Unifier) error {
	for k, v := range o.source.Variables {
		us := u.source.Variables[k]
		if err := doesConflict(us, v); err != nil {
			return err
		} else if us == nil || (us.Fn == nil && us.Base == nil && v.Base == nil && v.Fn == nil) {
			u.source.Variables[k] = v
			if us != nil {
				u.dest.Variables[us] = k
			}
		} else if v.Base != nil {
			if us != nil {
				if err := doesConflict(u.dest.Variables[us], v); err != nil {
					return err
				}
			}
			u.dest.Variables[us] = v
			u.source.Variables[k] = v
		}
	}
	for k, v := range o.dest.Variables {
		ud := u.dest.Variables[k]
		if err := doesConflict(ud, v); err != nil {
			return err
		} else if ud == nil || (ud.Fn == nil && ud.Base == nil && v.Base == nil && v.Fn == nil) {
			u.dest.Variables[k] = v
			if ud != nil {
				u.source.Variables[ud] = k
			}
		} else if v.Base != nil {
			if u.dest.Variables[k] != nil {
				if err := doesConflict(u.source.Variables[ud], v); err != nil {
					return err
				}
			}
			u.source.Variables[ud] = v
			u.dest.Variables[k] = v
		}
	}
	return nil
}

func Unify(a *types.TypePtr, b *types.TypePtr) (*Unifier, error) {
	if a.IsVariable() && b.IsVariable() {
		u := NewUnifier()
		u.source.Variables[a.Def] = b.Def
		u.dest.Variables[b.Def] = a.Def
		return u, nil
	}
	if a.IsVariable() {
		u := NewUnifier()
		u.source.Variables[a.Def] = b.Def
		return u, nil
	}
	if b.IsVariable() {
		u := NewUnifier()
		u.dest.Variables[b.Def] = a.Def
		return u, nil
	}
	if a.IsBase() && b.IsBase() {
		err := doesConflict(a.Def, b.Def)
		if err != nil {
			return nil, err
		}
		return NewUnifier(), nil
	}
	if a.IsFunction() && b.IsFunction() {
		if len(a.Def.Fn) != len(b.Def.Fn) {
			return nil, errors.New("wrong number of function arguments " + strconv.Itoa(len(a.Def.Fn)) + " given " + strconv.Itoa(len(b.Def.Fn)) + " required")
		}
		unifier := NewUnifier()
		for i := range a.Def.Fn {
			u, err := Unify(a.Def.Fn[i], b.Def.Fn[i])
			if err != nil {
				return nil, err
			}
			err = unifier.combine(u)
			if err != nil {
				return nil, err
			}
		}
		return unifier, nil
	}
	if a.IsFunction() {
		return nil, errors.New("not a function")
	}
	if b.IsFunction() {
		return nil, errors.New("a function required")
	}
	panic("missing unification rule")
}
