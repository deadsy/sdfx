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

// InvoluteGearParms defines the parameters for an involute gear.
type InvoluteGearParms struct {
	NumberTeeth   int     // number of gear teeth
	Module        float64 // pitch circle diameter / number of gear teeth
	PressureAngle float64 // gear pressure angle (radians)
	Backlash      float64 // backlash expressed as per-tooth distance at pitch circumference
	Clearance     float64 // additional root clearance
	RingWidth     float64 // width of ring wall (from root circle)
	Facets        int     // number of facets for involute flank
}

// InvoluteGear returns an 2D polygon for an involute gear.
func InvoluteGear(k *InvoluteGearParms) (sdf.SDF2, error) {

	if k.NumberTeeth <= 0 {
		return nil, sdf.ErrMsg("NumberTeeth <= 0")
	}
	if k.Module <= 0 {
		return nil, sdf.ErrMsg("Module <= 0")
	}
	if k.PressureAngle <= 0 {
		return nil, sdf.ErrMsg("PressureAngle <= 0")
	}
	if k.Backlash < 0 {
		return nil, sdf.ErrMsg("Backlash <= 0")
	}
	if k.Clearance < 0 {
		return nil, sdf.ErrMsg("Clearance < 0")
	}
	if k.RingWidth < 0 {
		return nil, sdf.ErrMsg("RingWidth < 0")
	}
	if k.Facets <= 0 {
		return nil, sdf.ErrMsg("Facets <= 0")
	}

	// pitch radius
	pitchRadius := float64(k.NumberTeeth) * k.Module * 0.5

	// base circle radius
	baseRadius := pitchRadius * math.Cos(k.PressureAngle)

	// addendum: radial distance from pitch circle to outside circle
	addendum := k.Module * 1.0
	// dedendum: radial distance from pitch circle to root circle
	dedendum := addendum + k.Clearance

	outerRadius := pitchRadius + addendum
	rootRadius := pitchRadius - dedendum
	ringRadius := rootRadius - k.RingWidth

	tooth := involuteGearTooth(
		k.NumberTeeth,
		k.Module,
		rootRadius,
		baseRadius,
		outerRadius,
		k.Backlash,
		k.Facets,
	)

	gear := sdf.RotateCopy2D(tooth, k.NumberTeeth)
	root := sdf.Circle2D(rootRadius)
	ring := sdf.Circle2D(ringRadius)

	return sdf.Difference2D(sdf.Union2D(gear, root), ring), nil
}

//-----------------------------------------------------------------------------
