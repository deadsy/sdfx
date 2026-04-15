package sdf

import (
	"math"
	"testing"

	v3 "github.com/deadsy/sdfx/vec/v3"
)

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
