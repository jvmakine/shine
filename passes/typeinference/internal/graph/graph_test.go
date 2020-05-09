package graph

import (
	"reflect"
	"testing"

	. "github.com/jvmakine/shine/types"
)

func TestGraphBuilding(tes *testing.T) {
	v1 := &TypeVar{}
	v2 := &TypeVar{}
	v3 := &TypeVar{}
	tests := []struct {
		name   string
		typs   []Type
		substs map[*TypeVar]Type
	}{{
		name: "add function parameter -> result relations correctly",
		typs: []Type{
			MakeFunction(Type{Variable: v1}, Type{Variable: v1}, Type{Variable: v1}),
			MakeFunction(Type{Variable: v3}, IntP, Type{Variable: v2}),
		},
		substs: map[*TypeVar]Type{v1: IntP, v2: IntP, v3: IntP},
	}, {
		name: "add function parameter -> result relations correctly when going through a variable",
		typs: []Type{
			MakeFunction(Type{Variable: v1}, Type{Variable: v1}, Type{Variable: v1}),
			MakeVariable(),
			MakeFunction(Type{Variable: v3}, IntP, Type{Variable: v2}),
		},
		substs: map[*TypeVar]Type{v1: IntP, v2: IntP, v3: IntP},
	},
	}
	for _, tt := range tests {
		tes.Run(tt.name, func(t *testing.T) {
			graph := MakeTypeGraph()
			s := tt.typs[0]
			r := tt.typs[1:]
			for len(r) > 0 {
				n := r[0]
				r = r[1:]
				graph.Add(s, n)
				s = n
			}
			subst, err := graph.Substitutions()
			if err != nil {
				t.Fatal(err)
			} else {
				for k, v := range tt.substs {
					if !reflect.DeepEqual(subst[k], v) {
						t.Errorf("wrong substitution = %v, want %v", subst[k].Signature(), v.Signature())
					}
				}
			}

		})
	}
}
