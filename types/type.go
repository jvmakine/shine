package types

import (
	"strconv"
	"strings"
)

type Primitive = string

const (
	Int  Primitive = "int"
	Bool Primitive = "bool"
)

type TypeDef struct {
	Base *Primitive
	Fn   []*TypePtr
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

type TypeCopyCtx = *map[*TypeDef]*TypeDef

func NewTypeCopyCtx() TypeCopyCtx {
	return &map[*TypeDef]*TypeDef{}
}

func (t *TypePtr) Copy(ctx TypeCopyCtx) *TypePtr {
	var params []*TypePtr = nil
	var def *TypeDef = nil
	if (*ctx)[t.Def] != nil {
		return &TypePtr{Def: (*ctx)[t.Def]}
	}
	if t.Def.Fn != nil {
		params = make([]*TypePtr, len(t.Def.Fn))
		for i, p := range t.Def.Fn {
			if (*ctx)[p.Def] == nil {
				(*ctx)[p.Def] = p.Copy(ctx).Def
			}
			params[i] = &TypePtr{Def: (*ctx)[p.Def]}
		}
	}
	def = &TypeDef{
		Fn:   params,
		Base: t.Def.Base,
	}
	(*ctx)[t.Def] = def
	return &TypePtr{Def: def}
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
