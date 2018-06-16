//-----------------------------------------------------------------------------
/*

Render an SDF

SDF2 -> DXF file
SDF3 -> STL file

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"sync"
)

//-----------------------------------------------------------------------------

// Render an SDF3 as a STL file.
func RenderSTL(
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

// Render an SDF3 as a STL file.
func RenderSTL_New(
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
	MarchingCubes_Octree(s, resolution, output)

	// stop the STL writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
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
	MarchingSquares_Quadtree(s, resolution, output)

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
