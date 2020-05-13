package typeinference

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/jvmakine/shine/types"
)

func TestSimpleReplacements(t *testing.T) {
	lens := MakeTLens()
	vari := MakeVariable()
	fun := MakeFunction(vari, vari, IntP)

	lens.Update(vari.Variable, RealP)
	result := lens.Convert(fun)
	want := MakeFunction(RealP, RealP, IntP)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Convert() result = %v, want %v", result, want)
	}

	result = lens.Convert(vari)
	want = RealP
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Convert() result = %v, want %v", result, want)
	}

	err := lens.Update(vari.Variable, IntP)
	wantErr := errors.New("can not unify int with real")
	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Convert() error = %v, want %v", err, wantErr)
	}

	err = lens.Update(vari.Variable, RealP)
	wantErr = nil
	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Convert() error = %v, want %v", err, wantErr)
	}

	vari2 := MakeVariable()
	lens.Update(vari2.Variable, vari)

	fun = MakeFunction(vari2, vari2, IntP)
	result = lens.Convert(fun)
	want = MakeFunction(RealP, RealP, IntP)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Convert() result = %v, want %v", result, want)
	}

	err = lens.Update(vari2.Variable, IntP)
	wantErr = errors.New("can not unify int with real")
	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Convert() error = %v, want %v", err, wantErr)
	}

	vari3, vari4 := MakeVariable(), MakeVariable()
	lens.Update(vari3.Variable, IntP)
	lens.Update(vari4.Variable, RealP)

	err = lens.Update(vari3.Variable, vari4)
	wantErr = errors.New("can not unify int with real")
	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Convert() error = %v, want %v", err, wantErr)
	}
}

func TestFunctionReplacements1(t *testing.T) {
	lens := MakeTLens()
	vari, fvari := MakeVariable(), MakeVariable()
	fun := MakeFunction(fvari, IntP)

	lens.Update(vari.Variable, fun)
	lens.Update(fvari.Variable, IntP)

	result := lens.Convert(vari)
	want := MakeFunction(IntP, IntP)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Convert() result = %v, want %v", result, want)
	}
}

func TestFunctionReplacements2(t *testing.T) {
	lens := MakeTLens()
	vari, fvari := MakeVariable(), MakeVariable()
	fun := MakeFunction(fvari, IntP)

	lens.Update(fvari.Variable, IntP)
	lens.Update(vari.Variable, fun)

	result := lens.Convert(vari)
	want := MakeFunction(IntP, IntP)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Convert() result = %v, want %v", result, want)
	}
}

func TestFunctionReplacements3(t *testing.T) {
	lens := MakeTLens()
	vari1, vari2, fvari := MakeVariable(), MakeVariable(), MakeVariable()

	lens.Update(vari1.Variable, MakeFunction(fvari, IntP))
	lens.Update(vari2.Variable, MakeFunction(IntP, fvari))
	lens.Update(vari2.Variable, vari1)

	result := lens.Convert(fvari)
	want := IntP
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Convert() result = %v, want %v", result, want)
	}
}

func TestFunctionReplacements4(t *testing.T) {
	lens := MakeTLens()
	vari1, vari2, fvari := MakeVariable(), MakeVariable(), MakeVariable()

	lens.Update(vari1.Variable, MakeFunction(fvari, IntP))
	lens.Update(vari2.Variable, MakeFunction(RealP, fvari))
	err := lens.Update(vari2.Variable, vari1)

	wantErr := errors.New("can not unify int with real")
	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Convert() error = %v, want %v", err, wantErr)
	}
}
