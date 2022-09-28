//-----------------------------------------------------------------------------
/*

Pipe Connectors

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
)

//-----------------------------------------------------------------------------

const name = "sch40:1"
const units = "mm"
const length = 40.0

//-----------------------------------------------------------------------------

func main() {

	// 2-way
	s, err := obj.StdPipeConnector3D(name, units, length, [6]bool{false, false, false, false, true, true})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pipe_connector_2a.stl", render.NewMarchingCubesOctree(300))

	// 2-way
	s, err = obj.StdPipeConnector3D(name, units, length, [6]bool{true, false, false, false, true, false})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pipe_connector_2b.stl", render.NewMarchingCubesOctree(300))

	// 3-way
	s, err = obj.StdPipeConnector3D(name, units, length, [6]bool{true, false, false, false, true, true})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pipe_connector_3a.stl", render.NewMarchingCubesOctree(300))

	// 3-way
	s, err = obj.StdPipeConnector3D(name, units, length, [6]bool{true, false, true, false, true, false})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pipe_connector_3b.stl", render.NewMarchingCubesOctree(300))

	// 4-way
	s, err = obj.StdPipeConnector3D(name, units, length, [6]bool{true, true, true, true, false, false})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pipe_connector_4a.stl", render.NewMarchingCubesOctree(300))

	// 4-way
	s, err = obj.StdPipeConnector3D(name, units, length, [6]bool{true, false, true, true, true, false})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pipe_connector_4b.stl", render.NewMarchingCubesOctree(300))

	// 5-way
	s, err = obj.StdPipeConnector3D(name, units, length, [6]bool{true, true, true, true, true, false})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "pipe_connector_5a.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
