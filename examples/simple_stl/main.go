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
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func main() {

	side := 30.0

	a := v3.Vec{0, 0, 0}
	b := v3.Vec{side, 0, 0}
	c := v3.Vec{0, side, 0}
	d := v3.Vec{0, 0, side}

	t1 := &sdf.Triangle3{a, b, d}
	t2 := &sdf.Triangle3{a, c, b}
	t3 := &sdf.Triangle3{a, d, c}
	t4 := &sdf.Triangle3{b, c, d}

	err := render.SaveSTL("simple.stl", []*sdf.Triangle3{t1, t2, t3, t4})
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------
