//-----------------------------------------------------------------------------
/*

The Simplest Manifold STL object.

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {

	side := 30.0

	a := sdf.V3{0, 0, 0}
	b := sdf.V3{side, 0, 0}
	c := sdf.V3{0, side, 0}
	d := sdf.V3{0, 0, side}

	t1 := render.NewTriangle3(a, b, d)
	t2 := render.NewTriangle3(a, c, b)
	t3 := render.NewTriangle3(a, d, c)
	t4 := render.NewTriangle3(b, c, d)

	err := render.SaveSTL("simple.stl", []*render.Triangle3{t1, t2, t3, t4})
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------
