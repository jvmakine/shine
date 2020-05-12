package ast

import (
	"strconv"
	"strings"
)

func (e *Exp) String() string {
	b := strings.Builder{}
	e.stringer(&b)
	return b.String()
}

func (e *Exp) stringer(b *strings.Builder) {
	if e.Id != nil {
		b.WriteString("ID[name:")
		b.WriteString(e.Id.Name)
		b.WriteString(",type:")
		b.WriteString(e.Id.Type.Signature())
		b.WriteString("]")
	} else if e.Const != nil {
		b.WriteString("CONST[type:")
		b.WriteString(e.Const.Type.Signature())
		b.WriteString("]")
	} else if e.Call != nil {
		b.WriteString("CALL[fun:{")
		e.Call.Function.stringer(b)
		b.WriteString("},type:")
		b.WriteString(e.Call.Type.Signature())
		b.WriteString(",params:{")
		for i, p := range e.Call.Params {
			p.stringer(b)
			if i < len(e.Call.Params)-1 {
				b.WriteString(",")
			}
		}
		b.WriteString("}]")
	} else if e.Def != nil {
		b.WriteString("DEF[type:")
		b.WriteString(e.Type().Signature())
		b.WriteString(",body:{")
		e.Def.Body.stringer(b)
		b.WriteString("}]")
	} else if e.Block != nil {
		b.WriteString("BLOCK[id:")
		b.WriteString(strconv.Itoa(e.Block.ID))
		b.WriteString(",type:")
		b.WriteString(e.Type().Signature())
		b.WriteString(",value:{")
		e.Block.Value.stringer(b)
		b.WriteString("},assignments:{")
		for k, p := range e.Block.Assignments {
			b.WriteString(k)
			b.WriteString(":")
			p.stringer(b)
			b.WriteString(",")
		}
		b.WriteString("}]")
	}
}
