//-----------------------------------------------------------------------------
/*

Splines

*/
//-----------------------------------------------------------------------------

package sdf

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
// 1st and 2nd derivatives are equal across the intervals.
// 2nd derivatives == 0 at the endpoints.
// See: http://mathworld.wolfram.com/CubicSpline.html

type CubicSpline struct {
	x0, k      float64
	a, b, c, d float64
}

// NewCubicSpline returns n-1 cubic splines for n data points.
func NewCubicSpline(data []V2) []CubicSpline {
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
	// Reformat as the cubic polynomial constants.
	s := make([]CubicSpline, n-1)
	for i := 0; i < n-1; i++ {
		x0 := data[i].X
		x1 := data[i+1].X
		y0 := data[i].Y
		y1 := data[i+1].Y
		D0 := x[i]
		D1 := x[i+1]
		s[i].x0 = x0
		s[i].k = 1.0 / (x1 - x0)
		s[i].a = y0
		s[i].b = D0
		s[i].c = 3*(y1-y0) - 2*D0 - D1
		s[i].d = 2*(y0-y1) + D0 + D1
	}
	return s
}

//-----------------------------------------------------------------------------

func (s *CubicSpline) Interpolate(x float64) float64 {
	t := s.k * (x - s.x0)
	return s.a + t*(s.b+t*(s.c+s.d*t))
}

//-----------------------------------------------------------------------------
