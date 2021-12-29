package typeinference

import (
	"errors"

	"github.com/jvmakine/shine/ast"
	"github.com/jvmakine/shine/types"
	. "github.com/jvmakine/shine/types"
)

func fun(ts ...Type) *ast.Exp {
	return &ast.Exp{Op: &ast.Op{Type: function(ts...)}}
}

func union(un ...Primitive) Type {
	return Type{Variable: &TypeVar{Union: un}}
}

func function(ts ...Type) Type {
	return MakeFunction(ts...)
}

func withVar(v Type, f func(t Type) *ast.Exp) *ast.Exp {
	return f(v)
}

var global map[string]*ast.Exp = map[string]*ast.Exp{
	"+":  withVar(union(Int, Real, String), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"-":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"*":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"%":  fun(IntP, IntP, IntP),
	"/":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, t) }),
	"<":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	">":  withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	">=": withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"<=": withVar(union(Int, Real), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"||": fun(BoolP, BoolP, BoolP),
	"&&": fun(BoolP, BoolP, BoolP),
	"==": withVar(union(Int, Bool, String), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"!=": withVar(union(Int, Bool), func(t Type) *ast.Exp { return fun(t, t, BoolP) }),
	"if": withVar(MakeVariable(), func(t Type) *ast.Exp { return fun(BoolP, t, t, t) }),
}

func typeConstant(constant *ast.Const) {
	if constant.Int != nil {
		constant.Type = IntP
	} else if constant.Bool != nil {
		constant.Type = BoolP
	} else if constant.Real != nil {
		constant.Type = RealP
	} else if constant.String != nil {
		constant.Type = StringP
	} else {
		panic("invalid const")
	}
}

func typeId(id *ast.Id, ctx *ast.VisitContext) error {
	block := ctx.BlockOf(id.Name)
	if ctx.Path()[id.Name] {
		id.Type = MakeVariable()
	} else if block != nil {
		b := ctx.BlockOf(id.Name)
		ref := b.Assignments[id.Name]
		if ref != nil {
			id.Type = ref.Type().Copy(NewTypeCopyCtx())
			return nil
		}
		tc := b.TCFunctions[id.Name]
		if tc != nil {
			id.Type = tc.TypeClass.Functions[id.Name].TypeDecl
			return nil
		}
		tdef := b.TypeDefs[id.Name]
		if tdef == nil {
			panic("no id found: " + id.Name)
		}
		if tdef.Struct != nil {
			id.Type = tdef.Type().Copy(NewTypeCopyCtx())
		} else {
			return errors.New("invalid type def")
		}
	} else if p := ctx.ParamOf(id.Name); p != nil {
		id.Type = p.Type
	} else {
		return errors.New("undefined id " + id.Name)
	}
	return nil
}

func typeOp(op *ast.Op, ctx *ast.VisitContext) error {
	g := global[op.Name]
	if g == nil {
		panic("invalid op " + op.Name)
	}
	op.Type = g.Type().Copy(NewTypeCopyCtx())
	return nil
}

func typeCall(call *ast.FCall, unifier Substitutions) error {
	call.Type = MakeVariable()
	ftype := call.MakeFunType()
	s, err := ftype.Unifier(call.Function.Type())
	if err != nil {
		return err
	}
	call.Type = s.Apply(call.Type)
	for _, p := range call.Params {
		p.Convert(s)
	}
	unifier.Combine(s)
	return nil
}

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
		ts := make([]types.Type, len(def.Struct.Fields)+1)
		sf := make([]types.SField, len(def.Struct.Fields))
		for i, f := range def.Struct.Fields {
			typ := f.Type
			if !typ.IsDefined() {
				typ = MakeVariable()
			}
			f.Type = typ
			ts[i] = typ
			sf[i] = SField{
				Name: f.Name,
				Type: typ,
			}
		}

		stru := types.MakeStructure(name, sf...)
		ts[len(def.Struct.Fields)] = stru

		nt, err := rewriteNamedType(types.MakeFunction(ts...), ctx)
		if err != nil {
			return err
		}

		def.Struct.Type = nt
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
			td := ctx.TypeDef(t.Named.Name)
			if td != nil {
				if len(td.FreeVariables) != len(t.Named.TypeArguments) {
					return Type{}, errors.New("invalid number of type arguments for " + t.Named.Name)
				}
				nt := td.Type()
				for i, a := range t.Named.TypeArguments {
					o := td.FreeVariables[i]
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
	return exp.Visit(func(v *ast.Exp, ctx *ast.VisitContext) error {
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
		} else if v.Block != nil {
			ctx = ctx.SubBlock(v.Block)
			err := v.Block.CheckValueCycles()
			if err != nil {
				return err
			}
			names, err := topologicalSort(v.Block.TypeDefs)
			if err != nil {
				return err
			}
			for _, name := range names {
				value := v.Block.TypeDefs[name]
				err := rewriteNamedTypeDef(name, value, ctx)
				if err != nil {
					return err
				}
			}
			for _, name := range names {
				value := v.Block.TypeDefs[name]
				free := map[string]Type{}
				used := map[string]bool{}
				for _, n := range value.FreeVariables {
					free[n] = MakeVariable()
					used[n] = false
					d := ctx.BlockOf(n)
					if d != nil || v.Block.Assignments[n] != nil || v.Block.TypeDefs[n] != nil {
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
					typ, err := resolveTypeVariables(value.Struct.Type, free, used)
					if err != nil {
						return err
					}
					value.Struct.Type = typ
				} else if value.TypeClass != nil {
					for name, f := range value.TypeClass.Functions {
						if ctx.BlockOf(name) != nil {
							return errors.New("redefinition of " + name)
						}
						_, err := resolveTypeVariables(f.TypeDecl, free, used)
						if err != nil {
							return err
						}
						if v.Block.TCFunctions == nil {
							v.Block.TCFunctions = map[string]*ast.TypeDefinition{}
						}
						v.Block.TCFunctions[name] = value
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
		}
		return nil
	})
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
	} else if typ.Structure != nil {
		nf := make([]SField, len(typ.Structure.Fields))
		for i, f := range typ.Structure.Fields {
			na, err := resolveTypeVariables(f.Type, free, used)
			if err != nil {
				return types.Type{}, err
			}
			nf[i] = SField{
				Name: f.Name,
				Type: na,
			}
			result = MakeStructure(typ.Structure.Name, nf...)
		}
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

func rewriter(t Type, ctx *ast.VisitContext) (Type, error) {
	if t.IsNamed() {
		tdef, err := resolveNamed(t.Named.Name, ctx)
		if err != nil {
			return Type{}, err
		}
		if len(tdef.FreeVariables) != len(t.Named.TypeArguments) {
			return Type{}, errors.New("wrong number of type arguments for " + t.Named.Name)
		}
		var resolved types.Type
		if tdef.Struct != nil {
			resolved = createStructType(t.Named.Name, tdef)
		} else {
			resolved = tdef.Type()
		}
		unifier := MakeSubstitutions()
		for i, ta := range t.Named.TypeArguments {
			nt, err := rewriter(ta, ctx)
			if err != nil {
				return Type{}, err
			}
			v := tdef.VaribleMap[tdef.FreeVariables[i]]
			nt, err = nt.Unify(v)
			if err != nil {
				return Type{}, err
			}
			err = unifier.Update(v.Variable, nt)
			if err != nil {
				return Type{}, err
			}
			t.Named.TypeArguments[i] = nt
		}

		return unifier.Apply(resolved), nil
	}
	return t, nil
}

func createStructType(name string, tdef *ast.TypeDefinition) Type {
	fs := make([]types.SField, len(tdef.Struct.Fields))
	for i, f := range tdef.Struct.Fields {
		if !f.Type.IsDefined() {
			fs[i] = types.SField{
				Name: f.Name,
				Type: types.MakeVariable(),
			}
		} else {
			fs[i] = types.SField{
				Name: f.Name,
				Type: f.Type,
			}
		}
	}
	return types.MakeStructure(name, fs...)
}

func rewriteNamed(exp *ast.Exp) error {
	return exp.RewriteTypes(rewriter)
}

func Infer(exp *ast.Exp) error {
	blockCount := 0
	if err := initialiseVariables(exp); err != nil {
		return err
	}
	if err := rewriteNamed(exp); err != nil {
		return err
	}
	unifier := MakeSubstitutions()
	crawler := func(v *ast.Exp, ctx *ast.VisitContext) error {
		if v.Const != nil {
			typeConstant(v.Const)
		} else if v.Block != nil {
			blockCount++
			v.Block.ID = blockCount
			for name := range v.Block.Assignments {
				if ctx.BlockOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
			}
			for name := range v.Block.TypeDefs {
				if ctx.BlockOf(name) != nil {
					return errors.New("redefinition of " + name)
				}
			}
		} else if v.Id != nil {
			if err := typeId(v.Id, ctx); err != nil {
				return err
			}
		} else if v.Op != nil {
			if err := typeOp(v.Op, ctx); err != nil {
				return err
			}
		} else if v.Call != nil {
			if err := typeCall(v.Call, unifier); err != nil {
				return err
			}
		} else if v.Def != nil {
			v.Convert(unifier)
		} else if v.TDecl != nil {
			uni, err := v.TDecl.Type.Unifier(v.TDecl.Exp.Type())
			if err != nil {
				return err
			}
			unifier.Combine(uni)
			v.TDecl.Exp.Convert(unifier)
		} else if v.FAccess != nil {
			vari := MakeVariable()
			typ := MakeStructuralVar(map[string]Type{v.FAccess.Field: vari})
			uni, err := v.FAccess.Exp.Type().Unifier(typ)
			if err != nil {
				return err
			}
			unifier.Combine(uni)
			v.FAccess.Type = vari
		}
		return nil
	}
	// infer used code
	visited, err := exp.CrawlAfter(crawler)
	if err != nil {
		return err
	}
	// infer unused code
	err = exp.VisitAfter(func(v *ast.Exp, ctx *ast.VisitContext) error {
		if !visited[v] {
			err := crawler(v, ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
