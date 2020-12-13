//-----------------------------------------------------------------------------
/*

Involute Gears

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// return the involute coordinate for a given angle
func involuteXY(
	r float64, // base radius
	theta float64, // involute angle
) sdf.V2 {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return sdf.V2{
		r * (c + theta*s),
		r * (s - theta*c),
	}
}

// return the involute angle for a given radial distance
func involuteTheta(
	r float64, // base radius
	d float64, // involute radial distance
) float64 {
	x := d / r
	return math.Sqrt(x*x - 1)
}

//-----------------------------------------------------------------------------

// involuteGearTooth returns a 2D profile for a single involute tooth.
func involuteGearTooth(
	numberTeeth int, // number of gear teeth
	gearModule float64, // pitch circle diameter / number of gear teeth
	rootRadius float64, // radius at tooth root
	baseRadius float64, // radius at the base of the involute
	outerRadius float64, // radius at the outside of the tooth
	backlash float64, // backlash expressed as units of pitch circumference
	facets int, // number of facets for involute flank
) sdf.SDF2 {

	pitchRadius := float64(numberTeeth) * gearModule / 2.0

	// work out the angular extent of the tooth on the base radius
	pitchPoint := involuteXY(baseRadius, involuteTheta(baseRadius, pitchRadius))
	faceAngle := math.Atan2(pitchPoint.Y, pitchPoint.X)
	backlashAngle := backlash / (2.0 * pitchRadius)
	centerAngle := sdf.Pi/(2.0*float64(numberTeeth)) + faceAngle - backlashAngle

	// work out the angles over which the involute will be used
	startAngle := involuteTheta(baseRadius, math.Max(baseRadius, rootRadius))
	stopAngle := involuteTheta(baseRadius, outerRadius)
	dtheta := (stopAngle - startAngle) / float64(facets)

	v := make([]sdf.V2, 2*(facets+1)+1)

	// lower tooth face
	m := sdf.Rotate(-centerAngle)
	angle := startAngle
	for i := 0; i <= facets; i++ {
		v[i] = m.MulPosition(involuteXY(baseRadius, angle))
		angle += dtheta
	}

	// upper tooth face (mirror the lower point)
	for i := 0; i <= facets; i++ {
		p := v[facets-i]
		v[facets+1+i] = sdf.V2{p.X, -p.Y}
	}

	// add the origin to make the polygon a tooth wedge
	v[2*(facets+1)] = sdf.V2{0, 0}

	return sdf.Polygon2D(v)
}

//-----------------------------------------------------------------------------

// InvoluteGear returns an 2D polygon for an involute gear.
func InvoluteGear(
	numberTeeth int, // number of gear teeth
	gearModule float64, // pitch circle diameter / number of gear teeth
	pressureAngle float64, // gear pressure angle (radians)
	backlash float64, // backlash expressed as per-tooth distance at pitch circumference
	clearance float64, // additional root clearance
	ringWidth float64, // width of ring wall (from root circle)
	facets int, // number of facets for involute flank
) (sdf.SDF2, error) {

	// pitch radius
	pitchRadius := float64(numberTeeth) * gearModule / 2.0

	// base circle radius
	baseRadius := pitchRadius * math.Cos(pressureAngle)

	// addendum: radial distance from pitch circle to outside circle
	addendum := gearModule * 1.0
	// dedendum: radial distance from pitch circle to root circle
	dedendum := addendum + clearance

	outerRadius := pitchRadius + addendum
	rootRadius := pitchRadius - dedendum
	ringRadius := rootRadius - ringWidth

	tooth := involuteGearTooth(
		numberTeeth,
		gearModule,
		rootRadius,
		baseRadius,
		outerRadius,
		backlash,
		facets,
	)

	gear := sdf.RotateCopy2D(tooth, numberTeeth)
	root := sdf.Circle2D(rootRadius)
	ring := sdf.Circle2D(ringRadius)

	return sdf.Difference2D(sdf.Union2D(gear, root), ring), nil
}

//-----------------------------------------------------------------------------
