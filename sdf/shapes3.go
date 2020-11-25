//-----------------------------------------------------------------------------
/*

Common 3D shapes.

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

// CounterBoredHole3D returns the SDF3 for a counterbored hole.
func CounterBoredHole3D(
	l float64, // total length
	r float64, // hole radius
	cbRadius float64, // counter bore radius
	cbDepth float64, // counter bore depth
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cylinder3D(cbDepth, cbRadius, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - cbDepth) / 2}))
	return Union3D(s0, s1)
}

// ChamferedHole3D returns the SDF3 for a chamfered hole (45 degrees).
func ChamferedHole3D(
	l float64, // total length
	r float64, // hole radius
	chRadius float64, // chamfer radius
) SDF3 {
	s0 := Cylinder3D(l, r, 0)
	s1 := Cone3D(chRadius, r, r+chRadius, 0)
	s1 = Transform3D(s1, Translate3d(V3{0, 0, (l - chRadius) / 2}))
	return Union3D(s0, s1)
}

// CounterSunkHole3D returns the SDF3 for a countersunk hole (45 degrees).
func CounterSunkHole3D(
	l float64, // total length
	r float64, // hole radius
) SDF3 {
	return ChamferedHole3D(l, r, r)
}

//-----------------------------------------------------------------------------

// ChamferedCylinder intersects a chamfered cylinder with an SDF3.
func ChamferedCylinder(s SDF3, kb, kt float64) SDF3 {
	// get the length and radius from the bounding box
	l := s.BoundingBox().Max.Z
	r := s.BoundingBox().Max.X
	p := NewPolygon()
	p.Add(0, -l)
	p.Add(r, -l).Chamfer(r * kb)
	p.Add(r, l).Chamfer(r * kt)
	p.Add(0, l)
	return Intersect3D(s, Revolve3D(Polygon2D(p.Vertices())))
}

//-----------------------------------------------------------------------------
