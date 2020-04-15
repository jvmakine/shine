package typeinferer

import (
	"errors"
	"reflect"
	"testing"
)

func TestUnify(t *testing.T) {
	type args struct {
		a *TypePtr
		b *TypePtr
	}
	tests := []struct {
		name  string
		left  *TypePtr
		right *TypePtr
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
		left:  fun("A", base("int")).v.Type.(*TypePtr),
		right: fun(base("bool"), base("int")).v.Type.(*TypePtr),
		want:  "(bool,int)",
		err:   nil,
	}, {
		name:  "should unify repeating variables in functions",
		left:  fun("A", "A").v.Type.(*TypePtr),
		right: fun(base("bool"), "A").v.Type.(*TypePtr),
		want:  "(bool,bool)",
		err:   nil,
	}, {
		name:  "should fail repeating variables in conflict",
		left:  fun("A", "A").v.Type.(*TypePtr),
		right: fun(base("bool"), base("int")).v.Type.(*TypePtr),
		want:  "",
		err:   errors.New("can not unify int with bool"),
	}, {
		name:  "should unify two way parameter references",
		left:  fun("A", "A", "B", "A").v.Type.(*TypePtr),
		right: fun("C", base("int"), "D", "E").v.Type.(*TypePtr),
		want:  "(int,int,V1,int)",
		err:   nil,
	}, {
		name:  "should fail on two way parameter reference conflicts",
		left:  fun("A", "A", base("bool"), "A").v.Type.(*TypePtr),
		right: fun("C", base("int"), "C", "E").v.Type.(*TypePtr),
		want:  "",
		err:   errors.New("can not unify bool with int"),
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
