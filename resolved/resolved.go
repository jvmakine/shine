package resolved

import "github.com/jvmakine/shine/types"

type ClosureParam struct {
	Name string
	Type types.Type
}

type Closure []ClosureParam

type ResolvedFnCall struct {
	ID string
}

type ResolvedFnDef struct {
	Closure Closure
}
