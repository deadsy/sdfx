package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func marchingTet4(s sdf.SDF3, box sdf.Box3, step float64) []*Tet4 {
	fmt.Printf("marching tetrahedra, bbox center: %v , step: %v\n", s.BoundingBox().Center(), step)
	var tetrahedra []*Tet4

	// Constant hard-coded tetrahedra vertices to develop and debug the output API.
	// https://cs.stackexchange.com/a/90011/67985
	tetrahedra = append(tetrahedra, &Tet4{
		V: [4]v3.Vec{{X: 0, Y: 0, Z: 0}, {X: 0, Y: 0, Z: 1}, {X: 0, Y: 1, Z: 1}, {X: 1, Y: 1, Z: 1}},
	})
	tetrahedra = append(tetrahedra, &Tet4{
		V: [4]v3.Vec{{X: 0, Y: 0, Z: 0}, {X: 0, Y: 1, Z: 0}, {X: 0, Y: 1, Z: 1}, {X: 1, Y: 1, Z: 1}},
	})
	tetrahedra = append(tetrahedra, &Tet4{
		V: [4]v3.Vec{{X: 0, Y: 0, Z: 0}, {X: 0, Y: 0, Z: 1}, {X: 1, Y: 0, Z: 1}, {X: 1, Y: 1, Z: 1}},
	})
	tetrahedra = append(tetrahedra, &Tet4{
		V: [4]v3.Vec{{X: 0, Y: 0, Z: 0}, {X: 1, Y: 0, Z: 0}, {X: 1, Y: 0, Z: 1}, {X: 1, Y: 1, Z: 1}},
	})
	tetrahedra = append(tetrahedra, &Tet4{
		V: [4]v3.Vec{{X: 0, Y: 0, Z: 0}, {X: 0, Y: 1, Z: 0}, {X: 1, Y: 1, Z: 0}, {X: 1, Y: 1, Z: 1}},
	})
	tetrahedra = append(tetrahedra, &Tet4{
		V: [4]v3.Vec{{X: 0, Y: 0, Z: 0}, {X: 1, Y: 0, Z: 0}, {X: 1, Y: 1, Z: 0}, {X: 1, Y: 1, Z: 1}},
	})

	// TODO: Logic.

	return tetrahedra
}

//-----------------------------------------------------------------------------

// MarchingTet4Uniform renders using marching Tetrahedra with uniform space sampling.
type MarchingTet4Uniform struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingTet4Uniform returns a RenderFE object.
func NewMarchingTet4Uniform(meshCells int) *MarchingTet4Uniform {
	return &MarchingTet4Uniform{
		meshCells: meshCells,
	}
}

// Info returns a string describing the rendered volume.
func (r *MarchingTet4Uniform) Info(s sdf.SDF3) string {
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
func (r *MarchingTet4Uniform) Render(s sdf.SDF3, output chan<- []*Tet4) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	output <- marchingTet4(s, bb, meshInc)
}

//-----------------------------------------------------------------------------
