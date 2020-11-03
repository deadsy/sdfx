//-----------------------------------------------------------------------------
/*

Midget Motor Casting Patterns
Popular Mechanics 1936

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func main() {
	scale := shrink * sdf.MillimetresPerInch
	sdf.RenderSTL(sdf.ScaleUniform3D(cylinderPattern(true, true), scale), 330, "cylinder_pattern.stl")
	sdf.RenderSTL(sdf.ScaleUniform3D(crankCaseFrontPattern(), scale), 300, "crankcase_front.stl")

	//sdf.RenderSTL(sdf.ScaleUniform3D(cylinderCoreBox(), shrink), 330, "cylinder_corebox.stl")
}

//-----------------------------------------------------------------------------
