//-----------------------------------------------------------------------------
/*

Hex Heads for nuts and bolts.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// Hex2D returns a 2d hexagon with rounded corners.
func Hex2D(radius, round float64) (sdf.SDF2, error) {
	delta := 2 * round / math.Sqrt(3)
	hex, err := sdf.Polygon2D(sdf.Nagon(6, radius-delta))
	if err != nil {
		return nil, err
	}
	return sdf.Offset2D(hex, round), nil
}

// Hex3D returns a 3d hexagon with rounded corners.
func Hex3D(radius, height, round float64) (sdf.SDF3, error) {
	hex, err := Hex2D(radius, round)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(hex, height), nil
}

//-----------------------------------------------------------------------------

// HexHead3D returns the rounded hex head for a nut or bolt.
func HexHead3D(
	radius float64, // radius of hex head
	height float64, // height of hex head
	round string, // rounding control (t)top, (b)bottom, (tb)top/bottom
) (sdf.SDF3, error) {
	// basic hex body
	hex3d, err := Hex3D(radius, height, radius*0.08)
	if err != nil {
		return nil, err
	}
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
			hex3d = sdf.Intersect3D(hex3d, sdf.Transform3D(sphere3d, sdf.Translate3d(v3.Vec{0, 0, -zOfs})))
		}
		if round == "b" || round == "tb" {
			hex3d = sdf.Intersect3D(hex3d, sdf.Transform3D(sphere3d, sdf.Translate3d(v3.Vec{0, 0, zOfs})))
		}
	}
	return hex3d, nil
}

//-----------------------------------------------------------------------------
