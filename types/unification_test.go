package types

import (
	"errors"
	"reflect"
	"testing"
)

func TestType_Unify(t *testing.T) {
	var1 := NewVariable()
	var2 := NewVariable()
	tests := []struct {
		name string
		a    Type
		b    Type
		want Type
		err  error
	}{{
		name: "unifies same primitives",
		a:    Int,
		b:    Int,
		want: Int,
	}, {
		name: "fails to unify different primitives",
		a:    Int,
		b:    Bool,
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "unifies union variables to subsets",
		a:    NewUnionVariable(Int, Bool, Real),
		b:    NewUnionVariable(Bool, Real, String),
		want: NewUnionVariable(Bool, Real),
	}, {
		name: "unifies union variables to primitives",
		a:    NewUnionVariable(Int, Bool),
		b:    NewUnionVariable(Bool, Real),
		want: Bool,
	}, {
		name: "fails to unify disjoint restricted primitives",
		a:    NewUnionVariable(Int, Bool),
		b:    NewUnionVariable(String, Real),
		err:  errors.New("can not unify V1[int|bool] with V1[string|real]"),
	}, {
		name: "unifies restricted variables with primitives",
		a:    Bool,
		b:    NewUnionVariable(Int, Bool),
		want: Bool,
	}, {
		name: "fails to unify restricted variables with incompatible primitives",
		a:    Real,
		b:    NewUnionVariable(Int, Bool),
		err:  errors.New("can not unify V1[int|bool] with real"),
	}, {
		name: "unifies identical functions",
		a:    NewFunction(Real, Real),
		b:    NewFunction(Real, Real),
		want: NewFunction(Real, Real),
	}, {
		name: "fails to unify mismatching functions",
		a:    NewFunction(Real, Real),
		b:    NewFunction(Real, Int),
		err:  errors.New("can not unify int with real"),
	}, {
		name: "fails to unify a function with a primitive",
		a:    NewFunction(Real, Real),
		b:    Int,
		err:  errors.New("can not unify (real)=>real with int"),
	}, {
		name: "unifies variable functions with variables",
		a:    NewVariable(),
		b:    NewFunction(NewVariable(), Real),
		want: NewFunction(NewVariable(), Real),
	}, {
		name: "unifies variables within functions",
		a:    NewFunction(Int, NewVariable()),
		b:    NewFunction(NewVariable(), Real),
		want: NewFunction(Int, Real),
	}, {
		name: "unifies variables with restricted variables",
		a:    NewVariable(),
		b:    NewUnionVariable(Int, Real),
		want: NewUnionVariable(Int, Real),
	}, {
		name: "unifies functions with overlapping variables",
		a:    NewFunction(var1, Int, var1),
		b:    NewFunction(var2, var2, var2),
		want: NewFunction(Int, Int, Int),
	}, {
		name: "fails to unify mismatching functions",
		a:    NewFunction(var1, Int, var1),
		b:    NewFunction(var2, var2, Real),
		err:  errors.New("can not unify int with real"),
	}, {
		name: "unifies matching structures",
		a:    NewStructure(Named{"a", Int}, Named{"b", Bool}),
		b:    NewStructure(Named{"a", Int}, Named{"b", Bool}),
		want: NewStructure(Named{"a", Int}, Named{"b", Bool}),
	}, {
		name: "fails to unify structure with a function",
		a:    NewStructure(Named{"a", Int}, Named{"b", Bool}),
		b:    NewFunction(Int, Bool),
		err:  errors.New("can not unify (int)=>bool with {a:int,b:bool}"),
	}, {
		name: "fails to unify structure with a primitive",
		a:    NewStructure(Named{"a", Int}),
		b:    Int,
		err:  errors.New("can not unify int with {a:int}"),
	}, {
		name: "unifies a structure with a variable",
		a:    NewStructure(Named{"a", Int}, Named{"b", Bool}),
		b:    NewVariable(),
		want: NewStructure(Named{"a", Int}, Named{"b", Bool}),
	}, {
		name: "fails to unify on name mismatch",
		a:    NewNamed("s1", NewStructure(Named{"a", Int})),
		b:    NewNamed("s2", NewStructure(Named{"a", Int})),
		err:  errors.New("can not unify s1[{a:int}] with s2[{a:int}]"),
	}, {
		name: "unifies identical recursive structures",
		a:    recursiveStruct("data", "r", Named{"a", Int}),
		b:    recursiveStruct("data", "r", Named{"a", Int}),
		want: recursiveStruct("data", "r", Named{"a", Int}),
	}, {
		name: "unifies structures with generic variables",
		a:    NewVariable(),
		b:    NewStructure(NewNamed("x", Int)),
		want: NewStructure(NewNamed("x", Int)),
	}, {
		name: "fails to unify union var with a structural var",
		a:    NewStructure(NewNamed("x", Int)),
		b:    NewUnionVariable(Int, Real),
		err:  errors.New("can not unify V1[int|real] with {x:int}"),
	}, {
		name: "combines non conflicting structural variables",
		a:    NewStructuralVar(NewNamed("x", Int)),
		b:    NewStructuralVar(NewNamed("y", Real)),
		want: NewStructuralVar(NewNamed("x", Int), NewNamed("y", Real)),
	}, {
		name: "fails on conflicting structural variables",
		a:    NewStructuralVar(NewNamed("x", Int)),
		b:    NewStructuralVar(NewNamed("x", Real)),
		err:  errors.New("can not unify int with real"),
	}, {
		name: "fails on structural variables with primitives",
		a:    NewStructuralVar(NewNamed("x", Int)),
		b:    Int,
		err:  errors.New("can not unify V1{x:int} with int"),
	}, {
		name: "unifies structural variables with structures",
		a:    NewNamed("a", NewStructure(NewNamed("x", Int))),
		b:    NewStructuralVar(NewNamed("x", NewVariable())),
		want: NewNamed("a", NewStructure(NewNamed("x", Int))),
	}, {
		name: "fails to unify conflicting structural variables with structures",
		a:    NewNamed("a", NewStructure(NewNamed("x", Int))),
		b:    NewStructuralVar(NewNamed("y", NewVariable())),
		err:  errors.New("can not unify V1{y:V2} with a[{x:int}]"),
	}, {
		name: "unifies variables wthin structural variables",
		a:    NewStructuralVar(NewNamed("x", NewUnionVariable(Int, Bool))),
		b:    NewStructuralVar(NewNamed("x", NewUnionVariable(Int, Real))),
		want: NewStructuralVar(NewNamed("x", Int)),
	}, {
		name: "unifies union variables with structural variables",
		a: NewUnionVariable(
			NewStructuralVar(NewNamed("x", Int), NewNamed("y", Real)),
			NewStructuralVar(NewNamed("x", Real), NewNamed("y", Real)),
		),
		b: NewUnionVariable(
			NewStructuralVar(NewNamed("x", Int), NewNamed("a", Bool)),
			NewStructuralVar(NewNamed("b", Bool)),
		),
		want: NewUnionVariable(
			NewStructuralVar(NewNamed("x", Int), NewNamed("y", Real), NewNamed("a", Bool)),
			NewStructuralVar(NewNamed("x", Int), NewNamed("y", Real), NewNamed("b", Bool)),
			NewStructuralVar(NewNamed("x", Real), NewNamed("y", Real), NewNamed("b", Bool)),
		),
	}, {
		name: "removes duplicates in union variables",
		a: NewUnionVariable(
			NewStructuralVar(NewNamed("x", Int)),
			NewStructuralVar(NewNamed("y", Real)),
		),
		b: NewUnionVariable(
			NewStructuralVar(NewNamed("x", Int), NewNamed("y", Real)),
		),
		want: NewStructuralVar(NewNamed("x", Int), NewNamed("y", Real)),
	}, {
		name: "unifies functions with union variables",
		a: NewUnionVariable(
			Int,
			NewFunction(Int, Int),
			NewFunction(Int, Real),
		),
		b: NewFunction(Int, NewVariable()),
		want: NewUnionVariable(
			NewFunction(Int, Int),
			NewFunction(Int, Real),
		),
	}, {
		name: "unifies primitives with free variables in unions",
		a: NewUnionVariable(
			Real,
			NewFunction(NewVariable(), NewVariable()),
			NewStructuralVar(NewNamed("x", NewVariable())),
			NewVariable(),
		),
		b:    Int,
		want: Int,
	}, {
		name: "unifies functions with free variables with functions with union variables",
		a:    NewFunction(NewVariable(), NewUnionVariable(Int, Real)),
		b:    NewFunction(NewUnionVariable(Int, Real), NewUnionVariable(Bool, Real)),
		want: NewFunction(NewUnionVariable(Int, Real), Real),
	}, {
		name: "unifies dependent variables in unions",
		a:    NewFunction(var1, var1),
		b:    NewFunction(NewUnionVariable(Int, Real), NewUnionVariable(Bool, Real)),
		want: NewFunction(Real, Real),
	}, {
		name: "unifies functions based on a union variable",
		a:    NewFunction(Int, NewVariable()),
		b:    NewUnionVariable(NewFunction(Int, Int), NewFunction(Real, Real)),
		want: NewFunction(Int, Int),
	}, {
		name: "unifies variable functions with union variables",
		a:    NewFunction(NewVariable(), NewVariable()),
		b:    NewUnionVariable(NewFunction(Int, Int), NewFunction(Real, Real)),
		want: NewUnionVariable(NewFunction(Int, Int), NewFunction(Real, Real)),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unify(tt.a, tt.b)
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
				gotsign := Signature(got)
				wantsign := Signature(tt.want)
				if gotsign != wantsign {
					t.Error(gotsign + " did not equal " + wantsign)
				}
			}
		})
	}
}
