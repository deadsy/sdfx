package sdf

import (
	"math"
	"testing"

	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Test_Buttress_WrapContinuity verifies that asymmetric buttress thread
// profiles produce the same SDF at x = +pitch/2 and x = -pitch/2 — the
// SawTooth wrap boundary. A discontinuity here would cause octree marching
// cubes to skip cubes straddling the boundary and produce mesh holes; the
// 2-period polygon design is what makes the SDF wrap-continuous.
func Test_Buttress_WrapContinuity(t *testing.T) {
	cases := []struct {
		name   string
		make   func(r, p float64) (SDF2, error)
		radius float64
		pitch  float64
	}{
		{"ANSI", ANSIButtressThread, 5, 2},
		{"Plastic", PlasticButtressThread, 5, 2},
		{"ANSI_steep", ANSIButtressThread, 2.5, 3},
		{"Plastic_steep", PlasticButtressThread, 2.5, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := c.make(c.radius, c.pitch)
			if err != nil {
				t.Fatal(err)
			}
			hp := c.pitch / 2
			// Sample a vertical line of points at the wrap boundary, comparing
			// SDF(+pitch/2, y) vs SDF(-pitch/2, y) at radii covering the
			// thread region. Tolerance is generous enough to absorb numerical
			// noise from polygon edge intersections, but tight enough to
			// catch the ~0.4mm gap that the periodic union previously masked.
			tol := 1e-6
			for _, y := range []float64{c.radius - 0.1, c.radius - 0.3, c.radius - 0.5, c.radius} {
				dR := s.Evaluate(v2.Vec{X: hp, Y: y})
				dL := s.Evaluate(v2.Vec{X: -hp, Y: y})
				if math.Abs(dR-dL) > tol {
					t.Errorf("y=%.3f: SDF discontinuous at wrap: dR(%+.4f)=%+.6f dL(%+.4f)=%+.6f delta=%.6f",
						y, hp, dR, -hp, dL, dR-dL)
				}
			}
		})
	}
}

// Test_Screw3D_RejectsDiscontinuousProfile verifies that Screw3D refuses to
// build a screw from a thread profile whose SDF is discontinuous at the
// SawTooth wrap boundary (x=±pitch/2). Such profiles silently produce holes
// in octree-rendered meshes, so we fail fast at construction. A user-defined
// asymmetric profile with flat scaffolding extensions (which is exactly the
// shape the original buttress profiles had) is the canonical example.
func Test_Screw3D_RejectsDiscontinuousProfile(t *testing.T) {
	// Asymmetric profile shaped like a buttress: notch at x=+pitch/4 with a
	// near-vertical right flank and a 45° left flank, with flat-crest
	// scaffolding extending from the period boundaries to ±pitch. SDF at
	// x=+pitch/2 sees the near-vertical flank just to the left; SDF at
	// x=-pitch/2 sees nothing nearby — so the wrap is discontinuous.
	pitch := 2.0
	radius := 5.0
	depth := 0.6
	tp := NewPolygon()
	tp.Add(pitch, 0)
	tp.Add(pitch, radius)
	tp.Add(pitch/2+0.05, radius)  // gentle (right) flank top
	tp.Add(pitch/4, radius-depth) // valley root, asymmetric position
	tp.Add(pitch/4-depth, radius) // 45° (left) flank top
	tp.Add(-pitch, radius)
	tp.Add(-pitch, 0)
	bad, err := Polygon2D(tp.Vertices())
	if err != nil {
		t.Fatal(err)
	}

	_, err = Screw3D(bad, 10, 0, pitch, 1)
	if err == nil {
		t.Fatalf("Screw3D accepted a discontinuous thread profile; expected an error")
	}
	t.Logf("got expected rejection: %s", err)
}

// Test_Screw_TaperSlope verifies that Evaluate uses tan(taper) not atan(taper)
// for the taper slope computation. At 30° the error is 16.5%.
func Test_Screw_TaperSlope(t *testing.T) {
	taperAngle := DtoR(30)
	slope := math.Tan(taperAngle)
	slopeAtan := math.Atan(taperAngle)

	thread, err := ISOThread(5, 2, true)
	if err != nil {
		t.Fatal(err)
	}
	screw, err := Screw3D(thread, 20, taperAngle, 2, 1)
	if err != nil {
		t.Fatal(err)
	}

	// At z=-8 (pitch multiple), SawTooth maps to 0 (thread crest center).
	// The tan-predicted crest radius is larger than the atan-predicted one.
	testZ := -8.0
	rCrestTan := 5 + math.Abs(testZ)*slope
	rCrestAtan := 5 + math.Abs(testZ)*slopeAtan

	// A point between the two predicted radii should be inside if tan is used.
	testR := (rCrestTan + rCrestAtan) / 2
	d := screw.Evaluate(v3.Vec{X: testR, Y: 0, Z: testZ})
	if d > 0 {
		t.Errorf("point at r=%.3f (between tan=%.3f and atan=%.3f) is outside (d=%+.6f); expected inside with tan slope",
			testR, rCrestTan, rCrestAtan, d)
	}

	// Verify surface is near the tan-predicted radius, not the atan-predicted one.
	dInside := screw.Evaluate(v3.Vec{X: rCrestTan - 0.5, Y: 0, Z: testZ})
	dOutside := screw.Evaluate(v3.Vec{X: rCrestTan + 0.5, Y: 0, Z: testZ})
	if dInside >= 0 {
		t.Errorf("r_crest-0.5: expected inside (d<0), got d=%+.4f", dInside)
	}
	if dOutside <= 0 {
		t.Errorf("r_crest+0.5: expected outside (d>0), got d=%+.4f", dOutside)
	}

	t.Logf("taper=30° slope_tan=%.4f slope_atan=%.4f (%.1f%% error)",
		slope, slopeAtan, 100*(slope-slopeAtan)/slope)
}
