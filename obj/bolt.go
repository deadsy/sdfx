//-----------------------------------------------------------------------------
/*

Bolt: Simple Bolts for 3d printing.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// BoltParms defines the parameters for a bolt.
type BoltParms struct {
	Thread      string  // name of thread
	Style       string  // head style "hex" or "knurl"
	Tolerance   float64 // subtract from external thread radius
	TotalLength float64 // threaded length + shank length
	ShankLength float64 // non threaded length
}

// Bolt returns a simple bolt suitable for 3d printing.
func Bolt(k *BoltParms) (sdf.SDF3, error) {
	// validate parameters
	t, err := sdf.ThreadLookup(k.Thread)
	if err != nil {
		return nil, err
	}
	if k.TotalLength < 0 {
		return nil, sdf.ErrMsg("TotalLength < 0")
	}
	if k.ShankLength < 0 {
		return nil, sdf.ErrMsg("ShankLength < 0")
	}
	if k.Tolerance < 0 {
		return nil, sdf.ErrMsg("Tolerance < 0")
	}

	// head
	var head sdf.SDF3
	hr := t.HexRadius()
	hh := t.HexHeight()
	switch k.Style {
	case "hex":
		head, err = HexHead3D(hr, hh, "b")
	case "knurl":
		head, err = KnurledHead3D(hr, hh, hr*0.25)
	default:
		return nil, sdf.ErrMsg(fmt.Sprintf("unknown style \"%s\"", k.Style))
	}
	if err != nil {
		return nil, err
	}

	// shank
	shankLength := k.ShankLength + hh/2
	shankOffset := shankLength / 2
	shank, err := sdf.Cylinder3D(shankLength, t.Radius, hh*0.08)
	if err != nil {
		return nil, err
	}
	shank = sdf.Transform3D(shank, sdf.Translate3d(v3.Vec{0, 0, shankOffset}))

	// external thread
	threadLength := k.TotalLength - k.ShankLength
	if threadLength < 0 {
		threadLength = 0
	}
	var thread sdf.SDF3
	if threadLength != 0 {
		r := t.Radius - k.Tolerance
		threadOffset := threadLength/2 + shankLength
		isoThread, err := sdf.ISOThread(r, t.Pitch, true)
		if err != nil {
			return nil, err
		}
		thread, err = sdf.Screw3D(isoThread, threadLength, t.Taper, t.Pitch, 1)
		if err != nil {
			return nil, err
		}
		// chamfer the thread
		thread, err = ChamferedCylinder(thread, 0, 0.5)
		if err != nil {
			return nil, err
		}
		thread = sdf.Transform3D(thread, sdf.Translate3d(v3.Vec{0, 0, threadOffset}))
	}

	return sdf.Union3D(head, shank, thread), nil
}

//-----------------------------------------------------------------------------
