package types

import (
	"errors"
	"reflect"
	"testing"

	"github.com/roamz/deepdiff"
)

func TestType_Unify(t *testing.T) {
	var1 := MakeVariable()
	var2 := MakeVariable()
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
	}, {
		name: "fails to unify different primitives",
		a:    MakePrimitive("int"),
		b:    MakePrimitive("bool"),
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "unifies union variables to subsets",
		a:    MakeUnionVar(IntP, BoolP, RealP),
		b:    MakeUnionVar(BoolP, RealP, StringP),
		want: MakeUnionVar(BoolP, RealP),
	}, {
		name: "unifies union variables to primitives",
		a:    MakeUnionVar(IntP, BoolP),
		b:    MakeUnionVar(BoolP, RealP),
		want: MakePrimitive("bool"),
	}, {
		name: "fails to unify disjoint restricted primitives",
		a:    MakeUnionVar(IntP, BoolP),
		b:    MakeUnionVar(StringP, RealP),
		err:  errors.New("can not unify V1[int|bool] with V1[string|real]"),
	}, {
		name: "unifies restricted variables with primitives",
		a:    MakePrimitive("bool"),
		b:    MakeUnionVar(IntP, BoolP),
		want: MakePrimitive("bool"),
	}, {
		name: "fails to unify restricted variables with incompatible primitives",
		a:    MakePrimitive("real"),
		b:    MakeUnionVar(IntP, BoolP),
		err:  errors.New("can not unify V1[int|bool] with real"),
	}, {
		name: "unifies identical functions",
		a:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		b:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		want: MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
	}, {
		name: "fails to unify mismatching functions",
		a:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		b:    MakeFunction(MakePrimitive("real"), MakePrimitive("int")),
		err:  errors.New("can not unify int with real"),
	}, {
		name: "fails to unify a function with a primitive",
		a:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		b:    MakePrimitive("int"),
		err:  errors.New("can not unify (real)=>real with int"),
	}, {
		name: "unifies variable functions with variables",
		a:    MakeVariable(),
		b:    MakeFunction(MakeVariable(), MakePrimitive("real")),
		want: MakeFunction(MakeVariable(), MakePrimitive("real")),
	}, {
		name: "unifies variables within functions",
		a:    MakeFunction(MakePrimitive("int"), MakeVariable()),
		b:    MakeFunction(MakeVariable(), MakePrimitive("real")),
		want: MakeFunction(MakePrimitive("int"), MakePrimitive("real")),
	}, {
		name: "unifies variables with restricted variables",
		a:    MakeVariable(),
		b:    MakeUnionVar(IntP, RealP),
		want: MakeUnionVar(IntP, RealP),
	}, {
		name: "unifies functions with overlapping variables",
		a:    MakeFunction(var1, IntP, var1),
		b:    MakeFunction(var2, var2, var2),
		want: MakeFunction(IntP, IntP, IntP),
	}, {
		name: "fails to unify mismatching functions",
		a:    MakeFunction(var1, IntP, var1),
		b:    MakeFunction(var2, var2, RealP),
		want: Type{},
		err:  errors.New("can not unify int with real"),
	}, {
		name: "unifies matching structures",
		a:    MakeStructure("s1", SField{"a", IntP}, SField{"b", BoolP}),
		b:    MakeStructure("s1", SField{"a", IntP}, SField{"b", BoolP}),
		want: MakeStructure("s1", SField{"a", IntP}, SField{"b", BoolP}),
	}, {
		name: "unifies matching structures with variables",
		a:    MakeStructure("s1", SField{"a", MakeVariable()}, SField{"b", BoolP}),
		b:    MakeStructure("s1", SField{"a", IntP}, SField{"b", MakeVariable()}),
		want: MakeStructure("s1", SField{"a", IntP}, SField{"b", BoolP}),
	}, {
		name: "unifies matching structures with variables",
		a:    MakeStructure("s1", SField{"a", MakeVariable()}, SField{"b", BoolP}),
		b:    MakeStructure("s1", SField{"a", IntP}, SField{"b", MakeVariable()}),
		want: MakeStructure("s1", SField{"a", IntP}, SField{"b", BoolP}),
	}, {
		name: "fails to unify structure with a function",
		a:    MakeStructure("", SField{"a", IntP}, SField{"b", BoolP}),
		b:    MakeFunction(IntP, BoolP),
		want: Type{},
		err:  errors.New("can not unify (int)=>bool with {a:int,b:bool}"),
	}, {
		name: "fails to unify structure with a primitive",
		a:    MakeStructure("", SField{"a", IntP}),
		b:    IntP,
		want: Type{},
		err:  errors.New("can not unify int with {a:int}"),
	}, {
		name: "unifies a structure with a variable",
		a:    MakeStructure("", SField{"a", IntP}, SField{"b", BoolP}),
		b:    MakeVariable(),
		want: MakeStructure("", SField{"a", IntP}, SField{"b", BoolP}),
	}, {
		name: "fails to unify on name mismatch",
		a:    MakeStructure("s1", SField{"a", IntP}),
		b:    MakeStructure("s2", SField{"a", IntP}),
		want: Type{},
		err:  errors.New("can not unify s1{a:int} with s2{a:int}"),
	}, {
		name: "unifies identical recursive structures",
		a:    recursiveStruct("data", "r", SField{"a", IntP}),
		b:    recursiveStruct("data", "r", SField{"a", IntP}),
		want: recursiveStruct("data", "r", SField{"a", IntP}),
	}, {
		name: "unifies structural variables with generic variables",
		a:    MakeVariable(),
		b:    MakeStructuralVar(map[string]Type{"x": IntP}),
		want: MakeStructuralVar(map[string]Type{"x": IntP}),
	}, {
		name: "fails to unify union var with a structural var",
		a:    MakeStructuralVar(map[string]Type{"x": IntP}),
		b:    MakeUnionVar(IntP, RealP),
		want: Type{},
		err:  errors.New("can not unify V1[int|real] with V1{x:int}"),
	}, {
		name: "combines non conflicting structural variables",
		a:    MakeStructuralVar(map[string]Type{"x": IntP}),
		b:    MakeStructuralVar(map[string]Type{"y": RealP}),
		want: MakeStructuralVar(map[string]Type{"x": IntP, "y": RealP}),
	}, {
		name: "fails on conflicting structural variables",
		a:    MakeStructuralVar(map[string]Type{"x": IntP}),
		b:    MakeStructuralVar(map[string]Type{"x": RealP}),
		want: Type{},
		err:  errors.New("can not unify int with real"),
	}, {
		name: "fails on structural variables with primitives",
		a:    MakeStructuralVar(map[string]Type{"x": IntP}),
		b:    IntP,
		want: Type{},
		err:  errors.New("can not unify V1{x:int} with int"),
	}, {
		name: "unifies structural variables with structures",
		a:    MakeStructure("a", SField{"x", IntP}),
		b:    MakeStructuralVar(map[string]Type{"x": MakeVariable()}),
		want: MakeStructure("a", SField{"x", IntP}),
	}, {
		name: "fails to unify conflicting structural variables with structures",
		a:    MakeStructure("a", SField{"x", IntP}),
		b:    MakeStructuralVar(map[string]Type{"y": MakeVariable()}),
		want: Type{},
		err:  errors.New("can not unify V1{y:V2} with a{x:int}"),
	}, {
		name: "unifies variables wthin structural variables",
		a:    MakeStructuralVar(map[string]Type{"x": MakeUnionVar(IntP, BoolP)}),
		b:    MakeStructuralVar(map[string]Type{"x": MakeUnionVar(IntP, RealP)}),
		want: MakeStructuralVar(map[string]Type{"x": IntP}),
	}, {
		name: "unifies union variables with structural variables",
		a: MakeUnionVar(
			MakeStructuralVar(map[string]Type{"x": IntP, "y": RealP}),
			MakeStructuralVar(map[string]Type{"x": RealP, "y": RealP}),
		),
		b: MakeUnionVar(
			MakeStructuralVar(map[string]Type{"x": IntP, "a": BoolP}),
			MakeStructuralVar(map[string]Type{"b": BoolP}),
		),
		want: MakeUnionVar(
			MakeStructuralVar(map[string]Type{"x": IntP, "y": RealP, "a": BoolP}),
			MakeStructuralVar(map[string]Type{"x": IntP, "y": RealP, "b": BoolP}),
			MakeStructuralVar(map[string]Type{"x": RealP, "y": RealP, "b": BoolP}),
		),
	}, {
		name: "removes duplicates in union variables",
		a: MakeUnionVar(
			MakeStructuralVar(map[string]Type{"x": IntP}),
			MakeStructuralVar(map[string]Type{"y": RealP}),
		),
		b: MakeUnionVar(
			MakeStructuralVar(map[string]Type{"x": IntP, "y": RealP}),
		),
		want: MakeStructuralVar(map[string]Type{"x": IntP, "y": RealP}),
	}, {
		name: "unifies functions with union variables",
		a: MakeUnionVar(
			IntP,
			MakeFunction(IntP, IntP),
			MakeFunction(IntP, RealP),
		),
		b: MakeFunction(IntP, MakeVariable()),
		want: MakeUnionVar(
			MakeFunction(IntP, IntP),
			MakeFunction(IntP, RealP),
		),
	}, {
		name: "unifies primitives with free variables in unions",
		a: MakeUnionVar(
			RealP,
			MakeFunction(MakeVariable(), MakeVariable()),
			MakeStructuralVar(map[string]Type{"x": MakeVariable()}),
			MakeVariable(),
		),
		b:    IntP,
		want: IntP,
	}, {
		name: "unifies functions with free variables with functions with union variables",
		a:    MakeFunction(MakeVariable(), MakeUnionVar(IntP, RealP)),
		b:    MakeFunction(MakeUnionVar(IntP, RealP), MakeUnionVar(BoolP, RealP)),
		want: MakeFunction(MakeUnionVar(IntP, RealP), RealP),
	}, {
		name: "unifies dependent variables in unions",
		a:    MakeFunction(var1, var1),
		b:    MakeFunction(MakeUnionVar(IntP, RealP), MakeUnionVar(BoolP, RealP)),
		want: MakeFunction(RealP, RealP),
	}, {
		name: "unifies functions based on a union variable",
		a:    MakeFunction(IntP, MakeVariable()),
		b:    MakeUnionVar(MakeFunction(IntP, IntP), MakeFunction(RealP, RealP)),
		want: MakeFunction(IntP, IntP),
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
				ok, _ := deepdiff.DeepDiff(got, tt.want)
				gotsign := got.Signature()
				wantsign := tt.want.Signature()
				if !ok {
					t.Error(gotsign + " did not equal " + wantsign)
				}
			}
		})
	}
}
