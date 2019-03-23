//-----------------------------------------------------------------------------
/*

KeyCaps for Cherry MX key switches

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

const stemX = 6.0
const stemY = 5.0

const crossDepth = 4.0
const crossWidth = 1.0
const crossX = 4.0
const stemRound = 0.05

// keyStem returns a keycap stem of a given length.
func keyStem(length float64) sdf.SDF3 {
	ofs := length - crossDepth
	s0 := sdf.Box3D(sdf.V3{crossX, crossWidth, length}, crossX*stemRound)
	s1 := sdf.Box3D(sdf.V3{crossWidth, stemY * (1.0 + 2.0*stemRound), length}, crossX*stemRound)
	cavity := sdf.Transform3D(sdf.Union3D(s0, s1), sdf.Translate3d(sdf.V3{0, 0, ofs}))
	stem := sdf.Box3D(sdf.V3{stemX, stemY, length}, stemX*stemRound)
	return sdf.Difference3D(stem, cavity)
}

//-----------------------------------------------------------------------------

// roundCap returns a round keycap.
func roundCap(diameter, height, wall float64) sdf.SDF3 {
	rOuter := 0.5 * diameter
	rInner := 0.5 * (diameter - (2.0 * wall))

	outer := sdf.Cylinder3D(height, rOuter, 0)
	inner := sdf.Cylinder3D(height, rInner, 0)
	inner = sdf.Transform3D(inner, sdf.Translate3d(sdf.V3{0, 0, wall}))

	keycap := sdf.Difference3D(outer, inner)
	stem := keyStem(height)

	return sdf.Union3D(keycap, stem)
}

//-----------------------------------------------------------------------------

func main() {
	sdf.RenderSTL(roundCap(18, 6, 1.5), 200, "round_cap.stl")
}

//-----------------------------------------------------------------------------
