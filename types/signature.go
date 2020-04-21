package types

import (
	"strconv"
	"strings"
)

type signctx struct {
	varc int
	varm map[*TypeVar]string
}

func sign(t Type, ctx *signctx) string {
	if t.IsPrimitive() {
		return *t.Primitive
	}
	if t.IsFunction() {
		var sb strings.Builder
		sb.WriteString("(")
		for i, p := range *t.Function {
			sb.WriteString(sign(p, ctx))
			if i < len(*t.Function)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		return sb.String()
	}
	if t.IsVariable() {
		if ctx.varm[t.Variable] == "" {
			ctx.varc++
			ctx.varm[t.Variable] = "V" + strconv.Itoa(ctx.varc)
			if len(t.Variable.Restrictions) > 0 {
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
	if !t.IsDefined() {
		panic("can not get signature from undefined type")
	}
	panic("invalid type")
}

func (t Type) Signature() string {
	varm := map[*TypeVar]string{}
	return sign(t, &signctx{varc: 0, varm: varm})
}
