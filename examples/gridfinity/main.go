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
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

func main() {

	kBase := obj.GfBaseParms{
		X: 3,
		Y: 4,
	}
	base := obj.GfBase(&kBase)
	render.ToSTL(base, "base.stl", render.NewMarchingCubesOctree(300))

	kBody := obj.GfBodyParms{
		Size:  v3i.Vec{1, 1, 3},
		Hole:  true,
		Empty: true,
	}
	body := obj.GfBody(&kBody)
	render.ToSTL(body, "body_1x1x3.stl", render.NewMarchingCubesOctree(300))

	kBody = obj.GfBodyParms{
		Size: v3i.Vec{1, 2, 1},
	}
	body = obj.GfBody(&kBody)
	render.ToSTL(body, "body_1x2x1.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
