package render

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// Specify the point to create the tetrahedra.
// Point can be on edges or corners.
// Index from 0 to 11 means an edge.
// Index from 12 to 19 means a corner.
// Corners were originally indexed from 0 to 7 but they are shifted by 12.
// So, corners are from 0+12 to 7+12 i.e. from 12 to 19.
func point(edges [12]v3.Vec, corners [8]v3.Vec, index int) v3.Vec {
	if 0 <= index && index < 12 {
		return edges[index]
	} else if index < 20 {
		return corners[index-12]
	} else {
		// Should never reach here.
		return v3.Vec{}
	}
}

// -----------------------------------------------------------------------------
// To avoid distorted tetrahedron with non-positive Jacobian or with flat shape or with degenerate faces.
// Regradless of input values at corners, return the point at the middle of the edge.
func mcInterpolateFE(p1, p2 v3.Vec, v1, v2, x float64) v3.Vec {
	// Pick the half way point
	t := 0.5

	return v3.Vec{
		X: p1.X + t*(p2.X-p1.X),
		Y: p1.Y + t*(p2.Y-p1.Y),
		Z: p1.Z + t*(p2.Z-p1.Z),
	}
}

//-----------------------------------------------------------------------------

// Is a tetrahedron almost flat?
// To avoid mathematical problem while FEA is run.
// MATHEMATICA script is available here:
// https://math.stackexchange.com/a/4709610/197913
func almostFlat(a, b, c, d v3.Vec) (bool, float64) {
	ab := b.Sub(a)
	ac := c.Sub(a)
	ad := d.Sub(a)

	// Note that the `Norm` function of MATHEMATICA is equivalent to our `Length()` function.
	nab := ab.Length()
	ncd := d.Sub(c).Length()
	nbd := d.Sub(b).Length()
	nbc := c.Sub(b).Length()
	nac := ac.Length()
	nad := ad.Length()

	// Check for 0 edge lengths
	if nab == 0 || ncd == 0 ||
		nbd == 0 || nbc == 0 ||
		nac == 0 || nad == 0 {
		return true, 0
	}

	volume := 1.0 / 6.0 * math.Abs(ab.Cross(ac).Dot(ad))
	denom := (nab + ncd) * (nac + nbd) * (nad + nbc)

	// Tolerance derived from here would be `480.0`:
	// https://math.stackexchange.com/a/4709610/197913
	// A different value is calibrated according to observations.
	// TODO: Could be further calibrated.
	tolerance := 1000.0

	rho := tolerance * volume / denom

	return rho < 1, volume
}

//-----------------------------------------------------------------------------

// Reference:
// CCX source code:
// ccx_2.20/src/shape4tet.f
func isBadGaussTet4(coords [4]v3.Vec, xi, et, ze float64) (bool, float64) {
	// Coordinates of the nodes.
	var xl [3][4]float64

	for i := 0; i < 4; i++ {
		xl[0][i] = coords[i].X
		xl[1][i] = coords[i].Y
		xl[2][i] = coords[i].Z
	}

	// Shape functions.
	var shp [4][4]float64

	shp[3][0] = 1.0 - xi - et - ze
	shp[3][1] = xi
	shp[3][2] = et
	shp[3][3] = ze

	// local derivatives of the shape functions: xi-derivative

	shp[0][0] = -1.0
	shp[0][1] = 1.0
	shp[0][2] = 0.0
	shp[0][3] = 0.0

	// local derivatives of the shape functions: eta-derivative

	shp[1][0] = -1.0
	shp[1][1] = 0.0
	shp[1][2] = 1.0
	shp[1][3] = 0.0

	// local derivatives of the shape functions: zeta-derivative

	shp[2][0] = -1.0
	shp[2][1] = 0.0
	shp[2][2] = 0.0
	shp[2][3] = 1.0

	// computation of the local derivative of the global coordinates (xs)
	xs := [3][3]float64{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			xs[i][j] = 0.0
			for k := 0; k < 4; k++ {
				xs[i][j] += xl[i][k] * shp[j][k]
			}
		}
	}

	// computation of the jacobian determinant
	xsj := xs[0][0]*(xs[1][1]*xs[2][2]-xs[1][2]*xs[2][1]) -
		xs[0][1]*(xs[1][0]*xs[2][2]-xs[1][2]*xs[2][0]) +
		xs[0][2]*(xs[1][0]*xs[2][1]-xs[1][1]*xs[2][0])

	// According to CCX source code to detect nonpositive jacobian determinant in element
	// Fortran threshold for non-positive Jacobian determinant is 1e-20.
	return xsj < 1e-20, xsj
}

// Reference:
// CCX source code:
// ccx_2.20/src/shape10tet.f
func isBadGaussTet10(coords [10]v3.Vec, xi, et, ze float64) (bool, float64) {
	// Coordinates of the nodes.
	var xl [3][10]float64

	for i := 0; i < 10; i++ {
		xl[0][i] = coords[i].X
		xl[1][i] = coords[i].Y
		xl[2][i] = coords[i].Z
	}

	// Shape functions.
	var shp [4][10]float64

	// Shape functions
	a := 1.0 - xi - et - ze
	shp[3][0] = (2.0*a - 1.0) * a
	shp[3][1] = xi * (2.0*xi - 1.0)
	shp[3][2] = et * (2.0*et - 1.0)
	shp[3][3] = ze * (2.0*ze - 1.0)
	shp[3][4] = 4.0 * xi * a
	shp[3][5] = 4.0 * xi * et
	shp[3][6] = 4.0 * et * a
	shp[3][7] = 4.0 * ze * a
	shp[3][8] = 4.0 * xi * ze
	shp[3][9] = 4.0 * et * ze

	// Local derivatives of the shape functions: xi-derivative
	shp[0][0] = 1.0 - 4.0*a
	shp[0][1] = 4.0*xi - 1.0
	shp[0][2] = 0.0
	shp[0][3] = 0.0
	shp[0][4] = 4.0 * (a - xi)
	shp[0][5] = 4.0 * et
	shp[0][6] = -4.0 * et
	shp[0][7] = -4.0 * ze
	shp[0][8] = 4.0 * ze
	shp[0][9] = 0.0

	// Local derivatives of the shape functions: eta-derivative
	shp[1][0] = 1.0 - 4.0*a
	shp[1][1] = 0.0
	shp[1][2] = 4.0*et - 1.0
	shp[1][3] = 0.0
	shp[1][4] = -4.0 * xi
	shp[1][5] = 4.0 * xi
	shp[1][6] = 4.0 * (a - et)
	shp[1][7] = -4.0 * ze
	shp[1][8] = 0.0
	shp[1][9] = 4.0 * ze

	// Local derivatives of the shape functions: zeta-derivative
	shp[2][0] = 1.0 - 4.0*a
	shp[2][1] = 0.0
	shp[2][2] = 0.0
	shp[2][3] = 4.0*ze - 1.0
	shp[2][4] = -4.0 * xi
	shp[2][5] = 0.0
	shp[2][6] = -4.0 * et
	shp[2][7] = 4.0 * (a - ze)
	shp[2][8] = 4.0 * xi
	shp[2][9] = 4.0 * et

	// Computation of the local derivative of the global coordinates (xs)
	var xs [3][3]float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			xs[i][j] = 0.0
			for k := 0; k < 10; k++ {
				xs[i][j] = xs[i][j] + xl[i][k]*shp[j][k]
			}
		}
	}

	// computation of the jacobian determinant
	xsj := xs[0][0]*(xs[1][1]*xs[2][2]-xs[1][2]*xs[2][1]) -
		xs[0][1]*(xs[1][0]*xs[2][2]-xs[1][2]*xs[2][0]) +
		xs[0][2]*(xs[1][0]*xs[2][1]-xs[1][1]*xs[2][0])

	// According to CCX source code to detect nonpositive jacobian determinant in element
	// Fortran threshold for non-positive Jacobian determinant is 1e-20.
	return xsj < 1e-20, xsj
}

//-----------------------------------------------------------------------------

// Exactly follow CCX source code leading to this error:
// *ERROR in e_c3d: nonpositive jacobian determinant in element
func isBadTet4(coords [4]v3.Vec) (bool, float64) {

	// xi, et, and ze are the coordinates of the Gauss point
	// in the integration scheme for the 4-node tetrahedral element.
	// For this element type, there is typically only 1 Gauss point used,
	// which is located at the centroid of the tetrahedron.
	// The coordinates of this Gauss point are (xi, et, ze) = (1/4, 1/4, 1/4).
	// Reference:
	// ccx_2.20/src/gauss.f
	var xi float64 = 0.25
	var et float64 = 0.25
	var ze float64 = 0.25

	return isBadGaussTet4(coords, xi, et, ze)
}

// Exactly follow CCX source code leading to this error:
// *ERROR in e_c3d: nonpositive jacobian determinant in element
func isBadTet10(coords [10]v3.Vec) (bool, float64) {
	// Gause points are according to CCX source code.
	// Reference:
	// ccx_2.20/src/gauss.f
	var gaussPoints [4]v3.Vec
	gaussPoints[0] = v3.Vec{0.138196601125011, 0.138196601125011, 0.138196601125011}
	gaussPoints[1] = v3.Vec{0.585410196624968, 0.138196601125011, 0.138196601125011}
	gaussPoints[2] = v3.Vec{0.138196601125011, 0.585410196624968, 0.138196601125011}
	gaussPoints[3] = v3.Vec{0.138196601125011, 0.138196601125011, 0.585410196624968}

	var bad bool
	var jacobianDeterminant float64

	for i := 0; i < 4; i++ {
		bad, jacobianDeterminant = isBadGaussTet10(coords, gaussPoints[i].X, gaussPoints[i].Y, gaussPoints[i].Z)
		if bad {
			return true, jacobianDeterminant
		}
	}

	return false, jacobianDeterminant
}

//-----------------------------------------------------------------------------

// If triangles are degenerate, then tetrahedra will be bad.
// To filter bad tetrahedra.
func degenerateTriangles(a, b, c, d v3.Vec) bool {
	// 4 triangles are possible.
	// Each triangle is a tetrahedron side.
	t := sdf.Triangle3{}
	t[0] = a
	t[1] = b
	t[2] = c
	// Use the epsilon value of `vertexbuffer.go`
	if t.Degenerate(0.0001) {
		return true
	}
	t[0] = a
	t[1] = b
	t[2] = d
	// Use the epsilon value of `vertexbuffer.go`
	if t.Degenerate(0.0001) {
		return true
	}
	t[0] = a
	t[1] = c
	t[2] = d
	// Use the epsilon value of `vertexbuffer.go`
	if t.Degenerate(0.0001) {
		return true
	}
	t[0] = b
	t[1] = c
	t[2] = d
	// Use the epsilon value of `vertexbuffer.go`
	return t.Degenerate(0.0001)
}

//-----------------------------------------------------------------------------
