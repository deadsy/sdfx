//-----------------------------------------------------------------------------
/*

Gridfinity Storage Parts

https://gridfinity.xyz/

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func gfShape(size v2.Vec, h0, h1, h2, round float64) sdf.SDF3 {
	// upper
	k := TruncRectPyramidParms{
		Size:       v3.Vec{size.X, size.Y, h2},
		BaseAngle:  sdf.DtoR(45),
		BaseRadius: round,
	}
	upper, _ := TruncRectPyramid3D(&k)

	// middle
	size = size.SubScalar(2.0 * h2)
	round -= h2
	m2d := sdf.Box2D(size, round)
	middle := sdf.Extrude3D(m2d, h1)
	middle = sdf.Transform3D(middle, sdf.Translate3d(v3.Vec{0, 0, h2 + 0.5*h1}))

	//lower
	k = TruncRectPyramidParms{
		Size:       v3.Vec{size.X, size.Y, h0},
		BaseAngle:  sdf.DtoR(45),
		BaseRadius: round,
	}
	lower, _ := TruncRectPyramid3D(&k)
	lower = sdf.Transform3D(lower, sdf.Translate3d(v3.Vec{0, 0, h2 + h1}))

	s := sdf.Union3D(upper, middle, lower)
	return sdf.Transform3D(s, sdf.RotateX(sdf.Pi))
}

//-----------------------------------------------------------------------------

const gfFemaleSize = 42.0
const gfFemaleRound = 8.0 * 0.5
const gfFemaleH0 = 0.7
const gfFemaleH1 = 1.8
const gfFemaleH2 = 2.15

func gfFemale() sdf.SDF3 {
	return gfShape(v2.Vec{gfFemaleSize, gfFemaleSize}, gfFemaleH0, gfFemaleH1, gfFemaleH2, gfFemaleRound)
}

const gfMaleSize = 41.5
const gfMaleRound = 7.5 * 0.5
const gfMaleH0 = 0.8
const gfMaleH1 = 1.8
const gfMaleH2 = 2.15
const gfMaleHeight = gfMaleH0 + gfMaleH1 + gfMaleH2

func gfMale() sdf.SDF3 {
	return gfShape(v2.Vec{gfMaleSize, gfMaleSize}, gfMaleH0, gfMaleH1, gfMaleH2, gfMaleRound)
}

const gfLipRound = 7.5 * 0.5
const gfLipH0 = 0.7
const gfLipH1 = 1.8
const gfLipH2 = 1.9
const gfLipHeight = gfLipH0 + gfLipH1 + gfLipH2

func gfLip(x, y float64) sdf.SDF3 {
	return gfShape(v2.Vec{x, y}, gfLipH0, gfLipH1, gfLipH2, gfLipRound)
}

//-----------------------------------------------------------------------------

// GfBaseParms are the gridfinity base parameters.
type GfBaseParms struct {
	X, Y   int     // size of base in gridfinity units
	Height float64 // height of base lattice
}

// GfBase returns a Gridfinity base grid.
func GfBase(k *GfBaseParms) sdf.SDF3 {
	if k.X <= 0 {
		k.X = 1
	}
	if k.Y <= 0 {
		k.Y = 1
	}
	const h = gfFemaleH0 + gfFemaleH1 + gfFemaleH2
	if k.Height < h {
		k.Height = h
	}

	// base body
	size := v2.Vec{float64(k.X), float64(k.Y)}.MulScalar(gfFemaleSize)
	b2d := sdf.Box2D(size, gfFemaleRound)
	base := sdf.Extrude3D(b2d, k.Height)

	// female holes
	posn := make([]v3.Vec, k.X*k.Y)
	xOfs := -0.5 * float64(k.X-1) * gfFemaleSize
	yOfs := -0.5 * float64(k.Y-1) * gfFemaleSize
	zOfs := k.Height * 0.5
	idx := 0
	for i := 0; i < k.X; i++ {
		for j := 0; j < k.Y; j++ {
			posn[idx] = v3.Vec{xOfs + float64(i)*gfFemaleSize, yOfs + float64(j)*gfFemaleSize, zOfs}
			idx++
		}
	}
	holes := sdf.Multi3D(gfFemale(), posn)

	return sdf.Difference3D(base, holes)
}

//-----------------------------------------------------------------------------

// GfBodyParms are the gridfinity body parameters.
type GfBodyParms struct {
	X, Y, Z int // size of body in gridfinity units
}

const gfHeightSize = 7.0

// GfBody returns a gridfinity body.
func GfBody(k *GfBodyParms) sdf.SDF3 {

	if k.X <= 0 {
		k.X = 1
	}
	if k.Y <= 0 {
		k.Y = 1
	}
	if k.Z <= 0 {
		k.Z = 1
	}

	// body
	size := v2.Vec{float64(k.X), float64(k.Y)}.MulScalar(gfFemaleSize).SubScalar(gfFemaleSize - gfMaleSize)
	b2d := sdf.Box2D(size, gfMaleRound)
	h := (float64(k.Z) * gfHeightSize) + gfLipH0 - gfMaleHeight
	body := sdf.Extrude3D(b2d, h)

	// base plugs
	posn := make([]v3.Vec, k.X*k.Y)
	xOfs := -0.5 * float64(k.X-1) * gfFemaleSize
	yOfs := -0.5 * float64(k.Y-1) * gfFemaleSize
	zOfs := -0.5 * h
	idx := 0
	for i := 0; i < k.X; i++ {
		for j := 0; j < k.Y; j++ {
			posn[idx] = v3.Vec{xOfs + float64(i)*gfFemaleSize, yOfs + float64(j)*gfFemaleSize, zOfs}
			idx++
		}
	}
	plugs := sdf.Multi3D(gfMale(), posn)

	// stacking lip
	lip := gfLip(size.X, size.Y)
	lip = sdf.Transform3D(lip, sdf.Translate3d(v3.Vec{0, 0, 0.5 * h}))

	return sdf.Difference3D(sdf.Union3D(body, plugs), lip)
}

//-----------------------------------------------------------------------------
