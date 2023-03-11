package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
)

//-----------------------------------------------------------------------------

// MarchingCubesFEUniform renders using marching cubes with uniform space sampling.
type MarchingCubesFEUniform struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingCubesFEUniform returns a RenderHex8 object.
func NewMarchingCubesFEUniform(meshCells int) *MarchingCubesFEUniform {
	return &MarchingCubesFEUniform{
		meshCells: meshCells,
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
// Finite elements are in the shape of tetrahedra.
func (r *MarchingCubesFEUniform) RenderTet4(s sdf.SDF3, output chan<- []*Tet4) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	output <- marchingCubesTet4(s, bb, meshInc)
}

// Render produces a finite elements mesh over the bounding volume of an sdf3.
// Finite elements are in the shape of hexahedra.
func (r *MarchingCubesFEUniform) RenderHex8(s sdf.SDF3, output chan<- []*Hex8) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	output <- marchingCubesHex8(s, bb, meshInc)
}

// Render produces a finite elements mesh over the bounding volume of an sdf3.
// Finite elements are in the shape of hexahedra.
func (r *MarchingCubesFEUniform) RenderHex20(s sdf.SDF3, output chan<- []*Hex20) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	output <- marchingCubesHex20(s, bb, meshInc)
}

//-----------------------------------------------------------------------------
