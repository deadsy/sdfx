//-----------------------------------------------------------------------------
/*

Truncated Rectangular Pyramid

This a rectangular base pyramid that has rounded edges and has been truncated.

It's an attractive object in its own right, but it's particularly useful for
sand-casting patterns because the slope implements a pattern draft and the
rounded edges minimise sand crumbling.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// TruncRectPyramidParms defines the parameters for a truncated rectangular pyramid.
type TruncRectPyramidParms struct {
	Size        v3.Vec  // size of truncated pyramid
	BaseAngle   float64 // base angle of pyramid (radians)
	BaseRadius  float64 // base corner radius
	RoundRadius float64 // edge rounding radius
}

// TruncRectPyramid3D returns a truncated rectangular pyramid with rounded edges.
func TruncRectPyramid3D(k *TruncRectPyramidParms) (sdf.SDF3, error) {
	if k.Size.LTZero() {
		return nil, sdf.ErrMsg("Size < 0")
	}
	if k.BaseAngle <= 0 || k.BaseAngle > sdf.DtoR(90) {
		return nil, sdf.ErrMsg("BaseAngle must be (0,90] degrees")
	}
	if k.BaseRadius < 0 {
		return nil, sdf.ErrMsg("BaseRadius < 0")
	}
	if k.RoundRadius < 0 {
		return nil, sdf.ErrMsg("RoundRadius < 0")
	}
	h := k.Size.Z
	dr := h / math.Tan(k.BaseAngle)
	rb := k.BaseRadius + dr
	rt := math.Max(k.BaseRadius-dr, 0)
	round := math.Min(0.5*rt, k.RoundRadius)
	s, err := sdf.Cone3D(2.0*h, rb, rt, round)
	if err != nil {
		return nil, err
	}
	wx := math.Max(k.Size.X-2.0*k.BaseRadius, 0)
	wy := math.Max(k.Size.Y-2.0*k.BaseRadius, 0)
	s = sdf.Elongate3D(s, v3.Vec{wx, wy, 0})
	s = sdf.Cut3D(s, v3.Vec{0, 0, 0}, v3.Vec{0, 0, 1})
	return s, nil
}

//-----------------------------------------------------------------------------
