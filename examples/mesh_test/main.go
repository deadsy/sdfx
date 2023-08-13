//-----------------------------------------------------------------------------
/*


 */
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

func testPolygon() (*sdf.Polygon, error) {
	b := sdf.NewBezier()
	b.Add(-788.571430, 666.647920)
	b.Add(-788.785400, 813.701340).Mid()
	b.Add(-759.449240, 1023.568700).Mid()
	b.Add(-588.571430, 1026.647900)
	b.Add(-417.693610, 1029.727200).Mid()
	b.Add(-583.793160, 507.272270).Mid()
	b.Add(-285.714290, 506.647920)
	b.Add(12.364584, 506.023560).Mid()
	b.Add(-137.634380, 1110.386900).Mid()
	b.Add(85.714281, 1115.219300)
	b.Add(309.062940, 1120.051800).Mid()
	b.Add(498.298980, 1086.587000).Mid()
	b.Add(491.428570, 903.790780)
	b.Add(484.558160, 720.994550).Mid()
	b.Add(79.128329, 547.886390).Mid()
	b.Add(62.857140, 292.362210)
	b.Add(46.585951, 36.838026).Mid()
	b.Add(367.678530, -5.375978).Mid()
	b.Add(374.285720, -179.066370)
	b.Add(380.892900, -352.756760).Mid()
	b.Add(273.020040, -521.481290).Mid()
	b.Add(131.428570, -521.923510)
	b.Add(-10.162890, -522.365730).Mid()
	b.Add(50.355420, -54.901413).Mid()
	b.Add(-134.285720, -59.066363)
	b.Add(-318.926860, -63.231312).Mid()
	b.Add(-304.285720, -429.542560).Mid()
	b.Add(-442.857150, -433.352080)
	b.Add(-581.428570, -437.161610).Mid()
	b.Add(-750.919960, -371.353320).Mid()
	b.Add(-748.571430, -221.923510)
	b.Add(-746.222890, -72.493698).Mid()
	b.Add(-413.586510, -77.312515).Mid()
	b.Add(-402.857140, 120.933630)
	b.Add(-424.396820, 260.368600).Mid()
	b.Add(-788.357460, 519.594510).Mid()
	b.Add(-788.571430, 666.647920)
	b.Close()
	return b.Polygon()
}

func getLines() []*sdf.Line2 {
	p, _ := testPolygon()
	l, _ := p.Lines()
	return l
}

//-----------------------------------------------------------------------------

func test1() error {

	m := getLines()

	s0, err := sdf.Mesh2D(m)
	if err != nil {
		return err
	}

	s1, err := sdf.Mesh2DSlow(m)
	if err != nil {
		return err
	}

	d := render.NewDXF("test.dxf")
	var ePoints []v2.Vec

	bb := s0.BoundingBox()
	for _, p := range bb.RandomSet(100000) {
		d0 := s0.Evaluate(p)
		d1 := s1.Evaluate(p)
		if !sdf.EqualFloat64(d0, d1, 1e-12) {
			ePoints = append(ePoints, p)
		}
	}

	d.Lines(m)

	boxes := s0.(*sdf.MeshSDF2).Boxes()
	for _, b := range boxes {
		d.Box(b)
	}

	log.Printf("%d distance errors", len(ePoints))
	if len(ePoints) != 0 {
		d.Points(ePoints, 0.2)
	}

	d.Save()

	return nil
}

//-----------------------------------------------------------------------------

func test2() error {

	m := getLines()

	s0, err := sdf.Mesh2D(m)
	if err != nil {
		return err
	}

	s1, err := sdf.Mesh2DSlow(m)
	if err != nil {
		return err
	}

	sdf.BenchmarkSDF2("Mesh2D", s0)
	sdf.BenchmarkSDF2("Mesh2DSlow", s1)

	return nil
}

//-----------------------------------------------------------------------------

func main() {
	err := test1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = test2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

}

//-----------------------------------------------------------------------------
