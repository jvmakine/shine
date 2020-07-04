package main

import (
	"errors"
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
									agg = (p2, p1, u, sum) => {
										cur = p1 + p2
										if (cur > u) sum else {
											nsum = if (cur % 2 == 0) sum + cur else sum
											agg(p1, cur, u, nsum)
										}
									}
									agg(1, 1, 100, 0)
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

									findLargestProd(10, 99, (x) => x % 2 == 0)
								`,
		err: nil,
	}, {
		name: "compile sequential functions",
		program: `
									a = (x) => (y) => (z) => x + y + z
									a(1)(2)(3)
								`,
		err: nil,
	}, {
		name: "cwork with unused values",
		program: `
									a = b
									b = 3
									3
								`,
		err: nil,
	}, {
		name: "compile values in right order",
		program: `
									one = 1
									two = 2
									a = one + two
									b = a + 1
									c = b + a
									d = c + two
									d + one
								`,
		err: nil,
	}, {
		name: "compiles closures with lambdas",
		program: `
								a = (f) => (z) => f(2, z)
								a((x, y) => {x + y})(1)
							`,
		err: nil,
	}, {
		name: "uses predefined types variable",
		program: `
								a = (x:int, y) => x + y
								a(1.0, 2.0)
							`,
		err: errors.New("can not unify int with real"),
	}, {
		name: "compile structure programs",
		program: `
					Person = (age, height, weight)
					bmi = (p) => p.weight / (p.height * p.height)
					person = Person(38, 1.73, 60.0)
					bmi(person)
				`,
		err: nil,
	}, {
		name: "compile interfaces",
		program: `
			a ~> { 
				add = (x) => a - x
				sub = (x) => a - x
			}
			
			b = {
				a ~> { add = (x) => a + x }
				f = (x:int) => x.add(7).add(3).sub(1)
				
				if (1.0.add(1.0) > 1.0) f(3) else f(0)
			}
			
			b
		`,
		err: nil,
	}, /*{
		name: "support closures in functions as values",
		program: `
			b = {
				add = (x, y) => x + y
				(z, w) => add(z, w)
			}
			b(1, 2)
		`,
		err: nil,
	}*/}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compile(tt.program)
			if err != nil {
				if reflect.DeepEqual(err, tt.err) {
					return
				}
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.err)
				return
			} else if tt.err != nil {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
