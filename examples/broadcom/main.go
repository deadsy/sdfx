//-----------------------------------------------------------------------------
/*

Broadcom Development Board Mounting Kits

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

var baseThickness = 3.0
var pillarHeight = 45.0

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

// cutout returns a rounded rectangle spanning (x0,y0) to (x1,y1).
func cutout(x0, y0, x1, y1 float64) sdf.SDF2 {
	s := sdf.Box2D(v2.Vec{x1 - x0, y1 - y0}, 3.0)
	return sdf.Transform2D(s, sdf.Translate2d(v2.Vec{0.5 * (x0 + x1), 0.5 * (y0 + y1)}))
}

//-----------------------------------------------------------------------------
// BCM947822EAP
// PCB: 200-WC000476-0T01-XX rev 01, 270 x 250 mm (panel rails removed)

// lower-left corner of the PCB in drill file coordinates (mils)
var bcm947822eapOrigin = v2.Vec{-6251.4, -4873.0}

// bcm947822eapPosition converts a drill file location (mils) to
// PCB relative coordinates (mm) at height z.
func bcm947822eapPosition(x, y, z float64) v3.Vec {
	return v3.Vec{
		(x - bcm947822eapOrigin.X) * sdf.Mil,
		(y - bcm947822eapOrigin.Y) * sdf.Mil,
		z,
	}
}

func bcm947822eapStandoffs() (sdf.SDF3, error) {

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes: NPTH 149.6 mil (M3)
	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 8.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.5, // M3 self-tapping
	}
	positions0 := v3.VecSet{
		bcm947822eapPosition(-5975.0, -4575.0, zOfs),
		bcm947822eapPosition(-5260.0, -2269.9, zOfs),
		bcm947822eapPosition(-5040.0, 781.0, zOfs),
		bcm947822eapPosition(-5040.0, 4484.0, zOfs),
		bcm947822eapPosition(4090.0, 4484.0, zOfs),
		bcm947822eapPosition(4090.0, -497.1, zOfs),
		bcm947822eapPosition(4100.0, -4594.0, zOfs),
	}
	s, err := obj.Standoff3D(k)
	if err != nil {
		return nil, err
	}
	s0 := sdf.Multi3D(s, positions0)

	// cpu supports: PTH 125.2 mil
	// The heatsink holes are a square of 4 at (+/-1771.7, +/-1771.7) mil,
	// but the lower 2 are occluded by metal shields so only the upper 2 are used.
	// The hole gives clearance for anything poking through the PCB.
	// Note: the PTH 126 mil holes are for a daughter board, not standoffs.
	positions1 := v3.VecSet{
		bcm947822eapPosition(-1771.7, 1771.7, zOfs),
		bcm947822eapPosition(1771.7, 1771.7, zOfs),
	}
	s1 := sdf.Multi3D(s, positions1)

	return sdf.Union3D(s0, s1), nil
}

func bcm947822eap() (sdf.SDF3, error) {

	pcbX := 270.0
	pcbY := 250.0

	// base
	pp := &obj.PanelParms{
		Size:         v2.Vec{pcbX, pcbY},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{".x.x", ".x.x", ".x.x", ".x.x"}, // 2 per side at 1/4 and 3/4, clear of the corner standoffs
	}
	s0, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}
	s0 = sdf.Transform2D(s0, sdf.Translate2d(v2.Vec{0.5 * pcbX, 0.5 * pcbY}))

	// cutouts (PCB coordinates), clear of the standoffs
	c := sdf.Union2D(
		cutout(45, 20, 250, 150),
		cutout(45, 162, 104, 225),
		cutout(124, 162, 194, 225),
		cutout(214, 162, 250, 225),
	)

	// extrude the base
	s1 := sdf.Extrude3D(sdf.Difference2D(s0, c), baseThickness)

	// add the standoffs
	s2, err := bcm947822eapStandoffs()
	if err != nil {
		return nil, err
	}
	s3 := sdf.Union3D(s1, s2)
	s3.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	// The board is mounted upside down (top side facing the base) so the
	// daughter board on the bottom side is accessible. Mirror the platform
	// about the PCB center line (x = pcbX/2) to match.
	m := sdf.Translate3d(v3.Vec{pcbX, 0, 0}).Mul(sdf.MirrorYZ())
	return sdf.Transform3D(s3, m), nil
}

//-----------------------------------------------------------------------------
// M3 self-tapping screw hole test piece.
// Standoffs with differing hole diameters to work out the best printed
// hole size. Left to right the holes increase in size and are marked with
// 1, 2, 3 small holes in the base.

func m3Test() (sdf.SDF3, error) {

	const height = 20.0
	const pitch = 15.0
	holes := []float64{2.4, 2.5, 2.6}

	n := float64(len(holes))
	baseX := n * pitch
	baseY := 24.0

	// base with marker holes
	base := sdf.Box2D(v2.Vec{baseX, baseY}, 3.0)
	marker, err := sdf.Circle2D(0.75)
	if err != nil {
		return nil, err
	}
	var markers []sdf.SDF2
	for i := range holes {
		x := (float64(i) - 0.5*(n-1)) * pitch
		for j := 0; j <= i; j++ {
			mx := x + (float64(j)-0.5*float64(i))*3.0
			markers = append(markers, sdf.Transform2D(marker, sdf.Translate2d(v2.Vec{mx, -8.5})))
		}
	}
	s0 := sdf.Extrude3D(sdf.Difference2D(base, sdf.Union2D(markers...)), baseThickness)

	// standoffs
	zOfs := 0.5 * (height + baseThickness)
	var standoffs []sdf.SDF3
	for i, d := range holes {
		k := &obj.StandoffParms{
			PillarHeight:   height,
			PillarDiameter: 8.0,
			HoleDepth:      10.0,
			HoleDiameter:   d,
		}
		s, err := obj.Standoff3D(k)
		if err != nil {
			return nil, err
		}
		x := (float64(i) - 0.5*(n-1)) * pitch
		standoffs = append(standoffs, sdf.Transform3D(s, sdf.Translate3d(v3.Vec{x, 0, zOfs})))
	}

	s1 := sdf.Union3D(append([]sdf.SDF3{s0}, standoffs...)...)
	s1.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s1, nil
}

//-----------------------------------------------------------------------------

func main() {

	bcm947822eap, err := bcm947822eap()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(bcm947822eap, shrink), "bcm947822eap.stl", render.NewMarchingCubesOctree(500))

	m3Test, err := m3Test()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(m3Test, shrink), "m3_test.stl", render.NewMarchingCubesOctree(600))
}

//-----------------------------------------------------------------------------
