//-----------------------------------------------------------------------------
/*

Render an SDF

SDF3 -> STL file
SDF2 -> DXF file
SDF2 -> SVG file

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"sync"
)

//-----------------------------------------------------------------------------

// RenderSTL renders an SDF3 as an STL file (uses octree sampling).
func RenderSTL(
	s SDF3, //sdf3 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {

	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV3i()

	fmt.Printf("rendering %s (%dx%dx%d, resolution %.2f)\n", path, cells[0], cells[1], cells[2], resolution)

	// write the triangles to an STL file
	var wg sync.WaitGroup
	output, err := WriteSTL(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	// run marching cubes to generate the triangle mesh
	marchingCubesOctree(s, resolution, output)

	// stop the STL writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

// RenderSTLSlow renders an SDF3 as an STL file (uses uniform grid sampling).
func RenderSTLSlow(
	s SDF3, //sdf3 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := bb1Size.ToV3i()
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := NewBox3(bb0.Center(), bb1Size)

	fmt.Printf("rendering %s (%dx%dx%d)\n", path, cells[0], cells[1], cells[2])

	// run marching cubes to generate the triangle mesh
	m := marchingCubes(s, bb, meshInc)
	err := SaveSTL(path, m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------

// RenderDXF renders an SDF2 as a DXF file. (uses quadtree sampling)
func RenderDXF(
	s SDF2, //sdf2 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {

	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV2i()

	fmt.Printf("rendering %s (%dx%d, resolution %.2f)\n", path, cells[0], cells[1], resolution)

	// write the line segments to a DXF file
	var wg sync.WaitGroup
	output, err := WriteDXF(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	// run marching squares to generate the line segments
	marchingSquaresQuadtree(s, resolution, output)

	// stop the DXF writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

// RenderDXFSlow renders an SDF2 as a DXF file. (uses uniform grid sampling)
func RenderDXFSlow(
	s SDF2, //sdf2 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := bb1Size.ToV2i()
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := NewBox2(bb0.Center(), bb1Size)

	fmt.Printf("rendering %s (%dx%d)\n", path, cells[0], cells[1])

	// run marching squares to generate the line segments
	m := marchingSquares(s, bb, meshInc)
	err := SaveDXF(path, m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------

// RenderSVG renders an SDF2 as an SVG file. (uses quadtree sampling)
func RenderSVG(
	s SDF2, // sdf2 to render
	meshCells int, // number of cells on the longest axis. e.g 200
	path string, // path to filename
	lineStyle string, // SVG line style
) error {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV2i()

	fmt.Printf("rendering %s (%dx%d, resolution %.2f)\n", path, cells[0], cells[1], resolution)

	// write the line segments to an SVG file
	var wg sync.WaitGroup
	output, err := WriteSVG(&wg, path, lineStyle)
	if err != nil {
		return err
	}

	// run marching squares to generate the line segments
	marchingSquaresQuadtree(s, resolution, output)

	// stop the SVG writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
	return nil
}

// RenderSVGSlow renders an SDF2 as an SVG file. (uses uniform grid sampling)
func RenderSVGSlow(
	s SDF2, // sdf2 to render
	meshCells int, // number of cells on the longest axis. e.g 200
	path string, // path to filename
	lineStyle string, // SVG line style
) error {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := bb1Size.ToV2i()
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := NewBox2(bb0.Center(), bb1Size)

	fmt.Printf("rendering %s (%dx%d)\n", path, cells[0], cells[1])

	// run marching squares to generate the line segments
	m := marchingSquares(s, bb, meshInc)
	return SaveSVG(path, lineStyle, m)
}

//-----------------------------------------------------------------------------
