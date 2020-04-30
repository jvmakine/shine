package resolved

import "github.com/jvmakine/shine/types"

type ClojureParam struct {
	Name string
	Type types.Primitive
}

type Clojure []ClojureParam

type ResolvedFn struct {
	ID      string
	Clojure Clojure
}
