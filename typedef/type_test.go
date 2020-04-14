package typedef

import (
	"reflect"
	"testing"
)

func TestSignature(t *testing.T) {
	tests := []struct {
		name string
		typ  *Type
		want TypeSignature
	}{{
		name: "return base type signature",
		typ:  primitive("int"),
		want: "int",
	}, {
		name: "return function type signature",
		typ:  function(primitive("int"), primitive("bool"), primitive("int")),
		want: "(int,bool,int)",
	}, {
		name: "return nested function type signature",
		typ:  function(primitive("int"), function(primitive("bool"), primitive("bool")), primitive("int")),
		want: "(int,(bool,bool),int)",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.typ
			got := tr.Signature()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type.Signature() = %v, want %v", got, tt.want)
			}
		})
	}
}
