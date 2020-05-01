package resolved

import "github.com/jvmakine/shine/types"

type ClojureParam struct {
	Name string
	Type types.Type
}

type Clojure []ClojureParam

type ResolvedFnCall struct {
	ID string
}

type ResolvedFnDef struct {
	Clojure Clojure
}
