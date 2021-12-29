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

func sign(t Type, ctx *signctx, level int) string {
	if t.IsPrimitive() {
		return *t.Primitive
	}
	if t.IsVariable() && !t.IsFunction() {
		if ctx.varm[t.Variable] == "" {
			ctx.varc++
			ctx.varm[t.Variable] = "V" + strconv.Itoa(ctx.varc)
			if t.IsUnionVar() {
				var sb strings.Builder
				sb.WriteString("[")
				for i, r := range t.Variable.Union {
					sb.WriteString(r)
					if i < len(t.Variable.Union)-1 {
						sb.WriteString("|")
					}
				}
				sb.WriteString("]")
				ctx.varm[t.Variable] += sb.String()
			} else if t.IsStructuralVar() {
				var sb strings.Builder
				sb.WriteString("{")
				i := 0
				for k, v := range t.Variable.Structural {
					sb.WriteString(k)
					sb.WriteString(":")
					sb.WriteString(sign(v, ctx, level))
					if i < len(t.Variable.Structural)-1 {
						sb.WriteString(",")
					}
					i++
				}
				sb.WriteString("}")
				ctx.varm[t.Variable] += sb.String()
			}
		}
		return ctx.varm[t.Variable]
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
				args[i] = "[[" + sign(a, ctx, level) + "]]"
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
