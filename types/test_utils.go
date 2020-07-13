package types

func recursiveStruct(name string, recField string, fields ...Named) Type {
	s := NewStructure(fields...)
	s.Fields = append(s.Fields, NewNamed(recField, NewNamed(name, s)))
	n := NewNamed(name, s)
	return n
}

func WithType(typ Type, f func(t Type) Type) Type {
	return f(typ)
}
