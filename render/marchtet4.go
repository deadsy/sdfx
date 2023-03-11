package render

import (
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
	result := mcToTriangles(p, v, x)

	// TODO: Create tetrahedra by composing proper tables.

	resultTet4 := make([]*Tet4, 0, len(result))
	for _, res := range result {
		t := Tet4{
			V:     [4]v3.Vec{},
			Layer: layerZ,
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
