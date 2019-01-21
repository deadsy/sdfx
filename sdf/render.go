//-----------------------------------------------------------------------------
/*

Render an SDF

SDF3 -> STL file
SDF2 -> DXF file

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"sync"
)

//-----------------------------------------------------------------------------

// Render an SDF3 as an STL file (octree sampling)
func RenderSTL(
	s SDF3, //sdf3 to render
	mesh_cells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {

	// work out the sampling resolution to use
	bb_size := s.BoundingBox().Size()
	resolution := bb_size.MaxComponent() / float64(mesh_cells)
	cells := bb_size.DivScalar(resolution).ToV3i()

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

// Render an SDF3 as an STL file.
func RenderSTL_Slow(
	s SDF3, //sdf3 to render
	mesh_cells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0_size := bb0.Size()
	mesh_inc := bb0_size.MaxComponent() / float64(mesh_cells)
	bb1_size := bb0_size.DivScalar(mesh_inc)
	bb1_size = bb1_size.Ceil().AddScalar(1)
	cells := bb1_size.ToV3i()
	bb1_size = bb1_size.MulScalar(mesh_inc)
	bb := NewBox3(bb0.Center(), bb1_size)

	fmt.Printf("rendering %s (%dx%dx%d)\n", path, cells[0], cells[1], cells[2])

	// run marching cubes to generate the triangle mesh
	m := MarchingCubes(s, bb, mesh_inc)
	err := SaveSTL(path, m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------

// Render an SDF2 as a DXF file. (quadtree sampling)
func RenderDXF(
	s SDF2, //sdf2 to render
	mesh_cells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {

	// work out the sampling resolution to use
	bb_size := s.BoundingBox().Size()
	resolution := bb_size.MaxComponent() / float64(mesh_cells)
	cells := bb_size.DivScalar(resolution).ToV2i()

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

// Render an SDF2 as a DXF file. (grid sampling)
func RenderDXF_Slow(
	s SDF2, //sdf2 to render
	mesh_cells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0_size := bb0.Size()
	mesh_inc := bb0_size.MaxComponent() / float64(mesh_cells)
	bb1_size := bb0_size.DivScalar(mesh_inc)
	bb1_size = bb1_size.Ceil().AddScalar(1)
	cells := bb1_size.ToV2i()
	bb1_size = bb1_size.MulScalar(mesh_inc)
	bb := NewBox2(bb0.Center(), bb1_size)

	fmt.Printf("rendering %s (%dx%d)\n", path, cells[0], cells[1])

	// run marching squares to generate the line segments
	m := MarchingSquares(s, bb, mesh_inc)
	err := SaveDXF(path, m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------

// Render an SDF2 as an SVG file. (quadtree sampling)
func RenderSVG(
	s SDF2, // sdf2 to render
	meshCells int, // number of cells on the longest axis. e.g 200
	path string, // path to filename
) error {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV2i()

	fmt.Printf("rendering %s (%dx%d, resolution %.2f)\n", path, cells[0], cells[1], resolution)

	// write the line segments to an SVG file
	var wg sync.WaitGroup
	output, err := WriteSVG(&wg, path)
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

// Render an SDF2 as an SVG file. (grid sampling)
func RenderSVG_Slow(
	s SDF2, // sdf2 to render
	meshCells int, // number of cells on the longest axis. e.g 200
	path string, // path to filename
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
	m := MarchingSquares(s, bb, meshInc)
	return SaveSVG(path, m)
}

//-----------------------------------------------------------------------------
