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
	n             int     // polynomial order
	a, b, c, d, e float64 // polynomial coefficients
}

// Return the bezier polynomial function value.
func (p *BezierPolynomial) f0(t float64) float64 {
	switch p.n {
	case 1:
		// linear
		return p.a + t*p.b
	case 2:
		// quadratic
		return p.a + t*(p.b+t*p.c)
	case 3:
		// cubic
		return p.a + t*(p.b+t*(p.c+t*p.d))
	case 4:
		// quintic
		return p.a + t*(p.b+t*(p.c+t*(p.d+t*p.e)))
	default:
		panic(fmt.Sprintf("bad polynomial order %d", p.n))
	}
}

// Given the end/control points calculate the polynomial coefficients.
func (p *BezierPolynomial) Set(x []float64) {
	p.n = len(x) - 1
	switch p.n {
	case 1:
		// linear
		p.a = x[0]
		p.b = -x[0] + x[1]
	case 2:
		// quadratic
		p.a = x[0]
		p.b = -2*x[0] + 2*x[1]
		p.c = x[0] - 2*x[1] + x[2]
	case 3:
		// cubic
		p.a = x[0]
		p.b = -3*x[0] + 3*x[1]
		p.c = 3*x[0] - 6*x[1] + 3*x[2]
		p.d = -x[0] + 3*x[1] - 3*x[2] + x[3]
	case 4:
		// quintic
		p.a = x[0]
		p.b = -4*x[0] + 4*x[1]
		p.c = 6*x[0] - 12*x[1] + 6*x[2]
		p.d = -4*x[0] + 12*x[1] - 12*x[2] + 4*x[3]
		p.e = x[0] - 4*x[1] + 6*x[2] - 4*x[3] + x[4]
	default:
		panic(fmt.Sprintf("bad polynomial order %d", p.n))
	}
	// zero out any very small coefficients
	sum := Abs(p.a) + Abs(p.b) + Abs(p.c) + Abs(p.d) + Abs(p.e)
	p.a = ZeroSmall(p.a, sum, POLY_EPSILON)
	p.b = ZeroSmall(p.b, sum, POLY_EPSILON)
	p.c = ZeroSmall(p.c, sum, POLY_EPSILON)
	p.d = ZeroSmall(p.d, sum, POLY_EPSILON)
	p.e = ZeroSmall(p.e, sum, POLY_EPSILON)
}

//-----------------------------------------------------------------------------

type BezierSpline struct {
	px, py BezierPolynomial // x/y bezier polynomials
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
