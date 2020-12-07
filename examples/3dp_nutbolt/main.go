//-----------------------------------------------------------------------------
/*

3D Printable Nuts and Bolts

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// Tolerance: Measured in mm. Typically 0.0 to 0.4. Larger is looser.
// Smaller is tighter. Heuristically it could be set to some fraction
// of an FDM nozzle size. It's worth experimenting to find out a good
// value for the specific application and printer.
// const mmTolerance = 0.4 // a bit loose
// const mmTolerance = 0.2 // very tight
// const mmTolerance = 0.3 // good plastic to plastic fit
const mmTolerance = 0.3
const inchTolerance = mmTolerance / sdf.MillimetresPerInch

// Quality: The long axis of the model is rendered with n cells. A larger
// value will take longer to generate, give a better resolution and a
// larger STL file size.
const quality = 200

//-----------------------------------------------------------------------------
// inch example

func inch() error {
	// bolt
	boltParms := obj.BoltParms{
		Thread:      "unc_5/8",
		Style:       "knurl",
		Tolerance:   inchTolerance,
		TotalLength: 2.0,
		ShankLength: 0.5,
	}
	bolt, err := obj.Bolt(&boltParms)
	if err != nil {
		return err
	}
	bolt = sdf.ScaleUniform3D(bolt, sdf.MillimetresPerInch)
	render.RenderSTLSlow(bolt, quality, "inch_bolt.stl")

	// nut
	nutParms := obj.NutParms{
		Thread:    "unc_5/8",
		Style:     "knurl",
		Tolerance: inchTolerance,
	}
	nut, err := obj.Nut(&nutParms)
	if err != nil {
		return err
	}
	nut = sdf.ScaleUniform3D(nut, sdf.MillimetresPerInch)
	render.RenderSTLSlow(nut, quality, "inch_nut.stl")

	return nil
}

//-----------------------------------------------------------------------------
// metric example

func metric() error {
	// bolt
	boltParms := obj.BoltParms{
		Thread:      "M16x2",
		Style:       "hex",
		Tolerance:   mmTolerance,
		TotalLength: 50.0,
		ShankLength: 10.0,
	}
	bolt, err := obj.Bolt(&boltParms)
	if err != nil {
		return err
	}
	render.RenderSTLSlow(bolt, quality, "metric_bolt.stl")

	// nut
	nutParms := obj.NutParms{
		Thread:    "M16x2",
		Style:     "hex",
		Tolerance: mmTolerance,
	}
	nut, err := obj.Nut(&nutParms)
	if err != nil {
		return err
	}
	render.RenderSTLSlow(nut, quality, "metric_nut.stl")

	return nil
}

//-----------------------------------------------------------------------------

func main() {
	err := inch()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	err = metric()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

//-----------------------------------------------------------------------------
