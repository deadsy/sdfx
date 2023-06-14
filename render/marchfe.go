package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
)

//-----------------------------------------------------------------------------

// A finite element can be linear or non-linear.
type Order int

const (
	Linear    Order = iota + 1 // 4-node tetrahedron and 8-node hexahedron
	Quadratic                  // 10-node tetrahedron and 20-node hexahedron
)

//-----------------------------------------------------------------------------

// Two shapes of finite element can be generated: tetrahedral and hexahedral.
type Shape int

const (
	Hexahedral Shape = iota + 1
	Tetrahedral
	Both
)

//-----------------------------------------------------------------------------

// MarchingCubesFEUniform renders using marching cubes with uniform space sampling.
type MarchingCubesFEUniform struct {
	meshCells int   // number of cells on the longest axis of bounding box. e.g 200
	order     Order //
	shape     Shape //
}

// NewMarchingCubesFEUniform returns a RenderHex8 object.
func NewMarchingCubesFEUniform(meshCells int, order Order, shape Shape) *MarchingCubesFEUniform {
	return &MarchingCubesFEUniform{
		meshCells: meshCells,
		order:     order,
		shape:     shape,
	}
}

// Info returns a string describing the rendered volume.
func (r *MarchingCubesFEUniform) Info(s sdf.SDF3) string {
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := conv.V3ToV3i(bb1Size)
	return fmt.Sprintf("%dx%dx%d", cells.X, cells.Y, cells.Z)
}

// To get the layer counts which are consistent with loops of marching algorithm.
func (r *MarchingCubesFEUniform) LayerCounts(s sdf.SDF3) (int, int, int) {
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	size := bb.Size()
	steps := conv.V3ToV3i(size.DivScalar(meshInc).Ceil())
	return steps.X, steps.Y, steps.Z
}

// Render produces a finite elements mesh over the bounding volume of an sdf3.
// Order and shape of finite elements are selectable.
func (r *MarchingCubesFEUniform) RenderFE(s sdf.SDF3, output chan<- []*Fe) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	output <- marchingCubesFE(s, bb, meshInc, r.order, r.shape)
}

//-----------------------------------------------------------------------------
