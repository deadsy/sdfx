//-----------------------------------------------------------------------------
/*

Extrusions

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func hex() SDF2 {
	return NewPolySDF2(Nagon(6, 20))
}

var height = 100.0

func TwistExtrude(sdf SDF2, p V3) float64 {
	m := Rotate(p.Z * 0.1)
	pnew := m.MulPosition(V2{p.X, p.Y})
	return sdf.Evaluate(pnew)
}

func extrude() SDF3 {
	s := NewExtrudeSDF3(NewOffsetSDF2(hex(), 5), height)
	s.(*ExtrudeSDF3).SetExtrude(TwistExtrude)
	return s
}

//-----------------------------------------------------------------------------

func main() {
	s := extrude()
	RenderSTL(s, "extrude.stl")
}

//-----------------------------------------------------------------------------
