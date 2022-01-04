package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name    string
		program string
	}{{
		name: "compiles functions as variables program without errors",
		program: `
			operate = (x, y, f) => { f(x, y) }
			add = (x, y) => { x + y }
			sub = (x, y) => { x - y }
			pick = (b) => { if (b) sub else add }

			operate(3, 1, pick(true)) + operate(5, 1, pick(false)) + operate(1, 1, (x, y) => { x + y })
		`,
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
	}, {
		name: "compile sequential functions",
		program: `
			a = (x) => (y) => (z) => x + y + z
			a(1)(2)(3)
		`,
	}, {
		name: "cwork with unused values",
		program: `
			a = b
			b = 3
			3
		`,
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
	}, {
		name: "compiles closures with lambdas",
		program: `
			a = (f) => (z) => f(2, z)
			a((x, y) => {x + y})(1)
		`,
	}, {
		name: "compile structure programs",
		program: `
			Person :: (age: int, height: real, weight: real)
			bmi = (p) => p.weight / (p.height * p.height)
			person = Person(38, 1.73, 60.0)
			bmi(person)
		`,
	}, {
		name: "compile named types with arguments",
		program: `
			Func[A,B] :: (A) => B
			f: Func[int,int] = (x) => x
			f(1)
		`,
	}, {
		name: "compile type classes",
		program: `
			Monad[F] :: { fmap[A,B] :: (F[A], (A) => F[B]) => F[B] }
			S[A] :: (value: A)
			Monad[S] -> { fmap = (s, f) => f(s.value) }
			
			fmap(S(1), (x) => S(x + 5)).value
		`,
	}, {
		name: "compile abstract functions based on type classes",
		program: `
			Monad[F] :: { fmap[A,B] :: (F[A], (A) => F[B]) => F[B] }
			S[A] :: (value: A)
			Monad[S] -> { fmap = (s, f) => f(s.value) }
			afun = (a, f) => fmap(a, f)
			
			afun(S(1), (y) => S(y + 1)).value
		`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := compileModule(tt.program)
			require.NoError(t, err)
		})
	}
}
