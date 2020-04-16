package inferer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jvmakine/shine/types"
)

func TestUnify(t *testing.T) {
	type args struct {
		a *types.TypePtr
		b *types.TypePtr
	}
	tests := []struct {
		name  string
		left  *types.TypePtr
		right *types.TypePtr
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
		err:   errors.New("can not unify int with bool"),
	}, {
		name:  "should unify same consts",
		left:  base("int"),
		right: base("int"),
		want:  "int",
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
	}, {
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unify(tt.left, tt.right)
			if err != nil {
				if !reflect.DeepEqual(err, tt.err) {
					t.Errorf("Unify() error = %v, want %v", err, tt.err)
				}
				return
			}
			got.ApplySource(tt.left)
			got.ApplyDest(tt.right)
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
