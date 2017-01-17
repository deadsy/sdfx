//-----------------------------------------------------------------------------
/*

Matrix Operations

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

type M44 struct {
	x00, x01, x02, x03 float64
	x10, x11, x12, x13 float64
	x20, x21, x22, x23 float64
	x30, x31, x32, x33 float64
}

type M33 struct {
	x00, x01, x02 float64
	x10, x11, x12 float64
	x20, x21, x22 float64
}

//-----------------------------------------------------------------------------

func RandomM33(a, b float64) M33 {
	m := M33{}
	m.x00 = random_range(a, b)
	m.x10 = random_range(a, b)
	m.x20 = random_range(a, b)
	m.x01 = random_range(a, b)
	m.x11 = random_range(a, b)
	m.x21 = random_range(a, b)
	m.x02 = random_range(a, b)
	m.x12 = random_range(a, b)
	m.x22 = random_range(a, b)
	return m
}

//-----------------------------------------------------------------------------

func Identity44() M44 {
	return M44{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func Identity33() M33 {
	return M33{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}
}

var Test33 = M33{2, 1, 1, 3, 2, 1, 2, 1, 2}
var Test33_Inv = M33{3, -1, -1, -4, 2, 1, -1, 0, 1}

//-----------------------------------------------------------------------------

func (a M44) MulPosition(b V3) V3 {
	x := a.x00*b.X + a.x01*b.Y + a.x02*b.Z + a.x03
	y := a.x10*b.X + a.x11*b.Y + a.x12*b.Z + a.x13
	z := a.x20*b.X + a.x21*b.Y + a.x22*b.Z + a.x23
	return V3{x, y, z}
}

func (a M33) MulPosition(b V2) V2 {
	// TODO
	return V2{0, 0}
}

//-----------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------
// http://dev.theomader.com/transform-bounding-boxes/

func (a M44) MulBox(box Box3) Box3 {
	r := V3{a.x00, a.x10, a.x20}
	u := V3{a.x01, a.x11, a.x21}
	b := V3{a.x02, a.x12, a.x22}
	t := V3{a.x03, a.x13, a.x23}
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

func (a M33) MulBox(box Box2) Box2 {
	// TODO
	return Box2{}
}

//-----------------------------------------------------------------------------

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

func (a M33) Determinant() float64 {
	return (a.x00*a.x11*a.x22 - a.x00*a.x12*a.x21 -
		a.x01*a.x10*a.x22 + a.x01*a.x12*a.x20 +
		a.x02*a.x10*a.x21 - a.x02*a.x11*a.x20)
}

//-----------------------------------------------------------------------------

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

func (a M33) Inverse() M33 {
	m := M33{}
	d := 1 / a.Determinant()
	m.x00 = (a.x11*a.x22 - a.x12*a.x21) * d
	m.x01 = (a.x12*a.x20 - a.x01*a.x22) * d
	m.x02 = (a.x01*a.x12 - a.x11*a.x02) * d
	m.x10 = (a.x12*a.x20 - a.x22*a.x10) * d
	m.x11 = (a.x22*a.x00 - a.x20*a.x02) * d
	m.x12 = (a.x02*a.x10 - a.x12*a.x00) * d
	m.x20 = (a.x10*a.x21 - a.x20*a.x11) * d
	m.x21 = (a.x20*a.x01 - a.x00*a.x21) * d
	m.x22 = (a.x00*a.x11 - a.x01*a.x10) * d
	return m
}

//-----------------------------------------------------------------------------
