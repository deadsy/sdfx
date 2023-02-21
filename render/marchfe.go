package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
)

//-----------------------------------------------------------------------------

func marchingTetrahedra(s sdf.SDF3, box sdf.Box3, step float64) []*Tetrahedron {
	// TODO: Logic.
	fmt.Printf("marching tetrahedra, bbox center: %v , step: %v\n", s.BoundingBox().Center(), step)
	return nil
}

//-----------------------------------------------------------------------------

// MarchingTetrahedraUniform renders using marching Tetrahedra with uniform space sampling.
type MarchingTetrahedraUniform struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingTetrahedraUniform returns a RenderFE object.
func NewMarchingTetrahedraUniform(meshCells int) *MarchingTetrahedraUniform {
	return &MarchingTetrahedraUniform{
		meshCells: meshCells,
	}
}

// Info returns a string describing the rendered volume.
func (r *MarchingTetrahedraUniform) Info(s sdf.SDF3) string {
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := conv.V3ToV3i(bb1Size)
	return fmt.Sprintf("%dx%dx%d", cells.X, cells.Y, cells.Z)
}

// Render produces a finite elements mesh over the bounding volume of an sdf3.
// Finite elements are in the shape of tetrahedra.
func (r *MarchingTetrahedraUniform) Render(s sdf.SDF3, output chan<- []*Tetrahedron) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	output <- marchingTetrahedra(s, bb, meshInc)
}

//-----------------------------------------------------------------------------
