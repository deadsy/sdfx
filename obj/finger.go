//-----------------------------------------------------------------------------
/*

2D Finger Button

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// FingerButtonParms defines the parameters for a 2D finger button.
type FingerButtonParms struct {
	Width  float64 // finger width
	Gap    float64 // gap between finger and body
	Length float64 // length of the finger
}

// FingerButton2D returns a 2D cutout for a finger button.
func FingerButton2D(k *FingerButtonParms) (sdf.SDF2, error) {
	r0 := 0.5 * k.Width
	r1 := r0 - k.Gap
	l := 2.0 * k.Length
	s := sdf.Difference2D(sdf.Line2D(l, r0), sdf.Line2D(l, r1))
	s = sdf.Cut2D(s, v2.Vec{0, 0}, v2.Vec{0, 1})
	return sdf.Transform2D(s, sdf.Translate2d(v2.Vec{-k.Length, 0})), nil
}

//-----------------------------------------------------------------------------
