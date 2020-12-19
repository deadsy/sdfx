//-----------------------------------------------------------------------------
/*

Hex Heads for nuts and bolts.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// HexHead3D returns the rounded hex head for a nut or bolt.
func HexHead3D(
	radius float64, // radius of hex head
	height float64, // height of hex head
	round string, // rounding control (t)top, (b)bottom, (tb)top/bottom
) (sdf.SDF3, error) {
	// basic hex body
	cornerRound := radius * 0.08
	hex2d, err := sdf.Polygon2D(sdf.Nagon(6, radius-cornerRound))
	if err != nil {
		return nil, err
	}
	hex2d = sdf.Offset2D(hex2d, cornerRound)
	hex3d := sdf.Extrude3D(hex2d, height)
	// round out the top and/or bottom as required
	if round != "" {
		topRound := radius * 1.6
		d := radius * math.Cos(sdf.DtoR(30))
		sphere3d, err := sdf.Sphere3D(topRound)
		if err != nil {
			return nil, err
		}
		zOfs := math.Sqrt(topRound*topRound-d*d) - height/2
		if round == "t" || round == "tb" {
			hex3d = sdf.Intersect3D(hex3d, sdf.Transform3D(sphere3d, sdf.Translate3d(sdf.V3{0, 0, -zOfs})))
		}
		if round == "b" || round == "tb" {
			hex3d = sdf.Intersect3D(hex3d, sdf.Transform3D(sphere3d, sdf.Translate3d(sdf.V3{0, 0, zOfs})))
		}
	}
	return hex3d, nil
}

//-----------------------------------------------------------------------------
