package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
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
		err:  nil,
	}, {
		name: "fails to unify different primitives",
		a:    MakePrimitive("int"),
		b:    MakePrimitive("bool"),
		want: Type{},
		err:  errors.New("can not unify bool with int"),
	}, {
		name: "unifies union variables to subsets",
		a:    MakeUnionVar("int", "bool", "real"),
		b:    MakeUnionVar("bool", "real", "foo"),
		want: MakeUnionVar("bool", "real"),
		err:  nil,
	}, {
		name: "unifies union variables to primitives",
		a:    MakeUnionVar("int", "bool"),
		b:    MakeUnionVar("bool", "real"),
		want: MakePrimitive("bool"),
		err:  nil,
	}, {
		name: "fails to unify disjoint restricted primitives",
		a:    MakeUnionVar("int", "bool"),
		b:    MakeUnionVar("bar", "foo"),
		want: Type{},
		err:  errors.New("can not unify V1[bar|foo] with V1[int|bool]"),
	}, {
		name: "unifies restricted variables with primitives",
		a:    MakePrimitive("bool"),
		b:    MakeUnionVar("int", "bool"),
		want: MakePrimitive("bool"),
		err:  nil,
	}, {
		name: "fails to unify restricted variables with incompatible primitives",
		a:    MakePrimitive("real"),
		b:    MakeUnionVar("int", "bool"),
		want: Type{},
		err:  errors.New("can not unify real with V1[int|bool]"),
	}, {
		name: "unifies identical functions",
		a:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		b:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		want: MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		err:  nil,
	}, {
		name: "fails to unify mismatching functions",
		a:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		b:    MakeFunction(MakePrimitive("real"), MakePrimitive("int")),
		want: Type{},
		err:  errors.New("can not unify int with real"),
	}, {
		name: "fails to unify a function with a primitive",
		a:    MakeFunction(MakePrimitive("real"), MakePrimitive("real")),
		b:    MakePrimitive("int"),
		want: Type{},
		err:  errors.New("can not unify (real)=>real with int"),
	}, {
		name: "unifies variable functions with variables",
		a:    MakeVariable(),
		b:    MakeFunction(MakeVariable(), MakePrimitive("real")),
		want: MakeFunction(MakeVariable(), MakePrimitive("real")),
		err:  nil,
	}, {
		name: "unifies variables within functions",
		a:    MakeFunction(MakePrimitive("int"), MakeVariable()),
		b:    MakeFunction(MakeVariable(), MakePrimitive("real")),
		want: MakeFunction(MakePrimitive("int"), MakePrimitive("real")),
		err:  nil,
	}, {
		name: "unifies variables with restricted variables",
		a:    MakeVariable(),
		b:    MakeUnionVar("int", "real"),
		want: MakeUnionVar("int", "real"),
		err:  nil,
	}, {
		name: "unifies functions with overlapping variables",
		a:    MakeFunction(var1, IntP, var1),
		b:    MakeFunction(var2, var2, var2),
		want: MakeFunction(IntP, IntP, IntP),
		err:  nil,
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
		err:  nil,
	}, {
		name: "unifies matching structures with variables",
		a:    MakeStructure("s1", SField{"a", MakeVariable()}, SField{"b", BoolP}),
		b:    MakeStructure("s1", SField{"a", IntP}, SField{"b", MakeVariable()}),
		want: MakeStructure("s1", SField{"a", IntP}, SField{"b", BoolP}),
		err:  nil,
	}, {
		name: "unifies matching structures with variables",
		a:    MakeStructure("s1", SField{"a", MakeVariable()}, SField{"b", BoolP}),
		b:    MakeStructure("s1", SField{"a", IntP}, SField{"b", MakeVariable()}),
		want: MakeStructure("s1", SField{"a", IntP}, SField{"b", BoolP}),
		err:  nil,
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
		err:  nil,
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
		err:  nil,
	}, {
		name: "unifies structural variables with generic variables",
		a:    MakeVariable(),
		b:    MakeStructuralVar(map[string]Type{"x": IntP}),
		want: MakeStructuralVar(map[string]Type{"x": IntP}),
		err:  nil,
	}, {
		name: "fails to unify union var with a structural var",
		a:    MakeStructuralVar(map[string]Type{"x": IntP}),
		b:    MakeUnionVar(Int, Real),
		want: Type{},
		err:  errors.New("can not unify V1[int|real] with V1{x:int}"),
	}, {
		name: "combines non conflicting structural variables",
		a:    MakeStructuralVar(map[string]Type{"x": IntP}),
		b:    MakeStructuralVar(map[string]Type{"y": RealP}),
		want: MakeStructuralVar(map[string]Type{"x": IntP, "y": RealP}),
		err:  nil,
	}, {
		name: "fails on conflicting structural variables",
		a:    MakeStructuralVar(map[string]Type{"x": IntP}),
		b:    MakeStructuralVar(map[string]Type{"x": RealP}),
		want: Type{},
		err:  errors.New("can not unify int with real"),
	}, {
		name: "unifies structural variables with structures",
		a:    MakeStructure("a", SField{"x", IntP}),
		b:    MakeStructuralVar(map[string]Type{"x": MakeVariable()}),
		want: MakeStructure("a", SField{"x", IntP}),
		err:  nil,
	}, {
		name: "fails to unify conflicting structural variables with structures",
		a:    MakeStructure("a", SField{"x", IntP}),
		b:    MakeStructuralVar(map[string]Type{"y": MakeVariable()}),
		want: Type{},
		err:  errors.New("can not unify V1{y:V2} with a{x:int}"),
	}, {
		name: "unifies variables wthin structural variables",
		a:    MakeStructuralVar(map[string]Type{"x": MakeUnionVar(Int, Bool)}),
		b:    MakeStructuralVar(map[string]Type{"x": MakeUnionVar(Int, Real)}),
		want: MakeStructuralVar(map[string]Type{"x": IntP}),
		err:  nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Unify(tt.b)
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				if err != nil {
					t.Errorf("Type.Unify() error = %v", err)
					return
				}
				require.Equal(t, tt.want, got)
			}
		})
	}
}
