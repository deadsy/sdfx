//-----------------------------------------------------------------------------
/*

Nuts and Bolts

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func nutAndBolt(
	name string, // name of thread
	totalLength float64, // threaded length + shank length
	shankLength float64, //  non threaded length
) (sdf.SDF3, error) {

	// bolt
	boltParms := obj.BoltParms{
		Thread:      name,
		Style:       "hex",
		TotalLength: totalLength,
		ShankLength: shankLength,
	}
	bolt, err := obj.Bolt(&boltParms)
	if err != nil {
		return nil, err
	}

	// nut
	nutParms := obj.NutParms{
		Thread: name,
		Style:  "hex",
	}
	nut, err := obj.Nut(&nutParms)
	if err != nil {
		return nil, err
	}

	zOffset := totalLength * 1.5
	nut = sdf.Transform3D(nut, sdf.Translate3d(sdf.V3{0, 0, zOffset}))

	return sdf.Union3D(nut, bolt), nil
}

//-----------------------------------------------------------------------------

func main() {

	xOffset := 1.5

	s0, err := nutAndBolt("unc_1/4", 2, 0.5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s0 = sdf.Transform3D(s0, sdf.Translate3d(sdf.V3{-0.6 * xOffset, 0, 0}))

	s1, err := nutAndBolt("unc_1/2", 2.0, 0.5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	//s1 = sdf.Transform3D(s1, sdf.Translate3d(sdf.V3{0, 0, 0}))

	s2, err := nutAndBolt("unc_1", 2.0, 0.5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s2 = sdf.Transform3D(s2, sdf.Translate3d(sdf.V3{xOffset, 0, 0}))

	render.RenderSTLSlow(sdf.Union3D(s0, s1, s2), 400, "nutandbolt.stl")
}

//-----------------------------------------------------------------------------
