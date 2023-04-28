//-----------------------------------------------------------------------------
/*

Gridfinity Storage Parts

https://gridfinity.xyz/

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------
// material shrinkage

const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func main() {

	kBase := obj.GfBaseParms{
		X: 3,
		Y: 4,
	}
	base := obj.GfBase(&kBase)
	render.ToSTL(sdf.ScaleUniform3D(base, shrink), "base.stl", render.NewMarchingCubesOctree(300))

	kBody := obj.GfBodyParms{
		X: 1,
		Y: 2,
		Z: 2,
	}
	body := obj.GfBody(&kBody)
	render.ToSTL(sdf.ScaleUniform3D(body, shrink), "body.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
