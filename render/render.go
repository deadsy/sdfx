//-----------------------------------------------------------------------------
/*

Top-Level Rendering Routines

*/
//-----------------------------------------------------------------------------

package render

import (
	"fmt"
	"sync"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// Render3 implementations produce a 3d triangle mesh over the bounding volume of an sdf3.
type Render3 interface {
	Render(sdf3 sdf.SDF3, meshCells int, output chan<- *Triangle3)
	Cells(s sdf.SDF3, meshCells int) (float64, sdf.V3i)
}

// DefaultRender3Cells is an internal function that avoids duplicate code. Used for Render3 implementations
func DefaultRender3Cells(s sdf.SDF3, meshCells int) (float64, sdf.V3i) {
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	return resolution, bbSize.DivScalar(resolution).ToV3i()
}

// ToSTL renders an SDF3 to an STL file.
func ToSTL(
	s sdf.SDF3, // sdf3 to render
	meshCells int, // number of cells on the longest axis of bounding box. e.g 200
	path string, // path to filename
	r Render3, // rendering method
) {
	resolution, cells := r.Cells(s, meshCells)
	fmt.Printf("rendering %s (%dx%dx%d, resolution %.2f)\n", path, cells[0], cells[1], cells[2], resolution)
	// write the triangles to an STL file
	var wg sync.WaitGroup
	output, err := WriteSTL(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	// run the renderer
	r.Render(s, meshCells, output)
	// stop the STL writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

//-----------------------------------------------------------------------------
// Legacy API (Use ToSTL for new designs) ...

// RenderSTL renders an SDF3 as an STL file (uses octree sampling).
func RenderSTL(
	s sdf.SDF3, //sdf3 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	ToSTL(s, meshCells, path, &MarchingCubesOctree{})
}

// RenderSTLSlow renders an SDF3 as an STL file (uses uniform grid sampling).
func RenderSTLSlow(
	s sdf.SDF3, //sdf3 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	ToSTL(s, meshCells, path, &MarchingCubesUniform{})
}

//-----------------------------------------------------------------------------
