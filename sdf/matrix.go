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
type M44 struct {
	x00, x01, x02, x03 float64
	x10, x11, x12, x13 float64
	x20, x21, x22, x23 float64
	x30, x31, x32, x33 float64
}

// M33 is a 3x3 matrix.
type M33 struct {
	x00, x01, x02 float64
	x10, x11, x12 float64
	x20, x21, x22 float64
}

// M22 is a 2x2 matrix.
type M22 struct {
	x00, x01 float64
	x10, x11 float64
}

//-----------------------------------------------------------------------------

// RandomM22 returns a 2x2 matrix with random elements.
func RandomM22(a, b float64) M22 {
	m := M22{randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b)}
	return m
}

// RandomM33 returns a 3x3 matrix with random elements.
func RandomM33(a, b float64) M33 {
	m := M33{randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b),
		randomRange(a, b)}
	return m
}

// RandomM44 returns a 4x4 matrix with random elements.
func RandomM44(a, b float64) M44 {
	m := M44{
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
	return m
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
		r.x00, r.x01, r.x02, 0,
		r.x10, r.x11, r.x12, 0,
		r.x20, r.x21, r.x22, 0,
		0, 0, 0, 1,
	}
}

//-----------------------------------------------------------------------------

// Equals tests the equality of 4x4 matrices.
func (a M44) Equals(b M44, tolerance float64) bool {
	return (math.Abs(a.x00-b.x00) < tolerance &&
		math.Abs(a.x01-b.x01) < tolerance &&
		math.Abs(a.x02-b.x02) < tolerance &&
		math.Abs(a.x03-b.x03) < tolerance &&
		math.Abs(a.x10-b.x10) < tolerance &&
		math.Abs(a.x11-b.x11) < tolerance &&
		math.Abs(a.x12-b.x12) < tolerance &&
		math.Abs(a.x13-b.x13) < tolerance &&
		math.Abs(a.x20-b.x20) < tolerance &&
		math.Abs(a.x21-b.x21) < tolerance &&
		math.Abs(a.x22-b.x22) < tolerance &&
		math.Abs(a.x23-b.x23) < tolerance &&
		math.Abs(a.x30-b.x30) < tolerance &&
		math.Abs(a.x31-b.x31) < tolerance &&
		math.Abs(a.x32-b.x32) < tolerance &&
		math.Abs(a.x33-b.x33) < tolerance)
}

// Equals tests the equality of 3x3 matrices.
func (a M33) Equals(b M33, tolerance float64) bool {
	return (math.Abs(a.x00-b.x00) < tolerance &&
		math.Abs(a.x01-b.x01) < tolerance &&
		math.Abs(a.x02-b.x02) < tolerance &&
		math.Abs(a.x10-b.x10) < tolerance &&
		math.Abs(a.x11-b.x11) < tolerance &&
		math.Abs(a.x12-b.x12) < tolerance &&
		math.Abs(a.x20-b.x20) < tolerance &&
		math.Abs(a.x21-b.x21) < tolerance &&
		math.Abs(a.x22-b.x22) < tolerance)
}

// Equals tests the equality of 2x2 matrices.
func (a M22) Equals(b M22, tolerance float64) bool {
	return (math.Abs(a.x00-b.x00) < tolerance &&
		math.Abs(a.x01-b.x01) < tolerance &&
		math.Abs(a.x10-b.x10) < tolerance &&
		math.Abs(a.x11-b.x11) < tolerance)
}

//-----------------------------------------------------------------------------

// MulPosition multiplies a v3.Vec position with a rotate/translate matrix.
func (a M44) MulPosition(b v3.Vec) v3.Vec {
	return v3.Vec{a.x00*b.X + a.x01*b.Y + a.x02*b.Z + a.x03,
		a.x10*b.X + a.x11*b.Y + a.x12*b.Z + a.x13,
		a.x20*b.X + a.x21*b.Y + a.x22*b.Z + a.x23}
}

// MulPosition multiplies a v2.Vec position with a rotate/translate matrix.
func (a M33) MulPosition(b v2.Vec) v2.Vec {
	return v2.Vec{a.x00*b.X + a.x01*b.Y + a.x02,
		a.x10*b.X + a.x11*b.Y + a.x12}
}

// MulPosition multiplies a v2.Vec position with a rotate matrix.
func (a M22) MulPosition(b v2.Vec) v2.Vec {
	return v2.Vec{a.x00*b.X + a.x01*b.Y,
		a.x10*b.X + a.x11*b.Y}
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
	m := M44{}
	m.x00 = a.x00*b.x00 + a.x01*b.x10 + a.x02*b.x20 + a.x03*b.x30
	m.x10 = a.x10*b.x00 + a.x11*b.x10 + a.x12*b.x20 + a.x13*b.x30
	m.x20 = a.x20*b.x00 + a.x21*b.x10 + a.x22*b.x20 + a.x23*b.x30
	m.x30 = a.x30*b.x00 + a.x31*b.x10 + a.x32*b.x20 + a.x33*b.x30
	m.x01 = a.x00*b.x01 + a.x01*b.x11 + a.x02*b.x21 + a.x03*b.x31
	m.x11 = a.x10*b.x01 + a.x11*b.x11 + a.x12*b.x21 + a.x13*b.x31
	m.x21 = a.x20*b.x01 + a.x21*b.x11 + a.x22*b.x21 + a.x23*b.x31
	m.x31 = a.x30*b.x01 + a.x31*b.x11 + a.x32*b.x21 + a.x33*b.x31
	m.x02 = a.x00*b.x02 + a.x01*b.x12 + a.x02*b.x22 + a.x03*b.x32
	m.x12 = a.x10*b.x02 + a.x11*b.x12 + a.x12*b.x22 + a.x13*b.x32
	m.x22 = a.x20*b.x02 + a.x21*b.x12 + a.x22*b.x22 + a.x23*b.x32
	m.x32 = a.x30*b.x02 + a.x31*b.x12 + a.x32*b.x22 + a.x33*b.x32
	m.x03 = a.x00*b.x03 + a.x01*b.x13 + a.x02*b.x23 + a.x03*b.x33
	m.x13 = a.x10*b.x03 + a.x11*b.x13 + a.x12*b.x23 + a.x13*b.x33
	m.x23 = a.x20*b.x03 + a.x21*b.x13 + a.x22*b.x23 + a.x23*b.x33
	m.x33 = a.x30*b.x03 + a.x31*b.x13 + a.x32*b.x23 + a.x33*b.x33
	return m
}

// Mul multiplies 3x3 matrices.
func (a M33) Mul(b M33) M33 {
	m := M33{}
	m.x00 = a.x00*b.x00 + a.x01*b.x10 + a.x02*b.x20
	m.x10 = a.x10*b.x00 + a.x11*b.x10 + a.x12*b.x20
	m.x20 = a.x20*b.x00 + a.x21*b.x10 + a.x22*b.x20
	m.x01 = a.x00*b.x01 + a.x01*b.x11 + a.x02*b.x21
	m.x11 = a.x10*b.x01 + a.x11*b.x11 + a.x12*b.x21
	m.x21 = a.x20*b.x01 + a.x21*b.x11 + a.x22*b.x21
	m.x02 = a.x00*b.x02 + a.x01*b.x12 + a.x02*b.x22
	m.x12 = a.x10*b.x02 + a.x11*b.x12 + a.x12*b.x22
	m.x22 = a.x20*b.x02 + a.x21*b.x12 + a.x22*b.x22
	return m
}

// Mul multiplies 2x2 matrices.
func (a M22) Mul(b M22) M22 {
	m := M22{}
	m.x00 = a.x00*b.x00 + a.x01*b.x10
	m.x01 = a.x00*b.x01 + a.x01*b.x11
	m.x10 = a.x10*b.x00 + a.x11*b.x10
	m.x11 = a.x10*b.x01 + a.x11*b.x11
	return m
}

//-----------------------------------------------------------------------------

// Add two 3x3 matrices.
func (a M33) Add(b M33) M33 {
	return M33{
		x00: a.x00 + b.x00,
		x10: a.x10 + b.x10,
		x20: a.x20 + b.x20,
		x01: a.x01 + b.x01,
		x11: a.x11 + b.x11,
		x21: a.x21 + b.x21,
		x02: a.x02 + b.x02,
		x12: a.x12 + b.x12,
		x22: a.x22 + b.x22,
	}
}

//-----------------------------------------------------------------------------

// MulScalar multiplies each 3x3 matrix component by a scalar.
func (a M33) MulScalar(k float64) M33 {
	return M33{
		x00: k * a.x00,
		x10: k * a.x10,
		x20: k * a.x20,
		x01: k * a.x01,
		x11: k * a.x11,
		x21: k * a.x21,
		x02: k * a.x02,
		x12: k * a.x12,
		x22: k * a.x22,
	}
}

//-----------------------------------------------------------------------------
// Transform bounding boxes - keep them axis aligned
// http://dev.theomader.com/transform-bounding-boxes/

// MulBox rotates/translates a 3d bounding box and resizes for axis-alignment.
func (a M44) MulBox(box Box3) Box3 {
	r := v3.Vec{a.x00, a.x10, a.x20}
	u := v3.Vec{a.x01, a.x11, a.x21}
	b := v3.Vec{a.x02, a.x12, a.x22}
	t := v3.Vec{a.x03, a.x13, a.x23}
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
	r := v2.Vec{a.x00, a.x10}
	u := v2.Vec{a.x01, a.x11}
	t := v2.Vec{a.x02, a.x12}
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
	return (a.x00*a.x11*a.x22*a.x33 - a.x00*a.x11*a.x23*a.x32 +
		a.x00*a.x12*a.x23*a.x31 - a.x00*a.x12*a.x21*a.x33 +
		a.x00*a.x13*a.x21*a.x32 - a.x00*a.x13*a.x22*a.x31 -
		a.x01*a.x12*a.x23*a.x30 + a.x01*a.x12*a.x20*a.x33 -
		a.x01*a.x13*a.x20*a.x32 + a.x01*a.x13*a.x22*a.x30 -
		a.x01*a.x10*a.x22*a.x33 + a.x01*a.x10*a.x23*a.x32 +
		a.x02*a.x13*a.x20*a.x31 - a.x02*a.x13*a.x21*a.x30 +
		a.x02*a.x10*a.x21*a.x33 - a.x02*a.x10*a.x23*a.x31 +
		a.x02*a.x11*a.x23*a.x30 - a.x02*a.x11*a.x20*a.x33 -
		a.x03*a.x10*a.x21*a.x32 + a.x03*a.x10*a.x22*a.x31 -
		a.x03*a.x11*a.x22*a.x30 + a.x03*a.x11*a.x20*a.x32 -
		a.x03*a.x12*a.x20*a.x31 + a.x03*a.x12*a.x21*a.x30)
}

// Determinant returns the determinant of a 3x3 matrix.
func (a M33) Determinant() float64 {
	return (a.x00*(a.x11*a.x22-a.x21*a.x12) -
		a.x01*(a.x10*a.x22-a.x20*a.x12) +
		a.x02*(a.x10*a.x21-a.x20*a.x11))
}

// Determinant returns the determinant of a 2x2 matrix.
func (a M22) Determinant() float64 {
	return a.x00*a.x11 - a.x01*a.x10
}

//-----------------------------------------------------------------------------

// Inverse returns the inverse of a 4x4 matrix.
func (a M44) Inverse() M44 {
	m := M44{}
	d := 1 / a.Determinant()
	m.x00 = (a.x12*a.x23*a.x31 - a.x13*a.x22*a.x31 + a.x13*a.x21*a.x32 - a.x11*a.x23*a.x32 - a.x12*a.x21*a.x33 + a.x11*a.x22*a.x33) * d
	m.x01 = (a.x03*a.x22*a.x31 - a.x02*a.x23*a.x31 - a.x03*a.x21*a.x32 + a.x01*a.x23*a.x32 + a.x02*a.x21*a.x33 - a.x01*a.x22*a.x33) * d
	m.x02 = (a.x02*a.x13*a.x31 - a.x03*a.x12*a.x31 + a.x03*a.x11*a.x32 - a.x01*a.x13*a.x32 - a.x02*a.x11*a.x33 + a.x01*a.x12*a.x33) * d
	m.x03 = (a.x03*a.x12*a.x21 - a.x02*a.x13*a.x21 - a.x03*a.x11*a.x22 + a.x01*a.x13*a.x22 + a.x02*a.x11*a.x23 - a.x01*a.x12*a.x23) * d
	m.x10 = (a.x13*a.x22*a.x30 - a.x12*a.x23*a.x30 - a.x13*a.x20*a.x32 + a.x10*a.x23*a.x32 + a.x12*a.x20*a.x33 - a.x10*a.x22*a.x33) * d
	m.x11 = (a.x02*a.x23*a.x30 - a.x03*a.x22*a.x30 + a.x03*a.x20*a.x32 - a.x00*a.x23*a.x32 - a.x02*a.x20*a.x33 + a.x00*a.x22*a.x33) * d
	m.x12 = (a.x03*a.x12*a.x30 - a.x02*a.x13*a.x30 - a.x03*a.x10*a.x32 + a.x00*a.x13*a.x32 + a.x02*a.x10*a.x33 - a.x00*a.x12*a.x33) * d
	m.x13 = (a.x02*a.x13*a.x20 - a.x03*a.x12*a.x20 + a.x03*a.x10*a.x22 - a.x00*a.x13*a.x22 - a.x02*a.x10*a.x23 + a.x00*a.x12*a.x23) * d
	m.x20 = (a.x11*a.x23*a.x30 - a.x13*a.x21*a.x30 + a.x13*a.x20*a.x31 - a.x10*a.x23*a.x31 - a.x11*a.x20*a.x33 + a.x10*a.x21*a.x33) * d
	m.x21 = (a.x03*a.x21*a.x30 - a.x01*a.x23*a.x30 - a.x03*a.x20*a.x31 + a.x00*a.x23*a.x31 + a.x01*a.x20*a.x33 - a.x00*a.x21*a.x33) * d
	m.x22 = (a.x01*a.x13*a.x30 - a.x03*a.x11*a.x30 + a.x03*a.x10*a.x31 - a.x00*a.x13*a.x31 - a.x01*a.x10*a.x33 + a.x00*a.x11*a.x33) * d
	m.x23 = (a.x03*a.x11*a.x20 - a.x01*a.x13*a.x20 - a.x03*a.x10*a.x21 + a.x00*a.x13*a.x21 + a.x01*a.x10*a.x23 - a.x00*a.x11*a.x23) * d
	m.x30 = (a.x12*a.x21*a.x30 - a.x11*a.x22*a.x30 - a.x12*a.x20*a.x31 + a.x10*a.x22*a.x31 + a.x11*a.x20*a.x32 - a.x10*a.x21*a.x32) * d
	m.x31 = (a.x01*a.x22*a.x30 - a.x02*a.x21*a.x30 + a.x02*a.x20*a.x31 - a.x00*a.x22*a.x31 - a.x01*a.x20*a.x32 + a.x00*a.x21*a.x32) * d
	m.x32 = (a.x02*a.x11*a.x30 - a.x01*a.x12*a.x30 - a.x02*a.x10*a.x31 + a.x00*a.x12*a.x31 + a.x01*a.x10*a.x32 - a.x00*a.x11*a.x32) * d
	m.x33 = (a.x01*a.x12*a.x20 - a.x02*a.x11*a.x20 + a.x02*a.x10*a.x21 - a.x00*a.x12*a.x21 - a.x01*a.x10*a.x22 + a.x00*a.x11*a.x22) * d
	return m
}

// Inverse returns the inverse of a 3x3 matrix.
func (a M33) Inverse() M33 {
	m := M33{}
	d := 1 / a.Determinant()
	m.x00 = (a.x11*a.x22 - a.x12*a.x21) * d
	m.x01 = (a.x21*a.x02 - a.x01*a.x22) * d
	m.x02 = (a.x01*a.x12 - a.x11*a.x02) * d
	m.x10 = (a.x12*a.x20 - a.x22*a.x10) * d
	m.x11 = (a.x22*a.x00 - a.x20*a.x02) * d
	m.x12 = (a.x02*a.x10 - a.x12*a.x00) * d
	m.x20 = (a.x10*a.x21 - a.x20*a.x11) * d
	m.x21 = (a.x20*a.x01 - a.x00*a.x21) * d
	m.x22 = (a.x00*a.x11 - a.x01*a.x10) * d
	return m
}

// Inverse returns the inverse of a 2x2 matrix.
func (a M22) Inverse() M22 {
	m := M22{}
	d := 1 / a.Determinant()
	m.x00 = a.x11 * d
	m.x01 = -a.x01 * d
	m.x10 = -a.x10 * d
	m.x11 = a.x00 * d
	return m
}

//-----------------------------------------------------------------------------
