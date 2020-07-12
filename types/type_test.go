package types

/*func TestType_AddToUnion(t *testing.T) {
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
		typ:  NewUnionVariable(Int, Real),
		arg:  NewVariable(),
		want: NewVariable(),
	}, {
			name: "adding a structures variable to a union results into a union variable",
			typ:  NewUnionVariable(Int, Real),
			arg:  MakeStructuralVar(map[string]Type{"a": NewVariable()}),
			want: NewUnionVariable(Int, Real, MakeStructuralVar(map[string]Type{"a": NewVariable()})),
		}, {
			name: "adding a duplicate value to a union does not change it",
			typ:  NewUnionVariable(Real, Int),
			arg:  Int,
			want: NewUnionVariable(Real, Int),
		}, {
			name: "adding a type to a non union type makes it a union type",
			typ:  Real,
			arg:  Int,
			want: NewUnionVariable(Real, Int),
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
*/
