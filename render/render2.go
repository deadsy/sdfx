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

// Render2 implementations produce a 2d triangle mesh over the bounding volume of a SDF2.
type Render2 interface {
	Render(sdf2 sdf.SDF2, meshCells int, output chan<- *Line)
	Info(sdf2 sdf.SDF2, meshCells int) string
}

// ToSVG renders an SDF2 to an SVG file.
func ToSVG(
	s sdf.SDF2, // sdf2 to render
	meshCells int, // number of cells on the longest axis of bounding box. e.g 200
	path, lineStyle string, // path to filename
	r Render2, // rendering method
) {
	if lineStyle == "" { // Set a default
		lineStyle = "fill:none;stroke:black;stroke-width:0.1"
	}
	fmt.Printf("rendering %s (%s)\n", path, r.Info(s, meshCells))
	// write the triangles to an SVG file
	var wg sync.WaitGroup
	output, err := WriteSVG(&wg, path, lineStyle)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	// run the renderer
	r.Render(s, meshCells, output)
	// stop the SVG writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

//-----------------------------------------------------------------------------
// Legacy API (Use ToSVG for new designs) ...

// Deprecated: RenderSVG renders an SDF2 as an SVG file (uses octree sampling).
func RenderSVG(
	s sdf.SDF2, //sdf2 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path, lineStyle string, //path to filename
) {
	ToSVG(s, meshCells, path, lineStyle, &MarchingSquaresQuadtree{})
}

// Deprecated: RenderSVGSlow renders an SDF2 as an SVG file (uses uniform grid sampling).
func RenderSVGSlow(
	s sdf.SDF2, //sdf2 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path, lineStyle string, //path to filename
) {
	ToSVG(s, meshCells, path, lineStyle, &MarchingSquaresUniform{})
}

//-----------------------------------------------------------------------------
