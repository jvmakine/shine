package typeinferer

import (
	"strconv"
	"strings"

	"github.com/jvmakine/shine/typedef"
)

type TypeDef struct {
	Base *typedef.Primitive
	Fn   []*TypePtr
}

type TypePtr struct {
	Def *TypeDef
}

func sign(t *TypePtr, varc *int, varm *map[*TypeDef]string) string {
	if t.IsBase() {
		return *t.Def.Base
	}
	if t.IsFunction() {
		var sb strings.Builder
		sb.WriteString("(")
		for i, p := range t.Def.Fn {
			sb.WriteString(p.Signature())
			if i < len(t.Def.Fn)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		return sb.String()
	}
	if (*varm)[t.Def] == "" {
		*varc = *varc + 1
		(*varm)[t.Def] = "V" + strconv.Itoa(*varc)
	}
	return (*varm)[t.Def]
}

func (t *TypePtr) Signature() string {
	varc := 0
	varm := map[*TypeDef]string{}
	return sign(t, &varc, &varm)
}

func base(t typedef.Primitive) *TypePtr {
	return &TypePtr{Def: &TypeDef{Base: &t}}
}

func function(ts ...*TypePtr) *TypePtr {
	return &TypePtr{&TypeDef{Fn: ts}}
}

func variable() *TypePtr {
	return &TypePtr{Def: &TypeDef{}}
}

func (t *TypePtr) Copy() *TypePtr {
	var params []*TypePtr = nil
	var def *TypeDef = nil
	if t.Def != nil {
		if t.Def.Fn != nil {
			params = make([]*TypePtr, len(t.Def.Fn))
			var seen map[*TypePtr]*TypePtr = map[*TypePtr]*TypePtr{}
			for i, p := range t.Def.Fn {
				if seen[p] == nil {
					seen[p] = p.Copy()
				}
				params[i] = seen[p]
			}
		}
		def = &TypeDef{
			Fn:   params,
			Base: t.Def.Base,
		}
	}

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

func (t *TypePtr) ReturnType() *TypePtr {
	return t.Def.Fn[len(t.Def.Fn)-1]
}
