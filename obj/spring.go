//-----------------------------------------------------------------------------
/*

Springs

3d printable plastic springs.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// SpringParms defines a 3d printable spring.
type SpringParms struct {
	Width         float64 // width of spring
	Height        float64 // height of spring (3d only)
	WallThickness float64 // thickness of wall
	Diameter      float64 // diameter of spring turn
	NumSections   int     // number of spring sections
}

//-----------------------------------------------------------------------------

// Spring2D returns a 2d spring.
func Spring2D(k *SpringParms) (sdf.SDF2, error) {
	outerRadius := 0.5 * k.Diameter
	innerRadius := outerRadius - k.WallThickness
	spacing := k.Diameter - k.WallThickness
	// check parameters
	if k.NumSections <= 0 {
		return nil, sdf.ErrMsg("NumSections <= 0")
	}
	if k.Width < 0 {
		return nil, sdf.ErrMsg("Width < 0")
	}
	if k.WallThickness < 0 {
		return nil, sdf.ErrMsg("WallThickness < 0")
	}
	if innerRadius <= 0 {
		return nil, sdf.ErrMsg("innerRadius <= 0")
	}
	// wall
	wall := sdf.Box2D(v2.Vec{k.WallThickness, k.Width}, 0)
	// left/right spring loops
	loop, err := Washer2D(&WasherParms{InnerRadius: innerRadius, OuterRadius: outerRadius})
	if err != nil {
		return nil, err
	}
	xOfs := 0.5 * spacing
	yOfs := 0.5 * k.Width
	rLoop := sdf.Cut2D(loop, v2.Vec{}, v2.Vec{-1, 0})
	rLoop = sdf.Transform2D(rLoop, sdf.Translate2d(v2.Vec{xOfs, yOfs}))
	lLoop := sdf.Cut2D(loop, v2.Vec{}, v2.Vec{1, 0})
	lLoop = sdf.Transform2D(lLoop, sdf.Translate2d(v2.Vec{xOfs, -yOfs}))
	// left/right sections
	lSection := sdf.Union2D(wall, lLoop)
	rSection := sdf.Union2D(wall, rLoop)
	// assemble the sections
	var parts []sdf.SDF2
	var posn v2.Vec
	for i := 0; i < k.NumSections; i++ {
		if i&1 == 0 {
			parts = append(parts, sdf.Transform2D(lSection, sdf.Translate2d(posn)))
		} else {
			parts = append(parts, sdf.Transform2D(rSection, sdf.Translate2d(posn)))
		}
		posn.X += spacing
	}
	// final wall
	parts = append(parts, sdf.Transform2D(wall, sdf.Translate2d(posn)))
	return sdf.Union2D(parts...), nil
}

// Spring3D returns a 3d spring.
func Spring3D(k *SpringParms) (sdf.SDF3, error) {
	if k.Height <= 0 {
		return nil, sdf.ErrMsg("Height <= 0")
	}
	s, err := Spring2D(k)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, k.Height), nil
}

//-----------------------------------------------------------------------------
