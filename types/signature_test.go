package types

import "testing"

func TestType_Signature(t *testing.T) {
	tests := []struct {
		name string
		typ  Type
		want string
	}{{
		name: "support structures without variables",
		typ:  MakeStructure("", SField{"a", IntP}, SField{"b", MakeFunction(RealP, RealP)}, SField{"c", BoolP}),
		want: "{a:int,b:(real)=>real,c:bool}",
	}, {
		name: "support structures with variables",
		typ: WithType(MakeVariable(), func(t Type) Type {
			return MakeStructure("", SField{"a", t}, SField{"b", MakeFunction(t, IntP)}, SField{"c", BoolP})
		}),
		want: "{a:V1,b:(V1)=>int,c:bool}",
	}, {
		name: "support named structures",
		typ:  MakeStructure("data", SField{"a", IntP}, SField{"b", BoolP}),
		want: "data{a:int,b:bool}",
	}, {
		name: "support recursive structures",
		typ:  recursiveStruct("data", "b", SField{"a", IntP}),
		want: "data{a:int,b:data}",
	}, {
		name: "support structural variables",
		typ:  MakeStructuralVar(map[string]Type{"x": IntP}),
		want: "V1{x:int}",
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
