//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
)

//-----------------------------------------------------------------------------

func main() {

	s, err := cc16a()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(s, 200, "cc16a.stl")

	s, err = cc16b()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(s, 200, "cc16b.stl")

	cc18a()

	s, err = cc18b()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(s, 200, "cc18b.stl")

	s, err = cc18c()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(s, 200, "cc18c.stl")
}

//-----------------------------------------------------------------------------
