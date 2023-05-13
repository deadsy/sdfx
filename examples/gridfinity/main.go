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
	"github.com/deadsy/sdfx/vec/v2i"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

func main() {

	kBase := obj.GfBaseParms{
		Size:   v2i.Vec{4, 4},
		Magnet: true,
		Hole:   true,
	}
	base := obj.GfBase(&kBase)
	render.ToSTL(base, "base_4x4.stl", render.NewMarchingCubesOctree(300))

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
