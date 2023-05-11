//-----------------------------------------------------------------------------
/*

Matrix Operations

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// M44 is a 4x4 matrix.
type M44 [16]float64

// M33 is a 3x3 matrix.
type M33 [9]float64

// M22 is a 2x2 matrix.
type M22 [4]float64

//-----------------------------------------------------------------------------

// RandomM22 returns a 2x2 matrix with random elements.
func RandomM22(a, b float64) M22 {
	return M22{randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b)}
}

// RandomM33 returns a 3x3 matrix with random elements.
func RandomM33(a, b float64) M33 {
	return M33{randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b)}
}

// RandomM44 returns a 4x4 matrix with random elements.
func RandomM44(a, b float64) M44 {
	return M44{
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b)}
}

//-----------------------------------------------------------------------------

// Identity3d returns a 4x4 identity matrix.
func Identity3d() M44 {
	return M44{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

// Identity2d returns a 3x3 identity matrix.
func Identity2d() M33 {
	return M33{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}
}

// Identity returns a 2x2 identity matrix.
func Identity() M22 {
	return M22{
		1, 0,
		0, 1}
}

// Translate3d returns a 4x4 translation matrix.
func Translate3d(v v3.Vec) M44 {
	return M44{
		1, 0, 0, v.X,
		0, 1, 0, v.Y,
		0, 0, 1, v.Z,
		0, 0, 0, 1}
}

// Translate2d returns a 3x3 translation matrix.
func Translate2d(v v2.Vec) M33 {
	return M33{
		1, 0, v.X,
		0, 1, v.Y,
		0, 0, 1}
}

// Scale3d returns a 4x4 scaling matrix.
// Scaling does not preserve distance. See: ScaleUniform3D()
func Scale3d(v v3.Vec) M44 {
	return M44{
		v.X, 0, 0, 0,
		0, v.Y, 0, 0,
		0, 0, v.Z, 0,
		0, 0, 0, 1}
}

// Scale2d returns a 3x3 scaling matrix.
// Scaling does not preserve distance. See: ScaleUniform2D().
func Scale2d(v v2.Vec) M33 {
	return M33{
		v.X, 0, 0,
		0, v.Y, 0,
		0, 0, 1}
}

// Rotate3d returns an orthographic 4x4 rotation matrix (right hand rule).
func Rotate3d(v v3.Vec, a float64) M44 {
	v = v.Normalize()
	s := math.Sin(a)
	c := math.Cos(a)
	m := 1 - c
	return M44{
		m*v.X*v.X + c, m*v.X*v.Y - v.Z*s, m*v.Z*v.X + v.Y*s, 0,
		m*v.X*v.Y + v.Z*s, m*v.Y*v.Y + c, m*v.Y*v.Z - v.X*s, 0,
		m*v.Z*v.X - v.Y*s, m*v.Y*v.Z + v.X*s, m*v.Z*v.Z + c, 0,
		0, 0, 0, 1}
}

// RotateX returns a 4x4 matrix with rotation about the X axis.
func RotateX(a float64) M44 {
	return Rotate3d(v3.Vec{1, 0, 0}, a)
}

// RotateY returns a 4x4 matrix with rotation about the Y axis.
func RotateY(a float64) M44 {
	return Rotate3d(v3.Vec{0, 1, 0}, a)
}

// RotateZ returns a 4x4 matrix with rotation about the Z axis.
func RotateZ(a float64) M44 {
	return Rotate3d(v3.Vec{0, 0, 1}, a)
}

// MirrorXY returns a 4x4 matrix with mirroring across the XY plane.
func MirrorXY() M44 {
	return M44{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, -1, 0,
		0, 0, 0, 1}
}

// MirrorXZ returns a 4x4 matrix with mirroring across the XZ plane.
func MirrorXZ() M44 {
	return M44{
		1, 0, 0, 0,
		0, -1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

// MirrorYZ returns a 4x4 matrix with mirroring across the YZ plane.
func MirrorYZ() M44 {
	return M44{
		-1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

// MirrorXeqY returns a 4x4 matrix with mirroring across the X == Y plane.
func MirrorXeqY() M44 {
	return M44{
		0, 1, 0, 0,
		1, 0, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

// MirrorX returns a 3x3 matrix with mirroring across the X axis.
func MirrorX() M33 {
	return M33{
		1, 0, 0,
		0, -1, 0,
		0, 0, 1}
}

// MirrorY returns a 3x3 matrix with mirroring across the Y axis.
func MirrorY() M33 {
	return M33{
		-1, 0, 0,
		0, 1, 0,
		0, 0, 1}
}

// Rotate2d returns an orthographic 3x3 rotation matrix (right hand rule).
func Rotate2d(a float64) M33 {
	s := math.Sin(a)
	c := math.Cos(a)
	return M33{
		c, -s, 0,
		s, c, 0,
		0, 0, 1}
}

// Rotate returns an orthographic 2x2 rotation matrix (right hand rule).
func Rotate(a float64) M22 {
	s := math.Sin(a)
	c := math.Cos(a)
	return M22{
		c, -s,
		s, c,
	}
}

// RotateToVector returns the rotation matrix that transforms a onto the same direction as b.
func RotateToVector(a, b v3.Vec) M44 {
	// is either vector == 0?
	if a.Equals(v3.Vec{}, epsilon) || b.Equals(v3.Vec{}, epsilon) {
		return Identity3d()
	}
	// normalize both vectors
	a = a.Normalize()
	b = b.Normalize()
	// are the vectors the same?
	if a.Equals(b, epsilon) {
		return Identity3d()
	}
	// are the vectors opposite (180 degrees apart)?
	if a.Neg().Equals(b, epsilon) {
		return M44{
			-1, 0, 0, 0,
			0, -1, 0, 0,
			0, 0, -1, 0,
			0, 0, 0, 1}
	}
	// general case
	// See:	https://math.stackexchange.com/questions/180418/calculate-rotation-matrix-to-align-vector-a-to-vector-b-in-3d
	v := a.Cross(b)
	k := 1 / (1 + a.Dot(b))
	vx := M33{0, -v.Z, v.Y, v.Z, 0, -v.X, -v.Y, v.X, 0}
	r := Identity2d().Add(vx).Add(vx.Mul(vx).MulScalar(k))
	return M44{
		r[0], r[1], r[2], 0,
		r[3], r[4], r[5], 0,
		r[6], r[7], r[8], 0,
		0, 0, 0, 1,
	}
}

//-----------------------------------------------------------------------------

// Equals tests the equality of 4x4 matrices.
func (a M44) Equals(b M44, tolerance float64) bool {
	return (math.Abs(a[0]-b[0]) < tolerance &&
		math.Abs(a[1]-b[1]) < tolerance &&
		math.Abs(a[2]-b[2]) < tolerance &&
		math.Abs(a[3]-b[3]) < tolerance &&
		math.Abs(a[4]-b[4]) < tolerance &&
		math.Abs(a[5]-b[5]) < tolerance &&
		math.Abs(a[6]-b[6]) < tolerance &&
		math.Abs(a[7]-b[7]) < tolerance &&
		math.Abs(a[8]-b[8]) < tolerance &&
		math.Abs(a[9]-b[9]) < tolerance &&
		math.Abs(a[10]-b[10]) < tolerance &&
		math.Abs(a[11]-b[11]) < tolerance &&
		math.Abs(a[12]-b[12]) < tolerance &&
		math.Abs(a[13]-b[13]) < tolerance &&
		math.Abs(a[14]-b[14]) < tolerance &&
		math.Abs(a[15]-b[15]) < tolerance)
}

// Equals tests the equality of 3x3 matrices.
func (a M33) Equals(b M33, tolerance float64) bool {
	return (math.Abs(a[0]-b[0]) < tolerance &&
		math.Abs(a[1]-b[1]) < tolerance &&
		math.Abs(a[2]-b[2]) < tolerance &&
		math.Abs(a[3]-b[3]) < tolerance &&
		math.Abs(a[4]-b[4]) < tolerance &&
		math.Abs(a[5]-b[5]) < tolerance &&
		math.Abs(a[6]-b[6]) < tolerance &&
		math.Abs(a[7]-b[7]) < tolerance &&
		math.Abs(a[8]-b[8]) < tolerance)
}

// Equals tests the equality of 2x2 matrices.
func (a M22) Equals(b M22, tolerance float64) bool {
	return (math.Abs(a[0]-b[0]) < tolerance &&
		math.Abs(a[1]-b[1]) < tolerance &&
		math.Abs(a[2]-b[2]) < tolerance &&
		math.Abs(a[3]-b[3]) < tolerance)
}

//-----------------------------------------------------------------------------

// MulPosition multiplies a v3.Vec position with a rotate/translate matrix.
func (a M44) MulPosition(b v3.Vec) v3.Vec {
	return v3.Vec{a[0]*b.X + a[1]*b.Y + a[2]*b.Z + a[3],
		a[4]*b.X + a[5]*b.Y + a[6]*b.Z + a[7],
		a[8]*b.X + a[9]*b.Y + a[10]*b.Z + a[11]}
}

// MulPosition multiplies a v2.Vec position with a rotate/translate matrix.
func (a M33) MulPosition(b v2.Vec) v2.Vec {
	return v2.Vec{a[0]*b.X + a[1]*b.Y + a[2], a[3]*b.X + a[4]*b.Y + a[5]}
}

// MulPosition multiplies a v2.Vec position with a rotate matrix.
func (a M22) MulPosition(b v2.Vec) v2.Vec {
	return v2.Vec{a[0]*b.X + a[1]*b.Y, a[2]*b.X + a[3]*b.Y}
}

//-----------------------------------------------------------------------------

// mulVertices2 multiples a set of v2.Vec vertices by a rotate/translate matrix.
func mulVertices2(v v2.VecSet, a M33) {
	for i := range v {
		v[i] = a.MulPosition(v[i])
	}
}

// mulVertices3 multiples a set of v3.Vec vertices by a rotate/translate matrix.
func mulVertices3(v v3.VecSet, a M44) {
	for i := range v {
		v[i] = a.MulPosition(v[i])
	}
}

//-----------------------------------------------------------------------------

// Mul multiplies 4x4 matrices.
func (a M44) Mul(b M44) M44 {
	return M44{
		a[0]*b[0] + a[1]*b[4] + a[2]*b[8] + a[3]*b[12],
		a[0]*b[1] + a[1]*b[5] + a[2]*b[9] + a[3]*b[13],
		a[0]*b[2] + a[1]*b[6] + a[2]*b[10] + a[3]*b[14],
		a[0]*b[3] + a[1]*b[7] + a[2]*b[11] + a[3]*b[15],
		a[4]*b[0] + a[5]*b[4] + a[6]*b[8] + a[7]*b[12],
		a[4]*b[1] + a[5]*b[5] + a[6]*b[9] + a[7]*b[13],
		a[4]*b[2] + a[5]*b[6] + a[6]*b[10] + a[7]*b[14],
		a[4]*b[3] + a[5]*b[7] + a[6]*b[11] + a[7]*b[15],
		a[8]*b[0] + a[9]*b[4] + a[10]*b[8] + a[11]*b[12],
		a[8]*b[1] + a[9]*b[5] + a[10]*b[9] + a[11]*b[13],
		a[8]*b[2] + a[9]*b[6] + a[10]*b[10] + a[11]*b[14],
		a[8]*b[3] + a[9]*b[7] + a[10]*b[11] + a[11]*b[15],
		a[12]*b[0] + a[13]*b[4] + a[14]*b[8] + a[15]*b[12],
		a[12]*b[1] + a[13]*b[5] + a[14]*b[9] + a[15]*b[13],
		a[12]*b[2] + a[13]*b[6] + a[14]*b[10] + a[15]*b[14],
		a[12]*b[3] + a[13]*b[7] + a[14]*b[11] + a[15]*b[15],
	}
}

// Mul multiplies 3x3 matrices.
func (a M33) Mul(b M33) M33 {
	return M33{
		a[0]*b[0] + a[1]*b[3] + a[2]*b[6],
		a[0]*b[1] + a[1]*b[4] + a[2]*b[7],
		a[0]*b[2] + a[1]*b[5] + a[2]*b[8],
		a[3]*b[0] + a[4]*b[3] + a[5]*b[6],
		a[3]*b[1] + a[4]*b[4] + a[5]*b[7],
		a[3]*b[2] + a[4]*b[5] + a[5]*b[8],
		a[6]*b[0] + a[7]*b[3] + a[8]*b[6],
		a[6]*b[1] + a[7]*b[4] + a[8]*b[7],
		a[6]*b[2] + a[7]*b[5] + a[8]*b[8],
	}
}

// Mul multiplies 2x2 matrices.
func (a M22) Mul(b M22) M22 {
	return M22{
		a[0]*b[0] + a[1]*b[2],
		a[0]*b[1] + a[1]*b[3],
		a[2]*b[0] + a[3]*b[2],
		a[2]*b[1] + a[3]*b[3],
	}
}

//-----------------------------------------------------------------------------

// Add two 3x3 matrices.
func (a M33) Add(b M33) M33 {
	return M33{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
		a[3] + b[3],
		a[4] + b[4],
		a[5] + b[5],
		a[6] + b[6],
		a[7] + b[7],
		a[8] + b[8],
	}
}

//-----------------------------------------------------------------------------

// MulScalar multiplies each 3x3 matrix component by a scalar.
func (a M33) MulScalar(k float64) M33 {
	return M33{
		k * a[0], k * a[1], k * a[2],
		k * a[3], k * a[4], k * a[5],
		k * a[6], k * a[7], k * a[8],
	}
}

//-----------------------------------------------------------------------------
// Transform bounding boxes - keep them axis aligned
// http://dev.theomader.com/transform-bounding-boxes/

// MulBox rotates/translates a 3d bounding box and resizes for axis-alignment.
func (a M44) MulBox(box Box3) Box3 {
	r := v3.Vec{a[0], a[4], a[8]}
	u := v3.Vec{a[1], a[5], a[9]}
	b := v3.Vec{a[2], a[6], a[10]}
	t := v3.Vec{a[3], a[7], a[11]}
	xa := r.MulScalar(box.Min.X)
	xb := r.MulScalar(box.Max.X)
	ya := u.MulScalar(box.Min.Y)
	yb := u.MulScalar(box.Max.Y)
	za := b.MulScalar(box.Min.Z)
	zb := b.MulScalar(box.Max.Z)
	xa, xb = xa.Min(xb), xa.Max(xb)
	ya, yb = ya.Min(yb), ya.Max(yb)
	za, zb = za.Min(zb), za.Max(zb)
	min := xa.Add(ya).Add(za).Add(t)
	max := xb.Add(yb).Add(zb).Add(t)
	return Box3{min, max}
}

// MulBox rotates/translates a 2d bounding box and resizes for axis-alignment.
func (a M33) MulBox(box Box2) Box2 {
	r := v2.Vec{a[0], a[3]}
	u := v2.Vec{a[1], a[4]}
	t := v2.Vec{a[2], a[5]}
	xa := r.MulScalar(box.Min.X)
	xb := r.MulScalar(box.Max.X)
	ya := u.MulScalar(box.Min.Y)
	yb := u.MulScalar(box.Max.Y)
	xa, xb = xa.Min(xb), xa.Max(xb)
	ya, yb = ya.Min(yb), ya.Max(yb)
	min := xa.Add(ya).Add(t)
	max := xb.Add(yb).Add(t)
	return Box2{min, max}
}

//-----------------------------------------------------------------------------

// Determinant returns the determinant of a 4x4 matrix.
func (a M44) Determinant() float64 {
	return (a[0]*a[5]*a[10]*a[15] - a[0]*a[5]*a[11]*a[14] +
		a[0]*a[6]*a[11]*a[13] - a[0]*a[6]*a[9]*a[15] +
		a[0]*a[7]*a[9]*a[14] - a[0]*a[7]*a[10]*a[13] -
		a[1]*a[6]*a[11]*a[12] + a[1]*a[6]*a[8]*a[15] -
		a[1]*a[7]*a[8]*a[14] + a[1]*a[7]*a[10]*a[12] -
		a[1]*a[4]*a[10]*a[15] + a[1]*a[4]*a[11]*a[14] +
		a[2]*a[7]*a[8]*a[13] - a[2]*a[7]*a[9]*a[12] +
		a[2]*a[4]*a[9]*a[15] - a[2]*a[4]*a[11]*a[13] +
		a[2]*a[5]*a[11]*a[12] - a[2]*a[5]*a[8]*a[15] -
		a[3]*a[4]*a[9]*a[14] + a[3]*a[4]*a[10]*a[13] -
		a[3]*a[5]*a[10]*a[12] + a[3]*a[5]*a[8]*a[14] -
		a[3]*a[6]*a[8]*a[13] + a[3]*a[6]*a[9]*a[12])
}

// Determinant returns the determinant of a 3x3 matrix.
func (a M33) Determinant() float64 {
	return (a[0]*(a[4]*a[8]-a[7]*a[5]) -
		a[1]*(a[3]*a[8]-a[6]*a[5]) +
		a[2]*(a[3]*a[7]-a[6]*a[4]))
}

// Determinant returns the determinant of a 2x2 matrix.
func (a M22) Determinant() float64 {
	return a[0]*a[3] - a[1]*a[2]
}

//-----------------------------------------------------------------------------

// Inverse returns the inverse of a 4x4 matrix.
func (a M44) Inverse() M44 {
	d := 1 / a.Determinant()
	return M44{
		(a[6]*a[11]*a[13] - a[7]*a[10]*a[13] + a[7]*a[9]*a[14] - a[5]*a[11]*a[14] - a[6]*a[9]*a[15] + a[5]*a[10]*a[15]) * d,
		(a[3]*a[10]*a[13] - a[2]*a[11]*a[13] - a[3]*a[9]*a[14] + a[1]*a[11]*a[14] + a[2]*a[9]*a[15] - a[1]*a[10]*a[15]) * d,
		(a[2]*a[7]*a[13] - a[3]*a[6]*a[13] + a[3]*a[5]*a[14] - a[1]*a[7]*a[14] - a[2]*a[5]*a[15] + a[1]*a[6]*a[15]) * d,
		(a[3]*a[6]*a[9] - a[2]*a[7]*a[9] - a[3]*a[5]*a[10] + a[1]*a[7]*a[10] + a[2]*a[5]*a[11] - a[1]*a[6]*a[11]) * d,
		(a[7]*a[10]*a[12] - a[6]*a[11]*a[12] - a[7]*a[8]*a[14] + a[4]*a[11]*a[14] + a[6]*a[8]*a[15] - a[4]*a[10]*a[15]) * d,
		(a[2]*a[11]*a[12] - a[3]*a[10]*a[12] + a[3]*a[8]*a[14] - a[0]*a[11]*a[14] - a[2]*a[8]*a[15] + a[0]*a[10]*a[15]) * d,
		(a[3]*a[6]*a[12] - a[2]*a[7]*a[12] - a[3]*a[4]*a[14] + a[0]*a[7]*a[14] + a[2]*a[4]*a[15] - a[0]*a[6]*a[15]) * d,
		(a[2]*a[7]*a[8] - a[3]*a[6]*a[8] + a[3]*a[4]*a[10] - a[0]*a[7]*a[10] - a[2]*a[4]*a[11] + a[0]*a[6]*a[11]) * d,
		(a[5]*a[11]*a[12] - a[7]*a[9]*a[12] + a[7]*a[8]*a[13] - a[4]*a[11]*a[13] - a[5]*a[8]*a[15] + a[4]*a[9]*a[15]) * d,
		(a[3]*a[9]*a[12] - a[1]*a[11]*a[12] - a[3]*a[8]*a[13] + a[0]*a[11]*a[13] + a[1]*a[8]*a[15] - a[0]*a[9]*a[15]) * d,
		(a[1]*a[7]*a[12] - a[3]*a[5]*a[12] + a[3]*a[4]*a[13] - a[0]*a[7]*a[13] - a[1]*a[4]*a[15] + a[0]*a[5]*a[15]) * d,
		(a[3]*a[5]*a[8] - a[1]*a[7]*a[8] - a[3]*a[4]*a[9] + a[0]*a[7]*a[9] + a[1]*a[4]*a[11] - a[0]*a[5]*a[11]) * d,
		(a[6]*a[9]*a[12] - a[5]*a[10]*a[12] - a[6]*a[8]*a[13] + a[4]*a[10]*a[13] + a[5]*a[8]*a[14] - a[4]*a[9]*a[14]) * d,
		(a[1]*a[10]*a[12] - a[2]*a[9]*a[12] + a[2]*a[8]*a[13] - a[0]*a[10]*a[13] - a[1]*a[8]*a[14] + a[0]*a[9]*a[14]) * d,
		(a[2]*a[5]*a[12] - a[1]*a[6]*a[12] - a[2]*a[4]*a[13] + a[0]*a[6]*a[13] + a[1]*a[4]*a[14] - a[0]*a[5]*a[14]) * d,
		(a[1]*a[6]*a[8] - a[2]*a[5]*a[8] + a[2]*a[4]*a[9] - a[0]*a[6]*a[9] - a[1]*a[4]*a[10] + a[0]*a[5]*a[10]) * d,
	}
}

// Inverse returns the inverse of a 3x3 matrix.
func (a M33) Inverse() M33 {
	d := 1 / a.Determinant()
	return M33{
		(a[4]*a[8] - a[5]*a[7]) * d,
		(a[7]*a[2] - a[1]*a[8]) * d,
		(a[1]*a[5] - a[4]*a[2]) * d,
		(a[5]*a[6] - a[8]*a[3]) * d,
		(a[8]*a[0] - a[6]*a[2]) * d,
		(a[2]*a[3] - a[5]*a[0]) * d,
		(a[3]*a[7] - a[6]*a[4]) * d,
		(a[6]*a[1] - a[0]*a[7]) * d,
		(a[0]*a[4] - a[1]*a[3]) * d,
	}
}

// Inverse returns the inverse of a 2x2 matrix.
func (a M22) Inverse() M22 {
	d := 1 / a.Determinant()
	return M22{
		a[3] * d,
		-a[1] * d,
		-a[2] * d,
		a[0] * d,
	}
}

//-----------------------------------------------------------------------------

// NewM44 returns a new matrix. Input is in row-major order.
func NewM44(x [16]float64) M44 {
	return M44{
		x[0], x[1], x[2], x[3],
		x[4], x[5], x[6], x[7],
		x[8], x[9], x[10], x[11],
		x[12], x[13], x[14], x[15],
	}
}

// Values returns the matrix values in row-major order.
func (a M44) Values() [16]float64 {
	return [16]float64{
		a[0], a[1], a[2], a[3],
		a[4], a[5], a[6], a[7],
		a[8], a[9], a[10], a[11],
		a[12], a[13], a[14], a[15],
	}
}

//-----------------------------------------------------------------------------

// NewM33 returns a new matrix. Input is in row-major order.
func NewM33(x [9]float64) M33 {
	return M33{
		x[0], x[1], x[2],
		x[3], x[4], x[5],
		x[6], x[7], x[8],
	}
}

// Values returns the matrix values in row-major order.
func (a M33) Values() [9]float64 {
	return [9]float64{
		a[0], a[1], a[2],
		a[3], a[4], a[5],
		a[6], a[7], a[8],
	}
}

//-----------------------------------------------------------------------------

// NewM22 returns a new matrix. Input is in row-major order.
func NewM22(x [4]float64) M22 {
	return M22{
		x[0], x[1],
		x[2], x[3]}
}

// Values returns the matrix values in row-major order.
func (a M22) Values() [4]float64 {
	return [4]float64{
		a[0], a[1],
		a[2], a[3]}
}

//-----------------------------------------------------------------------------
