//-----------------------------------------------------------------------------
/*

Simple Washer.

The washer can be partial and used to create circular wall segments.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// WasherParms defines the parameters for a washer.
type WasherParms struct {
	Thickness   float64 // thickness (3d only)
	InnerRadius float64 // inner radius
	OuterRadius float64 // outer radius
	Remove      float64 // fraction of complete washer removed
}

//-----------------------------------------------------------------------------

// Washer2D returns a 2d washer.
func Washer2D(k *WasherParms) (sdf.SDF2, error) {
	if k.InnerRadius >= k.OuterRadius {
		return nil, sdf.ErrMsg("InnerRadius >= OuterRadius")
	}
	if k.Remove != 0 {
		return nil, sdf.ErrMsg("TODO support Remove != 0")
	}
	outer, err := sdf.Circle2D(k.OuterRadius)
	if err != nil {
		return nil, err
	}
	inner, err := sdf.Circle2D(k.InnerRadius)
	if err != nil {
		return nil, err
	}
	return sdf.Difference2D(outer, inner), nil
}

//-----------------------------------------------------------------------------

// Washer3D returns a 3d washer.
// This can also be used to create circular walls.
func Washer3D(k *WasherParms) (sdf.SDF3, error) {
	if k.Thickness <= 0 {
		return nil, sdf.ErrMsg("Thickness <= 0")
	}
	if k.InnerRadius >= k.OuterRadius {
		return nil, sdf.ErrMsg("InnerRadius >= OuterRadius")
	}
	if k.Remove < 0 || k.Remove >= 1.0 {
		return nil, sdf.ErrMsg("Remove must be [0..1)")
	}

	if k.Remove == 0 {
		// difference of cylinders
		outer, err := sdf.Cylinder3D(k.Thickness, k.OuterRadius, 0)
		if err != nil {
			return nil, err
		}
		inner, err := sdf.Cylinder3D(k.Thickness, k.InnerRadius, 0)
		if err != nil {
			return nil, err
		}
		return sdf.Difference3D(outer, inner), nil
	}

	// build a 2d profile box
	dx := k.OuterRadius - k.InnerRadius
	dy := k.Thickness
	xofs := 0.5 * (k.InnerRadius + k.OuterRadius)
	b := sdf.Box2D(v2.Vec{dx, dy}, 0)
	b = sdf.Transform2D(b, sdf.Translate2d(v2.Vec{xofs, 0}))
	// rotate about the z-axis
	theta := sdf.Tau * (1.0 - k.Remove)
	s, err := sdf.RevolveTheta3D(b, theta)
	if err != nil {
		return nil, err
	}
	// center the removed portion on the x-axis
	dtheta := 0.5 * (sdf.Tau - theta)
	return sdf.Transform3D(s, sdf.RotateZ(dtheta)), nil
}

//-----------------------------------------------------------------------------
