package types

import (
	"strconv"
	"strings"
)

type signctx struct {
	varc int
	varm map[*TypeVar]string
}

func (f Function) sign(ctx *signctx) string {
	var sb strings.Builder
	sb.WriteString("(")
	if len(f) > 1 {
		for i, p := range f {
			sb.WriteString(sign(p, ctx))
			if i < len(f)-2 {
				sb.WriteString(",")
			} else if i < len(f)-1 {
				sb.WriteString(")=>")
			}
		}
	} else {
		sb.WriteString(")=>")
		sb.WriteString(sign(f[0], ctx))
	}
	return sb.String()
}

func (s Structure) sign(ctx *signctx) string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, p := range s {
		sb.WriteString(p.Name)
		sb.WriteString(":")
		sb.WriteString(sign(p.Type, ctx))
		if i < len(s)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func sign(t Type, ctx *signctx) string {
	if t.IsPrimitive() {
		return *t.Primitive
	}
	if t.IsVariable() && !t.IsFunction() {
		if ctx.varm[t.Variable] == "" {
			ctx.varc++
			ctx.varm[t.Variable] = "V" + strconv.Itoa(ctx.varc)
			if t.IsRestrictedVariable() {
				var sb strings.Builder
				sb.WriteString("[")
				for i, r := range t.Variable.Restrictions {
					sb.WriteString(r)
					if i < len(t.Variable.Restrictions)-1 {
						sb.WriteString("|")
					}
				}
				sb.WriteString("]")
				ctx.varm[t.Variable] += sb.String()
			}
		}
		return ctx.varm[t.Variable]
	}
	if t.IsFunction() {
		return t.Function.sign(ctx)
	}
	if t.IsStructure() {
		return t.Structure.sign(ctx)
	}
	if !t.IsDefined() {
		return "<undefined>"
	}
	panic("invalid type")
}

func (t Type) Signature() string {
	varm := map[*TypeVar]string{}
	return sign(t, &signctx{varc: 0, varm: varm})
}
