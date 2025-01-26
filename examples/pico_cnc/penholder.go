//-----------------------------------------------------------------------------
/*

Pen Holder for Path Testing

Inspired by: https://www.thingiverse.com/thing:2625750)

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------
// pen holder

func penHolder() (sdf.SDF3, error) {

	const holderHeight = 20.0
	const holderWidth = 25.0
	const shaftRadius = 8.0 * 0.5
	const penRadius = 13.0 * 0.5

	// spring
	k := &obj.SpringParms{
		Width:         holderWidth,       // width of spring
		Height:        holderHeight,      // height of spring (3d only)
		WallThickness: 1,                 // thickness of wall
		Diameter:      5,                 // diameter of spring turn
		NumSections:   3,                 // number of spring sections
		Boss:          [2]float64{12, 8}, // boss sizes
	}
	spring, err := k.Spring3D()
	if err != nil {
		return nil, err
	}

	return spring, nil
}

//-----------------------------------------------------------------------------
