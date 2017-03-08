//-----------------------------------------------------------------------------
/*

Splines

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
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
// Interpolate using cubic splines.
// interval: y(t) = a + bt + ct^2 + dt^3 for t in [0,1]
// 1st and 2nd derivatives are equal across intervals.
// 2nd derivatives == 0 at the endpoints (natural splines).
// See: http://mathworld.wolfram.com/CubicSpline.html

type CS struct {
	x0, k      float64
	a, b, c, d float64
}

type CubicSpline struct {
	xmin, xmax float64
	spline     []CS
}

// NewCubicSpline returns n-1 cubic splines for n x-ordered data points.
func NewCubicSpline(data []V2) CubicSpline {
	// Build and solve the tridiagonal matrix
	n := len(data)
	m := make([]V3, n)
	d := make([]float64, n)
	for i := 1; i < n-1; i++ {
		m[i] = V3{1, 4, 1}
		d[i] = 3 * (data[i+1].Y - data[i-1].Y)
	}
	// Special case the end splines.
	// Assume the 2nd derivative at the end points is 0.
	m[0] = V3{0, 2, 1}
	d[0] = 3 * (data[1].Y - data[0].Y)
	m[n-1] = V3{1, 2, 0}
	d[n-1] = 3 * (data[n-1].Y - data[n-2].Y)
	x := TriDiagonal(m, d)
	// The solution data are the first derivatives.
	// Reformat as the cubic polynomial coefficients.
	spline := make([]CS, n-1)
	for i := 0; i < n-1; i++ {
		x0 := data[i].X
		x1 := data[i+1].X
		y0 := data[i].Y
		y1 := data[i+1].Y
		D0 := x[i]
		D1 := x[i+1]
		spline[i].x0 = x0
		spline[i].k = 1.0 / (x1 - x0)
		spline[i].a = y0
		spline[i].b = D0
		spline[i].c = 3*(y1-y0) - 2*D0 - D1
		spline[i].d = 2*(y0-y1) + D0 + D1
	}
	return CubicSpline{data[0].X, data[n-1].X, spline}
}

//-----------------------------------------------------------------------------

// Interpolate a function value on a single cubic spline.
func (s *CS) Interpolate(x float64) float64 {
	t := s.k * (x - s.x0)
	return s.a + t*(s.b+t*(s.c+s.d*t))
}

// Return the first derivate for a single spline
func (s *CS) FirstDerivative(x float64) float64 {
	t := s.k * (x - s.x0)
	return s.b + t*(2*s.c+3*s.d*t)
}

// Return the spline index for a given value of x.
func (s CubicSpline) Index(x float64) *CS {
	// sanity checking
	n := len(s.spline)
	if n == 0 {
		panic("no splines")
	}
	// check x is within the range of the data points
	if x < s.xmin || x > s.xmax {
		panic("x is out of range")
	}
	// find the spline corresponding to the x value
	lo := 0
	hi := n
	for hi-lo > 1 {
		mid := (lo + hi) >> 1
		if s.spline[mid].x0 < x {
			lo = mid
		} else {
			hi = mid
		}
	}
	return &s.spline[lo]
}

// Interpolate a function value on a set of cubic splines.
func (s CubicSpline) Interpolate(x float64) float64 {
	return s.Index(x).Interpolate(x)
}

// Return the first derivate for a set of cubic splines.
func (s CubicSpline) FirstDerivative(x float64) float64 {
	return s.Index(x).FirstDerivative(x)
}

//-----------------------------------------------------------------------------
// WIP - distance minimisation

const N_SAMPLES = 10000

// return distance squared between point and spline
func (s *CubicSpline) Dist2(x float64, p V2) float64 {
	dx := x - p.X
	dy := s.Interpolate(x) - p.Y
	return dx*dx + dy*dy
}

// Dumb search for the minimum point/spline distance
func (s *CubicSpline) Min1(p V2) {
	delta := (s.xmax - s.xmin) / float64(N_SAMPLES)
	x := s.xmin

	dmin2 := s.Dist2(s.xmin, p)
	xmin := s.xmin

	for i := 0; i < N_SAMPLES; i++ {
		d2 := s.Dist2(x, p)
		if d2 < dmin2 {
			dmin2 = d2
			xmin = x
		}
		x += delta
	}

	fmt.Printf("xmin %f dmin %f\n", xmin, math.Sqrt(dmin2))
	m0 := V2{p.X - xmin, p.Y - s.Interpolate(xmin)}
	m1 := V2{1, s.FirstDerivative(xmin)}
	fmt.Printf("m0 %v m1 %v m0.m1 %f\n", m0, m1, m0.Dot(m1))

}

//-----------------------------------------------------------------------------
