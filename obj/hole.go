//-----------------------------------------------------------------------------
/*

Holes

*/
//-----------------------------------------------------------------------------

package obj

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// CounterBoredHole3D returns the SDF3 for a counterbored hole.
func CounterBoredHole3D(
	l float64, // total length
	r float64, // hole radius
	cbRadius float64, // counter bore radius
	cbDepth float64, // counter bore depth
) sdf.SDF3 {
	s0 := sdf.Cylinder3D(l, r, 0)
	s1 := sdf.Cylinder3D(cbDepth, cbRadius, 0)
	s1 = sdf.Transform3D(s1, sdf.Translate3d(sdf.V3{0, 0, (l - cbDepth) / 2}))
	return sdf.Union3D(s0, s1)
}

// ChamferedHole3D returns the SDF3 for a chamfered hole (45 degrees).
func ChamferedHole3D(
	l float64, // total length
	r float64, // hole radius
	chRadius float64, // chamfer radius
) sdf.SDF3 {
	s0 := sdf.Cylinder3D(l, r, 0)
	s1 := sdf.Cone3D(chRadius, r, r+chRadius, 0)
	s1 = sdf.Transform3D(s1, sdf.Translate3d(sdf.V3{0, 0, (l - chRadius) / 2}))
	return sdf.Union3D(s0, s1)
}

// CounterSunkHole3D returns the SDF3 for a countersunk hole (45 degrees).
func CounterSunkHole3D(
	l float64, // total length
	r float64, // hole radius
) sdf.SDF3 {
	return ChamferedHole3D(l, r, r)
}

//-----------------------------------------------------------------------------
