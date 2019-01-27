//-----------------------------------------------------------------------------
/*

Polygon Building Code

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"
)

//-----------------------------------------------------------------------------

// Polygon stores a set of 2d polygon vertices.
type Polygon struct {
	closed  bool            // is the polygon closed or open?
	reverse bool            // return the vertices in reverse order
	vlist   []PolygonVertex // list of polygon vertices
}

// PolygonVertex is a polygon vertex.
type PolygonVertex struct {
	relative bool    // vertex position is relative to previous vertex
	vtype    pvType  // type of polygon vertex
	vertex   V2      // vertex coordinates
	facets   int     // number of polygon facets to create when smoothing
	radius   float64 // radius of smoothing (0 == none)
}

// pvType is the type of a polygon vertex.
type pvType int

const (
	pvNormal pvType = iota // normal vertex
	pvHide                 // hide the line segment in rendering
	pvSmooth               // smooth the vertex
	pvArc                  // replace the line segment with an arc
)

//-----------------------------------------------------------------------------
// Operations on Polygon Vertices

// Rel positions the polygon vertex relative to the prior vertex.
func (v *PolygonVertex) Rel() *PolygonVertex {
	v.relative = true
	return v
}

// Polar treats the polygon vertex values as polar coordinates (r, theta).
func (v *PolygonVertex) Polar() *PolygonVertex {
	v.vertex = PolarToXY(v.vertex.X, v.vertex.Y)
	return v
}

// Hide hides the line segment for this vertex in the dxf render.
func (v *PolygonVertex) Hide() *PolygonVertex {
	v.vtype = pvHide
	return v
}

// Smooth marks the polygon vertex for smoothing.
func (v *PolygonVertex) Smooth(radius float64, facets int) *PolygonVertex {
	if radius != 0 && facets != 0 {
		v.radius = radius
		v.facets = facets
		v.vtype = pvSmooth
	}
	return v
}

// Chamfer marks the polygon vertex for chamfering.
func (v *PolygonVertex) Chamfer(size float64) *PolygonVertex {
	// Fake it with a 1 facet smoothing.
	// The size will be inaccurate for anything other than
	// 90 degree segments, but this is easy, and I'm lazy ...
	if size != 0 {
		v.radius = size * sqrtHalf
		v.facets = 1
		v.vtype = pvSmooth
	}
	return v
}

// Arc replaces a line segment with a circular arc.
func (v *PolygonVertex) Arc(radius float64, facets int) *PolygonVertex {
	if radius != 0 && facets != 0 {
		v.radius = radius
		v.facets = facets
		v.vtype = pvArc
	}
	return v
}

//-----------------------------------------------------------------------------

// nextVertex returns the next vertex in the polygon.
func (p *Polygon) nextVertex(i int) *PolygonVertex {
	if i == len(p.vlist)-1 {
		if p.closed {
			return &p.vlist[0]
		}
		return nil
	}
	return &p.vlist[i+1]
}

// prevVertex returns the previous vertex in the polygon.
func (p *Polygon) prevVertex(i int) *PolygonVertex {
	if i == 0 {
		if p.closed {
			return &p.vlist[len(p.vlist)-1]
		}
		return nil
	}
	return &p.vlist[i-1]
}

//-----------------------------------------------------------------------------
// convert line segments to arcs

// arcVertex replaces a line segment with a circular arc.
func (p *Polygon) arcVertex(i int) bool {
	// check the vertex
	v := &p.vlist[i]
	if v.vtype != pvArc {
		return false
	}
	// now it's a normal vertex
	v.vtype = pvNormal
	// check for the previous vertex
	pv := p.prevVertex(i)
	if pv == nil {
		return false
	}
	// The sign of the radius indicates which side of the chord the arc is on.
	side := Sign(v.radius)
	radius := Abs(v.radius)
	// two points on the chord
	a := pv.vertex
	b := v.vertex
	// Normal to chord
	ba := b.Sub(a).Normalize()
	n := V2{ba.Y, -ba.X}.MulScalar(side)
	// midpoint
	mid := a.Add(b).MulScalar(0.5)
	// distance from a to midpoint
	dMid := mid.Sub(a).Length()
	// distance from midpoint to center of arc
	dCenter := math.Sqrt((radius * radius) - (dMid * dMid))
	// center of arc
	c := mid.Add(n.MulScalar(dCenter))
	// work out the angle
	ac := a.Sub(c).Normalize()
	bc := b.Sub(c).Normalize()
	dtheta := -side * math.Acos(ac.Dot(bc)) / float64(v.facets)
	// rotation matrix
	m := Rotate(dtheta)
	// radius vector
	rv := m.MulPosition(a.Sub(c))
	// work out the new vertices
	vlist := make([]PolygonVertex, v.facets-1)
	for j := range vlist {
		vlist[j] = PolygonVertex{vertex: c.Add(rv)}
		rv = m.MulPosition(rv)
	}
	// insert the new vertices between the arc endpoints
	p.vlist = append(p.vlist[:i], append(vlist, p.vlist[i:]...)...)
	return true
}

// createArcs converts polygon line segments to arcs.
func (p *Polygon) createArcs() {
	done := false
	for done == false {
		done = true
		for i := range p.vlist {
			if p.arcVertex(i) {
				done = false
			}
		}
	}
}

//-----------------------------------------------------------------------------
// vertex smoothing

// Smooth the i-th vertex, return true if we smoothed it.
func (p *Polygon) smoothVertex(i int) bool {
	// check the vertex
	v := p.vlist[i]
	if v.vtype != pvSmooth {
		// fixed point
		return false
	}
	// get the next and previous points
	vn := p.nextVertex(i)
	vp := p.prevVertex(i)
	if vp == nil || vn == nil {
		// can't smooth the endpoints of an open polygon
		return false
	}
	// work out the angle
	v0 := vp.vertex.Sub(v.vertex).Normalize()
	v1 := vn.vertex.Sub(v.vertex).Normalize()
	theta := math.Acos(v0.Dot(v1))
	// distance from vertex to circle tangent
	d1 := v.radius / math.Tan(theta/2.0)
	if d1 > vp.vertex.Sub(v.vertex).Length() || d1 > vn.vertex.Sub(v.vertex).Length() {
		// unable to smooth - radius is too large
		return false
	}
	// tangent points
	p0 := v.vertex.Add(v0.MulScalar(d1))
	// distance from vertex to circle center
	d2 := v.radius / math.Sin(theta/2.0)
	// center of circle
	vc := v0.Add(v1).Normalize()
	c := v.vertex.Add(vc.MulScalar(d2))
	// rotation angle
	dtheta := Sign(v1.Cross(v0)) * (Pi - theta) / float64(v.facets)
	// rotation matrix
	rm := Rotate(dtheta)
	// radius vector
	rv := p0.Sub(c)
	// work out the new points
	points := make([]PolygonVertex, v.facets+1)
	for j := range points {
		points[j] = PolygonVertex{vertex: c.Add(rv)}
		rv = rm.MulPosition(rv)
	}
	// replace the old point with the new points
	p.vlist = append(p.vlist[:i], append(points, p.vlist[i+1:]...)...)
	return true
}

// smoothVertices smoothes the vertices of a polygon.
func (p *Polygon) smoothVertices() {
	done := false
	for done == false {
		done = true
		for i := range p.vlist {
			if p.smoothVertex(i) {
				done = false
			}
		}
	}
}

//-----------------------------------------------------------------------------

// relToAbs converts relative vertices to absolute vertices.
func (p *Polygon) relToAbs() {
	for i := range p.vlist {
		v := &p.vlist[i]
		if v.relative {
			pv := p.prevVertex(i)
			if pv.relative {
				panic("relative vertex needs an absolute reference")
			}
			v.vertex = v.vertex.Add(pv.vertex)
			v.relative = false
		}
	}
}

//-----------------------------------------------------------------------------

func (p *Polygon) fixups() {
	p.relToAbs()
	p.createArcs()
	p.smoothVertices()
}

//-----------------------------------------------------------------------------
// Public API for polygons

// Close closes the polygon.
func (p *Polygon) Close() {
	p.closed = true
}

// Reverse reverses the order the vertices are returned.
func (p *Polygon) Reverse() {
	p.reverse = true
}

// NewPolygon returns an empty polygon.
func NewPolygon() *Polygon {
	return &Polygon{}
}

// AddV2 adds a V2 vertex to a polygon.
func (p *Polygon) AddV2(x V2) *PolygonVertex {
	v := PolygonVertex{}
	v.vertex = x
	v.vtype = pvNormal
	p.vlist = append(p.vlist, v)
	return &p.vlist[len(p.vlist)-1]
}

// AddV2Set adds a set of V2 vertices to a polygon.
func (p *Polygon) AddV2Set(x []V2) {
	for _, v := range x {
		p.AddV2(v)
	}
}

// Add an x,y vertex to a polygon.
func (p *Polygon) Add(x, y float64) *PolygonVertex {
	return p.AddV2(V2{x, y})
}

// Drop the last vertex from the list.
func (p *Polygon) Drop() {
	p.vlist = p.vlist[:len(p.vlist)-1]
}

// Vertices returns the vertices of the polygon.
func (p *Polygon) Vertices() []V2 {
	if p.vlist == nil {
		return nil
	}
	p.fixups()
	n := len(p.vlist)
	v := make([]V2, n)
	if p.reverse {
		for i, pv := range p.vlist {
			v[n-1-i] = pv.vertex
		}
	} else {
		for i, pv := range p.vlist {
			v[i] = pv.vertex
		}
	}
	return v
}

// Render outputs a polygon as a 2D DXF file.
func (p *Polygon) Render(path string) error {
	if p.vlist == nil {
		return fmt.Errorf("no vertices")
	}
	p.fixups()
	fmt.Printf("rendering %s\n", path)
	d := NewDXF(path)
	for i := 0; i < len(p.vlist)-1; i++ {
		if p.vlist[i+1].vtype != pvHide {
			p0 := p.vlist[i].vertex
			p1 := p.vlist[i+1].vertex
			d.Line(p0, p1)
		}
	}
	// close the polygon if needed
	if p.closed {
		p0 := p.vlist[len(p.vlist)-1].vertex
		p1 := p.vlist[0].vertex
		if !p0.Equals(p1, tolerance) {
			d.Line(p0, p1)
		}
	}
	err := d.Save()
	if err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------

// Nagon return the vertices of a N sided regular polygon.
func Nagon(n int, radius float64) V2Set {
	if n < 3 {
		return nil
	}
	m := Rotate(Tau / float64(n))
	v := make(V2Set, n)
	p := V2{radius, 0}
	for i := 0; i < n; i++ {
		v[i] = p
		p = m.MulPosition(p)
	}
	return v
}

//-----------------------------------------------------------------------------
