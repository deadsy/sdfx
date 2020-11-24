//-----------------------------------------------------------------------------
/*

Simple Washer.

The washer can be partial and used to create circular wall segments.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"errors"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// WasherParms defines the parameters for a washer.
type WasherParms struct {
	Thickness   float64 // thickness
	InnerRadius float64 // inner radius
	OuterRadius float64 // outer radius
	Remove      float64 // fraction of complete washer removed
}

// Washer3D returns a washer.
// This is also used to create circular walls.
func Washer3D(k *WasherParms) (sdf.SDF3, error) {
	if k.Thickness <= 0 {
		return nil, errors.New("Thickness <= 0")
	}
	if k.InnerRadius >= k.OuterRadius {
		return nil, errors.New("InnerRadius >= OuterRadius")
	}
	if k.Remove < 0 || k.Remove >= 1.0 {
		return nil, errors.New("Remove must be [0..1)")
	}

	var s sdf.SDF3
	if k.Remove == 0 {
		// difference of cylinders
		outer := sdf.Cylinder3D(k.Thickness, k.OuterRadius, 0)
		inner := sdf.Cylinder3D(k.Thickness, k.InnerRadius, 0)
		s = sdf.Difference3D(outer, inner)
	} else {
		// build a 2d profile box
		dx := k.OuterRadius - k.InnerRadius
		dy := k.Thickness
		xofs := 0.5 * (k.InnerRadius + k.OuterRadius)
		b := sdf.Box2D(sdf.V2{dx, dy}, 0)
		b = sdf.Transform2D(b, sdf.Translate2d(sdf.V2{xofs, 0}))
		// rotate about the z-axis
		theta := sdf.Tau * (1.0 - k.Remove)
		s = sdf.RevolveTheta3D(b, theta)
		// center the removed portion on the x-axis
		dtheta := 0.5 * (sdf.Tau - theta)
		s = sdf.Transform3D(s, sdf.RotateZ(dtheta))
	}
	return s, nil
}

//-----------------------------------------------------------------------------
