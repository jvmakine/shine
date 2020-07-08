package types

type Union []Type

func (u Union) deduplicate() Union {
	res := Union{}
	for _, tt := range u {
		if len(res) == 0 {
			res = Union{tt}
		} else {
			shouldAppend := false
			for i, rt := range res {
				if tt.IsGeneralisationOf(rt) && !rt.HasFreeVars() {
					if len(res) > 1 {
						res = append(res[:i], res[i+1:]...)
					} else {
						res = Union{}
					}
					shouldAppend = true
				} else if !rt.IsGeneralisationOf(tt) {
					shouldAppend = true
				}
			}
			if shouldAppend {
				res = append(res, tt)
			}
		}
	}
	return res
}

func (r Union) Unify(o Type) (Type, error) {
	union := Union{}
	ous := Union{o}
	if o.IsUnionVar() {
		ous = o.Variable.Union
	}
	for _, rp := range r {
		for _, op := range ous {
			un, err := rp.Unify(op)
			if err == nil {
				union = append(union, un)
			}
		}
	}
	res := union.deduplicate()
	if len(res) == 0 {
		return Type{}, UnificationError(MakeUnionVar(r...), o)
	}
	if len(res) == 1 {
		return res[0], nil
	}
	return MakeUnionVar(res...), nil
}

func (t Type) AddToUnion(o Type) Type {
	var union Union
	if t.IsUnionVar() {
		union = append(t.Variable.Union, o).deduplicate()
	} else {
		ts := Union{t}
		ts = append(ts, o)
		union = ts.deduplicate()
	}
	if len(union) == 1 {
		return union[0]
	}
	return MakeUnionVar(union...)
}
