package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func marchingCubesTet4(s sdf.SDF3, box sdf.Box3, step float64) []*Tet4 {

	var tetrahedra []*Tet4
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
				tetrahedra = append(tetrahedra, mcToTet4(corners, values, 0, z)...)
				p.Z += dz
			}
			p.Y += dy
		}
		p.X += dx
	}

	return tetrahedra
}

//-----------------------------------------------------------------------------

func mcToTet4(p [8]v3.Vec, v [8]float64, x float64, layerZ int) []*Tet4 {
	// which of the 0..255 patterns do we have?
	index := 0
	for i := 0; i < 8; i++ {
		if v[i] < x {
			index |= 1 << uint(i)
		}
	}
	// do we have any triangles to create?
	if mcEdgeTable[index] == 0 {
		return nil
	}
	// work out the interpolated points on the edges
	var points [12]v3.Vec
	for i := 0; i < 12; i++ {
		bit := 1 << uint(i)
		if mcEdgeTable[index]&bit != 0 {
			a := mcPairTable[i][0]
			b := mcPairTable[i][1]
			points[i] = mcInterpolate(p[a], p[b], v[a], v[b], x)
		}
	}
	// create the triangles
	table := mcTriangleTable[index]
	count := len(table) / 3
	result := make([]*Triangle3, 0, count)
	for i := 0; i < count; i++ {
		t := Triangle3{}
		t.V[2] = points[table[i*3+0]]
		t.V[1] = points[table[i*3+1]]
		t.V[0] = points[table[i*3+2]]
		if !t.Degenerate(0) {
			result = append(result, &t)
		}
	}

	// TODO: Create tetrahedra by composing proper tables.
	resultTet4 := make([]*Tet4, 0, count)
	for _, res := range result {
		t := Tet4{
			V:     [4]v3.Vec{},
			layer: layerZ,
		}
		t.V[3] = v3.Vec{X: 0, Y: 0, Z: 0}
		t.V[2] = res.V[2]
		t.V[1] = res.V[1]
		t.V[0] = res.V[0]
		resultTet4 = append(resultTet4, &t)
	}

	return resultTet4
}

//-----------------------------------------------------------------------------

// MarchingTet4Uniform renders using marching Tetrahedra with uniform space sampling.
type MarchingTet4Uniform struct {
	meshCells int // number of cells on the longest axis of bounding box. e.g 200
}

// NewMarchingTet4Uniform returns a RenderTet4 object.
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

// To get the layer counts which are consistent with loops of marching algorithm.
func (r *MarchingTet4Uniform) LayerCounts(s sdf.SDF3) (int, int, int) {
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
func (r *MarchingTet4Uniform) Render(s sdf.SDF3, output chan<- []*Tet4) {
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

//-----------------------------------------------------------------------------
