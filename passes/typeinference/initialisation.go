package typeinference

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
)

func namedVariables(t *ast.TypeDefinition) map[string]bool {
	res := map[string]bool{}
	if t.Struct != nil {
		for _, f := range t.Struct.Fields {
			u := f.Type.NamedTypes()
			for n := range u {
				res[n] = true
			}
		}
	} else {
		res = t.TypeDecl.NamedTypes()
	}
	return res
}

func topologicalSort(defs map[string]*ast.TypeDefinition) ([]string, error) {
	result := []string{}
	todo := map[string]*ast.TypeDefinition{}
	for n, t := range defs {
		todo[n] = t
	}
	for len(todo) > 0 {
		prevLen := len(todo)
		for n, t := range todo {
			used := namedVariables(t)
			hit := false
			for u := range used {
				if todo[u] != nil {
					hit = true
					break
				}
			}
			if !hit {
				result = append(result, n)
				delete(todo, n)
			}
		}
		if len(todo) == prevLen {
			return nil, errors.New("recursive type declaration")
		}
	}
	return result, nil
}

func rewriteNamedTypeDef(name string, def *ast.TypeDefinition, ctx *ast.VisitContext) error {
	if def.Struct != nil {
		for _, f := range def.Struct.Fields {
			nt, err := rewriteNamedType(f.Type, ctx)
			if err != nil {
				return err
			}
			f.Type = nt
		}
		sf := make([]types.SField, len(def.Struct.Fields))
		for i, f := range def.Struct.Fields {
			typ := f.Type
			if !typ.IsDefined() {
				typ = MakeVariable()
			}
			f.Type = typ
			nt, err := rewriteNamedType(typ, ctx)
			if err != nil {
				return err
			}
			sf[i] = SField{
				Name: f.Name,
				Type: nt,
			}
		}

		stru := types.MakeStructure(name, sf...)
		def.Struct.Type = stru
	} else {
		nt, err := rewriteNamedType(def.TypeDecl, ctx)
		if err != nil {
			return err
		}
		def.TypeDecl = nt
	}
	return nil
}

func rewriteNamedType(from Type, ctx *ast.VisitContext) (Type, error) {
	return from.Rewrite(func(t Type) (Type, error) {
		if t.IsNamed() {
			tdef := ctx.TypeDef(t.Named.Name)
			if tdef != nil {
				if len(tdef.FreeVariables) != len(t.Named.TypeArguments) {
					return Type{}, errors.New("invalid number of type arguments for " + t.Named.Name)
				}
				nt := tdef.Type()
				for i, a := range t.Named.TypeArguments {
					o := tdef.FreeVariables[i]
					nt, _ = nt.Rewrite(func(t Type) (Type, error) {
						if t.IsNamed() && t.Named.Name == o {
							return a, nil
						}
						return t, nil
					})
				}
				return nt, nil
			}
			return t, nil
		}
		return t, nil
	})
}

func initialiseVariables(exp *ast.Exp) error {
	err := exp.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Def != nil {
			for _, p := range v.Def.Params {
				name := p.Name
				if ctx.BlockOf(name) != nil || ctx.ParamOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
				if !p.Type.IsDefined() {
					p.Type = MakeVariable()
				}
			}
		} else if block := v.Block; block != nil {
			return initialiseBlock(block, ctx)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return exp.RewriteTypes(nameRewriter)
}

func initialiseBlock(block *ast.Block, ctx *ast.VisitContext) error {
	ctx = ctx.SubBlock(block)
	err := block.CheckValueCycles()
	if err != nil {
		return err
	}
	names, err := topologicalSort(block.TypeDefs)
	if err != nil {
		return err
	}
	for _, name := range names {
		tdef := block.TypeDefs[name]
		if err := rewriteNamedTypeDef(name, tdef, ctx); err != nil {
			return err
		}
	}
	for _, name := range names {
		value := block.TypeDefs[name]
		free := map[string]Type{}
		used := map[string]bool{}
		index := map[string]int{}
		freeFields := make([]Type, len(value.FreeVariables))
		for i, n := range value.FreeVariables {
			free[n] = MakeVariable()
			used[n] = false
			index[n] = i
			freeFields[i] = free[n]
			d := ctx.BlockOf(n)
			if d != nil || block.Assignments[n] != nil || block.TypeDefs[n] != nil {
				return errors.New("redefinition of " + n)
			}
		}

		if value.Struct != nil {
			for _, f := range value.Struct.Fields {
				typ := f.Type
				typ, err = resolveTypeVariables(typ, free, used)
				if err != nil {
					return err
				}
				f.Type = typ
			}

			value.Struct.Type.Structure.TypeArguments = make([]Type, len(value.FreeVariables))
			for i, n := range value.FreeVariables {
				value.Struct.Type.Structure.TypeArguments[i] = free[n]
			}

			typ, err := resolveTypeVariables(value.Struct.Type, free, used)
			if err != nil {
				return err
			}
			value.Struct.Type = types.MakeStructure(name, typ.Structure.Fields...)
		} else if value.TypeClass != nil {
			for fname, f := range value.TypeClass.Functions {
				if ctx.BlockOf(fname) != nil {
					return errors.New("redefinition of " + fname)
				}

				nt, err := f.TypeDecl.Rewrite(func(a Type) (Type, error) {
					if a.IsNamed() && free[a.Named.Name].IsDefined() {
						used[a.Named.Name] = true
						idx := index[a.Named.Name]

						it := free[a.Named.Name]
						if len(a.Named.TypeArguments) > 0 {
							rws := make([]Type, len(a.Named.TypeArguments))
							for i, t := range a.Named.TypeArguments {
								nt, err := resolveTypeVariables(t, free, used)
								if err != nil {
									return Type{}, err
								}
								rws[i] = nt
							}

							it = types.MakeHierarchicalVar(it.Variable, rws...)
						}

						ff := make([]Type, len(freeFields))
						for i, f := range freeFields {
							ff[i] = f
						}
						ff[idx] = it

						nt := types.MakeTypeClassRef(name, idx, ff...)
						return nt, nil
					}
					return a, nil
				})
				if err != nil {
					return err
				}
				f.TypeDecl = nt

				nfree := map[string]Type{}
				nused := map[string]bool{}
				for n, t := range free {
					nfree[n] = t
					nused[n] = used[n]
				}
				for _, n := range f.FreeVariables {
					if nfree[n].IsDefined() {
						return errors.New("redefinition of " + n)
					}
					nfree[n] = types.MakeVariable()
				}

				nt, err = resolveTypeVariables(f.TypeDecl, nfree, nused)
				if err != nil {
					return err
				}
				f.TypeDecl = nt

				for n := range nused {
					if free[n].IsDefined() {
						used[n] = used[n] || nused[n]
					}
				}

				if block.TCFunctions == nil {
					block.TCFunctions = map[string]*ast.TypeDefinition{}
				}
				block.TCFunctions[fname] = value
			}
		} else {
			typ, err := resolveTypeVariables(value.TypeDecl, free, used)
			if err != nil {
				return err
			}
			value.TypeDecl = typ
		}
		value.VaribleMap = free
		for n, b := range used {
			if !b {
				return errors.New("unused free type " + n)
			}
		}
	}
	return nil
}

func resolveTypeVariables(typ types.Type, free map[string]Type, used map[string]bool) (types.Type, error) {
	result := typ
	if typ.Named != nil {
		if fv, ok := free[typ.Named.Name]; ok {
			used[typ.Named.Name] = true
			result = fv
		}
		for _, ta := range typ.Named.TypeArguments {
			_, err := resolveTypeVariables(ta, free, used)
			if err != nil {
				return types.Type{}, err
			}
		}
	} else if typ.Function != nil {
		args := make([]types.Type, len(*typ.Function))
		for i, a := range *typ.Function {
			na, err := resolveTypeVariables(a, free, used)
			if err != nil {
				return types.Type{}, err
			}
			args[i] = na
		}
		result = MakeFunction(args...)
	} else if s := typ.Structure; s != nil {
		nf := make([]SField, len(s.Fields))
		for i, f := range s.Fields {
			na, err := resolveTypeVariables(f.Type, free, used)
			if err != nil {
				return types.Type{}, err
			}
			nf[i] = SField{
				Name: f.Name,
				Type: na,
			}
			result = Type{Structure: &Structure{
				Name:           s.Name,
				TypeArguments:  s.TypeArguments,
				Fields:         nf,
				OrginalVars:    s.OrginalVars,
				OriginalFields: s.OriginalFields,
			}}
		}
	} else if typ.HVariable != nil {
		args := make([]types.Type, len(typ.HVariable.Params))
		for i, a := range typ.HVariable.Params {
			na, err := resolveTypeVariables(a, free, used)
			if err != nil {
				return types.Type{}, err
			}
			args[i] = na
		}
		result = MakeHierarchicalVar(typ.HVariable.Root, args...)
	} else if typ.TCRef != nil {
		args := make([]types.Type, len(typ.TCRef.TypeClassVars))
		for i, a := range typ.TCRef.TypeClassVars {
			na, err := resolveTypeVariables(a, free, used)
			if err != nil {
				return types.Type{}, err
			}
			args[i] = na
		}
		result = MakeTypeClassRef(typ.TCRef.TypeClass, typ.TCRef.Place, args...)
	}
	return result, nil
}

func resolveNamed(name string, ctx *ast.VisitContext) (*ast.TypeDefinition, error) {
	tdef := ctx.TypeDef(name)
	if tdef == nil {
		return nil, errors.New("type " + name + " is undefined")
	}
	return tdef, nil
}

func nameRewriter(t Type, ctx *ast.VisitContext) (Type, error) {
	if t.IsNamed() {
		tdef, err := resolveNamed(t.Named.Name, ctx)
		if err != nil {
			return Type{}, err
		}
		if len(t.Named.TypeArguments) != 0 && len(tdef.FreeVariables) != len(t.Named.TypeArguments) {
			return Type{}, errors.New("wrong number of type arguments for " + t.Named.Name)
		}
		resolved := tdef.Type()
		unifier := MakeSubstitutions()
		for i, ta := range t.Named.TypeArguments {
			nt, err := nameRewriter(ta, ctx)
			if err != nil {
				return Type{}, err
			}
			v := tdef.VaribleMap[tdef.FreeVariables[i]]
			nt, err = nt.Unify(v, ctx)
			if err != nil {
				return Type{}, err
			}
			err = unifier.Update(v.Variable, nt, ctx)
			if err != nil {
				return Type{}, err
			}
			t.Named.TypeArguments[i] = nt
		}

		return unifier.Apply(resolved, ctx), nil
	}
	return t, nil
}
