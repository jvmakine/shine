package types

import (
	"strconv"
	"strings"
)

type signctx struct {
	varc            int
	varm            map[*TypeVar]string
	definingStructs map[string]bool
}

func (f Function) sign(ctx *signctx, level int) string {
	var sb strings.Builder
	sb.WriteString("(")
	if len(f) > 1 {
		for i, p := range f {
			sb.WriteString(sign(p, ctx, level+1))
			if i < len(f)-2 {
				sb.WriteString(",")
			} else if i < len(f)-1 {
				sb.WriteString(")=>")
			}
		}
	} else {
		sb.WriteString(")=>")
		sb.WriteString(sign(f[0], ctx, level+1))
	}
	return sb.String()
}

func (s Structure) sign(ctx *signctx, level int) string {
	if ctx.definingStructs[s.Name] {
		return s.Name
	}
	if s.Name != "" {
		ctx.definingStructs[s.Name] = true
	}
	var sb strings.Builder
	if s.Name != "" {
		sb.WriteString(s.Name)
	}
	sb.WriteString("{")
	for i, p := range s.Fields {
		sb.WriteString(p.Name)
		sb.WriteString(":")
		sb.WriteString(sign(p.Type, ctx, level+1))
		if i < len(s.Fields)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func signVar(v *TypeVar, ctx *signctx, level int) string {
	if ctx.varm[v] == "" {
		ctx.varc++
		ctx.varm[v] = "V" + strconv.Itoa(ctx.varc)
		if v.Union != nil {
			var sb strings.Builder
			sb.WriteString("[")
			for i, r := range v.Union {
				sb.WriteString(r)
				if i < len(v.Union)-1 {
					sb.WriteString("|")
				}
			}
			sb.WriteString("]")
			ctx.varm[v] += sb.String()
		} else if len(v.Structural) > 0 {
			var sb strings.Builder
			sb.WriteString("{")
			i := 0
			for k, va := range v.Structural {
				sb.WriteString(k)
				sb.WriteString(":")
				sb.WriteString(sign(va, ctx, level))
				if i < len(v.Structural)-1 {
					sb.WriteString(",")
				}
				i++
			}
			sb.WriteString("}")
			ctx.varm[v] += sb.String()
		}
	}
	return ctx.varm[v]
}

func sign(t Type, ctx *signctx, level int) string {
	if t.IsPrimitive() {
		return *t.Primitive
	}
	if t.IsVariable() {
		return signVar(t.Variable, ctx, level)
	}
	if t.IsHVariable() {
		r := signVar(t.HVariable.Root, ctx, level)
		if len(t.HVariable.Params) > 0 {
			strs := make([]string, len(t.HVariable.Params))
			for i, p := range t.HVariable.Params {
				strs[i] = sign(p, ctx, level)
			}
			r += "[" + strings.Join(strs, ",") + "]"
		}
		return r
	}
	if t.IsFunction() {
		return t.Function.sign(ctx, level)
	}
	if t.IsStructure() {
		return t.Structure.sign(ctx, level)
	}
	if t.IsNamed() {
		post := ""
		if len(t.Named.TypeArguments) > 0 {
			strs := make([]string, len(t.Named.TypeArguments))
			for i, a := range t.Named.TypeArguments {
				strs[i] = sign(a, ctx, level)
			}
			post = "[" + strings.Join(strs, ",") + "]"
		}
		return t.Named.Name + post
	}
	if t.IsTypeClassRef() {
		res := t.TCRef.TypeClass + "["
		args := make([]string, len(t.TCRef.TypeClassVars))
		for i, a := range t.TCRef.TypeClassVars {
			if i == t.TCRef.Place && len(t.TCRef.TypeClassVars) > 1 {
				args[i] = "{" + sign(a, ctx, level) + "}"
			} else {
				args[i] = sign(a, ctx, level)
			}
		}
		return res + strings.Join(args, ",") + "]"
	}
	if !t.IsDefined() {
		return "<undefined>"
	}
	panic("invalid type")
}

func (t Type) Signature() string {
	varm := map[*TypeVar]string{}
	ds := map[string]bool{}
	ctx := signctx{varc: 0, varm: varm, definingStructs: ds}
	return sign(t, &ctx, 0)
}

// TODO: Remove
func (t Type) TSignature() string {
	return t.Signature()
}
