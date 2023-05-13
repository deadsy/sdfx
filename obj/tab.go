//-----------------------------------------------------------------------------
/*

Tabs for Connecting Objects.

Tab objects are used to align and connect upper and lower objects.
Generally the upper/lower boundary is the XY plane.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// Tab is the interface to a tab object.
type Tab interface {
	Body(upper bool, m sdf.M44) sdf.SDF3     // + to connected body
	Envelope(upper bool, m sdf.M44) sdf.SDF3 // - from connected body
}

// AddTabs to an upper or lower object.
func AddTabs(s sdf.SDF3, tab Tab, upper bool, mset []sdf.M44) sdf.SDF3 {
	bSet := make([]sdf.SDF3, len(mset))
	eSet := make([]sdf.SDF3, len(mset))
	for i := range mset {
		bSet[i] = tab.Body(upper, mset[i])
		eSet[i] = tab.Envelope(upper, mset[i])
	}
	body := sdf.Union3D(bSet...)
	envelope := sdf.Union3D(eSet...)
	return sdf.Union3D(sdf.Difference3D(s, envelope), body)
}

//-----------------------------------------------------------------------------
// simple straight tabs

// StraightTab contains the straight tab parameters.
type StraightTab struct {
	size      v3.Vec  // size of tab
	clearance float64 // clearance between male and female elements
}

// NewStraightTab returns a new straight tab object.
func NewStraightTab(size v3.Vec, clearance float64) (Tab, error) {
	return &StraightTab{
		size:      size,
		clearance: clearance,
	}, nil
}

// Body returns the upper/lower body of a straight tab.
func (t *StraightTab) Body(upper bool, m sdf.M44) sdf.SDF3 {
	if upper {
		return nil
	}
	s, _ := sdf.Box3D(t.size, 0)
	return sdf.Transform3D(s, m.Mul(sdf.Translate3d(v3.Vec{0, 0, 0.5 * t.size.Z})))
}

// Envelope returns the upper/lower envelope of a straight tab.
func (t *StraightTab) Envelope(upper bool, m sdf.M44) sdf.SDF3 {
	if upper {
		size := t.size.Add(v3.Vec{2.0 * t.clearance, 2.0 * t.clearance, t.clearance})
		s, _ := sdf.Box3D(size, 0)
		return sdf.Transform3D(s, m.Mul(sdf.Translate3d(v3.Vec{0, 0, 0.5 * size.Z})))
	}
	return nil
}

//-----------------------------------------------------------------------------
// 45 degreee angled tabs

// AngleTab contains the angle tab parameters.
type AngleTab struct {
	size      v3.Vec  // size of tab
	clearance float64 // clearance between male and female elements
}

// NewAngleTab returns a new angle tab object.
func NewAngleTab(size v3.Vec, clearance float64) (Tab, error) {
	return &AngleTab{
		size:      size,
		clearance: clearance,
	}, nil
}

// Body returns the upper/lower body of an angle tab.
func (t *AngleTab) Body(upper bool, m sdf.M44) sdf.SDF3 {
	if upper {
		return nil
	}
	size := t.size
	s, _ := sdf.Box3D(size, 0)
	xCut := size.X*0.5 - size.Z*0.5
	s = sdf.Cut3D(s, v3.Vec{xCut, 0, 0}, v3.Vec{-1, 0, 1})
	s = sdf.Cut3D(s, v3.Vec{-xCut, 0, 0}, v3.Vec{1, 0, -1})
	return sdf.Transform3D(s, m.Mul(sdf.Translate3d(v3.Vec{0, 0, 0.5 * size.Z})))
}

// Envelope returns the upper/lower envelope of an angle tab.
func (t *AngleTab) Envelope(upper bool, m sdf.M44) sdf.SDF3 {
	if upper {
		size := t.size.Add(v3.Vec{2.0 * t.clearance, 2.0 * t.clearance, t.clearance})
		s, _ := sdf.Box3D(size, 0)
		xCut := size.X*0.5 - size.Z*0.5
		s = sdf.Cut3D(s, v3.Vec{xCut, 0, 0}, v3.Vec{-1, 0, 1})
		s = sdf.Cut3D(s, v3.Vec{-xCut, 0, 0}, v3.Vec{1, 0, -1})
		return sdf.Transform3D(s, m.Mul(sdf.Translate3d(v3.Vec{0, 0, 0.5 * size.Z})))
	}
	return nil
}

//-----------------------------------------------------------------------------
// screw pillar tab

// ScrewTab contains the screw tab parameters.
type ScrewTab struct {
	Length     float64 // length of pillar
	Radius     float64 // radius of pillar
	Round      bool    // round the bottom of the pillar
	HoleUpper  float64 // length of upper hole
	HoleLower  float64 // length of lower hole
	HoleRadius float64 // radius of hole
}

func (t *ScrewTab) screwBody() sdf.SDF3 {
	var round float64
	if t.Round {
		round = t.Radius
	}
	s, _ := sdf.Cylinder3D(2.0*t.Length, t.Radius, round)
	return sdf.Cut3D(s, v3.Vec{0, 0, 0}, v3.Vec{0, 0, -1})
}

func (t *ScrewTab) screwHole() sdf.SDF3 {
	l := t.HoleUpper + t.HoleLower
	zOfs := t.HoleUpper - (0.5 * l)
	s, _ := sdf.Cylinder3D(l, t.HoleRadius, 0)
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, zOfs}))
}

// NewScrewTab returns a new screw tab object.
func NewScrewTab(k *ScrewTab) (Tab, error) {
	return k, nil
}

// Body returns the upper/lower body of an angle tab.
func (t *ScrewTab) Body(upper bool, m sdf.M44) sdf.SDF3 {
	if upper {
		return nil
	}
	return sdf.Transform3D(sdf.Difference3D(t.screwBody(), t.screwHole()), m)
}

// Envelope returns the upper/lower envelope of an angle tab.
func (t *ScrewTab) Envelope(upper bool, m sdf.M44) sdf.SDF3 {
	if upper {
		return sdf.Transform3D(t.screwHole(), m)
	}
	return sdf.Transform3D(t.screwBody(), m)
}

//-----------------------------------------------------------------------------
