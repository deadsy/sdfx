//-----------------------------------------------------------------------------
/*

Nuts: Simple Nut for 3d printing.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// NutParms defines the parameters for a nut.
type NutParms struct {
	Thread    string  // name of thread
	Style     string  // head style "hex" or "knurl"
	Tolerance float64 // add to internal thread radius
}

// Nut returns a simple nut suitable for 3d printing.
func Nut(k *NutParms) (sdf.SDF3, error) {
	// validate parameters
	t, err := sdf.ThreadLookup(k.Thread)
	if err != nil {
		return nil, err
	}
	if k.Tolerance < 0 {
		return nil, sdf.ErrMsg("Tolerance < 0")
	}

	// nut body
	var nut sdf.SDF3
	nr := t.HexRadius()
	nh := t.HexHeight()
	switch k.Style {
	case "hex":
		nut, err = HexHead3D(nr, nh, "tb")
	case "knurl":
		nut, err = KnurledHead3D(nr, nh, nr*0.25)
	default:
		return nil, sdf.ErrMsg(fmt.Sprintf("unknown style \"%s\"", k.Style))
	}
	if err != nil {
		return nil, err
	}

	// internal thread
	isoThread, err := sdf.ISOThread(t.Radius+k.Tolerance, t.Pitch, false)
	if err != nil {
		return nil, err
	}
	thread, err := sdf.Screw3D(isoThread, nh, t.Pitch, 1, t.Taper)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(nut, thread), nil
}

//-----------------------------------------------------------------------------
