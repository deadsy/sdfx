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
	closed bool // is the polygon closed or open?
	vlist  []PV // list of polygon vertices
}

// polygon vertex
type PV struct {
	prev   *PV     // previous vertex
	empty  bool    // is there a line segment for this vertex?
	vertex V2      // vertex coordinates
	facets int     // number of polygon facets to create when smoothing
	radius float64 // radius of smoothing (0 == none)
}

//-----------------------------------------------------------------------------
// Operations on Polygon Vertices

// Rel positions the polygon vertex relative to the prior vertex.
func (v *PV) Rel() *PV {
	if v.prev != nil {
		v.vertex = v.vertex.Add(v.prev.vertex)
	}
	return v
}

// Polar treats the polygon vertex values as polar coordinates (r, theta).
func (v *PV) Polar() *PV {
	v.vertex = PolarToXY(v.vertex.X, v.vertex.Y)
	return v
}

// Break treats this vertex as a break in the polygon.
func (v *PV) Break() *PV {
	v.empty = true
	return v
}

// Smooth marks the polygon vertex for smoothing.
func (v *PV) Smooth(radius float64, facets int) *PV {
	v.radius = radius
	v.facets = facets
	return v
}

//-----------------------------------------------------------------------------
// Operations on Polygons

// next_vertex return the next vertex in the polygon
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

// prev_vertex returns the previous vertex in the polygon
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

// smooth_vertex smoothes the i-th vertex, return true if we smoothed it
func (p *Polygon) smooth_vertex(i int) bool {

	v := p.vlist[i]
	if v.radius == 0 {
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

// Smooth does vertex smoothing on a polygon.
func (p *Polygon) Smooth() {
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

// Close closes the polygon.
func (p *Polygon) Close() {
	if p.vlist != nil {
		v_first := p.vlist[0]
		v_last := p.vlist[len(p.vlist)-1]
		if !v_first.vertex.Equals(v_last.vertex, 0) {
			p.vlist = append(p.vlist, v_first)
		}
	}
	p.closed = true
}

// NewPolygon returns an empty polygon.
func NewPolygon() *Polygon {
	return &Polygon{}
}

// Add adds a polygon vertex to a polygon.
func (p *Polygon) Add(x, y float64) *PV {
	v := PV{}
	v.vertex.X = x
	v.vertex.Y = y
	if p.vlist != nil {
		v.prev = &p.vlist[len(p.vlist)-1]
	}
	p.vlist = append(p.vlist, v)
	return &p.vlist[len(p.vlist)-1]
}

// Vertices returns the vertices of the polygon.
func (p *Polygon) Vertices() []V2 {
	v := make([]V2, len(p.vlist))
	for i, pv := range p.vlist {
		v[i] = pv.vertex
	}
	return v
}

// Render outputs a polygon as a 2D DXF file.
func (p *Polygon) Render(path string) {
	d := dxf.NewDrawing()
	for i := 0; i < len(p.vlist)-1; i++ {
		if !p.vlist[i+1].empty {
			p0 := p.vlist[i].vertex
			p1 := p.vlist[i+1].vertex
			d.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
		}
	}
	err := d.SaveAs(path)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
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
