package ast

import (
	"strconv"
	"strings"
)

type Options struct {
	Types bool
}

func (e *Exp) String() string {
	options := &Options{Types: true}
	return e.Format(options)
}

func (e *Exp) Format(options *Options) string {
	b := strings.Builder{}
	e.stringer(&b, 0, options)
	return b.String()
}

func (e *Exp) stringer(b *strings.Builder, indent int, options *Options) {
	if e.Id != nil {
		b.WriteString("\"")
		b.WriteString(e.Id.Name)
		b.WriteString("\"")
	} else if e.Op != nil {
		b.WriteString(e.Op.Name)
	} else if e.Const != nil {
		if e.Const.Int != nil {
			b.WriteString(strconv.FormatInt(*e.Const.Int, 10))
		} else if e.Const.Real != nil {
			b.WriteString(strconv.FormatFloat(*e.Const.Real, 'f', -1, 64))
		} else if e.Const.Bool != nil {
			if *e.Const.Bool {
				b.WriteString("true")
			} else {
				b.WriteString("false")
			}
		} else {
			panic("invalid const")
		}
	} else if e.Call != nil {
		e.Call.Function.stringer(b, indent, options)
		b.WriteString("(")
		for i, p := range e.Call.Params {
			p.stringer(b, indent, options)
			if i < len(e.Call.Params)-1 {
				b.WriteString(",")
			}
		}
		b.WriteString(")")
	} else if e.Def != nil {
		b.WriteString("(")
		for i, p := range e.Def.Params {
			b.WriteString(p.Name)
			if options.Types {
				b.WriteString(":")
				b.WriteString(p.Type.Signature())
			}
			if i < len(e.Def.Params)-1 || e.Def.HasClosure() {
				b.WriteString(",")
			}
		}
		if e.Def.HasClosure() {
			b.WriteString("[")
			for i, p := range e.Def.Closure.Fields {
				b.WriteString(p.Name)
				if options.Types {
					b.WriteString(":")
					b.WriteString(p.Type.Signature())
				}
				if i < len(e.Def.Closure.Fields)-1 {
					b.WriteString(",")
				}
			}
			b.WriteString("]")
		}
		b.WriteString(") => ")
		e.Def.Body.stringer(b, indent, options)
	} else if e.Block != nil {
		b.WriteString("{")
		newline(b, indent+2)
		for k, p := range e.Block.Assignments {
			b.WriteString(k)
			b.WriteString(" = ")
			p.stringer(b, indent+2, options)
			newline(b, indent+2)
		}
		e.Block.Value.stringer(b, indent+2, options)
		newline(b, indent)
		b.WriteString("}")
	} else if e.FAccess != nil {
		e.FAccess.Exp.stringer(b, indent, options)
		b.WriteString(".")
		b.WriteString(e.FAccess.Field)
	}
	if e.Op == nil && options.Types {
		b.WriteString(":")
		b.WriteString(e.Type().Signature())
	}
}

func newline(b *strings.Builder, indent int) {
	b.WriteString("\n")
	i := 0
	for i < indent {
		b.WriteString(" ")
		i++
	}
}
