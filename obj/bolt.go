//-----------------------------------------------------------------------------
/*

Bolt: Simple Bolts for 3d printing.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"errors"
	"fmt"

	"github.com/deadsy/sdfx/sdf"
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
		return nil, errors.New("total length < 0")
	}
	if k.ShankLength < 0 {
		return nil, errors.New("shank length < 0")
	}
	if k.Tolerance < 0 {
		return nil, errors.New("tolerance < 0")
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
		return nil, fmt.Errorf("unknown style \"%s\"", k.Style)
	}
	if err != nil {
		return nil, err
	}

	// shank
	shankLength := k.ShankLength + hh/2
	shankOffset := shankLength / 2
	shank := sdf.Cylinder3D(shankLength, t.Radius, hh*0.08)
	shank = sdf.Transform3D(shank, sdf.Translate3d(sdf.V3{0, 0, shankOffset}))

	// external thread
	threadLength := k.TotalLength - k.ShankLength
	if threadLength < 0 {
		threadLength = 0
	}
	var thread sdf.SDF3
	if threadLength != 0 {
		r := t.Radius - k.Tolerance
		threadOffset := threadLength/2 + shankLength
		isoThread := sdf.ISOThread(r, t.Pitch, true)
		thread = sdf.Screw3D(isoThread, threadLength, t.Pitch, 1)
		// chamfer the thread
		thread = ChamferedCylinder(thread, 0, 0.5)
		thread = sdf.Transform3D(thread, sdf.Translate3d(sdf.V3{0, 0, threadOffset}))
	}

	return sdf.Union3D(head, shank, thread), nil
}

//-----------------------------------------------------------------------------
