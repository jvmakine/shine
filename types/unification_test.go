package types

import (
	"errors"
	"reflect"
	"testing"
)

func TestType_Unify(t *testing.T) {
	tests := []struct {
		name string
		a    Type
		b    Type
		want Type
		err  error
	}{{
		name: "unifies same primitives",
		a:    MakePrimitive("int"),
		b:    MakePrimitive("int"),
		want: MakePrimitive("int"),
		err:  nil,
	}, {
		name: "fails to unify different primitives",
		a:    MakePrimitive("int"),
		b:    MakePrimitive("bool"),
		want: Type{},
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "unifies restricted primitives to subsets",
		a:    MakeRestricted("int", "bool", "real"),
		b:    MakeRestricted("bool", "real", "foo"),
		want: MakeRestricted("bool", "real"),
		err:  nil,
	}, {
		name: "fails to unify disjoint restricted primitives",
		a:    MakeRestricted("int", "bool"),
		b:    MakeRestricted("bar", "foo"),
		want: Type{},
		err:  errors.New("can not unify V1[int|bool] with V1[bar|foo]"),
	}, {
		name: "unifies restricted variables with primitives",
		a:    MakePrimitive("bool"),
		b:    MakeRestricted("int", "bool"),
		want: MakePrimitive("bool"),
		err:  nil,
	}, {
		name: "fails to unify restricted variables with incompatible primitives",
		a:    MakePrimitive("real"),
		b:    MakeRestricted("int", "bool"),
		want: Type{},
		err:  errors.New("can not unify real with V1[int|bool]"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Unify(tt.b)
			if tt.err != nil {
				if !reflect.DeepEqual(err, tt.err) {
					t.Errorf("Type.Unify() error = %v, wantErr %v", err, tt.err)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Type.Unify() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Type.Unify() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}
