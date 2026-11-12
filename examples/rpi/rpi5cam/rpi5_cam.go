//-----------------------------------------------------------------------------
/*

Raspberry Pi 5 with Module 3 Camera Board Support

*/
//-----------------------------------------------------------------------------

package rpi5cam

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

const baseThickness = 3.0
const pillarHeight = 15.0

//-----------------------------------------------------------------------------
// RPi Camera Module 3

const cameraX = 21.0
const cameraY = 12.5

func cameraStandoffs(yOfs float64) (sdf.SDF3, error) {

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 5.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}

	const xOfs = 0.0
	const zOfs = 0.5 * (pillarHeight + baseThickness)

	const dX = cameraX * 0.5
	const dY = cameraY * 0.5

	positions := v3.VecSet{
		{xOfs + dX, yOfs + dY, zOfs},
		{xOfs - dX, yOfs - dY, zOfs},
		{xOfs + dX, yOfs - dY, zOfs},
		{xOfs - dX, yOfs + dY, zOfs},
	}

	s, _ := obj.Standoff3D(k)
	return sdf.Multi3D(s, positions), nil
}

//-----------------------------------------------------------------------------
// RPi 5

const rpiX = 49.0
const rpiY = 58.0

func rpiStandoffs(yOfs float64) (sdf.SDF3, error) {

	// standoffs with screw holes
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}

	const xOfs = 0.0
	const zOfs = 0.5 * (pillarHeight + baseThickness)

	const dX = rpiX * 0.5
	const dY = rpiY * 0.5

	positions := v3.VecSet{
		{xOfs + dX, yOfs + dY, zOfs},
		{xOfs - dX, yOfs - dY, zOfs},
		{xOfs + dX, yOfs - dY, zOfs},
		{xOfs - dX, yOfs + dY, zOfs},
	}

	s, _ := obj.Standoff3D(k)
	return sdf.Multi3D(s, positions), nil
}

//-----------------------------------------------------------------------------

func Support() (sdf.SDF3, error) {

	s0, _ := cameraStandoffs(-55.0)
	s1, _ := rpiStandoffs(10.0)

	const baseX = 65.0
	const baseY = 150.0

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{baseX, baseY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "xx", "x", "xx"},
	}
	p0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}

	// cutout
	c0 := sdf.Box2D(v2.Vec{baseX * 0.5, baseY * 0.65}, 3.0)
	c0 = sdf.Transform2D(c0, sdf.Translate2d(v2.Vec{0, 15.0}))

	// extrude the base
	s4 := sdf.Extrude3D(sdf.Difference2D(p0, c0), baseThickness)

	// filleting
	s5 := sdf.Union3D(s0, s1, s4)
	s5.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))
	return s5, nil
}

//-----------------------------------------------------------------------------
