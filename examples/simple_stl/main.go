//-----------------------------------------------------------------------------
/*

The Simplest Manifold STL object.

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {

	side := 30.0

	a := V3{0, 0, 0}
	b := V3{side, 0, 0}
	c := V3{0, side, 0}
	d := V3{0, 0, side}

	t1 := NewTriangle3(a, b, d)
	t2 := NewTriangle3(a, c, b)
	t3 := NewTriangle3(a, d, c)
	t4 := NewTriangle3(b, c, d)

	m := NewMesh([]*Triangle3{t1, t2, t3, t4})
	err := SaveSTL("simple.stl", m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------
