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

type TypeUnion struct {
	Types []*TypePtr
}

type TypeDef struct {
	Base  *Primitive
	Fn    []*TypePtr
	Union *TypeUnion
}

type TypePtr struct {
	Def *TypeDef
}

func MakeFun(ts ...*TypePtr) *TypePtr {
	return &TypePtr{&TypeDef{Fn: ts}}
}

type signctx struct {
	varc int
	varm map[*TypeDef]string
}

func sign(t *TypePtr, ctx *signctx) string {
	if t.IsBase() {
		return *t.Def.Base
	}
	if t.IsFunction() {
		var sb strings.Builder
		sb.WriteString("(")
		for i, p := range t.Def.Fn {
			sb.WriteString(sign(p, ctx))
			if i < len(t.Def.Fn)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		return sb.String()
	}
	if t.IsUnion() {
		var sb strings.Builder
		sb.WriteString("(")
		for i, p := range t.Def.Union.Types {
			sb.WriteString(sign(p, ctx))
			if i < len(t.Def.Union.Types)-1 {
				sb.WriteString("|")
			}
		}
		sb.WriteString(")")
		return sb.String()
	}
	if ctx.varm[t.Def] == "" {
		ctx.varc++
		ctx.varm[t.Def] = "V" + strconv.Itoa(ctx.varc)
	}
	return ctx.varm[t.Def]
}

func (t *TypePtr) Signature() string {
	varm := map[*TypeDef]string{}
	return sign(t, &signctx{varc: 0, varm: varm})
}

type TypeCopyCtx struct {
	defs map[*TypeDef]*TypeDef
	ptrs map[*TypePtr]*TypePtr
}

func NewTypeCopyCtx() *TypeCopyCtx {
	return &TypeCopyCtx{
		defs: map[*TypeDef]*TypeDef{},
		ptrs: map[*TypePtr]*TypePtr{},
	}

}

func (t *TypePtr) Copy(ctx *TypeCopyCtx) *TypePtr {
	var params []*TypePtr = nil
	var def *TypeDef = nil
	if ctx.ptrs[t] != nil {
		return ctx.ptrs[t]
	}
	if ctx.defs[t.Def] != nil {
		res := &TypePtr{Def: ctx.defs[t.Def]}
		ctx.ptrs[t] = res
		return res
	}
	if t.Def.Fn != nil {
		params = make([]*TypePtr, len(t.Def.Fn))
		for i, p := range t.Def.Fn {
			if ctx.ptrs[p] != nil {
				params[i] = ctx.ptrs[p]
			} else {
				if ctx.defs[p.Def] == nil {
					ctx.defs[p.Def] = p.Copy(ctx).Def
				}
				res := &TypePtr{Def: ctx.defs[p.Def]}
				ctx.ptrs[p] = res
				params[i] = res
			}
		}
	}
	def = &TypeDef{
		Fn:   params,
		Base: t.Def.Base,
	}
	ctx.defs[t.Def] = def
	res := &TypePtr{Def: def}
	ctx.ptrs[t] = res
	return res
}

func (t *TypePtr) IsFunction() bool {
	return t.Def.Fn != nil
}

func (t *TypePtr) IsBase() bool {
	return t.Def.Base != nil
}

func (t *TypePtr) IsVariable() bool {
	return !t.IsBase() && !t.IsFunction()
}

func (t *TypePtr) IsUnion() bool {
	return t.Def.Union != nil
}

func (t *TypePtr) IsDefined() bool {
	if t.IsVariable() {
		return false
	}
	if t.IsBase() {
		return true
	}
	for _, p := range t.Def.Fn {
		if !p.IsDefined() {
			return false
		}
	}
	return true
}

func (t *TypePtr) ReturnType() *TypePtr {
	return t.Def.Fn[len(t.Def.Fn)-1]
}
