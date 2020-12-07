//-----------------------------------------------------------------------------
/*

PCB Standoffs, Mounting Pillars

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// StandoffParms defines the parameters for a board standoff pillar.
type StandoffParms struct {
	PillarHeight   float64
	PillarDiameter float64
	HoleDepth      float64 // > 0 is a hole, < 0 is a support stub
	HoleDiameter   float64
	NumberWebs     int // number of triangular gussets around the standoff base
	WebHeight      float64
	WebDiameter    float64
	WebWidth       float64
}

// pillarWeb returns a single pillar web
func pillarWeb(k *StandoffParms) sdf.SDF3 {
	w := sdf.NewPolygon()
	w.Add(0, 0)
	w.Add(0.5*k.WebDiameter, 0)
	w.Add(0, k.WebHeight)
	s := sdf.Extrude3D(sdf.Polygon2D(w.Vertices()), k.WebWidth)
	m := sdf.Translate3d(sdf.V3{0, 0, -0.5 * k.PillarHeight}).Mul(sdf.RotateX(sdf.DtoR(90.0)))
	return sdf.Transform3D(s, m)
}

// pillarWebs returns a set of pillar webs
func pillarWebs(k *StandoffParms) sdf.SDF3 {
	if k.NumberWebs == 0 {
		return nil
	}
	return sdf.RotateCopy3D(pillarWeb(k), k.NumberWebs)
}

// pillar returns a cylindrical pillar
func pillar(k *StandoffParms) sdf.SDF3 {
	return sdf.Cylinder3D(k.PillarHeight, 0.5*k.PillarDiameter, 0)
}

// pillarHole returns a pillar screw hole (or support stub)
func pillarHole(k *StandoffParms) sdf.SDF3 {
	if k.HoleDiameter == 0.0 || k.HoleDepth == 0.0 {
		return nil
	}
	s := sdf.Cylinder3D(math.Abs(k.HoleDepth), 0.5*k.HoleDiameter, 0)
	zOfs := 0.5 * (k.PillarHeight - k.HoleDepth)
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, 0, zOfs}))
}

// Standoff3D returns a single board standoff.
func Standoff3D(k *StandoffParms) sdf.SDF3 {
	s0 := sdf.Union3D(pillar(k), pillarWebs(k))
	if k.NumberWebs != 0 {
		// Cut off any part of the webs that protrude from the top of the pillar
		s0 = sdf.Intersect3D(s0, sdf.Cylinder3D(k.PillarHeight, k.WebDiameter, 0))
	}
	// Add the pillar hole/stub
	if k.HoleDepth >= 0.0 {
		// hole
		s0 = sdf.Difference3D(s0, pillarHole(k))
	} else {
		// support stub
		s0 = sdf.Union3D(s0, pillarHole(k))
	}
	return s0
}

//-----------------------------------------------------------------------------
