package types

import (
	"reflect"
	"testing"
)

func TestType_AddToUnion(t *testing.T) {
	type args struct {
		o Type
	}
	tests := []struct {
		name string
		typ  Type
		arg  Type
		want Type
	}{{
		name: "adding a free variable to a union results into a free variable",
		typ:  MakeUnionVar(IntP, RealP),
		arg:  MakeVariable(),
		want: MakeVariable(),
	}, {
		name: "adding a structures variable to a union results into a union variable",
		typ:  MakeUnionVar(IntP, RealP),
		arg:  MakeStructuralVar(map[string]Type{"a": MakeVariable()}),
		want: MakeUnionVar(IntP, RealP, MakeStructuralVar(map[string]Type{"a": MakeVariable()})),
	}, {
		name: "adding a duplicate value to a union does not change it",
		typ:  MakeUnionVar(RealP, IntP),
		arg:  IntP,
		want: MakeUnionVar(RealP, IntP),
	}, {
		name: "adding a type to a non union type makes it a union type",
		typ:  RealP,
		arg:  IntP,
		want: MakeUnionVar(RealP, IntP),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.typ
			if got := tr.AddToUnion(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type.AddToUnion() = %v, want %v", got, tt.want)
			}
		})
	}
}
