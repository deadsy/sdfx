//-----------------------------------------------------------------------------
/*

Finite elements - tetrahedra - from triangle mesh.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
)

func main() {
	stl := "../../files/teapot.stl"

	// read the stl file.
	file, err := os.OpenFile(stl, os.O_RDONLY, 0400)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// create the SDF from the STL mesh
	teapotSdf, err := obj.ImportSTL(file, 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// Render SDF3 to finite elements.
	// Output file can be used by ABAQUS or CalculiX.
	render.ToFE(teapotSdf, "teapot.inp", render.NewMarchingTetrahedraUniform(200))
}
