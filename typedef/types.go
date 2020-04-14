package typedef

import (
	"strings"
)

type TypeSignature = string

type Primitive = string

const (
	Int  Primitive = "int"
	Bool Primitive = "bool"
)

type Function = []*Type

type Type struct {
	Primitive *Primitive
	Function  *Function
}

func (t *Type) IsPrimitive() bool {
	return t.Primitive != nil
}

func (t *Type) IsFunction() bool {
	return t.Function != nil
}

func (t *Type) Signature() TypeSignature {
	if t.IsPrimitive() {
		return *t.Primitive
	}
	if t.IsFunction() {
		var sb strings.Builder
		sb.WriteString("(")
		for i, p := range *t.Function {
			sb.WriteString(p.Signature())
			if i < len(*t.Function)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		return sb.String()
	}
	panic("invalid type")
}

func primitive(t Primitive) *Type {
	return &Type{Primitive: &t}
}

func function(ts ...*Type) *Type {
	return &Type{Function: &ts}
}
