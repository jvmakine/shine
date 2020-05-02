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

					operate(3, 1, pick(true)) + operate(5, 1, pick(false)) + operate(1, 1, (x, y) => { x + y })
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
	}, {
		name: "compile nested function args without errors",
		program: `
				findLargestProd = (start, stop, cond) => {
					agg = (b, e, le, ri, max, c) => {
						if (le > e) {
							agg(b, e, b, ri + 1, max, c)
						} else if (ri > e) {
							max
						} else if (le * ri > max && c(le * ri)) {
							agg(b, e, le + 1, ri, le * ri, c)
						} else {
							agg(b, e, le + 1, ri, max, c)
						}
					}
					agg(start, stop, start, start, 0, cond)
				}

				findLargestProd(10, 99, (x) => { x % 2 == 0 })
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
