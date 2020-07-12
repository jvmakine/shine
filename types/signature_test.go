package types

import "testing"

func TestType_Signature(t *testing.T) {
	tests := []struct {
		name string
		typ  Type
		want string
	}{{
		name: "support structures without variables",
		typ:  NewStructure(Named{"a", Int}, Named{"b", NewFunction(Real, Real)}, Named{"c", Bool}),
		want: "{a:int,b:(real)=>real,c:bool}",
	}, {
		name: "support structural variables with variables",
		typ: WithType(NewVariable(), func(t Type) Type {
			return NewStructuralVar(Named{"a", t}, Named{"b", NewFunction(t, Int)}, Named{"c", Bool})
		}),
		want: "V1{a:V2,b:(V2)=>int,c:bool}",
	}, {
		name: "support named structures",
		typ:  NewNamed("data", NewStructure(Named{"a", Int}, Named{"b", Bool})),
		want: "data[{a:int,b:bool}]",
	}, {
		name: "support recursive structures",
		typ:  recursiveStruct("data", "b", Named{"a", Int}),
		want: "data[{a:int,b:data}]",
	}, {
		name: "support structural variables",
		typ:  NewStructuralVar(NewNamed("x", Int)),
		want: "V1{x:int}",
	}, {
		name: "supports union variables in functions",
		typ:  NewFunction(NewUnionVariable(Int, Real), NewUnionVariable(Int, Bool)),
		want: "(V1[int|real])=>V2[int|bool]",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.typ
			if got := Signature(tr); got != tt.want {
				t.Errorf("Type.Signature() = %v, want %v", got, tt.want)
			}
		})
	}
}
