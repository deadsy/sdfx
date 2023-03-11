package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func marchingCubesHex8(s sdf.SDF3, box sdf.Box3, step float64) []*Hex8 {

	var tetrahedra []*Hex8
	size := box.Size()
	base := box.Min
	steps := conv.V3ToV3i(size.DivScalar(step).Ceil())
	inc := size.Div(conv.V3iToV3(steps))

	// start the evaluation routines
	evalRoutines()

	// create the SDF layer cache
	l := newLayerYZ(base, inc, steps)
	// evaluate the SDF for x = 0
	l.Evaluate(s, 0)

	nx, ny, nz := steps.X, steps.Y, steps.Z
	dx, dy, dz := inc.X, inc.Y, inc.Z

	var p v3.Vec
	p.X = base.X
	for x := 0; x < nx; x++ {
		// read the x + 1 layer
		l.Evaluate(s, x+1)
		// process all cubes in the x and x + 1 layers
		p.Y = base.Y
		for y := 0; y < ny; y++ {
			p.Z = base.Z
			for z := 0; z < nz; z++ {
				x0, y0, z0 := p.X, p.Y, p.Z
				x1, y1, z1 := x0+dx, y0+dy, z0+dz
				corners := [8]v3.Vec{
					{x0, y0, z0},
					{x1, y0, z0},
					{x1, y1, z0},
					{x0, y1, z0},
					{x0, y0, z1},
					{x1, y0, z1},
					{x1, y1, z1},
					{x0, y1, z1}}
				values := [8]float64{
					l.Get(0, y, z),
					l.Get(1, y, z),
					l.Get(1, y+1, z),
					l.Get(0, y+1, z),
					l.Get(0, y, z+1),
					l.Get(1, y, z+1),
					l.Get(1, y+1, z+1),
					l.Get(0, y+1, z+1)}
				tetrahedra = append(tetrahedra, mcToHex8(corners, values, 0, z)...)
				p.Z += dz
			}
			p.Y += dy
		}
		p.X += dx
	}

	return tetrahedra
}

//-----------------------------------------------------------------------------

func mcToHex8(p [8]v3.Vec, v [8]float64, x float64, layerZ int) []*Hex8 {
	result := make([]*Hex8, 0)

	anyPositive := false
	for i := 0; i < 8; i++ {
		if v[i] > 0 {
			anyPositive = true
			break
		}
	}

	// Create a finite element if all 8 values are non-positive.
	// Finite element is inside the 3D model if all values are non-positive.
	// Of course, some spaces are missed by this approach.
	//
	// TODO: Come up with a more sophisticated approach?

	if !anyPositive {
		fe := Hex8{
			V:     [8]v3.Vec{},
			Layer: layerZ,
		}
		fe.V[7] = p[7]
		fe.V[6] = p[6]
		fe.V[5] = p[5]
		fe.V[4] = p[4]
		fe.V[3] = p[3]
		fe.V[2] = p[2]
		fe.V[1] = p[1]
		fe.V[0] = p[0]
		result = append(result, &fe)
	}

	return result
}

//-----------------------------------------------------------------------------

// MarchingHex8Uniform renders using marching cubes with uniform space sampling.
type MarchingHex8Uniform struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingHex8Uniform returns a RenderHex8 object.
func NewMarchingHex8Uniform(meshCells int) *MarchingHex8Uniform {
	return &MarchingHex8Uniform{
		meshCells: meshCells,
	}
}

// Info returns a string describing the rendered volume.
func (r *MarchingHex8Uniform) Info(s sdf.SDF3) string {
	bb0 := s.BoundingBox()
	bb0Size := bb0.Size()
	meshInc := bb0Size.MaxComponent() / float64(r.meshCells)
	bb1Size := bb0Size.DivScalar(meshInc)
	bb1Size = bb1Size.Ceil().AddScalar(1)
	cells := conv.V3ToV3i(bb1Size)
	return fmt.Sprintf("%dx%dx%d", cells.X, cells.Y, cells.Z)
}

// To get the layer counts which are consistent with loops of marching algorithm.
func (r *MarchingHex8Uniform) LayerCounts(s sdf.SDF3) (int, int, int) {
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
// Finite elements are in the shape of hexahedra.
func (r *MarchingHex8Uniform) Render(s sdf.SDF3, output chan<- []*Hex8) {
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

//-----------------------------------------------------------------------------
