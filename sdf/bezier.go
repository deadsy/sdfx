//-----------------------------------------------------------------------------
/*

Interpolate using Bezier Curves

*/
//-----------------------------------------------------------------------------

package sdf

import "fmt"

//-----------------------------------------------------------------------------

const POLY_EPSILON = 1e-12

//-----------------------------------------------------------------------------

type BezierPolynomial struct {
	n          int     // polynomial order
	a, b, c, d float64 // polynomial coefficients
}

func (p *BezierPolynomial) f0(t float64) float64 {
	switch p.n {
	case 1:
		return p.a + t*p.b
	case 2:
		return p.a + t*(p.b+t*p.c)
	case 3:
		return p.a + t*(p.b+t*(p.c+p.d*t))
	default:
		panic(fmt.Sprintf("bad polynomial order %d", p.n))
	}
	return 0
}

func (p *BezierPolynomial) Set(x []float64) {
	p.n = len(x) - 1
	switch p.n {
	case 1:
		p.a = x[0]
		p.b = x[1] - x[0]
	case 2:
		p.a = x[0]
		p.b = 2 * (x[1] - x[0])
		p.c = x[2] - 2*x[1] + x[0]
	case 3:
		p.a = x[0]
		p.b = 3 * (x[1] - x[0])
		p.c = 3 * (x[2] - 2*x[1] + x[0])
		p.d = x[3] - 3*x[2] + 3*x[1] - x[0]
	default:
		panic(fmt.Sprintf("bad polynomial order %d", p.n))
	}
	sum := Abs(p.a) + Abs(p.b) + Abs(p.c) + Abs(p.d)
	p.a = ZeroSmall(p.a, sum, POLY_EPSILON)
	p.b = ZeroSmall(p.b, sum, POLY_EPSILON)
	p.c = ZeroSmall(p.c, sum, POLY_EPSILON)
	p.d = ZeroSmall(p.d, sum, POLY_EPSILON)
}

//-----------------------------------------------------------------------------

type BezierSpline struct {
	px, py BezierPolynomial // bezier polynomial
}

// Return the function value for a given t value.
func (s *BezierSpline) f0(t float64) V2 {
	return V2{s.px.f0(t), s.py.f0(t)}
}

func (s *BezierSpline) Set(p []V2) {
	x := make([]float64, len(p))
	y := make([]float64, len(p))
	for i, v := range p {
		x[i] = v.X
		y[i] = v.Y
	}
	s.px.Set(x)
	s.py.Set(y)
}

// Return a polygon approximating the bezier spline.
func (s *BezierSpline) Polygonize(n int) *Polygon {
	p := NewPolygon()
	dt := 1.0 / float64(n-1)
	t := 0.0
	for i := 0; i < n; i++ {
		p.AddV2(s.f0(t))
		t += dt
	}
	return p
}

//-----------------------------------------------------------------------------
