/*
 * This file is just a test of the logic used in some of the test files in fixtures/
 * and not a test of the package itself.
 */

package subleqha

import (
	"math"
	"testing"
)

func TestAND(t *testing.T) {
	for a := 0; a <= math.MaxUint8; a++ {
		for b := 0; b <= math.MaxUint8; b++ {
			want := a & b
			got := op_AND(a, b, 8)
			if want != got {
				t.Errorf("op_AND  a: %8b, b: %8b, got: %8b, want: %8b", a, b, got, want)
			}
		}
	}
}

func TestOR(t *testing.T) {
	for a := 0; a <= math.MaxUint8; a++ {
		for b := 0; b <= math.MaxUint8; b++ {
			want := a | b
			got := op_OR(a, b, 8)
			if want != got {
				t.Errorf("op_OR  a: %8b, b: %8b, got: %8b, want: %8b", a, b, got, want)
			}
		}
	}
}

// This is just here to test logic of routine suitable for running under SUBLEQ
func op_AND(a, b, n int) int {
	hbitval := int(math.Pow(2, float64(n-1)))
	res := 0
	for x := 0; x < n; x++ {
		m := 0
		res += res
		if a >= hbitval {
			m++
			a -= hbitval
		}
		if b >= hbitval {
			b -= hbitval
			if m == 1 {
				res++
			}
		}

		a += a
		b += b
	}
	return res
}

// This is just here to test logic of routine suitable for running under SUBLEQ
func op_OR(a, b, n int) int {
	hbitval := int(math.Pow(2, float64(n-1)))
	res := 0
	for x := 0; x < n; x++ {
		m := 0
		res += res
		if a >= hbitval {
			m++
			a -= hbitval
		}
		if b >= hbitval {
			m++
			b -= hbitval
		}

		if m > 0 {
			res++
		}

		a += a
		b += b
	}
	return res
}
