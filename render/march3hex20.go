package render

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func marchingCubesHex20(s sdf.SDF3, box sdf.Box3, step float64) []*Hex20 {

	var fes []*Hex20
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
				fes = append(fes, mcToHex20(corners, values, 0, z)...)
				p.Z += dz
			}
			p.Y += dy
		}
		p.X += dx
	}

	return fes
}

//-----------------------------------------------------------------------------

func mcToHex20(p [8]v3.Vec, v [8]float64, x float64, layerZ int) []*Hex20 {
	result := make([]*Hex20, 0)

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
		fe := Hex20{
			V:     [20]v3.Vec{},
			Layer: layerZ,
		}

		// Refer to CalculiX solver documentation:
		// http://www.dhondt.de/ccx_2.20.pdf

		// Points on cube corners:
		fe.V[7] = p[7]
		fe.V[6] = p[6]
		fe.V[5] = p[5]
		fe.V[4] = p[4]
		fe.V[3] = p[3]
		fe.V[2] = p[2]
		fe.V[1] = p[1]
		fe.V[0] = p[0]

		// Points on cube edges:
		fe.V[8] = p[0].Add(p[1]).MulScalar(0.5)
		fe.V[9] = p[1].Add(p[2]).MulScalar(0.5)
		fe.V[10] = p[2].Add(p[3]).MulScalar(0.5)
		fe.V[11] = p[3].Add(p[0]).MulScalar(0.5)

		fe.V[12] = p[4].Add(p[5]).MulScalar(0.5)
		fe.V[13] = p[5].Add(p[6]).MulScalar(0.5)
		fe.V[14] = p[6].Add(p[7]).MulScalar(0.5)
		fe.V[15] = p[7].Add(p[4]).MulScalar(0.5)

		fe.V[16] = p[0].Add(p[4]).MulScalar(0.5)
		fe.V[17] = p[1].Add(p[5]).MulScalar(0.5)
		fe.V[18] = p[2].Add(p[6]).MulScalar(0.5)
		fe.V[19] = p[3].Add(p[7]).MulScalar(0.5)

		result = append(result, &fe)
	}

	return result
}

//-----------------------------------------------------------------------------
