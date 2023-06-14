package render

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func marchingCubesHex20Tet10(s sdf.SDF3, box sdf.Box3, step float64) []*Fe {
	var fes []*Fe
	size := box.Size()
	base := box.Min
	steps := conv.V3ToV3i(size.DivScalar(step).Ceil())
	inc := size.Div(conv.V3iToV3(steps))

	// start the evaluation routines
	evalRoutines()

	// create the SDF layer cache
	l := newLayerXY(base, inc, steps)
	// evaluate the SDF for z = 0
	l.Evaluate(s, 0)

	nx, ny, nz := steps.X, steps.Y, steps.Z
	dx, dy, dz := inc.X, inc.Y, inc.Z

	var p v3.Vec
	p.Z = base.Z
	for z := 0; z < nz; z++ {
		// read the z + 1 layer
		l.Evaluate(s, z+1)
		// process all cubes in the z and z + 1 layers
		p.X = base.X
		for x := 0; x < nx; x++ {
			p.Y = base.Y
			for y := 0; y < ny; y++ {
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
					l.Get(x, y, 0),
					l.Get(x+1, y, 0),
					l.Get(x+1, y+1, 0),
					l.Get(x, y+1, 0),
					l.Get(x, y, 1),
					l.Get(x+1, y, 1),
					l.Get(x+1, y+1, 1),
					l.Get(x, y+1, 1),
				}
				fes = append(fes, mcToHex20Tet10(corners, values, 0, x, y, z)...)
				p.Y += dy
			}
			p.X += dx
		}
		p.Z += dz
	}

	return fes
}

//-----------------------------------------------------------------------------

func mcToHex20Tet10(p [8]v3.Vec, v [8]float64, x float64, layerX, layerY, layerZ int) []*Fe {
	result := mcToHex20(p, v, x, layerX, layerY, layerZ)

	if len(result) < 1 {
		result = mcToTet10(p, v, x, layerX, layerY, layerZ)
	}

	return result
}

//-----------------------------------------------------------------------------
