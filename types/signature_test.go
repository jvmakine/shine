package types

import "testing"

func TestType_Signature(t *testing.T) {
	tests := []struct {
		name string
		typ  Type
		want string
	}{{
		name: "support structures without variables",
		typ:  MakeStructure("S", []Type{}, SField{"a", IntP}, SField{"b", MakeFunction(RealP, RealP)}, SField{"c", BoolP}),
		want: "S",
	}, {
		name: "support structures with variables",
		typ: WithType(MakeVariable(), func(t Type) Type {
			return MakeStructure("S", []Type{t}, SField{"a", t}, SField{"b", MakeFunction(t, IntP)}, SField{"c", BoolP})
		}),
		want: "S[V1]",
	}, {
		name: "support structural variables",
		typ:  MakeStructuralVar(map[string]Type{"x": IntP}),
		want: "V1{x:int}",
	}, {
		name: "support hierarchical variables",
		typ:  MakeHierarchicalVar(MakeVariable().Variable, MakeVariable()),
		want: "V1[V2]",
	}, {
		name: "support variable functions",
		typ:  MakeFunction(MakeVariable(), MakeVariable()),
		want: "(V1)=>V2",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.typ
			if got := tr.Signature(); got != tt.want {
				t.Errorf("Type.Signature() = %v, want %v", got, tt.want)
			}
		})
	}
}
