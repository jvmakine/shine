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
		ctx  MockUnificationCtx
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
		name: "unifies functions with overlapping variables",
		a:    NewFunction(var1, Int, var1),
		b:    NewFunction(var2, var2, var2),
		want: NewFunction(Int, Int, Int),
	}, {
		name: "fails to unify mismatching functions with variables",
		a:    NewFunction(var1, Int, var1),
		b:    NewFunction(var2, var2, Real),
		err:  errors.New("can not unify int with real"),
	}, {
		name: "unifies matching structures",
		a:    NewStructure(NewNamed("a", Int), NewNamed("b", Bool)),
		b:    NewStructure(NewNamed("a", Int), NewNamed("b", Bool)),
		want: NewStructure(NewNamed("a", Int), NewNamed("b", Bool)),
	}, {
		name: "fails to unify structure with a function",
		a:    NewStructure(NewNamed("a", Int), NewNamed("b", Bool)),
		b:    NewFunction(Int, Bool),
		err:  errors.New("can not unify (int)=>bool with {a:int,b:bool}"),
	}, {
		name: "fails to unify structure with a primitive",
		a:    NewStructure(NewNamed("a", Int)),
		b:    Int,
		err:  errors.New("can not unify int with {a:int}"),
	}, {
		name: "unifies a structure with a variable",
		a:    NewStructure(NewNamed("a", Int), NewNamed("b", Bool)),
		b:    NewVariable(),
		want: NewStructure(NewNamed("a", Int), NewNamed("b", Bool)),
	}, {
		name: "fails to unify on name mismatch",
		a:    NewNamed("s1", NewStructure(NewNamed("a", Int))),
		b:    NewNamed("s2", NewStructure(NewNamed("a", Int))),
		err:  errors.New("can not unify s1 with s2"),
	}, {
		name: "unifies identical recursive structures",
		a:    recursiveStruct("data", "r", NewNamed("a", Int)),
		b:    recursiveStruct("data", "r", NewNamed("a", Int)),
		want: recursiveStruct("data", "r", NewNamed("a", Int)),
	}, {
		name: "unifies structures with generic variables",
		a:    NewVariable(),
		b:    NewStructure(NewNamed("x", Int)),
		want: NewStructure(NewNamed("x", Int)),
	}, {
		name: "combines non conflicting structural variables",
		a:    NewVariable(NewNamed("x", Int)),
		b:    NewVariable(NewNamed("y", Real)),
		want: NewVariable(NewNamed("x", Int), NewNamed("y", Real)),
	}, {
		name: "fails on conflicting structural variables",
		a:    NewVariable(NewNamed("x", Int)),
		b:    NewVariable(NewNamed("x", Real)),
		err:  errors.New("can not unify int with real"),
	}, {
		name: "fails on structural variables with primitives",
		a:    NewVariable(NewNamed("x", Int)),
		b:    Int,
		err:  errors.New("can not unify V1{x:int} with int"),
	}, {
		name: "unifies structural variables with structures",
		a:    NewNamed("a", NewStructure(NewNamed("x", Int))),
		b:    NewVariable(NewNamed("x", NewVariable())),
		want: NewNamed("a", NewStructure(NewNamed("x", Int))),
	}, {
		name: "fails to unify conflicting structural variables with structures",
		a:    NewNamed("a", NewStructure(NewNamed("x", Int))),
		b:    NewVariable(NewNamed("y", NewVariable())),
		err:  errors.New("can not unify V1{y:V2} with a"),
	}, {
		name: "unifies structural variables based on the unification context",
		a:    Int,
		b:    NewVariable(NewNamed("a", NewFunction(Int, Int))),
		ctx:  MockUnificationCtx{"a": NewFunction(Int, NewFunction(Int, Int))},
		want: Int,
	}, {
		name: "unifies structural variables based on the unification context for functions",
		a:    NewFunction(Int, String),
		b:    NewVariable(NewNamed("a", NewFunction(Int, String))),
		ctx:  MockUnificationCtx{"a": NewFunction(NewVariable(), NewFunction(Int, String))},
		want: NewFunction(Int, String),
	}, {
		name: "unifies variables in structures",
		a:    NewVariable(NewNamed("a", Int), NewNamed("b", NewVariable())),
		b:    NewVariable(NewNamed("a", NewVariable()), NewNamed("b", Real)),
		want: NewVariable(NewNamed("a", Int), NewNamed("b", Real)),
	}, {
		name: "unifies structural variables with functions",
		a:    NewVariable(NewNamed("a", Int), NewNamed("b", var1)),
		b:    NewFunction(Int, var1),
		want: NewVariable(NewNamed("a", Int), NewNamed("b", var1), NewNamed("%call", NewFunction(Int, var1))),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unify(tt.a, tt.b, tt.ctx)
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

func TestType_RecursiveUnify(t *testing.T) {
	var1 := NewVariable()
	var2 := NewVariable(NewNamed("a", var1))
	ctx := MockUnificationCtx{"a": NewFunction(Int, Int)}

	result, err := Unifier(var1, var2, ctx)
	if err != nil {
		t.Errorf("Type.Unifier() error = %v", err)
		return
	}
	err = result.Update(var1.ID, Int, ctx)
	if err != nil {
		t.Errorf("Type.Update() error = %v", err)
		return
	}
	typ, _ := var2.convert(result, newSubstCtx())
	p, isP := typ.(Primitive)
	if !isP || p.ID != "int" {
		t.Errorf("result is not an int")
		return
	}
}

func TestType_RecursiveVariableUnify(t *testing.T) {
	var1 := NewVariable()
	var1.Fields["a"] = var1

	var2 := NewVariable()
	var2.Fields["a"] = var2

	ctx := MockUnificationCtx{"a": NewFunction(Int, Int)}

	result, err := Unifier(var1, var2, ctx)
	if err != nil {
		t.Errorf("Type.Unifier() error = %v", err)
		return
	}
	err = result.Add(var1, Int, ctx)
	if err != nil {
		t.Errorf("Type.Update() error = %v", err)
		return
	}
	typ, _ := var2.convert(result, newSubstCtx())
	if Signature(typ) != "int" {
		t.Errorf("result is not an int")
		return
	}
}

func TestType_RecursiveVariableFunctionUnify(t *testing.T) {
	var1 := NewVariable()
	var1.Fields["a"] = NewFunction(var1, var1)

	var2 := NewVariable()
	var2.Fields["a"] = NewFunction(var2, var2)

	ctx := MockUnificationCtx{"a": NewFunction(Int, NewFunction(Int, Int))}

	result, err := Unifier(var1, var2, ctx)
	if err != nil {
		t.Errorf("Type.Unifier() error = %v", err)
		return
	}
	err = result.Add(var1, Int, ctx)
	if err != nil {
		t.Errorf("Type.Update() error = %v", err)
		return
	}
	typ, _ := var2.convert(result, newSubstCtx())
	if Signature(typ) != "int" {
		t.Errorf("result is not an int")
		return
	}
}

func TestType_ChainUnify(tt *testing.T) {
	ctx := MockUnificationCtx{}
	vars := make([]Variable, 10)
	for i := range vars {
		vars[i] = NewVariable()
	}
	res := MakeSubstitutions()
	for i := range vars {
		if i > 0 {
			res.Update(vars[i].ID, vars[i-1], ctx)
		}
	}
	res.Add(vars[0], Int, ctx)

	for i := range vars {
		t := res.Apply(vars[i])
		if t != Int {
			tt.Error("int expected, got " + Signature(t))
		}
	}
}

func TestType_UnifyStructures(tt *testing.T) {
	ctx := MockUnificationCtx{"a": NewFunction(Int, Int)}

	var1 := NewVariable()
	var2 := NewVariable(NewNamed("a", NewFunction(var1, var1)))
	var3 := NewVariable(NewNamed("a", NewFunction(Int, var1)))

	res := MakeSubstitutions()
	if err := res.Add(var2, var3, ctx); err != nil {
		tt.Error(err)
	}

	t := res.Apply(var1)

	signt := Signature(t)
	singexp := Signature(Int)
	if signt != singexp {
		tt.Errorf("Type.Unify() got = %s", signt)
	}
}

func TestType_UnifyComplexFunctions(tt *testing.T) {
	ctx := MockUnificationCtx{"a": NewFunction(Int, Int)}

	var1 := NewVariable()
	var2 := NewVariable(NewNamed("a", var1))
	var3 := NewVariable()

	fun1 := NewFunction(var2, var2, NewFunction(var2, var3), var3)
	fun3 := NewFunction(NewVariable(), Real)
	fun4 := NewFunction(Int, Real)
	fun2 := NewFunction(NewVariable(), NewVariable(), fun3, NewVariable())

	res := MakeSubstitutions()
	if err := res.Add(var1, Int, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Add(fun2, fun1, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Add(fun3, fun4, ctx); err != nil {
		tt.Error(err)
	}

	signt := Signature(res.Apply(fun1))
	singexp := Signature(NewFunction(Int, Int, NewFunction(Int, Real), Real))
	if signt != singexp {
		tt.Errorf("Type.Unify() got = %s, expected %s", signt, singexp)
	}
}

func TestType_UnifyAllTypes(tt *testing.T) {
	ctx := MockUnificationCtx{"a": NewFunction(Real, Int)}

	var1 := NewVariable()
	var2 := NewVariable(NewNamed("a", Int))
	var3 := NewVariable()

	res := MakeSubstitutions()
	if err := res.Add(var3, Real, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Add(var1, var2, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Add(var2, var3, ctx); err != nil {
		tt.Error(err)
	}

	singexp := Signature(Real)

	signt := Signature(res.Apply(var2))
	if signt != singexp {
		tt.Errorf("Type.Unify() got = %s", signt)
	}
	signt = Signature(res.Apply(var1))
	if signt != singexp {
		tt.Errorf("Type.Unify() got = %s", signt)
	}
	signt = Signature(res.Apply(var3))
	if signt != singexp {
		tt.Errorf("Type.Unify() got = %s", signt)
	}
}

func TestType_SavesReferencesCorrectly(tt *testing.T) {
	ctx := MockUnificationCtx{}
	var1 := NewVariable()
	var2 := NewVariable(NewNamed("a", Int))
	var3 := NewVariable()
	res := MakeSubstitutions()
	if err := res.Add(var3, var1, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Add(var1, var2, ctx); err != nil {
		tt.Error(err)
	}
	if !(*res.references)[var2.ID][var1.ID] {
		tt.Error("no reference to " + var1.ID)
	}
	if !(*res.references)[var2.ID][var3.ID] {
		tt.Error("no reference to " + var3.ID)
	}
}

func TestType_PropagateContextsOnCombine(tt *testing.T) {
	ctx := MockUnificationCtx{"a": NewFunction(Int, Int)}

	var1 := NewVariable()
	var2 := NewVariable(NewNamed("a", Int))
	var3 := NewVariable()
	var4 := NewVariable()

	res := MakeSubstitutions()
	res2 := MakeSubstitutions()
	if err := res2.Add(var4, var1, ctx); err != nil {
		tt.Error(err)
	}
	if err := res2.Add(var1, var2, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Add(var3, var4, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Add(var4, Int, ctx); err != nil {
		tt.Error(err)
	}
	if err := res.Combine(res2, ctx); err != nil {
		tt.Error(err)
	}

	typ, _ := var3.convert(res, newSubstCtx())
	if typ.(Contextual).GetContext() == nil {
		tt.Error("nil context")
	}
}

type MockUnificationCtx map[string]Type

func (ctx MockUnificationCtx) StructuralTypeFor(name string, typ Type) Type {
	return ctx[name]
}
