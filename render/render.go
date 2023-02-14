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
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// Render3 renders a 3D triangle mesh over the bounding volume of an sdf3.
type Render3 interface {
	Render(sdf3 sdf.SDF3, output chan<- []*Triangle3)
	Info(sdf3 sdf.SDF3) string
}

// Render2 renders a 2D line set over the bounding area of an sdf2.
type Render2 interface {
	Render(s sdf.SDF2, output chan<- []*Line)
	Info(s sdf.SDF2) string
}

//-----------------------------------------------------------------------------

// Renders an SDF3 to a triangle mesh and returns the vertex buffer.
func ToVertices(
	s sdf.SDF3, // sdf3 to render
	r Render3, // rendering method
) []v3.Vec {
	vertices := make([]v3.Vec, 0)

	var wg sync.WaitGroup

	// Write the triangles to corresponding vertices.
	writer := writeVertices(&wg, &vertices)

	// Run the renderer.
	r.Render(s, writer)

	// Stop the writer reading on the channel.
	close(writer)

	// Wait for the write to complete.
	wg.Wait()

	return vertices
}

// ToSTL renders an SDF3 to an STL file.
func ToSTL(
	s sdf.SDF3, // sdf3 to render
	path string, // path to filename
	r Render3, // rendering method
) {
	fmt.Printf("rendering %s (%s)\n", path, r.Info(s))
	// write the triangles to an STL file
	var wg sync.WaitGroup
	output, err := WriteSTL(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	// run the renderer
	r.Render(s, output)
	// stop the STL writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

//-----------------------------------------------------------------------------

// To3MF renders an SDF3 to a 3MF file.
func To3MF(
	s sdf.SDF3, // sdf3 to render
	path string, // path to filename
	r Render3, // rendering method
) {
	fmt.Printf("rendering %s (%s)\n", path, r.Info(s))
	// write the triangles to a 3MF file
	var wg sync.WaitGroup
	output, err := Write3MF(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	// run the renderer
	r.Render(s, output)
	// stop the STL writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

//-----------------------------------------------------------------------------

// ToDXF renders an SDF2 to a DXF file.
func ToDXF(
	s sdf.SDF2, // sdf2 to render
	path string, // path to filename
	r Render2, // rendering method
) {
	fmt.Printf("rendering %s (%s)\n", path, r.Info(s))
	// write the line segments to a DXF file
	var wg sync.WaitGroup
	output, err := WriteDXF(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	// run the renderer
	r.Render(s, output)
	// stop the DXF writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

//-----------------------------------------------------------------------------

const svgLineStyle = "fill:none;stroke:black;stroke-width:0.1"

// ToSVG renders an SDF2 to an SVG file.
func ToSVG(
	s sdf.SDF2, // sdf2 to render
	path string, // path to filename
	r Render2, // rendering method
) {
	fmt.Printf("rendering %s (%s)\n", path, r.Info(s))
	// write the line segments to an SVG file
	var wg sync.WaitGroup
	output, err := WriteSVG(&wg, path, svgLineStyle)
	if err != nil {
		fmt.Printf("%s", err)
	}
	// run the renderer
	r.Render(s, output)
	// stop the SVG writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

//-----------------------------------------------------------------------------
