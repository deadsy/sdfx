//-----------------------------------------------------------------------------
/*

Fidget Spinner

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var bearing_od = 22.0
var bearing_id = 8.0
var bearing_w = 7.0

var bearing_or = bearing_od / 2
var bearing_ir = bearing_id / 2

//-----------------------------------------------------------------------------

// Return an N petal bezier flower.
func flower(n int, r0, r1, r2 float64) SDF2 {

	theta := TAU / float64(n)
	b := NewBezier()

	p0 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(-135)))
	p1 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(-45)))
	p2 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(45)))
	p3 := V2{r1, 0}.Add(PolarToXY(r0, DtoR(135)))
	p4 := PolarToXY(r2, theta/2)

	m := Rotate(theta)

	for i := 0; i < n; i++ {
		ofs := float64(i) * theta

		b.AddV2(p0).Handle(ofs+DtoR(-45), r0/2, r0/2)
		b.AddV2(p1).Handle(ofs+DtoR(45), r0/2, r0/2)
		b.AddV2(p2).Handle(ofs+DtoR(135), r0/2, r0/2)
		b.AddV2(p3).Handle(ofs+DtoR(225), r0/2, r0/2)
		b.AddV2(p4).Handle(ofs+theta/2+DtoR(90), r2/1.5, r2/1.5)

		p0 = m.MulPosition(p0)
		p1 = m.MulPosition(p1)
		p2 = m.MulPosition(p2)
		p3 = m.MulPosition(p3)
		p4 = m.MulPosition(p4)
	}

	b.Close()
	p := b.Polygon()
	return Polygon2D(p.Vertices())
}

//-----------------------------------------------------------------------------

func body() SDF3 {

	n := 3
	r0 := bearing_or + 4.0
	r1 := 45.0 - r0
	r2 := bearing_or + 4.0

	// body
	s1 := ExtrudeRounded3D(flower(n, r0, r1, r2), bearing_w, bearing_w/4.0)
	// periphery holes
	s2 := MakeBoltCircle3D(bearing_w, bearing_or, r1, n)
	// center hole
	s3 := Cylinder3D(bearing_w, bearing_or, 0)

	return Difference3D(s1, Union3D(s2, s3))
}

//-----------------------------------------------------------------------------

func spincap() SDF3 {

	t := 3.0
	sx := (12.0 - bearing_id) / 2.0
	sy := 1.0
	h := t + sy + (bearing_w-1.0)/2.0

	p := NewPolygon()
	p.Add(0, 0)
	p.Add(bearing_or, 0).Smooth(t/1.5,6)
	p.Add(bearing_or, t)
	p.Add(bearing_ir+sx, t)
	p.Add(bearing_ir+sx, t+sy)
	p.Add(bearing_ir, t+sy)
	p.Add(bearing_ir, h)
	p.Add(0, h)

	s0 := Polygon2D(p.Vertices())
	return Revolve3D(s0)
}

//-----------------------------------------------------------------------------

func main() {
	RenderSTL(body(), 300, "body.stl")
	RenderSTL(spincap(), 300, "cap.stl")
}

//-----------------------------------------------------------------------------
