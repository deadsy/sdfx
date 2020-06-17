//-----------------------------------------------------------------------------
/*

Nuts and Bolts

*/
//-----------------------------------------------------------------------------

package main

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

func nutAndBolt(
	name string, // name of thread
	totalLength float64, // threaded length + shank length
	shankLength float64, //  non threaded length
) sdf.SDF3 {

	// bolt
	boltParms := sdf.BoltParms{
		Thread:      name,
		Style:       "hex",
		TotalLength: totalLength,
		ShankLength: shankLength,
	}
	bolt, _ := sdf.Bolt(&boltParms)

	// nut
	nutParms := sdf.NutParms{
		Thread: name,
		Style:  "hex",
	}
	nut, _ := sdf.Nut(&nutParms)
	zOffset := totalLength * 1.5
	nut = sdf.Transform3D(nut, sdf.Translate3d(sdf.V3{0, 0, zOffset}))

	return sdf.Union3D(nut, bolt)
}

//-----------------------------------------------------------------------------

func main() {

	xOffset := 1.5

	s0 := nutAndBolt("unc_1/4", 2, 0.5)
	s0 = sdf.Transform3D(s0, sdf.Translate3d(sdf.V3{-0.6 * xOffset, 0, 0}))

	s1 := nutAndBolt("unc_1/2", 2.0, 0.5)
	s1 = sdf.Transform3D(s1, sdf.Translate3d(sdf.V3{0, 0, 0}))

	s2 := nutAndBolt("unc_1", 2.0, 0.5)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(sdf.V3{xOffset, 0, 0}))

	sdf.RenderSTLSlow(sdf.Union3D(s0, s1, s2), 400, "nutandbolt.stl")
}

//-----------------------------------------------------------------------------
