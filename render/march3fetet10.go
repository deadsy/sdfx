package render

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func marchingCubesTet10(s sdf.SDF3, box sdf.Box3, step float64) []*Tet10 {

	var fes []*Tet10
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
				fes = append(fes, mcToTet10(corners, values, 0, z)...)
				p.Z += dz
			}
			p.Y += dy
		}
		p.X += dx
	}

	return fes
}

//-----------------------------------------------------------------------------

func mcToTet10(p [8]v3.Vec, v [8]float64, x float64, layerZ int) []*Tet10 {
	// which of the 0..255 patterns do we have?
	index := 0
	for i := 0; i < 8; i++ {
		if v[i] < x {
			index |= 1 << uint(i)
		}
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

	// Create the tetrahedra.
	table := mcTetrahedronTable[index]
	count := len(table) / 4
	result := make([]*Tet10, 0, count)
	for i := 0; i < count; i++ {
		t := Tet10{
			V:     [10]v3.Vec{},
			Layer: layerZ,
		}
		// Points on tetrahedron corners.
		t.V[0] = point(points, p, table[i*4+0])
		t.V[1] = point(points, p, table[i*4+1])
		t.V[2] = point(points, p, table[i*4+2])
		t.V[3] = point(points, p, table[i*4+3])
		// Points on tetrahedron edges.
		// Followoing CalculiX node numbering.
		t.V[5-1] = t.V[1-1].Add(t.V[2-1]).MulScalar(0.5)
		t.V[6-1] = t.V[2-1].Add(t.V[3-1]).MulScalar(0.5)
		t.V[7-1] = t.V[1-1].Add(t.V[3-1]).MulScalar(0.5)
		t.V[8-1] = t.V[1-1].Add(t.V[4-1]).MulScalar(0.5)
		t.V[9-1] = t.V[2-1].Add(t.V[4-1]).MulScalar(0.5)
		t.V[10-1] = t.V[3-1].Add(t.V[4-1]).MulScalar(0.5)
		result = append(result, &t)
	}

	return result
}

//-----------------------------------------------------------------------------
