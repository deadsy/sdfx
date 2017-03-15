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
// Operations on individual splines

// Convert an x value to a t value.
func (s *CS) XtoT(x float64) float64 {
	return s.k * (x - s.x0)
}

// Convert a t value to an x value.
func (s *CS) TtoX(t float64) float64 {
	return s.x0 + (t / s.k)
}

// Return the function value for a given t value.
func (s *CS) Function(t float64) float64 {
	return s.a + t*(s.b+t*(s.c+s.d*t))
}

// Return the first derivative for a given t value.
func (s *CS) FirstDerivative(t float64) float64 {
	return s.b + t*(2*s.c+3*s.d*t)
}

// Return the second derivative for a given t value.
func (s *CS) SecondDerivative(t float64) float64 {
	return 2*s.c + 6*s.d*t
}

//-----------------------------------------------------------------------------

// Return the spline used for a given value of x.
func (s CubicSpline) Find(x float64) *CS {
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

// Return the function value on a set of cubic splines.
func (s CubicSpline) Function(x float64) float64 {
	cs := s.Find(x)
	return cs.Function(cs.XtoT(x))
}

//-----------------------------------------------------------------------------

const N_SAMPLES = 1000

// Return a 2D polygon approximating the cubic spline.
func (s *CubicSpline) Polygonize() SDF2 {
	p := NewPolygon()
	p.Add(s.xmin, 0)
	p.Add(s.xmax, 0)
	dx := (s.xmax - s.xmin) / float64(N_SAMPLES-1)
	x := s.xmax
	for i := 0; i < N_SAMPLES; i++ {
		p.Add(x, s.Function(x))
		x -= dx
	}
	p.Render("spline.dxf")
	return Polygon2D(p.Vertices())
}

//-----------------------------------------------------------------------------
// WIP - distance minimisation

// return distance squared between point and spline
func (ss *CubicSpline) Dist2(x float64, p V2) float64 {
	dx := x - p.X
	dy := ss.Function(x) - p.Y
	return dx*dx + dy*dy
}

// Dumb search for the minimum point/spline distance
func (s *CubicSpline) Min1(p V2) float64 {
	delta := (s.xmax - s.xmin) / float64(N_SAMPLES)
	x := s.xmin

	xmin := s.xmin

	dmin2 := s.Dist2(s.xmin, p)
	for i := 0; i < N_SAMPLES; i++ {
		d2 := s.Dist2(x, p)
		if d2 < dmin2 {
			dmin2 = d2
			xmin = x
		}
		x += delta
	}

	dmin := math.Sqrt(dmin2)
	fmt.Printf("dumb %v to %v %f\n", p, V2{xmin, s.Function(xmin)}, dmin)
	return dmin
}

//-----------------------------------------------------------------------------

func (s *CS) D0(t0, y0, t float64) float64 {
	dy := s.Function(t) - y0
	dt := t - t0
	return dt*dt + dy*dy
}

func (s *CS) D1(t0, y0, t float64) float64 {
	dy := s.Function(t) - y0
	dt := t - t0
	y1 := s.FirstDerivative(t)
	return 2 * (dt + y1*dy)
}

func (s *CS) D2(t0, y0, t float64) float64 {
	dy := s.Function(t) - y0
	y1 := s.FirstDerivative(t)
	y2 := s.SecondDerivative(t)
	return 2 * (1 + y1*y1 + y2*dy)
}

// Return a new t estimate for minimum distance using the Newton Raphson method.
func (s *CS) NR_Iterate(t0, y0, t float64) float64 {

	// We are minimising the distance squared function.
	// We are looking for the zeroes of the first derivative of this function.
	dy := s.Function(t) - y0
	dt := t - t0
	y1 := s.FirstDerivative(t)
	y2 := s.SecondDerivative(t)

	// d0 := dt * dt + dy * dy // distance2
	// d1 := 2 * (dt + y1*dy) // first derivative
	// d2 := 2 * (1 + y1*y1 + y2*dy) // second derivative
	// tnew = t - d1 / d2

	return t - (dt+y1*dy)/(1+y1*y1+y2*dy)
}

// Newton Raphson search for the minimum point/spline distance
func (ss *CubicSpline) Min2(p V2) float64 {

	s := ss.Find(p.X)
	t0 := s.XtoT(p.X)
	y0 := p.Y

	t := t0

	for i := 0; i < 10; i++ {
		fmt.Printf("t %f x %f y %f\n", t, s.TtoX(t), s.Function(t))
		t = s.NR_Iterate(t0, y0, t)
	}

	xmin := s.TtoX(t)
	return math.Sqrt(ss.Dist2(xmin, p))
}

//-----------------------------------------------------------------------------
