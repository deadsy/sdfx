//-----------------------------------------------------------------------------
/*

Midget Motor Casting Patterns
Popular Mechanics 1936

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func main() {
	const scale = shrink * sdf.MillimetresPerInch

	cp, err := cylinderPattern(true, true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(cp, scale), "cylinder_pattern.stl", render.NewMarchingCubesOctree(330))

	ccfp, err := ccFrontPattern()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(ccfp, scale), "crankcase_front.stl", render.NewMarchingCubesOctree(300))

	//render.ToSTL(sdf.ScaleUniform3D(cylinderCoreBox(), shrink), "cylinder_corebox.stl", render.NewMarchingCubesOctree(330))
}

//-----------------------------------------------------------------------------
