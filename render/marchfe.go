package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
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
	order     Order // Linear or quadratic.
	shape     Shape // Hexahedral, tetrahedral, or both.
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

// Render produces a finite elements mesh over the bounding volume of an sdf3.
// Order and shape of finite elements are selectable.
func (r *MarchingCubesFEUniform) RenderFE(s sdf.SDF3, output sdf.FeWriter) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)
	marchingCubesFE(s, bb, meshInc, r.order, r.shape, output)
}

//-----------------------------------------------------------------------------

// To get the voxel count, dimension, and min/max corner which are consistent with loops of marching algorithm.
// This func loops are exactly like `marchingCubesFE` loops. We have to be consistant.
func (r *MarchingCubesFEUniform) Voxels(s sdf.SDF3) (v3i.Vec, v3.Vec, []v3.Vec, []v3.Vec) {
	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	bb1Size = bb1Size.MulScalar(meshInc)
	bb := sdf.NewBox3(bb0.Center(), bb1Size)

	size := bb.Size()
	base := bb.Min
	steps := conv.V3ToV3i(size.DivScalar(meshInc).Ceil())
	inc := size.Div(conv.V3iToV3(steps))

	nx, ny, nz := steps.X, steps.Y, steps.Z
	dx, dy, dz := inc.X, inc.Y, inc.Z

	mins := make([]v3.Vec, 0, nz*nx*ny)
	maxs := make([]v3.Vec, 0, nz*nx*ny)

	var p v3.Vec
	p.X = base.X
	for x := 0; x < nx; x++ {
		p.Y = base.Y
		for y := 0; y < ny; y++ {
			p.Z = base.Z
			for z := 0; z < nz; z++ {
				x0, y0, z0 := p.X, p.Y, p.Z
				x1, y1, z1 := x0+dx, y0+dy, z0+dz

				mins = append(mins, v3.Vec{X: x0, Y: y0, Z: z0})
				maxs = append(maxs, v3.Vec{X: x1, Y: y1, Z: z1})

				p.Z += dz
			}
			p.Y += dy
		}
		p.X += dx
	}

	return v3i.Vec{X: nx, Y: ny, Z: nz}, v3.Vec{X: dx, Y: dy, Z: dz}, mins, maxs
}

//-----------------------------------------------------------------------------
