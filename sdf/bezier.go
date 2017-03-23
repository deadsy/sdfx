//-----------------------------------------------------------------------------
/*

Interpolate using Bezier Curves

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

const POLY_EPSILON = 1e-12

//-----------------------------------------------------------------------------

// Cubic polynomial
type Poly3 struct {
	a, b, c, d float64
}

// Return the polynomial function value for a given t value.
func (p *Poly3) f0(t float64) float64 {
	return p.a + t*(p.b+t*(p.c+p.d*t))
}

// Set polynomial coefficent values.
func (p *Poly3) Set(x0, x1, x2, x3 float64) {
	p.a = x0
	// TODO
	// Zero out any coefficients that are small relative to the others.
	sum := Abs(p.a) + Abs(p.b) + Abs(p.c) + Abs(p.d)
	p.a = ZeroSmall(p.a, sum, POLY_EPSILON)
	p.b = ZeroSmall(p.b, sum, POLY_EPSILON)
	p.c = ZeroSmall(p.c, sum, POLY_EPSILON)
	p.d = ZeroSmall(p.d, sum, POLY_EPSILON)
}

//-----------------------------------------------------------------------------

// Quadratic polynomial
type Poly2 struct {
	a, b, c float64
}

// Return the polynomial function value for a given t value.
func (p *Poly2) f0(t float64) float64 {
	return p.a + t*(p.b+t*p.c)
}

// Set polynomial coefficent values.
func (p *Poly2) Set(x0, x1, x2 float64) {
	p.a = x0
	p.b = 2 * (x1 - x0)
	p.c = x0 - 2*x1 + x2
	// Zero out any coefficients that are small relative to the others.
	sum := Abs(p.a) + Abs(p.b) + Abs(p.c)
	p.a = ZeroSmall(p.a, sum, POLY_EPSILON)
	p.b = ZeroSmall(p.b, sum, POLY_EPSILON)
	p.c = ZeroSmall(p.c, sum, POLY_EPSILON)
}

//-----------------------------------------------------------------------------

// Linear polynomial
type Poly1 struct {
	a, b float64
}

// Return the polynomial function value for a given t value.
func (p *Poly1) f0(t float64) float64 {
	return p.a + t*p.b
}

// Set polynomial coefficent values.
func (p *Poly1) Set(x0, x1 float64) {
	p.a = x0
	p.b = x1 - x0
	// Zero out any coefficients that are small relative to the others.
	sum := Abs(p.a) + Abs(p.b)
	p.a = ZeroSmall(p.a, sum, POLY_EPSILON)
	p.b = ZeroSmall(p.b, sum, POLY_EPSILON)
}

//-----------------------------------------------------------------------------
