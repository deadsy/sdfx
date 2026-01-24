//-----------------------------------------------------------------------------
/*

Pico RX Housing


*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// display0Bezel : 320x240 TJCTM24028-SPI
func display0Bezel(thickness float64, positive bool) (sdf.SDF3, error) {

	const displayX = 60.0
	const displayY = 45.0
	const cornerRounding = 2.0

	if positive == false {
		// return a panel hole for the display
		s0 := sdf.Box2D(v2.Vec{displayX, displayY}, cornerRounding)
		return sdf.Extrude3D(s0, thickness), nil
	}

	// return the display support standoffs
	const standOffZ = 4.0

	// standoffs with screw holes
	k0 := &obj.StandoffParms{
		PillarHeight:   standOffZ,
		PillarDiameter: 5.0,
		HoleDepth:      standOffZ,
		HoleDiameter:   2.4, // #4 screw
	}

	s0, err := obj.Standoff3D(k0)
	if err != nil {
		return nil, err
	}

	const xOfs = 0.5 * (86.0 - 3.0 - 6.92)
	const yOfs = 0.5 * (50.0 - 3.0 - 3.0)
	zOfs := 0.5 * (standOffZ + thickness)
	skew := xOfs - 0.5*displayX - 4.5

	positions := v3.VecSet{
		{xOfs + skew, yOfs, zOfs},
		{xOfs + skew, -yOfs, zOfs},
		{-xOfs + skew, yOfs, zOfs},
		{-xOfs + skew, -yOfs, zOfs},
	}
	return sdf.Multi3D(s0, positions), nil
}

func bezel0() (sdf.SDF3, error) {

	const panelX = 100.0
	const panelY = 60.0
	const thickness = 3.0

	p0 := sdf.Box2D(v2.Vec{panelX, panelY}, 2.0)
	panel := sdf.Extrude3D(p0, thickness)

	b0, err := display0Bezel(thickness, false)
	if err != nil {
		return nil, err
	}

	b1, err := display0Bezel(thickness, true)
	if err != nil {
		return nil, err
	}

	return sdf.Union3D(b1, sdf.Difference3D(panel, b0)), nil
}

//-----------------------------------------------------------------------------

// display1Bezel : GME12864-11 (128x64 SSD1306)
func display1Bezel(thickness float64, positive bool) (sdf.SDF3, error) {

	const displayX = 26.0
	const displayY = 14.0
	const cornerRounding = 1.0

	if positive == false {
		// return a panel hole for the display
		s0 := sdf.Box2D(v2.Vec{displayX, displayY}, cornerRounding)
		return sdf.Extrude3D(s0, thickness), nil
	}

	// return the display support standoffs
	const standOffZ = 2.1

	// standoffs with screw holes
	k0 := &obj.StandoffParms{
		PillarHeight:   standOffZ,
		PillarDiameter: 4.0,
		HoleDepth:      standOffZ,
		HoleDiameter:   2.4, // #4 screw
	}

	s0, err := obj.Standoff3D(k0)
	if err != nil {
		return nil, err
	}

	const xOfs = 0.5 * 23.5
	const yOfs = 0.5 * 23.8
	zOfs := 0.5 * (standOffZ + thickness)
	skew := yOfs - 0.5*displayY - 2.0

	positions := v3.VecSet{
		{xOfs, yOfs + skew, zOfs},
		{xOfs, -yOfs + skew, zOfs},
		{-xOfs, yOfs + skew, zOfs},
		{-xOfs, -yOfs + skew, zOfs},
	}
	return sdf.Multi3D(s0, positions), nil
}

func bezel1() (sdf.SDF3, error) {

	const panelX = 40.0
	const panelY = 40.0
	const thickness = 3.0

	p0 := sdf.Box2D(v2.Vec{panelX, panelY}, 2.0)
	panel := sdf.Extrude3D(p0, thickness)

	b0, err := display1Bezel(thickness, false)
	if err != nil {
		return nil, err
	}

	b1, err := display1Bezel(thickness, true)
	if err != nil {
		return nil, err
	}

	return sdf.Union3D(b1, sdf.Difference3D(panel, b0)), nil
}

//-----------------------------------------------------------------------------

func bezel2() (sdf.SDF3, error) {

	const panelX = 80.0
	const panelY = 40.0
	const panelThickness = 2.5

	p0 := sdf.Box2D(v2.Vec{panelX, panelY}, 2.0)
	panel := sdf.Extrude3D(p0, panelThickness)

	// push button
	pb, err := sdf.Box3D(v3.Vec{13.2, 10.8, panelThickness}, 0)
	if err != nil {
		return nil, err
	}
	const xOfs = 24.0
	pb0 := sdf.Transform3D(pb, sdf.Translate3d(v3.Vec{xOfs, 0, 0}))
	pb1 := sdf.Transform3D(pb, sdf.Translate3d(v3.Vec{-xOfs, 0, 0}))

	// rotary encoder
	k := obj.PanelHoleParms{
		Diameter:  9.8,
		Thickness: panelThickness,
	}
	r0, err := obj.PanelHole3D(&k)
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(panel, sdf.Union3D(pb0, pb1, r0)), nil
}

//-----------------------------------------------------------------------------

func main() {
	b0, err := bezel0()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(b0, "bezel0.stl", render.NewMarchingCubesOctree(500))

	b1, err := bezel1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(b1, "bezel1.stl", render.NewMarchingCubesOctree(500))

	b2, err := bezel2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(b2, "bezel2.stl", render.NewMarchingCubesOctree(500))

}

//-----------------------------------------------------------------------------
