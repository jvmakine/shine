package inferer

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/jvmakine/shine/types"
)

func TestUnify(t *testing.T) {
	type args struct {
		a Type
		b Type
	}
	tests := []struct {
		name  string
		left  Type
		right Type
		want  string
		err   error
	}{{
		name:  "should unify a single variable",
		left:  variable(),
		right: base("int"),
		want:  "int",
		err:   nil,
	}, {
		name:  "should fail to unify differing consts",
		left:  base("bool"),
		right: base("int"),
		want:  "",
		err:   errors.New("can not unify bool with int"),
	}, {
		name:  "should unify same consts",
		left:  base("int"),
		right: base("int"),
		want:  "int",
		err:   nil,
	}, {
		name:  "should unify union variables with a base type",
		left:  base("int"),
		right: union(Int, Real),
		want:  "int",
		err:   nil,
	}, {
		name:  "should fail to unify union variables with invalid type",
		left:  union(Int, Real),
		right: base("bool"),
		want:  "",
		err:   errors.New("can not unify bool with V1[int|real]"),
	}, {
		name:  "should unify two different unions",
		left:  union(Bool, Real),
		right: union(Bool, Int),
		want:  "bool",
		err:   nil,
	}, {
		name:  "should unify simple variables in functions",
		left:  fun("A", base("int")).v.Type,
		right: fun(base("bool"), base("int")).v.Type,
		want:  "(bool,int)",
		err:   nil,
	}, {
		name:  "should unify repeating variables in functions",
		left:  fun("A", "A").v.Type,
		right: fun(base("bool"), "A").v.Type,
		want:  "(bool,bool)",
		err:   nil,
	}, {
		name:  "should fail repeating variables in conflict",
		left:  fun("A", "A").v.Type,
		right: fun(base("bool"), base("int")).v.Type,
		want:  "",
		err:   errors.New("can not unify int with bool"),
	}, {
		name:  "should unify two way parameter references",
		left:  fun("A", "A", "B", "A").v.Type,
		right: fun("C", base("int"), "D", "E").v.Type,
		want:  "(int,int,V1,int)",
		err:   nil,
	}, {
		name:  "should fail on two way parameter reference conflicts",
		left:  fun("A", "A", base("bool"), "A").v.Type,
		right: fun("C", base("int"), "C", "E").v.Type,
		want:  "",
		err:   errors.New("can not unify bool with int"),
	}, {
		name:  "unify to sets of variables into one",
		left:  fun("A", "A", "B", "A").v.Type,
		right: fun("A", "B", "A", "A").v.Type,
		want:  "(V1,V1,V1,V1)",
		err:   nil,
	}, /*{
		name:  "unify to sets of variables into one base value",
		left:  fun("A", "A", "B", "A").v.Type,
		right: fun("A", "B", "A", base("int")).v.Type,
		want:  "(int,int,int,int)",
		err:   nil,
	}, {
		name:  "fail to unify sets of variables into two base values",
		left:  fun("A", "A", base("bool"), "A").v.Type,
		right: fun("A", "B", "A", base("int")).v.Type,
		want:  "",
		err:   errors.New("can not unify int with bool"),
	}, {
		name:  "fail to unify functions of mismaching number of arguments",
		left:  fun("A", "A").v.Type,
		right: fun("A", "A", "A").v.Type,
		want:  "",
		err:   errors.New("wrong number of function arguments 2 given 3 required"),
	}, {
		name:  "fail to unify functions with values",
		left:  base("int"),
		right: fun(base("int")).v.Type,
		want:  "",
		err:   errors.New("a function required"),
	}, {
		name:  "unify union variable functions",
		left:  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, t) }).v.Type,
		right: withVar(union(Int, Real), func(t Type) *excon { return fun(base("int"), t, t) }).v.Type,
		want:  "(int,int,int)",
		err:   nil,
	}, {
		name:  "fail to unify conflictingunion variable functions",
		left:  withVar(union(Int, Real), func(t Type) *excon { return fun(t, t, base("real")) }).v.Type,
		right: withVar(union(Int, Real), func(t Type) *excon { return fun(base("int"), t, t) }).v.Type,
		want:  "",
		err:   errors.New("can not unify real with int"),
	}*/}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unify(tt.left, tt.right)
			if err != nil {
				if !reflect.DeepEqual(err, tt.err) {
					t.Errorf("Unify() error = %v, want %v", err, tt.err)
				}
				return
			}
			got.Apply(&tt.left)
			got.Apply(&tt.right)
			ls := tt.left.Signature()
			rs := tt.right.Signature()
			if rs != ls {
				t.Errorf("left and right mismatch (%v != %v)", ls, rs)
				return
			}
			if rs != tt.want {
				t.Errorf("got %v, want %v", rs, tt.want)
			}
		})
	}
}
