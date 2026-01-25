//-----------------------------------------------------------------------------
/*

LCD/OLED Displays

4 mounting posts with through holes for screws.
1 hole for the display window.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// DisplayParms defines the parameters for a display mount.
type DisplayParms struct {
	Window          v2.Vec  // display window
	Rounding        float64 // window corner rounding
	Supports        v2.Vec  // support positions
	SupportHeight   float64 // support height
	SupportDiameter float64 // support diameter
	HoleDiameter    float64 // suport hole diameter
	Offset          v2.Vec  // offset between window and supports
	Thickness       float64 // panel thickness
	Countersunk     bool    // counter sink screws on panel face
}

func displayPositions(k *DisplayParms, zOfs float64) v3.VecSet {

	xOfs := 0.5 * k.Supports.X
	yOfs := 0.5 * k.Supports.Y

	return v3.VecSet{
		{xOfs + k.Offset.X, yOfs + k.Offset.Y, zOfs},
		{xOfs + k.Offset.X, -yOfs + k.Offset.Y, zOfs},
		{-xOfs + k.Offset.X, yOfs + k.Offset.Y, zOfs},
		{-xOfs + k.Offset.X, -yOfs + k.Offset.Y, zOfs},
	}
}

func Display(k *DisplayParms, negative bool) (sdf.SDF3, error) {

	if negative {
		// display window
		w0 := sdf.Box2D(v2.Vec{k.Window.X, k.Window.Y}, k.Rounding)
		window := sdf.Extrude3D(w0, k.Thickness)

		// support holes
		length := k.SupportHeight + k.Thickness
		radius := 0.5 * k.HoleDiameter

		var hole sdf.SDF3
		var err error
		if k.Countersunk {
			hole, err = ChamferedHole3D(length, radius, 1.5*radius)
			if err != nil {
				return nil, err
			}
			hole = sdf.Transform3D(hole, sdf.MirrorXY())
		} else {
			hole, err = sdf.Cylinder3D(length, radius, 0)
			if err != nil {
				return nil, err
			}
		}
		zOfs := 0.5 * (length - k.Thickness)
		positions := displayPositions(k, zOfs)
		holes := sdf.Multi3D(hole, positions)
		return sdf.Union3D(window, holes), nil
	}

	// supports
	s, err := sdf.Cylinder3D(k.SupportHeight, k.SupportDiameter*0.5, 0)
	if err != nil {
		return nil, err
	}
	zOfs := 0.5 * (k.SupportHeight + k.Thickness)
	positions := displayPositions(k, zOfs)
	return sdf.Multi3D(s, positions), nil
}

//-----------------------------------------------------------------------------
