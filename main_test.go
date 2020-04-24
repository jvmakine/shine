package main

import (
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name    string
		program string
		err     error
	}{{
		name: "compiles functions as variables program without errors",
		program: `
				operate = (x, y, f) => { f(x, y) }
				add = (x, y) => { x + y }
				sub = (x, y) => { x - y }
				pick = (b) => { if (b) sub else add }

				operate(1, 2, pick(true)) + operate(5, 2, pick(false))
			`,
		err: nil,
	}, {
		name: "compile euler2 without errors",
		program: `
			sum_even_fib = (upto) => {
				agg = (p2, p1, u, sum) => {
					cur = p1 + p2
					if (cur > u) sum else {
						nsum = if (cur % 2 == 0) sum + cur else sum
						agg(p1, cur, u, nsum)
					}
				}
				agg(1, 1, upto, 0)
			}
			
			sum_even_fib(100)
		`,
		err: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compile(tt.program)
			if err != nil {
				if reflect.DeepEqual(err, tt.err) {
					return
				}
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}
