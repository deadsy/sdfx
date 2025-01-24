//-----------------------------------------------------------------------------
/*

Springs

3d printable plastic springs.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
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
	if k.NumSections <= 0 {
		return nil, sdf.ErrMsg("NumSections <= 0")
	}
	if k.Width < 0 {
		return nil, sdf.ErrMsg("Width < 0")
	}
	if k.WallThickness < 0 {
		return nil, sdf.ErrMsg("WallThickness < 0")
	}
	radius := 0.5 * k.Diameter
	if radius < k.WallThickness {
		return nil, sdf.ErrMsg("radius < k.WallThickness")
	}

	return nil, nil
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
