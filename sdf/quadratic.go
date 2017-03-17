//-----------------------------------------------------------------------------
/*

Quadratic Solver

*/
//-----------------------------------------------------------------------------

package sdf

import "math"

//-----------------------------------------------------------------------------

type QSoln int

const (
	ZERO_SOLN QSoln = iota
	ONE_SOLN
	TWO_SOLN
	INF_SOLN
)

func quadratic(a, b, c float64) ([]float64, QSoln) {
	if a == 0 {
		if b == 0 {
			if c == 0 {
				// a = 0, b = 0, c = 0
				return nil, INF_SOLN
			} else {
				// a = 0, b = 0, c != 0
				return nil, ZERO_SOLN
			}
		} else {
			return []float64{-c / b}, ONE_SOLN
		}
	}
	det := b*b - 4*a*c
	if det == 0 {
		return []float64{-b / (2 * a)}, ONE_SOLN
	}
	d := math.Sqrt(det) / (2 * a)
	x := -b / (2 * a)
	return []float64{x + d, x - d}, TWO_SOLN
}

//-----------------------------------------------------------------------------
