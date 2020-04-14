package types

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
		typ:  base("int"),
		want: "int",
	}, {
		name: "return function type signature",
		typ:  function(base("int"), base("bool"), base("int")),
		want: "(int,bool,int)",
	}, {
		name: "return nested function type signature",
		typ:  function(base("int"), function(base("bool"), base("bool")), base("int")),
		want: "(int,(bool,bool),int)",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.typ
			got, err := tr.Signature()
			if err != nil {
				t.Errorf("Type.Signature() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type.Signature() = %v, want %v", got, tt.want)
			}
		})
	}
}
