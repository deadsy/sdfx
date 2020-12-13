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
	render.RenderSTL(sdf.ScaleUniform3D(cp, scale), 330, "cylinder_pattern.stl")

	ccfp, err := ccFrontPattern()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(ccfp, scale), 300, "crankcase_front.stl")

	//render.RenderSTL(sdf.ScaleUniform3D(cylinderCoreBox(), shrink), 330, "cylinder_corebox.stl")
}

//-----------------------------------------------------------------------------
