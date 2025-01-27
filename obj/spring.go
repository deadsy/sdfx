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
	Width         float64    // width of spring
	Height        float64    // height of spring (3d only)
	WallThickness float64    // thickness of wall
	Diameter      float64    // diameter of spring turn
	NumSections   int        // number of spring sections
	Boss          [2]float64 // boss sizes
}

// springLength returns the total spring length.
func (k *SpringParms) SpringLength() float64 {
	length := k.Boss[0] + k.Boss[1]
	length += k.WallThickness * (float64(k.NumSections) - 1)
	length += (k.Diameter - 2.0*k.WallThickness) * float64(k.NumSections)
	return length
}

// Spring2D returns a 2d spring.
func (k *SpringParms) Spring2D() (sdf.SDF2, error) {

	outerRadius := 0.5 * k.Diameter
	innerRadius := outerRadius - k.WallThickness
	spacing := 2.0 * innerRadius

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
	if k.Boss[0] < k.WallThickness {
		k.Boss[0] = k.WallThickness
	}
	if k.Boss[1] < k.WallThickness {
		k.Boss[1] = k.WallThickness
	}

	// wall thickness
	var wt []float64
	wt = append(wt, k.Boss[0])
	for i := 0; i < k.NumSections-1; i++ {
		wt = append(wt, k.WallThickness)
	}
	wt = append(wt, k.Boss[1])

	// left/right spring loops
	loop, err := Washer2D(&WasherParms{InnerRadius: innerRadius, OuterRadius: outerRadius})
	if err != nil {
		return nil, err
	}
	rLoop := sdf.Cut2D(loop, v2.Vec{}, v2.Vec{-1, 0})
	lLoop := sdf.Cut2D(loop, v2.Vec{}, v2.Vec{1, 0})

	// build the spring
	var parts []sdf.SDF2
	posn := v2.Vec{-0.5 * k.SpringLength(), 0}
	for i, t := range wt {
		wall := sdf.Box2D(v2.Vec{t, k.Width + k.WallThickness}, 0.5*k.WallThickness)
		posn.X += 0.5 * t
		parts = append(parts, sdf.Transform2D(wall, sdf.Translate2d(posn)))
		if i != len(wt)-1 {
			xOfs := 0.5 * (t + spacing)
			yOfs := 0.5 * k.Width
			if i&1 == 0 {
				parts = append(parts, sdf.Transform2D(lLoop, sdf.Translate2d(posn.Add(v2.Vec{xOfs, -yOfs}))))
			} else {
				parts = append(parts, sdf.Transform2D(rLoop, sdf.Translate2d(posn.Add(v2.Vec{xOfs, yOfs}))))
			}
		}
		posn.X += 0.5*t + spacing
	}

	return sdf.Union2D(parts...), nil
}

// Spring3D returns a 3d spring.
func (k *SpringParms) Spring3D() (sdf.SDF3, error) {
	if k.Height <= 0 {
		return nil, sdf.ErrMsg("Height <= 0")
	}
	s, err := k.Spring2D()
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, k.Height), nil
}

//-----------------------------------------------------------------------------
