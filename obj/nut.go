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

type ThreadedCylinderParms struct {
	Height    float64 // height of cylinder
	Diameter  float64 // diameter of cylinder
	Thread    string  // name of thread
	Tolerance float64 // add to internal thread radius
}

// Object returns a cylinder with an internal thread.
func (k *ThreadedCylinderParms) Object() (sdf.SDF3, error) {
	// validate parameters
	t, err := sdf.ThreadLookup(k.Thread)
	if err != nil {
		return nil, err
	}
	if k.Diameter < 0 {
		return nil, sdf.ErrMsg("Diameter < 0")
	}
	if k.Height < 0 {
		return nil, sdf.ErrMsg("Height < 0")
	}
	if k.Tolerance < 0 {
		return nil, sdf.ErrMsg("Tolerance < 0")
	}
	body, err := sdf.Cylinder3D(k.Height, 0.5*k.Diameter, 0)
	if err != nil {
		return nil, err
	}
	// internal thread
	t = t.ToMillimetre()
	isoThread, err := sdf.ISOThread(t.Radius+k.Tolerance, t.Pitch, false)
	if err != nil {
		return nil, err
	}
	thread, err := sdf.Screw3D(isoThread, k.Height, t.Taper, t.Pitch, 1)
	if err != nil {
		return nil, err
	}
	return sdf.Difference3D(body, thread), nil
}

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
	thread, err := sdf.Screw3D(isoThread, nh, t.Taper, t.Pitch, 1)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(nut, thread), nil
}

//-----------------------------------------------------------------------------
