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
	"github.com/deadsy/sdfx/vec/v2i"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------

func gfShape(size v2.Vec, h0, h1, h2, h3, round float64) sdf.SDF3 {

	// upper (h0)
	k := TruncRectPyramidParms{
		Size:       v3.Vec{size.X, size.Y, h0},
		BaseAngle:  sdf.DtoR(45),
		BaseRadius: round,
	}
	upper, _ := TruncRectPyramid3D(&k)

	// middle (h1)
	size = size.SubScalar(2.0 * h0)
	round -= h0
	m2d := sdf.Box2D(size, round)
	middle := sdf.Extrude3D(m2d, h1)
	middle = sdf.Transform3D(middle, sdf.Translate3d(v3.Vec{0, 0, h0 + 0.5*h1}))

	// lower (h2)
	k = TruncRectPyramidParms{
		Size:       v3.Vec{size.X, size.Y, h2},
		BaseAngle:  sdf.DtoR(45),
		BaseRadius: round,
	}
	lower, _ := TruncRectPyramid3D(&k)
	lower = sdf.Transform3D(lower, sdf.Translate3d(v3.Vec{0, 0, h0 + h1}))

	// extension (h3)
	var ext sdf.SDF3
	if h3 > 0 {
		size = size.SubScalar(2.0 * h2)
		round -= h2
		ext2d := sdf.Box2D(size, round)
		ext = sdf.Extrude3D(ext2d, h3)
		ext = sdf.Transform3D(ext, sdf.Translate3d(v3.Vec{0, 0, h0 + h1 + h2 + 0.5*h3}))
	}

	return sdf.Transform3D(sdf.Union3D(upper, middle, lower, ext), sdf.RotateX(sdf.Pi))
}

func gfGrid(x, y int, zOfs float64) []v3.Vec {
	grid := make([]v3.Vec, x*y)
	xOfs := -0.5 * float64(x-1) * gfFemaleSize
	yOfs := -0.5 * float64(y-1) * gfFemaleSize
	idx := 0
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			grid[idx] = v3.Vec{xOfs + float64(i)*gfFemaleSize, yOfs + float64(j)*gfFemaleSize, zOfs}
			idx++
		}
	}
	return grid
}

//-----------------------------------------------------------------------------

const gfHoleOffset = 4.8
const gfHoleMinor = 0.5 * 3.0
const gfHoleMajor = 0.5 * 6.5
const gfHoleHeight = 2.0

func gfHoles(r, h, zOfs float64) sdf.SDF3 {
	const ofs = 0.5*gfMaleSize - (gfMaleH0 + gfMaleH2 + gfHoleOffset)
	hole, _ := sdf.Cylinder3D(h, r, 0)
	posn := []v3.Vec{
		{ofs, ofs, zOfs},
		{-ofs, ofs, zOfs},
		{ofs, -ofs, zOfs},
		{-ofs, -ofs, zOfs},
	}
	return sdf.Multi3D(hole, posn)
}

func gfThruHoles(h float64) sdf.SDF3 {
	zOfs := 0.5*h - gfMaleHeight + gfHoleHeight
	return gfHoles(gfHoleMinor, h, zOfs)
}

const gfFemaleSize = 42.0
const gfFemaleRound = 0.5 * 8.0
const gfFemaleH0 = 2.15
const gfFemaleH1 = 1.8
const gfFemaleH2 = 0.7
const gfFemaleHeight = gfFemaleH0 + gfFemaleH1 + gfFemaleH2

func gfFemale(ext float64) sdf.SDF3 {
	return gfShape(v2.Vec{gfFemaleSize, gfFemaleSize}, gfFemaleH0, gfFemaleH1, gfFemaleH2, ext, gfFemaleRound)
}

const gfMaleSize = 41.5
const gfMaleRound = 0.5 * 7.5
const gfMaleH0 = 2.15
const gfMaleH1 = 1.8
const gfMaleH2 = 0.8
const gfMaleHeight = gfMaleH0 + gfMaleH1 + gfMaleH2

func gfMale() sdf.SDF3 {
	plug := gfShape(v2.Vec{gfMaleSize, gfMaleSize}, gfMaleH0, gfMaleH1, gfMaleH2, 0, gfMaleRound)
	holes := gfHoles(gfHoleMajor, gfHoleHeight, 0.5*gfHoleHeight-gfMaleHeight)
	return sdf.Difference3D(plug, holes)
}

const gfLipRound = 0.5 * 7.5
const gfLipH0 = 1.9
const gfLipH1 = 1.8
const gfLipH2 = 0.7
const gfLipHeight = gfLipH0 + gfLipH1 + gfLipH2

func gfLip(x, y, empty float64) sdf.SDF3 {
	return gfShape(v2.Vec{x, y}, gfLipH0, gfLipH1, gfLipH2, empty, gfLipRound)
}

const gfHeightSize = 7.0

// values not in the specifications
const gfFloor = 1.0      // floor thickness for an empty container
const gfBaseHeight = 4.0 // extra base height (for magnet mounts, side attachments)

//-----------------------------------------------------------------------------

// GfBaseParms are the gridfinity base parameters.
type GfBaseParms struct {
	Size   v2i.Vec // size of base in gridfinity units
	Magnet bool    // add magnet mounts
	Hole   bool    // add mounting holes
}

// GfBase returns a Gridfinity base grid.
func GfBase(k *GfBaseParms) sdf.SDF3 {
	if k.Size.X <= 0 {
		k.Size.X = 1
	}
	if k.Size.Y <= 0 {
		k.Size.Y = 1
	}

	h := gfFemaleHeight
	if k.Magnet || k.Hole {
		h += gfBaseHeight
	}

	// base body
	size := v2.Vec{float64(k.Size.X), float64(k.Size.Y)}.MulScalar(gfFemaleSize)
	b2d := sdf.Box2D(size, gfFemaleRound)
	base := sdf.Extrude3D(b2d, h)

	// main holes
	grid := gfGrid(k.Size.X, k.Size.Y, 0.5*h)
	holes := sdf.Multi3D(gfFemale(h-gfFemaleHeight), grid)

	// magnet mounts
	var magnets sdf.SDF3
	if k.Magnet || k.Hole {
		const r = gfMaleH0 + gfMaleH2 + gfHoleOffset
		magnets = sdf.Multi3D(gfHoles(r, gfBaseHeight, -h+0.5*gfBaseHeight), grid)
		zOfs := -0.5*gfHoleHeight - h + gfBaseHeight
		magnetHoles := sdf.Multi3D(gfHoles(gfHoleMajor, gfHoleHeight, zOfs), grid)
		magnets = sdf.Difference3D(magnets, magnetHoles)
	}

	// mounting holes
	if k.Hole {
		mountHoles := sdf.Multi3D(gfHoles(gfHoleMinor, h, -0.5*h), grid)
		magnets = sdf.Difference3D(magnets, mountHoles)
	}

	return sdf.Union3D(sdf.Difference3D(base, holes), magnets)
}

//-----------------------------------------------------------------------------

// GfBodyParms are the gridfinity body parameters.
type GfBodyParms struct {
	Size  v3i.Vec // size of body in gridfinity units
	Empty bool    // return an empty container
	Hole  bool    // add through holes to the body
}

// GfBody returns a gridfinity body.
func GfBody(k *GfBodyParms) sdf.SDF3 {

	if k.Size.X <= 0 {
		k.Size.X = 1
	}
	if k.Size.Y <= 0 {
		k.Size.Y = 1
	}
	if k.Size.Z <= 0 {
		k.Size.Z = 1
	}

	// body
	size := v2.Vec{float64(k.Size.X), float64(k.Size.Y)}.MulScalar(gfFemaleSize).SubScalar(gfFemaleSize - gfMaleSize)
	b2d := sdf.Box2D(size, gfMaleRound)
	h := (float64(k.Size.Z) * gfHeightSize) + gfLipHeight - gfMaleHeight
	body := sdf.Extrude3D(b2d, h)

	// grid positions
	grid := gfGrid(k.Size.X, k.Size.Y, -0.5*h)

	// base plugs
	plugs := sdf.Multi3D(gfMale(), grid)

	// through holes
	var holes sdf.SDF3
	if k.Hole {
		holes = sdf.Multi3D(gfThruHoles(h+gfMaleHeight-gfHoleHeight), grid)
	}

	// stacking lip
	empty := 0.0
	if k.Empty {
		empty = h - gfLipHeight - gfFloor
	}
	lip := gfLip(size.X, size.Y, empty)
	lip = sdf.Transform3D(lip, sdf.Translate3d(v3.Vec{0, 0, 0.5 * h}))

	return sdf.Difference3D(sdf.Union3D(body, plugs), sdf.Union3D(lip, holes))
}

//-----------------------------------------------------------------------------
