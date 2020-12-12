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
) (sdf.SDF3, error) {
	s0, err := sdf.Cylinder3D(l, r, 0)
	if err != nil {
		return nil, err
	}
	s1, err := sdf.Cylinder3D(cbDepth, cbRadius, 0)
	if err != nil {
		return nil, err
	}
	s1 = sdf.Transform3D(s1, sdf.Translate3d(sdf.V3{0, 0, (l - cbDepth) / 2}))
	return sdf.Union3D(s0, s1), nil
}

// ChamferedHole3D returns the SDF3 for a chamfered hole (45 degrees).
func ChamferedHole3D(
	l float64, // total length
	r float64, // hole radius
	chRadius float64, // chamfer radius
) (sdf.SDF3, error) {
	s0, err := sdf.Cylinder3D(l, r, 0)
	if err != nil {
		return nil, err
	}
	s1, err := sdf.Cone3D(chRadius, r, r+chRadius, 0)
	if err != nil {
		return nil, err
	}
	s1 = sdf.Transform3D(s1, sdf.Translate3d(sdf.V3{0, 0, (l - chRadius) / 2}))
	return sdf.Union3D(s0, s1), nil
}

// CounterSunkHole3D returns the SDF3 for a countersunk hole (45 degrees).
func CounterSunkHole3D(
	l float64, // total length
	r float64, // hole radius
) (sdf.SDF3, error) {
	return ChamferedHole3D(l, r, r)
}

//-----------------------------------------------------------------------------

// BoltCircle2D returns a 2D profile for a flange bolt circle.
func BoltCircle2D(
	holeRadius float64, // radius of bolt holes
	circleRadius float64, // radius of bolt circle
	numHoles int, // number of bolts
) sdf.SDF2 {
	s := sdf.Circle2D(holeRadius)
	s = sdf.Transform2D(s, sdf.Translate2d(sdf.V2{circleRadius, 0}))
	s = sdf.RotateCopy2D(s, numHoles)
	return s
}

// BoltCircle3D returns a 3D object for a flange bolt circle.
func BoltCircle3D(
	holeDepth float64, // depth of bolt holes
	holeRadius float64, // radius of bolt holes
	circleRadius float64, // radius of bolt circle
	numHoles int, // number of bolts
) sdf.SDF3 {
	s := BoltCircle2D(holeRadius, circleRadius, numHoles)
	return sdf.Extrude3D(s, holeDepth)
}

//-----------------------------------------------------------------------------
