//-----------------------------------------------------------------------------
/*

Splines

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
)

//-----------------------------------------------------------------------------

// Solve the tridiagonal matrix equation m.x = d, return x
// See: https://en.wikipedia.org/wiki/Tridiagonal_matrix_algorithm
func TriDiagonal(m []V3, d []float64) []float64 {
	// Sanity checks
	n := len(m)
	if len(d) != n {
		panic("bad sizes rows(m) != rows(d)")
	}
	if m[0].X != 0 || m[n-1].Z != 0 {
		panic("bad values for tridiagonal matrix")
	}
	if m[0].Y == 0 {
		panic("m[0].Y == 0")
	}
	cp := make([]float64, n) // c-prime
	x := make([]float64, n)  // d-prime -> x solution
	// elimination
	cp[0] = m[0].Z / m[0].Y
	x[0] = d[0] / m[0].Y
	for i := 1; i < n; i++ {
		denom := m[i].Y - m[i].X*cp[i-1]
		if denom == 0 {
			panic("denom == 0")
		}
		cp[i] = m[i].Z / denom
		x[i] = (d[i] - m[i].X*x[i-1]) / denom
	}
	// back substitution
	for i := n - 2; i >= 0; i-- {
		x[i] -= cp[i] * x[i+1]
	}
	return x
}

//-----------------------------------------------------------------------------
/*

y(t) = a + bt + ct^2 + dt^3
y'(t) = b + 2ct + 3dt^2

a = y(0)
b = y'(0)
c = 3 (y(1) - y(0)) - 2 y'(0) - y'(1)
d = 2 (y(0) - y(1)) + y'(0) + y'(1)

*/
//-----------------------------------------------------------------------------

type CubicSpline struct {
	data []V2
}

func NewCubicSpline(data []V2) *CubicSpline {
	s := CubicSpline{}
	s.data = data
	return &s
}

func (s *CubicSpline) Dump() {
	fmt.Printf("%+v\n", s.data)
}

//-----------------------------------------------------------------------------
