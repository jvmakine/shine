package types

import (
	"strconv"
	"strings"
)

type Primitive = string

const (
	Int  Primitive = "int"
	Bool Primitive = "bool"
	Real Primitive = "real"
)

type TVarID = int64

type TypeVar struct {
	Restrictions []Primitive
}

type Type struct {
	Function  *[]Type
	Variable  *TypeVar
	Primitive *Primitive
}

func MakeFun(ts ...Type) Type {
	return Type{Function: &ts}
}

func (t Type) FreeVars() []*TypeVar {
	if t.Variable != nil {
		return []*TypeVar{t.Variable}
	}
	if t.Function != nil {
		res := []*TypeVar{}
		for _, p := range *t.Function {
			res = append(res, p.FreeVars()...)
		}
		return res
	}
	return []*TypeVar{}
}

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

type TypeCopyCtx struct {
	vars map[*TypeVar]*TypeVar
}

func NewTypeCopyCtx() *TypeCopyCtx {
	return &TypeCopyCtx{
		vars: map[*TypeVar]*TypeVar{},
	}

}

func (t Type) Copy(ctx *TypeCopyCtx) Type {
	if t.Variable != nil {
		if ctx.vars[t.Variable] == nil {
			ctx.vars[t.Variable] = t.Variable.Copy()
		}
		return Type{Variable: ctx.vars[t.Variable]}
	}
	if t.Function != nil {
		ps := make([]Type, len(*t.Function))
		for i, p := range *t.Function {
			ps[i] = p.Copy(ctx)
		}
		return Type{Function: &ps}
	}
	return t
}

func (t *TypeVar) Copy() *TypeVar {
	return &TypeVar{
		Restrictions: t.Restrictions,
	}
}

func (t Type) IsFunction() bool {
	return t.Function != nil
}

func (t Type) IsPrimitive() bool {
	return t.Primitive != nil
}

func (t Type) IsVariable() bool {
	return t.Variable != nil
}

func (t Type) IsDefined() bool {
	return t.Function != nil || t.Variable != nil || t.Primitive != nil
}

func (t Type) HasFreeVars() bool {
	return len(t.FreeVars()) == 0
}

func (t Type) AsPrimitive() Primitive {
	if !t.IsPrimitive() {
		panic("type not primitive type")
	}
	return *t.Primitive
}

func (t Type) ReturnType() Type {
	return (*t.Function)[len(*t.Function)-1]
}
