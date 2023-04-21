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
	s = sdf.Transform3D(s, m.Mul(sdf.Translate3d(v3.Vec{0, 0, 0.5 * t.size.Z})))
	return s
}

// Envelope returns the upper/lower envelope of a straight tab.
func (t *StraightTab) Envelope(upper bool, m sdf.M44) sdf.SDF3 {
	if upper {
		s, _ := sdf.Box3D(v3.Vec{t.size.X + t.clearance, t.size.Y + t.clearance, t.size.Z}, 0)
		s = sdf.Transform3D(s, m.Mul(sdf.Translate3d(v3.Vec{0, 0, 0.5 * t.size.Z})))
		return s
	}
	return nil
}

//-----------------------------------------------------------------------------
