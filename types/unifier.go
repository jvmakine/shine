package types

import (
	"errors"
	"strconv"
)

type Subs struct {
	Variables map[*TypePtr]*TypePtr
}

func apply(t *TypePtr, s *Subs, d *Subs) {
	if fv := s.Variables[t]; fv != nil {
		if dv := d.Variables[fv]; dv != nil {
			t.Def = dv.Def
		} else {
			if fv.IsVariable() {
				panic("invalid substitution")
			}
			t.Def = fv.Def
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

func (u *Unifier) ApplySource(t *TypePtr) {
	apply(t, &u.source, &u.dest)
}

func (u *Unifier) ApplyDest(t *TypePtr) {
	apply(t, &u.dest, &u.source)
}

func NewUnifier() *Unifier {
	return &Unifier{
		source: Subs{Variables: map[*TypePtr]*TypePtr{}},
		dest:   Subs{Variables: map[*TypePtr]*TypePtr{}},
	}
}

func doesConflict(x *TypePtr, y *TypePtr) error {
	if x == nil || y == nil || x.Def.Base == nil || y.Def.Base == nil || x.Def == nil || y.Def == nil {
		return nil
	}
	a := *(x.Def.Base)
	b := *(y.Def.Base)
	if a != b {
		return errors.New("can not unify " + b + " with " + a)
	}
	return nil
}

func (u *Unifier) combine(o *Unifier) error {
	for k, v := range o.source.Variables {
		if u.source.Variables[k] == nil {
			u.source.Variables[k] = v
		} else if err := doesConflict(u.source.Variables[k], v); err != nil {
			return err
		} else if v.IsBase() {
			if u.source.Variables[k] != nil {
				if err := doesConflict(u.dest.Variables[u.source.Variables[k]], v); err != nil {
					return err
				}
			}
			u.source.Variables[k] = v
		}
	}
	for k, v := range o.dest.Variables {
		if u.dest.Variables[k] == nil {
			u.dest.Variables[k] = v
		} else if err := doesConflict(u.dest.Variables[k], v); err != nil {
			return err
		} else if v.IsBase() {
			if u.dest.Variables[k] != nil {
				if err := doesConflict(u.source.Variables[u.dest.Variables[k]], v); err != nil {
					return err
				}
			}
			u.dest.Variables[k] = v
		}
	}
	return nil
}

func Unify(a *TypePtr, b *TypePtr) (*Unifier, error) {
	if a.IsVariable() && b.IsVariable() {
		u := NewUnifier()
		u.source.Variables[a] = b
		u.dest.Variables[b] = a
		return u, nil
	}
	if a.IsVariable() {
		u := NewUnifier()
		u.source.Variables[a] = b
		return u, nil
	}
	if b.IsVariable() {
		u := NewUnifier()
		u.dest.Variables[b] = a
		return u, nil
	}
	if a.IsBase() && b.IsBase() {
		err := doesConflict(a, b)
		if err != nil {
			return nil, err
		}
		return NewUnifier(), nil
	}
	if a.IsFunction() && b.IsFunction() {
		if len(a.Def.Fn) != len(b.Def.Fn) {
			return nil, errors.New("wrong number of function arguments " + strconv.Itoa(len(a.Def.Fn)) + "given " + strconv.Itoa(len(b.Def.Fn)) + "required")
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
