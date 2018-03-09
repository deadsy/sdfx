//-----------------------------------------------------------------------------
/*

Text Example

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {

	f, err := LoadFont("cmr10.ttf")
	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}

	s2d, err := TextSDF2(f, "Hello World!")
	RenderDXF(s2d, 200, "shape.dxf")

	s3d := ExtrudeRounded3D(s2d, 200, 20)
	RenderSTL(s3d, 200, "shape.stl")
}

//-----------------------------------------------------------------------------
