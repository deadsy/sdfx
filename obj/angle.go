//-----------------------------------------------------------------------------
/*

Angle: Create profiles for steel/aluminum angle.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// AngleLeg defines the parameters for one leg of a piece of angle.
type AngleLeg struct {
	Length    float64
	Thickness float64
}

// AngleParms defines the parameters for a piece of angle.
type AngleParms struct {
	X, Y       AngleLeg // angle legs
	RootRadius float64  // radius of inside fillet
	Length     float64  // length (3d only)
}

// Angle2D returns a 2d angle profile.
func Angle2D(k *AngleParms) (sdf.SDF2, error) {
	if k.X.Length <= 0 {
		return nil, sdf.ErrMsg("k.X.Length <= 0")
	}
	if k.X.Thickness <= 0 {
		return nil, sdf.ErrMsg("k.X.Thickness <= 0")
	}
	if k.Y.Length <= 0 {
		return nil, sdf.ErrMsg("k.Y.Length <= 0")
	}
	if k.Y.Thickness <= 0 {
		return nil, sdf.ErrMsg("k.Y.Thickness <= 0")
	}
	if k.Y.Thickness >= k.X.Length {
		return nil, sdf.ErrMsg("k.Y.Thickness >= k.X.Length")
	}
	if k.X.Thickness >= k.Y.Length {
		return nil, sdf.ErrMsg("k.X.Thickness >= k.Y.Length")
	}
	if k.RootRadius < 0 {
		return nil, sdf.ErrMsg("k.RootRadius < 0")
	}
	if k.RootRadius > (k.X.Length - k.Y.Thickness) {
		return nil, sdf.ErrMsg("k.RootRadius > (k.X.LengthA - k.Y.Thickness)")
	}
	if k.RootRadius > (k.Y.Length - k.X.Thickness) {
		return nil, sdf.ErrMsg("k.RootRadius > (k.Y.Length - k.X.Thickness)")
	}

	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(k.X.Length, 0)
	p.Add(k.X.Length, k.X.Thickness)
	p.Add(k.Y.Thickness, k.X.Thickness).Smooth(k.RootRadius, 6)
	p.Add(k.Y.Thickness, k.Y.Length)
	p.Add(0, k.Y.Length)

	return sdf.Polygon2D(p.Vertices())
}

// Angle3D returns a piece of 3d angle.
func Angle3D(k *AngleParms) (sdf.SDF3, error) {
	if k.Length <= 0 {
		return nil, sdf.ErrMsg("k.Length <= 0")
	}
	s, err := Angle2D(k)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, k.Length), nil
}

//-----------------------------------------------------------------------------
