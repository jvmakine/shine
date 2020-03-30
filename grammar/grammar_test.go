package grammar

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Block
	}{{
		name:  "parse a simple numeric expression",
		input: "42",
		want:  &Block{Value: &Expression{Term: &TermExpression{Left: intTerm(42)}}},
	}, {
		name:  "parse + term expression",
		input: "1 + 2",
		want:  &Block{Value: &Expression{Term: binTerm("+", intTerm(1), intTerm(2))}},
	}, {
		name:  "parse - term expression",
		input: "1 - 2",
		want:  &Block{Value: &Expression{Term: binTerm("-", intTerm(1), intTerm(2))}},
	}, {
		name:  "parse * factor expression",
		input: "2 * 3",
		want:  &Block{Value: &Expression{Term: &TermExpression{Left: binFact("*", intVal(2), intVal(3))}}},
	}, {
		name:  "parse / factor expression",
		input: "2 / 3",
		want:  &Block{Value: &Expression{Term: &TermExpression{Left: binFact("/", intVal(2), intVal(3))}}},
	}, {
		name:  "maintain right precedence with + and *",
		input: "2 + 3 * 4",
		want:  &Block{Value: &Expression{Term: binTerm("+", intTerm(2), binFact("*", intVal(3), intVal(4)))}},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := Parse(tt.input)
			got := prog.Body
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func binFact(op string, left *Value, right *Value) *Term {
	return &Term{Left: left, Right: []*OpFactor{&OpFactor{Operation: &op, Right: right}}}
}

func binTerm(op string, left *Term, right *Term) *TermExpression {
	return &TermExpression{Left: left, Right: []*OpTerm{&OpTerm{Operation: &op, Right: right}}}
}

func intTerm(i int) *Term {
	return &Term{Left: intVal(i)}
}

func intVal(i int) *Value {
	integer := i
	return &Value{Int: &integer}
}
