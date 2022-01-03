package types

func recursiveStruct(name string, recField string, fields ...SField) Type {
	s := MakeStructure(name, fields...)
	s.Structure.Fields = append(s.Structure.Fields, SField{recField, s})
	return s
}
