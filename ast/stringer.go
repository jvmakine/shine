package ast

import (
	"strings"
)

type FormatOptions struct {
	Types bool
}

func String(e Expression) string {
	options := &FormatOptions{Types: true}
	return Format(e, options)
}

func Format(e Expression, options *FormatOptions) string {
	b := strings.Builder{}
	e.Format(&b, 0, options)
	return b.String()
}

func newline(b *strings.Builder, indent int) {
	b.WriteString("\n")
	i := 0
	for i < indent {
		b.WriteString(" ")
		i++
	}
}
