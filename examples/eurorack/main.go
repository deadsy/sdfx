//-----------------------------------------------------------------------------
/*

Create Eurorack Panels

http://www.doepfer.de/a100_man/a100m_e.htm

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
)

//-----------------------------------------------------------------------------

func main() {

	s0, err := obj.EuroRackPanel(3, 7)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	s1, err := obj.EuroRackPanel(3, 12)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	render.RenderDXF(s0, 300, "er_7hp.dxf")
	render.RenderDXF(s1, 300, "er_12hp.dxf")
}

//-----------------------------------------------------------------------------
