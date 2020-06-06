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

func TestType_TSignature(t *testing.T) {
	tests := []struct {
		name string
		typ  Type
		want string
	}{{
		name: "support recursive structures",
		typ:  recursiveStruct("data", "b", SField{"a", IntP}),
		want: "data{a:int,b:<st>}",
	}, {
		name: "support functions",
		typ:  MakeFunction(IntP, MakeFunction(IntP, IntP), MakeFunction(RealP)),
		want: "(int,<fn>)=><fn>",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.typ
			if got := tr.TSignature(); got != tt.want {
				t.Errorf("Type.TSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
