// Package resolved contains an AST for program where
// all the types have been resolved, functions converted to global
// names and different type instantations of functions converted to
// different functions
package resolved

type FunctionId = string

type FunctionRef struct {
	Name      string
	BlockId   string
	Signature string
}

func (f *FunctionRef) Id() FunctionId {
	return f.Name + "%%" + f.BlockId + "%%" + f.Signature
}

type ResolvedExp struct {
}

type Resolved struct {
	Functions map[FunctionId]bool
}
