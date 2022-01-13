//-----------------------------------------------------------------------------
/*

Spirals

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/Yeicor/surreal/surreal2"
	"log"
	"math"
	"math/rand"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	s, err := sdf.ArcSpiral2D(1.0, 20.0, 0.25*sdf.Pi, 8*sdf.Tau, 1.0)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	render.RenderDXF(s, 400, "spiral.dxf")
	render.ToSVG(s, -1, "spiral.svg", "", &render.Surreal2{Algorithm: surreal2.New( // Defaults with increased epsilon values (should be as low as possible for each SDF)
		math.Pi/30, 1e-3, 1e-8, sdf.V2i{1, 1}, 0.1, 1, 1e-8, 1, 100, rand.NewSource(0))})
}

//-----------------------------------------------------------------------------
