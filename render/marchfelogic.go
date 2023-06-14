package render

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func marchingCubesFE(s sdf.SDF3, box sdf.Box3, step float64, order Order, shape Shape) []*Fe {
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
				fes = append(fes, mcToFE(corners, values, x, y, z, order, shape)...)
				p.Y += dy
			}
			p.X += dx
		}
		p.Z += dz
	}

	return fes
}

//-----------------------------------------------------------------------------

func mcToFE(corners [8]v3.Vec, values [8]float64, x, y, z int, order Order, shape Shape) []*Fe {
	var fes []*Fe
	switch order {
	case Linear:
		{
			switch shape {
			case Hexahedral:
				{
					fes = append(fes, mcToHex8(corners, values, 0, x, y, z)...)
				}
			case Tetrahedral:
				{
					fes = append(fes, mcToTet4(corners, values, 0, x, y, z)...)
				}
			case Both:
				{
					tmp := mcToHex8(corners, values, 0, x, y, z)
					if len(tmp) < 1 {
						tmp = mcToTet4(corners, values, 0, x, y, z)
					}
					fes = append(fes, tmp...)
				}
			}
		}
	case Quadratic:
		{
			switch shape {
			case Hexahedral:
				{
					fes = append(fes, mcToHex20(corners, values, 0, x, y, z)...)
				}
			case Tetrahedral:
				{
					fes = append(fes, mcToTet10(corners, values, 0, x, y, z)...)
				}
			case Both:
				{
					tmp := mcToHex20(corners, values, 0, x, y, z)
					if len(tmp) < 1 {
						tmp = mcToTet10(corners, values, 0, x, y, z)
					}
					fes = append(fes, tmp...)
				}
			}
		}
	}

	return fes
}

//-----------------------------------------------------------------------------

func mcToHex8(p [8]v3.Vec, v [8]float64, x float64, layerX, layerY, layerZ int) []*Fe {
	result := make([]*Fe, 0)

	anyPositive := false
	for i := 0; i < 8; i++ {
		if v[i] > 0 {
			anyPositive = true
			break
		}
	}

	// Create a finite element if all 8 values are non-positive.
	// Finite element is inside the 3D model if all values are non-positive.

	if !anyPositive {
		fe := Fe{
			V: make([]v3.Vec, 8),
			X: layerX,
			Y: layerY,
			Z: layerZ,
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

func mcToHex20(p [8]v3.Vec, v [8]float64, x float64, layerX, layerY, layerZ int) []*Fe {
	result := make([]*Fe, 0)

	anyPositive := false
	for i := 0; i < 8; i++ {
		if v[i] > 0 {
			anyPositive = true
			break
		}
	}

	// Create a finite element if all 8 values are non-positive.
	// Finite element is inside the 3D model if all values are non-positive.

	if !anyPositive {
		fe := Fe{
			V: make([]v3.Vec, 20),
			X: layerX,
			Y: layerY,
			Z: layerZ,
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

func mcToTet4(p [8]v3.Vec, v [8]float64, x float64, layerX, layerY, layerZ int) []*Fe {
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
	result := make([]*Fe, 0, count)
	for i := 0; i < count; i++ {
		t := Fe{
			V: make([]v3.Vec, 4),
			X: layerX,
			Y: layerY,
			Z: layerZ,
		}

		t.V[0] = point(points, p, table[i*4+0])
		t.V[1] = point(points, p, table[i*4+1])
		t.V[2] = point(points, p, table[i*4+2])
		t.V[3] = point(points, p, table[i*4+3])
		degenerated := degenerateTriangles(t.V[0], t.V[1], t.V[2], t.V[3])
		flat, volume := almostFlat(t.V[0], t.V[1], t.V[2], t.V[3])

		// In the case of marching cubes algorithm to generate triangle, it's avoiding zero-area triangles by `!t.Degenerate(0)` check.
		// In our case of marching cubes algorithm to generate tetrahedron, we can do a check too:
		bad, jacobianDeterminant := isBadTet4([4]v3.Vec{t.V[0], t.V[1], t.V[2], t.V[3]})
		if !degenerated && !bad && !flat {
			result = append(result, &t)
		} else {
			fmt.Println("Bad element: tet4")
			fmt.Println("Non-positive Jacobian determinant? ", bad)
			fmt.Println("Jacobian determinant: ", jacobianDeterminant)
			fmt.Println("Almost flat? ", flat)
			fmt.Println("Volume: ", volume)
			fmt.Println("Degenerated? ", degenerated)
		}
	}

	return result
}

//-----------------------------------------------------------------------------

func mcToTet10(p [8]v3.Vec, v [8]float64, x float64, layerX, layerY, layerZ int) []*Fe {
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
	result := make([]*Fe, 0, count)
	for i := 0; i < count; i++ {
		t := Fe{
			V: make([]v3.Vec, 10),
			X: layerX,
			Y: layerY,
			Z: layerZ,
		}

		// Points on tetrahedron corners.
		t.V[0] = point(points, p, table[i*4+0])
		t.V[1] = point(points, p, table[i*4+1])
		t.V[2] = point(points, p, table[i*4+2])
		t.V[3] = point(points, p, table[i*4+3])
		degenerated := degenerateTriangles(t.V[0], t.V[1], t.V[2], t.V[3])
		flat, volume := almostFlat(t.V[0], t.V[1], t.V[2], t.V[3])
		// Points on tetrahedron edges.
		// Followoing CalculiX node numbering.
		t.V[5-1] = t.V[1-1].Add(t.V[2-1]).MulScalar(0.5)
		t.V[6-1] = t.V[2-1].Add(t.V[3-1]).MulScalar(0.5)
		t.V[7-1] = t.V[1-1].Add(t.V[3-1]).MulScalar(0.5)
		t.V[8-1] = t.V[1-1].Add(t.V[4-1]).MulScalar(0.5)
		t.V[9-1] = t.V[2-1].Add(t.V[4-1]).MulScalar(0.5)
		t.V[10-1] = t.V[3-1].Add(t.V[4-1]).MulScalar(0.5)
		// In the case of marching cubes algorithm to generate triangle, it's avoiding zero-area triangles by `!t.Degenerate(0)` check.
		// In our case of marching cubes algorithm to generate tetrahedron, we can do a check too:
		bad, jacobianDeterminant := isBadTet10([10]v3.Vec{t.V[0], t.V[1], t.V[2], t.V[3], t.V[4], t.V[5], t.V[6], t.V[7], t.V[8], t.V[9]})
		if !degenerated && !bad && !flat {
			result = append(result, &t)
		} else {
			fmt.Println("Bad element: tet10")
			fmt.Println("Non-positive Jacobian determinant? ", bad)
			fmt.Println("Jacobian determinant: ", jacobianDeterminant)
			fmt.Println("Almost flat? ", flat)
			fmt.Println("Volume: ", volume)
			fmt.Println("Degenerated? ", degenerated)
		}
	}

	return result
}

//-----------------------------------------------------------------------------
