package inferer

import (
	"errors"
	"strconv"

	. "github.com/jvmakine/shine/types"
)

type Subs struct {
	Variables map[*TypeVar]Type
}

type Unifier struct {
	graph   map[*TypeVar][]Type
	replace map[*TypeVar]Type
}

func NewUnifier() *Unifier {
	return &Unifier{}
}

func doesConflict(x Type, y Type) error {
	if x.IsPrimitive() && y.IsPrimitive() && *x.Primitive != *y.Primitive {
		return errors.New("can not unify " + *x.Primitive + " with " + *y.Primitive)
	}
	return nil
}

func (u *Unifier) buildGraph(a Type, b Type) error {
	u.graph = map[*TypeVar][]Type{}
	if err := u.addToGraph(a, b); err != nil {
		return err
	}
	return nil
}

func (u *Unifier) addToGraph(a Type, b Type) error {
	if err := doesConflict(a, b); err != nil {
		return err
	}
	if a.IsVariable() {
		u.graph[a.Variable] = append(u.graph[a.Variable], b)
		return nil
	}
	if b.IsVariable() {
		u.graph[b.Variable] = append(u.graph[b.Variable], a)
		return nil
	}
	// TODO: functions to variables unification
	if a.IsFunction() && b.IsFunction() {
		al := len(*a.Function)
		bl := len(*b.Function)
		if al != bl {
			return errors.New("wrong number of function arguments: " + strconv.Itoa(al) + " != " + strconv.Itoa(bl))
		}
		for i := range *a.Function {
			err := u.addToGraph((*a.Function)[i], (*b.Function)[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (u *Unifier) traverse(a *TypeVar) []Type {
	todo := []*TypeVar{a}
	result := []Type{}
	inResult := map[Type]bool{}
	visited := map[*TypeVar]bool{}
	for len(todo) > 0 {
		next := todo[0]
		todo = todo[1:]
		if !visited[next] {
			visited[next] = true
			r := Type{Variable: next}
			if !inResult[r] {
				result = append(result, r)
				inResult[r] = true
			}
			for _, f := range u.graph[next] {
				if f.IsVariable() {
					todo = append(todo, f.Variable)
				} else {
					if !inResult[f] {
						result = append(result, f)
						inResult[f] = true
					}
				}
			}
		}
	}
	return result
}

func solveRestrictions(ts []Type) (Type, error) {
	var ps []Primitive = nil
	for _, t := range ts {
		if t.IsPrimitive() {
			if ps == nil || len(ps) == 0 {
				ps = []Primitive{*t.Primitive}
			} else {
				found := false
				for _, p := range ps {
					if p == *t.Primitive {
						found = true
						break
					}
				}
				if !found {
					sig := ps[0]
					if len(ps) > 1 {
						sig = Type{Variable: &TypeVar{ps}}.Signature()
					}
					return Type{}, errors.New("can not unify " + *t.Primitive + " with " + sig)
				}
				ps = []Primitive{*t.Primitive}
			}
		} else if t.IsVariable() && len(t.Variable.Restrictions) > 0 {
			if ps == nil || len(ps) == 0 {
				ps = t.Variable.Restrictions
			} else {
				res := []Primitive{}
				found := map[Primitive]bool{}
				for _, p := range t.Variable.Restrictions {
					found[p] = true
				}
				for _, p := range ps {
					if found[p] {
						res = append(res, p)
					}
				}
				ps = res
			}
		}
	}
	if len(ps) == 1 {
		return Type{Primitive: &ps[0]}, nil
	}
	return Type{Variable: &TypeVar{ps}}, nil
}

func (u *Unifier) buildReplace() error {
	u.replace = map[*TypeVar]Type{}
	visited := map[*TypeVar]bool{}
	for k := range u.graph {
		if !visited[k] {
			visited[k] = true
			trv := u.traverse(k)
			solved, err := solveRestrictions(trv)
			if err != nil {
				return err
			}
			for _, t := range trv {
				if t.IsVariable() {
					u.replace[t.Variable] = solved
				}
			}
		}
	}
	// TODO
	return nil
}

func (u *Unifier) Apply(a *Type) {
	if a.IsVariable() && u.replace[a.Variable].IsDefined() {
		to := u.replace[a.Variable]
		a.Variable = to.Variable
		a.Primitive = to.Primitive
		a.Function = to.Function
	}
	if a.IsFunction() {
		pars := make([]Type, len(*a.Function))
		for i, p := range *a.Function {
			u.Apply(&p)
			pars[i] = p
		}
		a.Function = &pars
	}
}

func Unify(a Type, b Type) (*Unifier, error) {
	uni := NewUnifier()
	if err := uni.buildGraph(a, b); err != nil {
		return nil, err
	}
	if err := uni.buildReplace(); err != nil {
		return nil, err
	}
	return uni, nil
}
