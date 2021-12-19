//-----------------------------------------------------------------------------
/*

STL Load/Save

*/
//-----------------------------------------------------------------------------

package render

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// STLHeader defines the STL file header.
type STLHeader struct {
	_     [80]uint8 // Header
	Count uint32    // Number of triangles
}

// STLTriangle defines the triangle data within an STL file.
type STLTriangle struct {
	Normal, Vertex1, Vertex2, Vertex3 [3]float32
	_                                 uint16 // Attribute byte count
}

//-----------------------------------------------------------------------------

// SaveSTL writes a triangle mesh to an STL file.
func SaveSTL(path string, mesh []*Triangle3) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	header := STLHeader{}
	header.Count = uint32(len(mesh))
	if err := binary.Write(buf, binary.LittleEndian, &header); err != nil {
		return err
	}

	var d STLTriangle
	for _, triangle := range mesh {
		n := triangle.Normal()
		d.Normal[0] = float32(n.X)
		d.Normal[1] = float32(n.Y)
		d.Normal[2] = float32(n.Z)
		d.Vertex1[0] = float32(triangle.V[0].X)
		d.Vertex1[1] = float32(triangle.V[0].Y)
		d.Vertex1[2] = float32(triangle.V[0].Z)
		d.Vertex2[0] = float32(triangle.V[1].X)
		d.Vertex2[1] = float32(triangle.V[1].Y)
		d.Vertex2[2] = float32(triangle.V[1].Z)
		d.Vertex3[0] = float32(triangle.V[2].X)
		d.Vertex3[1] = float32(triangle.V[2].Y)
		d.Vertex3[2] = float32(triangle.V[2].Z)
		if err := binary.Write(buf, binary.LittleEndian, &d); err != nil {
			return err
		}
	}

	return buf.Flush()
}

//-----------------------------------------------------------------------------

// WriteSTL writes a stream of triangles to an STL file.
func WriteSTL(wg *sync.WaitGroup, path string) (chan<- *Triangle3, error) {

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	// Use buffered IO for optimal IO writes.
	// The default buffer size doesn't appear to limit performance.
	buf := bufio.NewWriter(f)

	// write an empty header
	hdr := STLHeader{}
	if err := binary.Write(buf, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}

	// External code writes triangles to this channel.
	// This goroutine reads the channel and writes triangles to the file.
	c := make(chan *Triangle3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer f.Close()

		var count uint32
		var d STLTriangle
		// read triangles from the channel and write them to the file
		for t := range c {
			n := t.Normal()
			d.Normal[0] = float32(n.X)
			d.Normal[1] = float32(n.Y)
			d.Normal[2] = float32(n.Z)
			d.Vertex1[0] = float32(t.V[0].X)
			d.Vertex1[1] = float32(t.V[0].Y)
			d.Vertex1[2] = float32(t.V[0].Z)
			d.Vertex2[0] = float32(t.V[1].X)
			d.Vertex2[1] = float32(t.V[1].Y)
			d.Vertex2[2] = float32(t.V[1].Z)
			d.Vertex3[0] = float32(t.V[2].X)
			d.Vertex3[1] = float32(t.V[2].Y)
			d.Vertex3[2] = float32(t.V[2].Z)
			if err := binary.Write(buf, binary.LittleEndian, &d); err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			count++
		}
		// flush the triangles
		buf.Flush()

		// back to the start of the file
		if _, err := f.Seek(0, 0); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		// rewrite the header with the correct mesh count
		hdr.Count = count
		if err := binary.Write(f, binary.LittleEndian, &hdr); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	return c, nil
}

//-----------------------------------------------------------------------------

// STLRenderer is an isosurface generator capable of producing 3D triangles forming a mesh from a SDF
type STLRenderer interface {
	// Render builds triangles from the specified SDF
	Render(sdf3 sdf.SDF3, meshCells int, output chan<- *Triangle3)
}

// STLRendererMarchingCubesUniform renders using marching cubes (uniform grid sampling)
type STLRendererMarchingCubesUniform struct{}

func (m *STLRendererMarchingCubesUniform) Render(s sdf.SDF3, meshCells int, output chan<- *Triangle3) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	for _, tri := range marchingCubes(s, bb, meshInc) {
		output <- tri
	}
}

// STLRendererMarchingCubesOctree renders using marching cubes (octree sampling)
type STLRendererMarchingCubesOctree struct{}

func (m *STLRendererMarchingCubesOctree) Render(s sdf.SDF3, meshCells int, output chan<- *Triangle3) {
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)

	//cells := bbSize.DivScalar(resolution).ToV3i()
	//fmt.Printf("rendering %s (%dx%dx%d, resolution %.2f)\n", cells[0], cells[1], cells[2], resolution)

	marchingCubesOctree(s, resolution, output)
}

// STLRendererDualContouring renders using dual contouring (octree sampling, sharp edges!, automatic simplification)
type STLRendererDualContouring struct {
	// Simplify: how much to simplify (if >=0).
	// NOTE: Meshing might fail with simplification enabled and greater than 0 (FIXME),
	// but the mesh might can still simplified later using external tools (the main benefit of dual contouring is sharp edges).
	Simplify float64
	// RCond [0, 1) is the parameter that controls the accuracy of sharp edges, with lower being more accurate
	// but it can cause instability leading to large wrong triangles. Leave the default if unsure.
	RCond float64
	// LockVertices makes sure each vertex stays in its voxel, avoiding small or bad triangles that may be generated
	// otherwise, but it also may remove some sharp edges.
	LockVertices bool
}

func (m *STLRendererDualContouring) Render(s sdf.SDF3, meshCells int, output chan<- *Triangle3) {
	if m.RCond == 0 {
		m.RCond = 1e-3
	}
	// work out the sampling resolution to use
	bbSize := s.BoundingBox().Size()
	resolution := bbSize.MaxComponent() / float64(meshCells)
	cells := bbSize.DivScalar(resolution).ToV3i()
	// Build the octree
	dcOctreeRootNode := dcNewOctree(cells, m.RCond, m.LockVertices)
	dcOctreeRootNode.Populate(s)
	// Simplify it
	if m.Simplify >= 0 {
		dcOctreeRootNode.Simplify(s, m.Simplify)
	}
	// Generate the final mesh
	dcOctreeRootNode.GenerateMesh(output)
}

// RenderSTL renders an SDF3 as an STL file (uses octree sampling).
func RenderSTL(
	s sdf.SDF3, //sdf3 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	// Default to marching cubes for backwards compatibility (and speed)
	RenderSTLCustom(s, meshCells, path, &STLRendererMarchingCubesOctree{})
	//RenderSTLCustom(s, meshCells, path, &STLRendererDualContouring{})
}

// RenderSTLSlow renders an SDF3 as an STL file (uses uniform grid sampling).
func RenderSTLSlow(
	s sdf.SDF3, //sdf3 to render
	meshCells int, //number of cells on the longest axis. e.g 200
	path string, //path to filename
) {
	RenderSTLCustom(s, meshCells, path, &STLRendererMarchingCubesUniform{})
}

// RenderSTLCustom renders an SDF3 as an STL file (using any STLRenderer for mesh generation).
func RenderSTLCustom(s sdf.SDF3, meshCells int, path string, renderer STLRenderer) {
	// write the triangles to an STL file
	var wg sync.WaitGroup
	output, err := WriteSTL(&wg, path)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	// run marching cubes to generate the triangle mesh
	renderer.Render(s, meshCells, output)

	// stop the STL writer reading on the channel
	close(output)
	// wait for the file write to complete
	wg.Wait()
}

//-----------------------------------------------------------------------------
