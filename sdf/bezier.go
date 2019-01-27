//-----------------------------------------------------------------------------
/*

Create curves using Bezier splines.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math/rand"
)

//-----------------------------------------------------------------------------

// BezierPolynomial contains the bezier polynomial parameters.
type BezierPolynomial struct {
	n             int     // polynomial order
	a, b, c, d, e float64 // polynomial coefficients
}

// Return the bezier polynomial function value.
func (p *BezierPolynomial) f0(t float64) float64 {
	switch p.n {
	case 0:
		// point
		return p.a
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
		// quartic
		return p.a + t*(p.b+t*(p.c+t*(p.d+t*p.e)))
	default:
		panic(fmt.Sprintf("bad polynomial order %d", p.n))
	}
}

// Return the 1st derivative of the bezier polynomial.
func (p *BezierPolynomial) f1(t float64) float64 {
	switch p.n {
	case 0:
		// point
		return 0
	case 1:
		// linear
		return p.b
	case 2:
		// quadratic
		return p.b + t*2*p.c
	case 3:
		// cubic
		return p.b + t*(2*p.c+t*3*p.d)
	case 4:
		// quartic
		return p.b + t*(2*p.c+t*(3*p.d+t*4*p.e))
	default:
		panic(fmt.Sprintf("bad polynomial order %d", p.n))
	}
}

// Return the 2nd derivative of the bezier polynomial.
func (p *BezierPolynomial) f2(t float64) float64 {
	switch p.n {
	case 0:
		// point
		return 0
	case 1:
		// linear
		return 0
	case 2:
		// quadratic
		return 2 * p.c
	case 3:
		// cubic
		return 2 * (p.c + t*3*p.d)
	case 4:
		// quartic
		return 2 * (p.c + t*3*(p.d+t*2*p.e))
	default:
		panic(fmt.Sprintf("bad polynomial order %d", p.n))
	}
}

// Set calculates bezier polynomial coefficients given the end/control points.
func (p *BezierPolynomial) Set(x []float64) {
	p.n = len(x) - 1
	switch p.n {
	case 0:
		// point
		p.a = x[0]
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
		// quartic
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
	p.a = ZeroSmall(p.a, sum, epsilon)
	p.b = ZeroSmall(p.b, sum, epsilon)
	p.c = ZeroSmall(p.c, sum, epsilon)
	p.d = ZeroSmall(p.d, sum, epsilon)
	p.e = ZeroSmall(p.e, sum, epsilon)
	// reduce the polynomial to the lowest order
	if p.n == 4 && p.e == 0 {
		p.n = 3
	}
	if p.n == 3 && p.d == 0 {
		p.n = 2
	}
	if p.n == 2 && p.c == 0 {
		p.n = 1
	}
	if p.n == 1 && p.b == 0 {
		p.n = 0
	}
}

//-----------------------------------------------------------------------------

// BezierSpline contains the x/y bezier curves for a 2D spline.
type BezierSpline struct {
	tolerance float64          // tolerance for adaptive sampling
	px, py    BezierPolynomial // x/y bezier polynomials
}

// Return the function value for a given t value.
func (s *BezierSpline) f0(t float64) V2 {
	return V2{s.px.f0(t), s.py.f0(t)}
}

// Sample generates polygon samples for a bezier spline.
func (s *BezierSpline) Sample(p *Polygon, t0, t1 float64, p0, p1 V2, n int) {

	// test the midpoint
	tmid := (t0 + t1) / 2
	pmid := s.f0(tmid)
	if colinearSlow(pmid, p0, p1, s.tolerance) {
		// the curve could be periodic so perturb the midpoint
		// pick a t value in [0.45,0.55]
		k := 0.45 + 0.1*rand.Float64()
		t2 := t0 + k*(t1-t0)
		p2 := s.f0(t2)
		if colinearSlow(p2, p0, p1, s.tolerance) {
			// looks flat enough, add the line segment
			if t0 == 0 {
				// add p0 for the first point on the spline
				p.AddV2(p0)
			}
			p.AddV2(p1)
			return
		}
	}
	// have we hit the recursion limit?
	if n > 8 {
		fmt.Printf("warn: bezier spline resursion limit %v\n", s)
		if t0 == 0 {
			// add p0 for the first point on the spline
			p.AddV2(p0)
		}
		p.AddV2(p1)
		return
	}
	// not flat enough, subdivide and recurse
	s.Sample(p, t0, tmid, p0, pmid, n+1)
	s.Sample(p, tmid, t1, pmid, p1, n+1)
}

// NewBezierSpline returns a bezier spline from the provided control/end points.
func NewBezierSpline(p []V2) *BezierSpline {
	//fmt.Printf("%v\n", p)
	s := BezierSpline{}
	// closer to 0, more polygon line segments
	s.tolerance = 0.02 // sin(theta)
	// work out the polynomials
	x := make([]float64, len(p))
	y := make([]float64, len(p))
	for i, v := range p {
		x[i] = v.X
		y[i] = v.Y
	}
	s.px.Set(x)
	s.py.Set(y)
	return &s
}

//-----------------------------------------------------------------------------

// bezierVertexType specifies the type of bezier control/endpoint.
type bezierVertexType int

const (
	endpoint bezierVertexType = iota // endpoint
	midpoint                         // midpoint
)

// BezierVertex specifies the vertex for a bezier curve.
type BezierVertex struct {
	vtype     bezierVertexType // type of bezier vertex
	vertex    V2               // vertex coordinates
	handleFwd V2               // polar coordinates of forward handle
	handleRev V2               // polar coordinates of reverse handle
}

// Bezier curve specification..
type Bezier struct {
	closed bool           // is the curve closed or open?
	vlist  []BezierVertex // list of bezier vertices
}

//-----------------------------------------------------------------------------

// Convert handles to control points.
func (b *Bezier) handles() {
	// new control vertex list
	var vlist []BezierVertex
	for _, v := range b.vlist {
		fwd := v.handleFwd
		rev := v.handleRev
		v.handleFwd = V2{}
		v.handleRev = V2{}
		// add a control midpoint for the reverse handle
		if rev.X != 0 {
			cp := BezierVertex{}
			cp.vtype = midpoint
			cp.vertex = PolarToXY(rev.X, rev.Y).Add(v.vertex)
			vlist = append(vlist, cp)
		}
		// add the original curve end point.
		vlist = append(vlist, v)
		// add a control midpoint for the forward handle
		if fwd.X != 0 {
			cp := BezierVertex{}
			cp.vtype = midpoint
			cp.vertex = PolarToXY(fwd.X, fwd.Y).Add(v.vertex)
			vlist = append(vlist, cp)
		}
	}
	// find the first endpoint control vertex
	i := 0
	for i = range vlist {
		if vlist[i].vtype == endpoint {
			break
		}
	}
	// move any leading midpoints to the end of the list
	if i != 0 {
		vlist = append(vlist[i:], vlist[:i]...)
	}
	// replace the original control vertex list
	b.vlist = vlist
}

// Take care of curve closure.
func (b *Bezier) closure() {
	// do we need to close the curve?
	if !b.closed {
		return
	}
	if len(b.vlist) == 0 || len(b.vlist) == 1 {
		panic("bad number of vertices")
	}
	first := b.vlist[0]
	last := b.vlist[len(b.vlist)-1]
	if first.vtype != endpoint {
		panic("first control vertex should be an endpoint")
	}
	if last.vtype == endpoint {
		if !last.vertex.Equals(first.vertex, tolerance) {
			// the first and last vertices aren't equal.
			// add the first vertex to close the curve
			b.vlist = append(b.vlist, first)
		}
	} else if last.vtype == midpoint {
		// add the first vertex to close the curve
		b.vlist = append(b.vlist, first)
	} else {
		panic("bad vertex type")
	}
}

// Do some validation checks on the control vertices.
func (b *Bezier) validate() {
	// basic checks
	n := len(b.vlist)
	if n < 2 {
		panic("bezier curve must have at least two points")
	}
	if b.vlist[0].vtype != endpoint {
		panic("bezier curve must start with an endpoint")
	}
	if !b.closed && b.vlist[n-1].vtype != endpoint {
		panic("non-closed bezier curve must end with an endpoint")
	}
}

// Post definition control point fixups.
func (b *Bezier) fixups() {
	b.handles()
	b.closure()
	b.validate()
}

//-----------------------------------------------------------------------------
// Public API for Bezier Curves.

// NewBezier returns an empty bezier curve.
func NewBezier() *Bezier {
	return &Bezier{}
}

// Close the bezier curve.
func (b *Bezier) Close() {
	b.closed = true
}

// AddV2 adds a V2 vertex to a polygon.
func (b *Bezier) AddV2(x V2) *BezierVertex {
	v := BezierVertex{}
	v.vertex = x
	v.vtype = endpoint
	b.vlist = append(b.vlist, v)
	return &b.vlist[len(b.vlist)-1]
}

// Add an x,y vertex to a polygon.
func (b *Bezier) Add(x, y float64) *BezierVertex {
	return b.AddV2(V2{x, y})
}

// Mid marks the vertex as a mid-curve control point.
func (v *BezierVertex) Mid() *BezierVertex {
	v.vtype = midpoint
	return v
}

// HandleFwd sets the slope handle in the forward direction.
func (v *BezierVertex) HandleFwd(theta, r float64) *BezierVertex {
	if v.vtype == midpoint {
		panic("can't place a handle on a curve midpoint")
	}
	v.handleFwd = V2{Abs(r), theta}
	return v
}

// HandleRev sets the slope handle in the reverse direction.
func (v *BezierVertex) HandleRev(theta, r float64) *BezierVertex {
	if v.vtype == midpoint {
		panic("can't place a handle on a curve midpoint")
	}
	v.handleRev = V2{Abs(r), theta}
	return v
}

// Handle marks the vertex with a slope control handle.
func (v *BezierVertex) Handle(theta, fwd, rev float64) *BezierVertex {
	v.HandleFwd(theta, fwd)
	v.HandleRev(theta+Pi, rev)
	return v
}

// Polygon returns a polygon approximating the bezier curve.
func (b *Bezier) Polygon() *Polygon {
	b.fixups()

	// generate the splines from the vertices
	var splines []*BezierSpline
	var vertices []V2

	n := len(b.vlist)
	state := endpoint
	i := 0
	for i < n {
		v := b.vlist[i]
		if state == endpoint {
			if v.vtype == endpoint {
				// start of spline
				vertices = []V2{v.vertex}
				// get the midpoints
				i++
				state = midpoint
			} else {
				panic("bad vertex type")
			}
		} else if state == midpoint {
			if v.vtype == endpoint {
				// end of spline
				vertices = append(vertices, v.vertex)
				splines = append(splines, NewBezierSpline(vertices))
				// this endpoint is the start of the next spline, don't advance
				state = endpoint
				// check for the last endpoint
				if i == n-1 {
					// end of the list
					break
				}
			} else if v.vtype == midpoint {
				// add a spline midpoint
				vertices = append(vertices, v.vertex)
				i++
			} else {
				panic("bad vertex type")
			}
		} else {
			panic("bad state")
		}
	}

	// render the splines to a polygon
	p := NewPolygon()
	n = len(splines)
	for i, s := range splines {
		if s.px.n == 0 && s.py.n == 0 {
			// This is a point, not a curve. Skip it.
			continue
		}
		// Add the spline vertices
		s.Sample(p, 0, 1, s.f0(0), s.f0(1), 0)
		if i != n-1 {
			// drop the last vertex since it is the first vertex of the next spline
			p.Drop()
		}
	}
	return p
}

//-----------------------------------------------------------------------------
