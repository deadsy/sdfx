//-----------------------------------------------------------------------------
/*

Bee Hive Parts

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

// entrance0 returns an open entrance
func entrance0(size sdf.V3) (sdf.SDF3, error) {
	r := size.Y * 0.5
	s0 := sdf.Line2D(size.X-(2*r), r)
	s1 := sdf.Extrude3D(s0, size.Z)
	return s1, nil
}

// entrance1 returns a vent entrance
func entrance1(size sdf.V3) (sdf.SDF3, error) {

	const rows = 3
	const cols = 16
	const holeRadius = 3.2 * 0.5

	hole, err := sdf.Circle2D(holeRadius)
	if err != nil {
		return nil, err
	}

	size.X -= 2 * holeRadius
	size.Y -= 2 * holeRadius
	dx := size.X / (cols - 1)
	dy := size.Y / (rows - 1)
	xOfs := -size.X / 2
	yOfs := size.Y / 2

	positions := []sdf.V2{}
	x := xOfs
	for i := 0; i < cols; i++ {
		y := yOfs
		for j := 0; j < rows; j++ {
			positions = append(positions, sdf.V2{x, y})
			y -= dy
		}
		x += dx
	}
	s := sdf.Multi2D(hole, positions)

	return sdf.Extrude3D(s, size.Z), nil
}

func entranceWheel() (sdf.SDF3, error) {
	const radius = 6.5 * 0.5 * sdf.MillimetresPerInch
	const thickness = 0.25 * sdf.MillimetresPerInch

	plate, err := sdf.Cylinder3D(thickness, radius, 0)
	if err != nil {
		return nil, err
	}

	hole, err := sdf.Cylinder3D(thickness, 2.5, 0)
	if err != nil {
		return nil, err
	}

	entranceSize := sdf.V3{
		4 * sdf.MillimetresPerInch,
		0.5 * sdf.MillimetresPerInch,
		thickness,
	}

	const k = 1.6
	ofs := k * entranceSize.X * 0.5 * math.Tan(sdf.DtoR(30))

	// open entrance
	e0, err := entrance0(entranceSize)
	if err != nil {
		return nil, err
	}
	e0 = sdf.Transform3D(e0, sdf.Translate3d(sdf.V3{0, ofs, 0}))

	// vent entrance
	e1, err := entrance1(entranceSize)
	if err != nil {
		return nil, err
	}
	e1 = sdf.Transform3D(e1, sdf.Translate3d(sdf.V3{0, ofs, 0}))
	e1 = sdf.Transform3D(e1, sdf.RotateZ(sdf.DtoR(120)))

	return sdf.Difference3D(plate, sdf.Union3D(e0, e1, hole)), nil
}

//-----------------------------------------------------------------------------

func holePattern(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = byte('x')
	}
	return string(s)
}

func entranceReducer() (sdf.SDF3, error) {

	const zSize = 0.25 * sdf.MillimetresPerInch
	const xSize = 6.0 * sdf.MillimetresPerInch
	const ySize = 1.9 * sdf.MillimetresPerInch

	k := obj.PanelParms{
		Size:         sdf.V2{xSize, ySize},
		CornerRadius: 5.0,
	}
	s, err := obj.Panel2D(&k)
	if err != nil {
		return nil, err
	}

	const holeRadius = (3.0 / 16.0) * sdf.MillimetresPerInch
	hole := sdf.Line2D(2*holeRadius, holeRadius)
	hole = sdf.Transform2D(hole, sdf.Rotate2d(sdf.DtoR(90)))

	const entranceSize = 4.0 * sdf.MillimetresPerInch
	const n = 6
	const gap = (entranceSize - (n * holeRadius)) / (n + 1)
	const yOfs = -ySize * 0.5
	const xOfs = (n - 1) * (holeRadius + gap) * 0.5
	p0 := sdf.V2{-xOfs, yOfs}
	p1 := sdf.V2{xOfs + holeRadius + gap, yOfs}
	hole = sdf.LineOf2D(hole, p0, p1, holePattern(n))

	return sdf.Extrude3D(sdf.Difference2D(s, hole), zSize), nil
}

//-----------------------------------------------------------------------------

func main() {

	p0, err := entranceReducer()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(p0, shrink), 300, "reducer.stl")

	p1, err := entranceWheel()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(p1, shrink), 300, "wheel.stl")
}

//-----------------------------------------------------------------------------
