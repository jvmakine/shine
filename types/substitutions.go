package types

type Substitutions struct {
	substitutions map[*TypeVar]Type
	references    map[*TypeVar]map[*TypeVar]bool
}

func MakeSubstitutions() Substitutions {
	return Substitutions{
		substitutions: map[*TypeVar]Type{},
		references:    map[*TypeVar]map[*TypeVar]bool{},
	}
}

type substCtx struct {
	visited map[*Structure]bool
}

func (s Substitutions) Apply(t Type) Type {
	return apply(s, t, &substCtx{visited: map[*Structure]bool{}})
}

func apply(s Substitutions, t Type, ctx *substCtx) Type {
	target := s.substitutions[t.Variable]
	if !target.IsDefined() {
		target = t
	}
	if target.IsFunction() {
		ntyps := make([]Type, len(target.FunctTypes()))
		for i, v := range target.FunctTypes() {
			ntyps[i] = apply(s, v, ctx)
		}
		return MakeFunction(ntyps...)
	}
	if struc := target.Structure; struc != nil {
		if ctx.visited[struc] {
			return target
		}
		ctx.visited[struc] = true
		fields := make([]SField, len(struc.Fields))
		for i, v := range struc.Fields {
			fields[i] = SField{
				Name: v.Name,
				Type: apply(s, v.Type, ctx),
			}
		}
		nt := make([]Type, len(struc.TypeArguments))
		for i, a := range struc.TypeArguments {
			nt[i] = apply(s, a, ctx)
		}
		return Type{Structure: &Structure{
			Name:          struc.Name,
			TypeArguments: nt,
			Fields:        fields,

			OrginalVars:    struc.OrginalVars,
			OriginalFields: struc.OriginalFields,
		}}
	}
	if target.IsStructuralVar() {
		for k, v := range target.Variable.Structural {
			target.Variable.Structural[k] = apply(s, v, ctx)
		}
	}
	if r := target.TCRef; r != nil {
		fields := make([]Type, len(r.TypeClassVars))
		for i, a := range r.TypeClassVars {
			fields[i] = apply(s, a, ctx)
		}
		nr := MakeTypeClassRef(r.TypeClass, r.Place, fields...)
		if !nr.TCRef.TypeClassVars[nr.TCRef.Place].HasFreeVars() {
			return nr.TCRef.TypeClassVars[nr.TCRef.Place]
		}
		return nr
	}
	if target.IsHVariable() {
		for i, a := range target.HVariable.Params {
			target.HVariable.Params[i] = apply(s, a, ctx)
		}
		if ss := s.substitutions[target.HVariable.Root]; ss.IsDefined() {
			res := s.substitutions[target.HVariable.Root]
			return res.Instantiate(target.HVariable.Params)
		}
	}
	return target
}

func (s Substitutions) Update(from *TypeVar, to Type) error {
	if from == to.Variable {
		return nil
	}

	result := s.Apply(to)

	if p := s.substitutions[from]; p.IsDefined() {
		if result != p {
			uni, err := result.Unifier(p)
			if err != nil {
				return err
			}
			s.Combine(uni)
			result = uni.Apply(result)
		}
	}

	s.substitutions[from] = result

	for _, fv := range result.FreeVars() {
		if s.references[fv] == nil {
			s.references[fv] = map[*TypeVar]bool{}
		}
		s.references[fv][from] = true
	}

	if rs := s.references[from]; rs != nil {
		s.references[from] = nil
		subs := MakeSubstitutions()
		subs.Update(from, result)
		for k := range rs {
			substit := s.substitutions[k]
			s.substitutions[k] = subs.Apply(substit)
			for _, fv := range s.substitutions[k].FreeVars() {
				if s.references[fv] == nil {
					s.references[fv] = map[*TypeVar]bool{}
				}
				s.references[fv][from] = true
			}
		}
	}

	return nil
}

func (s Substitutions) Combine(o Substitutions) error {
	for f, t := range o.substitutions {
		err := s.Update(f, t)
		if err != nil {
			return err
		}
	}
	return nil
}
