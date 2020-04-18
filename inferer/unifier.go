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
			if fv.Bases == nil && fv.Fn == nil {
				panic("invalid substitution")
			}
			t.Def.Fn = fv.Fn
			t.Def.Bases = fv.Bases
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

func unionConflict(a *types.TypeDef, b *types.TypeDef) bool {
	return len(unifyUnion(a, b).Bases) == 0
}

func doesConflict(x *types.TypeDef, y *types.TypeDef) error {
	if x == nil || y == nil {
		return nil
	} else if (x.Fn == nil) != (y.Fn == nil) {
		return errors.New("a function required")
	} else if x.Bases != nil && y.Bases != nil {
		conflicts := unionConflict(x, y)
		a := x.Signature()
		b := y.Signature()
		if conflicts {
			return errors.New("can not unify " + a + " with " + b)
		}
	}
	return nil
}

func combineSub(source Subs, dest Subs, from Subs) error {
	for k, v := range from.Variables {
		us := source.Variables[k]
		if err := doesConflict(us, v); err != nil {
			return err
		} else if us == nil || (us.IsVariable() && v.IsVariable()) {
			source.Variables[k] = v
			if us != nil {
				dest.Variables[us] = k
			}
		} else if v.IsPrimitive() {
			if us != nil {
				if err := doesConflict(dest.Variables[us], v); err != nil {
					return err
				}
			}
			dest.Variables[us] = v
			source.Variables[k] = v
		}
	}
	return nil
}

func (u *Unifier) combine(o *Unifier) error {
	if err := combineSub(u.source, u.dest, o.source); err != nil {
		return err
	}
	return combineSub(u.dest, u.source, o.dest)
}

func unifyUnion(ad *types.TypeDef, bd *types.TypeDef) *types.TypeDef {
	a := ad.Bases
	b := bd.Bases
	as := map[string]bool{}
	for _, p := range a {
		as[*p] = true
	}
	res := []*types.Primitive{}
	for _, p := range b {
		if as[*p] {
			res = append(res, p)
		}
	}
	if len(res) == len(a) {
		return ad
	}
	return &types.TypeDef{Bases: res}
}

func Unify(a *types.TypePtr, b *types.TypePtr) (*Unifier, error) {
	err := doesConflict(a.Def, b.Def)
	if err != nil {
		return nil, err
	}
	if a.IsVariable() && b.IsVariable() {
		u := NewUnifier()
		if a.IsUnion() && b.IsUnion() {
			u.source.Variables[a.Def] = unifyUnion(b.Def, a.Def)
			u.dest.Variables[b.Def] = unifyUnion(a.Def, b.Def)
		} else {
			u.source.Variables[a.Def] = b.Def
			u.dest.Variables[b.Def] = a.Def
		}
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
	return NewUnifier(), nil
}
