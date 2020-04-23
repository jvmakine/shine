package inferer

import (
	"errors"
	"strconv"

	. "github.com/jvmakine/shine/types"
)

type TypeGraph struct {
	graph map[*TypeVar][]Type
}

func MakeTypeGraph() TypeGraph {
	return TypeGraph{map[*TypeVar][]Type{}}
}

func (g TypeGraph) findFunct(v *TypeVar) Type {
	if vs := g.graph[v]; vs != nil {
		for _, v := range vs {
			if v.IsFunction() {
				return v
			}
		}
	}
	return Type{}
}

func (g TypeGraph) Add(a Type, b Type) error {
	if _, err := a.Unify(b); err != nil {
		return err
	}
	if a.IsVariable() {
		g.graph[a.Variable] = append(g.graph[a.Variable], b)
	}
	if b.IsVariable() {
		g.graph[b.Variable] = append(g.graph[b.Variable], a)
	}
	var fun1, fun2 Type
	if a.IsFunction() {
		fun1 = a
		if b.IsFunction() {
			fun2 = b
		} else if b.IsVariable() {
			fun2 = g.findFunct(b.Variable)
		}
	} else if b.IsFunction() {
		fun1 = b
		if a.IsFunction() {
			fun2 = a
		} else if a.IsVariable() {
			fun2 = g.findFunct(a.Variable)
		}
	}

	if fun1.IsFunction() && fun2.IsFunction() {
		ap := fun1.FunctTypes()
		bp := fun2.FunctTypes()
		al := len(ap)
		bl := len(bp)
		if al != bl {
			return errors.New("wrong number of function arguments: " + strconv.Itoa(al) + " != " + strconv.Itoa(bl))
		}
		for i := range ap {
			err := g.Add(ap[i], bp[i])
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
			for _, f := range g.graph[next] {
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
	for k := range g.graph {
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
