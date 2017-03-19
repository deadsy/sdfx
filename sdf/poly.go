//-----------------------------------------------------------------------------
/*

Polygon Building Code

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"math"

	"github.com/yofu/dxf"
)

//-----------------------------------------------------------------------------

type Polygon struct {
	closed  bool // is the polygon closed or open?
	reverse bool // return the vertices in reverse order
	vlist   []PV // list of polygon vertices
}

// polygon vertex
type PV struct {
	relative bool    // vertex position is relative to previous vertex
	vtype    PVType  // type of polygon vertex
	vertex   V2      // vertex coordinates
	facets   int     // number of polygon facets to create when smoothing
	radius   float64 // radius of smoothing (0 == none)
}

type PVType int

const (
	NORMAL PVType = iota // normal vertex
	HIDE                 // hide the line segment in rendering
	SMOOTH               // smooth the vertex
	ARC                  // replace the line segment with an arc
)

//-----------------------------------------------------------------------------
// Operations on Polygon Vertices

// Rel positions the polygon vertex relative to the prior vertex.
func (v *PV) Rel() *PV {
	v.relative = true
	return v
}

// Polar treats the polygon vertex values as polar coordinates (r, theta).
func (v *PV) Polar() *PV {
	v.vertex = PolarToXY(v.vertex.X, v.vertex.Y)
	return v
}

// Hide hides the line segment for this vertex in the dxf render.
func (v *PV) Hide() *PV {
	v.vtype = HIDE
	return v
}

// Smooth marks the polygon vertex for smoothing.
func (v *PV) Smooth(radius float64, facets int) *PV {
	v.radius = radius
	v.facets = facets
	v.vtype = SMOOTH
	return v
}

// Chamfer marks the polygon vertex for chamfering.
func (v *PV) Chamfer(size float64) *PV {
	// Fake it with a 1 facet smoothing.
	// The size will be inaacurate for anything other than
	// 90 degree segments, but this is easy, and I'm lazy ...
	v.radius = size / math.Sqrt(2)
	v.facets = 1
	v.vtype = SMOOTH
	return v
}

// Arc replaces a line segment with a circular arc.
func (v *PV) Arc(radius float64, facets int) *PV {
	v.radius = radius
	v.facets = facets
	v.vtype = ARC
	return v
}

//-----------------------------------------------------------------------------

// Return the next vertex in the polygon.
func (p *Polygon) next_vertex(i int) *PV {
	if i == len(p.vlist)-1 {
		if p.closed {
			return &p.vlist[0]
		} else {
			return nil
		}
	}
	return &p.vlist[i+1]
}

// Return the previous vertex in the polygon.
func (p *Polygon) prev_vertex(i int) *PV {
	if i == 0 {
		if p.closed {
			return &p.vlist[len(p.vlist)-1]
		} else {
			return nil
		}
	}
	return &p.vlist[i-1]
}

//-----------------------------------------------------------------------------
// convert line segments to arcs

// Replace a line segment with a circular arc.
func (p *Polygon) arc_vertex(i int) bool {
	// check the vertex
	v := &p.vlist[i]
	if v.vtype != ARC {
		return false
	}
	// now it's a normal vertex
	v.vtype = NORMAL
	// check for the previous vertex
	pv := p.prev_vertex(i)
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
	d_mid := mid.Sub(a).Length()
	// distance from midpoint to center of arc
	d_center := math.Sqrt((radius * radius) - (d_mid * d_mid))
	// center of arc
	c := mid.Add(n.MulScalar(d_center))
	// work out the angle
	ac := a.Sub(c).Normalize()
	bc := b.Sub(c).Normalize()
	dtheta := -side * math.Acos(ac.Dot(bc)) / float64(v.facets)
	// rotation matrix
	m := Rotate(dtheta)
	// radius vector
	rv := m.MulPosition(a.Sub(c))
	// work out the new vertices
	vlist := make([]PV, v.facets-1)
	for j, _ := range vlist {
		vlist[j] = PV{vertex: c.Add(rv)}
		rv = m.MulPosition(rv)
	}
	// insert the new vertices between the arc endpoints
	p.vlist = append(p.vlist[:i], append(vlist, p.vlist[i:]...)...)
	return true
}

// Convert polygon line segments to arcs.
func (p *Polygon) create_arcs() {
	done := false
	for done == false {
		done = true
		for i, _ := range p.vlist {
			if p.arc_vertex(i) {
				done = false
			}
		}
	}
}

//-----------------------------------------------------------------------------
// vertex smoothing

// Smooth the i-th vertex, return true if we smoothed it.
func (p *Polygon) smooth_vertex(i int) bool {
	// check the vertex
	v := p.vlist[i]
	if v.vtype != SMOOTH {
		// fixed point
		return false
	}
	// get the next and previous points
	vn := p.next_vertex(i)
	vp := p.prev_vertex(i)
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
	dtheta := Sign(v1.Cross(v0)) * (PI - theta) / float64(v.facets)
	// rotation matrix
	rm := Rotate(dtheta)
	// radius vector
	rv := p0.Sub(c)
	// work out the new points
	points := make([]PV, v.facets+1)
	for j, _ := range points {
		points[j] = PV{vertex: c.Add(rv)}
		rv = rm.MulPosition(rv)
	}
	// replace the old point with the new points
	p.vlist = append(p.vlist[:i], append(points, p.vlist[i+1:]...)...)
	return true
}

// Smooth the vertices of a polygon.
func (p *Polygon) smooth_vertices() {
	done := false
	for done == false {
		done = true
		for i, _ := range p.vlist {
			if p.smooth_vertex(i) {
				done = false
			}
		}
	}
}

//-----------------------------------------------------------------------------

// Converts relative vertices to absolute vertices.
func (p *Polygon) relative_to_absolute() {
	for i, _ := range p.vlist {
		v := &p.vlist[i]
		if v.relative {
			pv := p.prev_vertex(i)
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
	p.relative_to_absolute()
	p.create_arcs()
	p.smooth_vertices()
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

// Add a V2 vertex to a polygon.
func (p *Polygon) AddV2(x V2) *PV {
	v := PV{}
	v.vertex = x
	v.vtype = NORMAL
	p.vlist = append(p.vlist, v)
	return &p.vlist[len(p.vlist)-1]
}

// Add an x,y vertex to a polygon.
func (p *Polygon) Add(x, y float64) *PV {
	return p.AddV2(V2{x, y})
}

// Vertices returns the vertices of the polygon.
func (p *Polygon) Vertices() []V2 {
	if p.vlist == nil {
		return nil
	}
	p.fixups()
	v := make([]V2, len(p.vlist))
	for i, pv := range p.vlist {
		v[i] = pv.vertex
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
	d := dxf.NewDrawing()
	for i := 0; i < len(p.vlist)-1; i++ {
		if p.vlist[i+1].vtype != HIDE {
			p0 := p.vlist[i].vertex
			p1 := p.vlist[i+1].vertex
			d.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
		}
	}
	// close the polygon if needed
	if p.closed {
		p0 := p.vlist[len(p.vlist)-1].vertex
		p1 := p.vlist[0].vertex
		if !p0.Equals(p1, 0) {
			d.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
		}
	}
	err := d.SaveAs(path)
	if err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------

// Return the vertices of a N sided regular polygon
func Nagon(n int, radius float64) V2Set {
	if n < 3 {
		return nil
	}
	m := Rotate(TAU / float64(n))
	v := make(V2Set, n)
	p := V2{radius, 0}
	for i := 0; i < n; i++ {
		v[i] = p
		p = m.MulPosition(p)
	}
	return v
}

//-----------------------------------------------------------------------------
