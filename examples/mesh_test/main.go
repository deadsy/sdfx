//-----------------------------------------------------------------------------
/*


 */
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

func main() {

	x := -100.0
	const dx0 = 1
	const dx1 = 0.5

	lines := make([]*sdf.Line2, 200)
	for i := range lines {
		lines[i] = &sdf.Line2{v2.Vec{x, -1}, v2.Vec{x + dx0, 1}}
		x += dx1
	}

	s, err := sdf.Mesh2D(lines)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	s.Evaluate(v2.Vec{0, 0})

}

//-----------------------------------------------------------------------------
