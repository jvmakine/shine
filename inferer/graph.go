package inferer

import (
	"errors"
	"strconv"

	. "github.com/jvmakine/shine/types"
)

type TypeGraph map[*TypeVar][]Type

func (g TypeGraph) Combine(other TypeGraph) {
	for k, v := range other {
		g[k] = append(g[k], v...)
	}
}

func (g TypeGraph) Add(a Type, b Type) error {
	if _, err := a.Unify(b); err != nil {
		return err
	}
	if a.IsVariable() {
		g[a.Variable] = append(g[a.Variable], b)
	}
	if b.IsVariable() {
		g[b.Variable] = append(g[b.Variable], a)
	}
	// TODO: functions to variables unification
	if a.IsFunction() && b.IsFunction() {
		al := len(*a.Function)
		bl := len(*b.Function)
		if al != bl {
			return errors.New("wrong number of function arguments: " + strconv.Itoa(al) + " != " + strconv.Itoa(bl))
		}
		for i := range *a.Function {
			err := g.Add((*a.Function)[i], (*b.Function)[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (g TypeGraph) traverse(a *TypeVar) []Type {
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
			for _, f := range g[next] {
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

func (g TypeGraph) Substitutions() (Substitutions, error) {
	result := Substitutions{}
	done := map[*TypeVar]bool{}
	for k := range g {
		ts := g.traverse(k)
		vars := []*TypeVar{}
		ires := MakeVariable()

		for _, t := range ts {
			if done[t.Variable] {
				continue
			}
			if t.IsVariable() {
				done[t.Variable] = true
				vars = append(vars, t.Variable)
			}
			r, err := ires.Unify(t)
			ires = r
			if err != nil {
				return nil, err
			}
		}
		for _, v := range vars {
			result[v] = ires
		}
	}
	return result, nil
}
