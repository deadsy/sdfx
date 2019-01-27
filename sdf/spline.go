//-----------------------------------------------------------------------------
/*

Interpolate using Cubic Splines

x(t) = a + bt + ct^2 + dt^3 for t in [0,1]
y(t) = a + bt + ct^2 + dt^3 for t in [0,1]

1st and 2nd derivatives are continuous across intervals.
2nd derivatives == 0 at the endpoints (natural splines).
See: http://mathworld.wolfram.com/CubicSpline.html

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
)

//-----------------------------------------------------------------------------

// TriDiagonal solves the tridiagonal matrix equation m.x = d, returns x.
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

// CubicPolynomial is a cubic polynomial
type CubicPolynomial struct {
	a, b, c, d float64 // polynomial coefficients
}

// Return the function value for a given t value.
func (p *CubicPolynomial) f0(t float64) float64 {
	return p.a + t*(p.b+t*(p.c+p.d*t))
}

// Return the first derivative for a given t value.
func (p *CubicPolynomial) f1(t float64) float64 {
	return p.b + t*(2*p.c+3*p.d*t)
}

// Return the second derivative for a given t value.
func (p *CubicPolynomial) f2(t float64) float64 {
	return 2*p.c + 6*p.d*t
}

// Set cubic polynomial coefficient values.
func (p *CubicPolynomial) Set(y0, y1, D0, D1 float64) {
	p.a = y0
	p.b = D0
	p.c = 3*(y1-y0) - 2*D0 - D1
	p.d = 2*(y0-y1) + D0 + D1
	// Zero out any coefficients that are small relative to the others.
	sum := Abs(p.a) + Abs(p.b) + Abs(p.c) + Abs(p.d)
	p.a = ZeroSmall(p.a, sum, epsilon)
	p.b = ZeroSmall(p.b, sum, epsilon)
	p.c = ZeroSmall(p.c, sum, epsilon)
	p.d = ZeroSmall(p.d, sum, epsilon)
}

// Return the t values for f1 == 0 (local minima/maxima)
func (p *CubicPolynomial) f1Zeroes() []float64 {
	t, _ := quadratic(3*p.d, 2*p.c, p.b)
	return t
}

//-----------------------------------------------------------------------------

// CubicSpline is a 2d cubic spline.
type CubicSpline struct {
	idx    int             // index within spline set
	p0, p1 V2              // end points of cubic spline
	px, py CubicPolynomial // cubic polynomial
}

// Return the function value for a given t value.
func (s *CubicSpline) f0(t float64) V2 {
	return V2{s.px.f0(t), s.py.f0(t)}
}

// Return the first derivative for a given t value.
func (s *CubicSpline) f1(t float64) V2 {
	return V2{s.px.f1(t), s.py.f1(t)}
}

// Return the second derivative for a given t value.
func (s *CubicSpline) f2(t float64) V2 {
	return V2{s.px.f2(t), s.py.f2(t)}
}

// BoundingBox returns the bounding box for a cubic spline.
func (s *CubicSpline) BoundingBox() Box2 {
	p := V2Set{s.p0, s.p1}
	// x minima/maxima
	for _, t := range s.px.f1Zeroes() {
		p = append(p, s.f0(Clamp(t, 0, 1)))
	}
	// y minima/maxima
	for _, t := range s.py.f1Zeroes() {
		p = append(p, s.f0(Clamp(t, 0, 1)))
	}
	return Box2{p.Min(), p.Max()}
}

const nrTolerance = 0.0001
const nrMaxIters = 10

// nrIterate is Newton-Raphson Iteration for minimum spline distance.
func (s *CubicSpline) nrIterate(t float64, p V2) float64 {
	// We are minimising the distance squared function.
	// We are looking for the zeroes of the first derivative of this function.
	// dx = x0 - p.X
	// dy = y0 - p.Y
	// d0 = dx*dx + dy*dy // distance * distance
	// d1 = 2*(dx*x1 + dy*y1)
	// d2 = 2*(dx*x2 + x1*x1 + dy*y2 + y1*y1)
	// tnew = t - d1 / d2
	f0 := s.f0(t)
	f1 := s.f1(t)
	f2 := s.f2(t)
	dx := f0.X - p.X
	dy := f0.Y - p.Y
	return t - (dx*f1.X+dy*f1.Y)/(dx*f2.X+f1.X*f1.X+dy*f2.Y+f1.Y*f1.Y)
}

//-----------------------------------------------------------------------------

// CubicSplineSDF2 is an SDF2 made from a set of cubic splines.
type CubicSplineSDF2 struct {
	spline   []CubicSpline // cubic splines
	maxiters int           // max newton-raphson iterations
	bb       Box2          // bounding box
}

// find an individual spline and t value within the set of cubic splines making up the SDF2.
func (s *CubicSplineSDF2) find(t float64) (*CubicSpline, float64) {
	n := len(s.spline)
	t = Clamp(t, 0, float64(n))
	i := int(t)
	t -= float64(i)
	// correct for the last spline
	if i == n {
		i--
		t = 1
	}
	return &s.spline[i], t
}

// f0 returns the function value for a given t value.
func (s *CubicSplineSDF2) f0(t float64) V2 {
	cs, t := s.find(t)
	return cs.f0(t)
}

// f1 returns the first derivative for a given t value.
func (s *CubicSplineSDF2) f1(t float64) V2 {
	cs, t := s.find(t)
	return cs.f1(t)
}

// f2 returns the second derivative for a given t value.
func (s *CubicSplineSDF2) f2(t float64) V2 {
	cs, t := s.find(t)
	return cs.f2(t)
}

// d0 returns the distance squared between a point and a point on the splines curve.
func (s *CubicSplineSDF2) d0(t float64, p V2) float64 {
	f0 := s.f0(t)
	dx := f0.X - p.X
	dy := f0.Y - p.Y
	return dx*dx + dy*dy
}

func (s *CubicSplineSDF2) d1(t float64, p V2) float64 {
	f0 := s.f0(t)
	f1 := s.f1(t)
	dx := f0.X - p.X
	dy := f0.Y - p.Y
	return 2 * (dx*f1.X + dy*f1.Y)
}

func (s *CubicSplineSDF2) d2(t float64, p V2) float64 {
	f0 := s.f0(t)
	f1 := s.f1(t)
	f2 := s.f2(t)
	dx := f0.X - p.X
	dy := f0.Y - p.Y
	return 2 * (dx*f2.X + f1.X*f1.X + dy*f2.Y + f1.Y*f1.Y)
}

// CubicSpline2D returns an SDF2 made from a set of cubic splines.
func CubicSpline2D(knot []V2) SDF2 {
	if len(knot) < 2 {
		panic("cubic splines need at least 2 knots")
	}
	s := CubicSplineSDF2{}
	s.maxiters = nrMaxIters

	// Build and solve the tridiagonal matrices
	n := len(knot)
	m := make([]V3, n)
	dx := make([]float64, n)
	dy := make([]float64, n)
	for i := 1; i < n-1; i++ {
		m[i] = V3{1, 4, 1}
		dx[i] = 3 * (knot[i+1].X - knot[i-1].X)
		dy[i] = 3 * (knot[i+1].Y - knot[i-1].Y)
	}
	// Special case the end splines.
	// Assume the 2nd derivative at the end points is 0.
	m[0] = V3{0, 2, 1}
	dx[0] = 3 * (knot[1].X - knot[0].X)
	dy[0] = 3 * (knot[1].Y - knot[0].Y)
	m[n-1] = V3{1, 2, 0}
	dx[n-1] = 3 * (knot[n-1].X - knot[n-2].X)
	dy[n-1] = 3 * (knot[n-1].Y - knot[n-2].Y)
	// solve to give the first derivatives at the knot points
	xx := TriDiagonal(m, dx)
	xy := TriDiagonal(m, dy)

	// The solution data are the first derivatives.
	// Reformat as the cubic polynomial coefficients.
	s.spline = make([]CubicSpline, n-1)
	for i := 0; i < n-1; i++ {
		s.spline[i].idx = i
		s.spline[i].p0 = knot[i]
		s.spline[i].p1 = knot[i+1]
		s.spline[i].px.Set(knot[i].X, knot[i+1].X, xx[i], xx[i+1])
		s.spline[i].py.Set(knot[i].Y, knot[i+1].Y, xy[i], xy[i+1])
	}

	// work out the bounding box
	s.bb = s.spline[0].BoundingBox()
	for i := 1; i < n-1; i++ {
		s.bb = s.bb.Extend(s.spline[i].BoundingBox())
	}
	return &s
}

// Evaluate returns the minimum distance from a point to the cubic spline SDF2.
// Note: This uses Newton-Raphson minimisation and is unstable in some circumstances.
// A simple (and slower) solution is to convert the cubic spline SDF2 to a polygon
// SDF2 and use that for rendering.
func (s *CubicSplineSDF2) Evaluate(p V2) float64 {

	// initial estimate
	n := 9 // len(s.spline)
	cs, t := s.find(float64(n) / 2)

	var i int
	for i = 0; i < s.maxiters; i++ {

		tOld := t
		t = cs.nrIterate(t, p)
		fmt.Printf("%d tOld %f t %f\n", cs.idx, tOld, t)

		if t < 0 {
			// previous spline
			if cs.idx == 0 {
				// no previous splines
				t = 0
				break
			}
			// find the previous spline
			cs, t = s.find(float64(cs.idx) + t)
		} else if t > 1 {
			// next spline
			if cs.idx == n-1 {
				// on the last spline
				t = 1
				break
			}
			// find the next spline
			cs, t = s.find(float64(cs.idx) + t)
		} else {
			// on the same spline
			if Abs(t-tOld) < nrTolerance*Abs(t) {
				// The t estimate is within tolerance
				break
			}
		}
	}
	t += float64(cs.idx)
	dmin := math.Sqrt(s.d0(t, p))

	//if i == s.maxiters {
	//	// deliberately cause rendering problems
	//	dmin = 0
	//}

	fmt.Printf("p %v f0 %v t %f\n", p, s.f0(t), t)

	return dmin
}

// BoundingBox returns the 2d bounding box of a cubic spline.
func (s *CubicSplineSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// Polygonize returns a polygon approximating the cubic spline SDF2.
func (s *CubicSplineSDF2) Polygonize(n int) *Polygon {
	p := NewPolygon()
	dt := float64(len(s.spline)) / float64(n-1)
	t := 0.0
	for i := 0; i < n; i++ {
		p.AddV2(s.f0(t))
		t += dt
	}
	return p
}

// PolySpline2D returns a polygon SDF2 approximating a cubic spline SDF2.
func (s *CubicSplineSDF2) PolySpline2D(n int) SDF2 {
	p := s.Polygonize(n)
	return Polygon2D(p.Vertices())
}

//-----------------------------------------------------------------------------
