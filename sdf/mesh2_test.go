//-----------------------------------------------------------------------------
/*

Mesh 2D Testing and Benchmarking

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"testing"
)

//-----------------------------------------------------------------------------

func testPolygon() (*Polygon, error) {
	b := NewBezier()
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

func getLines() []*Line2 {
	p, _ := testPolygon()
	return VertexToLine(p.Vertices(), true)
}

//-----------------------------------------------------------------------------

func Benchmark_Mesh2D(b *testing.B) {
	m := getLines()
	s0, err := Mesh2D(m)
	if err != nil {
		b.Fatalf("error: %s", err)
	}
	bb := s0.BoundingBox()
	b.Run("Mesh2D", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = s0.Evaluate(bb.Random())
		}
	})
}

func Benchmark_Mesh2DSlow(b *testing.B) {
	m := getLines()
	s0, err := Mesh2DSlow(m)
	if err != nil {
		b.Fatalf("error: %s", err)
	}
	bb := s0.BoundingBox()
	b.Run("Mesh2DSlow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = s0.Evaluate(bb.Random())
		}
	})
}

//-----------------------------------------------------------------------------

const nPoints = 20000

func Test_Mesh2D(t *testing.T) {

	m := getLines()

	s0, err := Mesh2D(m)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	s1, err := Mesh2DSlow(m)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	bb := s0.BoundingBox()
	pSet := bb.RandomSet(nPoints)

	for _, p := range pSet {
		d0 := s0.Evaluate(p)
		d1 := s1.Evaluate(p)
		if !EqualFloat64(d0, d1, tolerance) {
			e := d0 - d1
			t.Errorf("%v fast %f slow %f error %f", p, d0, d1, e)
		}
	}
}

//-----------------------------------------------------------------------------

func Test_Mesh2D_Cache(t *testing.T) {

	m := getLines()

	s0, err := Mesh2D(m)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	s1 := Cache2D(s0)

	bb := s0.BoundingBox()
	pSet := bb.RandomSet(nPoints)

	for i := 0; i < 4; i++ {
		for _, p := range pSet {
			d0 := s0.Evaluate(p)
			d1 := s1.Evaluate(p)
			if !EqualFloat64(d0, d1, tolerance) {
				e := d0 - d1
				t.Errorf("%v no-cache %f cache %f error %f", p, d0, d1, e)
			}
		}
	}
}

//-----------------------------------------------------------------------------
