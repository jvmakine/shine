package types

type ClosureParam struct {
	Name string
	Type Type
}

type Closure []ClosureParam
