package render

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func marchingCubesHex8(s sdf.SDF3, box sdf.Box3, step float64) []*Hex8 {

	var fes []*Hex8
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
				fes = append(fes, mcToHex8(corners, values, 0, z)...)
				p.Z += dz
			}
			p.Y += dy
		}
		p.X += dx
	}

	return fes
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

		// Refer to CalculiX solver documentation:
		// http://www.dhondt.de/ccx_2.20.pdf

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
