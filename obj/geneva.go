//-----------------------------------------------------------------------------
/*

Geneva Drive

See: https://en.wikipedia.org/wiki/Geneva_drive

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// GenevaParms specfies the geneva drive parameters.
type GenevaParms struct {
	NumSectors     int     // number of sectors in the driven wheel
	CenterDistance float64 // center to center distance of driver/driven wheels
	DriverRadius   float64 // radius of lock portion of driver wheel
	DrivenRadius   float64 // radius of driven wheel
	PinRadius      float64 // radius of driver pin
	Clearance      float64 // pin/slot and wheel/wheel clearance
}

// Geneva2D makes 2d profiles for the driver/driven wheels of a geneva drive.
func Geneva2D(k *GenevaParms) (sdf.SDF2, sdf.SDF2, error) {

	if k.NumSectors < 2 {
		return nil, nil, sdf.ErrMsg("invalid number of sectors, must be > 2")
	}
	if k.CenterDistance <= 0 ||
		k.DrivenRadius <= 0 ||
		k.DriverRadius <= 0 ||
		k.PinRadius <= 0 {
		return nil, nil, sdf.ErrMsg("invalid dimensions, must be > 0")
	}
	if k.Clearance < 0 {
		return nil, nil, sdf.ErrMsg("invalid clearance, must be >= 0")
	}
	if k.CenterDistance > k.DrivenRadius+k.DriverRadius {
		return nil, nil, sdf.ErrMsg("center distance is too large")
	}

	// work out the pin offset from the center of the driver wheel
	theta := sdf.Tau / (2.0 * float64(k.NumSectors))
	d := k.CenterDistance
	r := k.DrivenRadius
	pinOffset := math.Sqrt((d * d) + (r * r) - (2 * d * r * math.Cos(theta)))

	// driven wheel
	sDriven, err := sdf.Circle2D(k.DrivenRadius - k.Clearance)
	if err != nil {
		return nil, nil, err
	}
	// cutouts for the driver wheel
	s, err := sdf.Circle2D(k.DriverRadius + k.Clearance)
	if err != nil {
		return nil, nil, err
	}
	s = sdf.Transform2D(s, sdf.Translate2d(v2.Vec{k.CenterDistance, 0}))
	s = sdf.RotateCopy2D(s, k.NumSectors)
	sDriven = sdf.Difference2D(sDriven, s)
	// cutouts for the pin slots
	slotLength := pinOffset + k.DrivenRadius - k.CenterDistance
	s = sdf.Line2D(2*slotLength, k.PinRadius+k.Clearance)
	s = sdf.Transform2D(s, sdf.Translate2d(v2.Vec{k.DrivenRadius, 0}))
	s = sdf.RotateCopy2D(s, k.NumSectors)
	s = sdf.Transform2D(s, sdf.Rotate2d(theta))
	sDriven = sdf.Difference2D(sDriven, s)

	// driver wheel
	sDriver, err := sdf.Circle2D(k.DriverRadius - k.Clearance)
	if err != nil {
		return nil, nil, err
	}
	// cutout for the driven wheel
	s, err = sdf.Circle2D(k.DrivenRadius + k.Clearance)
	if err != nil {
		return nil, nil, err
	}
	s = sdf.Transform2D(s, sdf.Translate2d(v2.Vec{k.CenterDistance, 0}))
	sDriver = sdf.Difference2D(sDriver, s)
	// driver pin
	s, err = sdf.Circle2D(k.PinRadius)
	if err != nil {
		return nil, nil, err
	}
	s = sdf.Transform2D(s, sdf.Translate2d(v2.Vec{pinOffset, 0}))
	sDriver = sdf.Union2D(sDriver, s)

	return sDriver, sDriven, nil
}

//-----------------------------------------------------------------------------
