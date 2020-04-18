package inferer

import (
	"errors"
	"strconv"

	"github.com/jvmakine/shine/types"
)

type Subs struct {
	Variables map[*types.TypeDef]*types.TypeDef
}

func (s Subs) replace(from *types.TypeDef, to *types.TypeDef, mirror Subs) error {
	err := doesConflict(s.Variables[from], to)
	if err != nil {
		return err
	}
	for k, v := range s.Variables {
		if v == from {
			s.Variables[k] = to
		}
	}
	if s.Variables[from] != nil {
		err := doesConflict(to, mirror.Variables[s.Variables[from]])
		if err != nil {
			return err
		}
	}
	if s.Variables[from] == nil || !to.IsVariable() {
		s.Variables[from] = to
	}
	return nil
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

func unionMatch(ad *types.TypeDef, bd *types.TypeDef) bool {
	as := map[string]bool{}
	hits := 0
	for _, p := range ad.Bases {
		as[*p] = true
	}
	for _, p := range bd.Bases {
		if as[*p] {
			hits++
		}
	}
	return hits == len(ad.Bases) && hits == len(bd.Bases)
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

func (u *Unifier) unify(a *types.TypePtr, b *types.TypePtr) error {
	if err := doesConflict(a.Def, b.Def); err != nil {
		return err
	} else if a.IsVariable() && b.IsVariable() {
		if err := u.source.replace(a.Def, b.Def, u.dest); err != nil {
			return err
		}
		if err := u.dest.replace(b.Def, a.Def, u.source); err != nil {
			return err
		}
	} else if a.IsVariable() {
		if err := u.source.replace(a.Def, b.Def, u.dest); err != nil {
			return err
		}
		if err := u.dest.replace(a.Def, b.Def, u.source); err != nil {
			return err
		}
	} else if b.IsVariable() {
		if err := u.dest.replace(b.Def, a.Def, u.source); err != nil {
			return err
		}
	} else if a.IsFunction() && b.IsFunction() {
		if len(a.Def.Fn) != len(b.Def.Fn) {
			return errors.New("wrong number of function arguments " + strconv.Itoa(len(a.Def.Fn)) + " given " + strconv.Itoa(len(b.Def.Fn)) + " required")
		}
		for i := range a.Def.Fn {
			if err := u.unify(a.Def.Fn[i], b.Def.Fn[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func Unify(a *types.TypePtr, b *types.TypePtr) (*Unifier, error) {
	uni := NewUnifier()
	err := uni.unify(a, b)
	if err != nil {
		return nil, err
	}
	return uni, err
}
