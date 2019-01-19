//-----------------------------------------------------------------------------
/*

Quadratic Solver

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

type qSoln int

const (
	zeroSoln qSoln = iota
	oneSoln
	twoSoln
	infSoln
)

// Return the real solutions of ax^2 + bx + c = 0
func quadratic(a, b, c float64) ([]float64, qSoln) {
	// TODO Fix all comparisons to 0
	if a == 0 {
		if b == 0 {
			if c == 0 {
				// a = 0, b = 0, c = 0
				return nil, infSoln
			}
			// a = 0, b = 0, c != 0
			return nil, zeroSoln
		}
		// a =0, b != 0, c != 0
		return []float64{-c / b}, oneSoln
	}
	det := b*b - 4*a*c
	if det < 0 {
		return nil, zeroSoln
	}
	x := -b / (2 * a)
	if det == 0 {
		return []float64{x}, oneSoln
	}
	d := math.Sqrt(det) / (2 * a)
	return []float64{x + d, x - d}, twoSoln
}

//-----------------------------------------------------------------------------
