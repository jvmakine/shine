package typeinference

import (
	. "github.com/jvmakine/shine/types"
)

type TLens struct {
	substitutions map[*TypeVar]Type
	references    map[*TypeVar]map[*TypeVar]bool
}

func MakeTLens() *TLens {
	return &TLens{substitutions: map[*TypeVar]Type{}, references: map[*TypeVar]map[*TypeVar]bool{}}
}

func (l *TLens) Substitutions() Substitutions {
	return l.substitutions
}

func (l *TLens) Convert(t Type) Type {
	return l.Substitutions().Apply(t)
}

func (l *TLens) Update(from *TypeVar, to Type) error {
	result := l.Convert(to)

	if p := l.substitutions[from]; p.IsDefined() {
		uni, err := result.Unifier(p)
		if err != nil {
			return err
		}
		for f, t := range uni {
			l.Update(f, t)
		}
		result = uni.Apply(result)
	}

	l.substitutions[from] = result

	for _, fv := range result.FreeVars() {
		if l.references[fv] == nil {
			l.references[fv] = map[*TypeVar]bool{}
		}
		l.references[fv][from] = true
	}

	if rs := l.references[from]; rs != nil {
		l.references[from] = nil
		s := Substitutions{from: result}
		for k := range rs {
			substit := l.substitutions[k]
			l.substitutions[k] = s.Apply(substit)
			for _, fv := range l.substitutions[k].FreeVars() {
				if l.references[fv] == nil {
					l.references[fv] = map[*TypeVar]bool{}
				}
				l.references[fv][from] = true
			}
		}
	}
	return nil
}
