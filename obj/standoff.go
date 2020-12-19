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
func pillarWeb(k *StandoffParms) (sdf.SDF3, error) {
	w := sdf.NewPolygon()
	w.Add(0, 0)
	w.Add(0.5*k.WebDiameter, 0)
	w.Add(0, k.WebHeight)
	p, err := sdf.Polygon2D(w.Vertices())
	if err != nil {
		return nil, err
	}
	s := sdf.Extrude3D(p, k.WebWidth)
	m := sdf.Translate3d(sdf.V3{0, 0, -0.5 * k.PillarHeight}).Mul(sdf.RotateX(sdf.DtoR(90.0)))
	return sdf.Transform3D(s, m), nil
}

// pillarWebs returns a set of pillar webs
func pillarWebs(k *StandoffParms) (sdf.SDF3, error) {
	if k.NumberWebs == 0 {
		// no webs
		return nil, nil
	}
	web, err := pillarWeb(k)
	if err != nil {
		return nil, err
	}
	return sdf.RotateCopy3D(web, k.NumberWebs), nil
}

// pillar returns a cylindrical pillar
func pillar(k *StandoffParms) (sdf.SDF3, error) {
	return sdf.Cylinder3D(k.PillarHeight, 0.5*k.PillarDiameter, 0)
}

// pillarHole returns a pillar screw hole (or support stub)
func pillarHole(k *StandoffParms) (sdf.SDF3, error) {
	if k.HoleDiameter == 0.0 || k.HoleDepth == 0.0 {
		// no hole
		return nil, nil
	}
	s, err := sdf.Cylinder3D(math.Abs(k.HoleDepth), 0.5*k.HoleDiameter, 0)
	if err != nil {
		return nil, err
	}
	zOfs := 0.5 * (k.PillarHeight - k.HoleDepth)
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, 0, zOfs})), nil
}

// Standoff3D returns a single board standoff.
func Standoff3D(k *StandoffParms) (sdf.SDF3, error) {
	pillar, err := pillar(k)
	if err != nil {
		return nil, err
	}
	webs, err := pillarWebs(k)
	if err != nil {
		return nil, err
	}
	s := sdf.Union3D(pillar, webs)
	if k.NumberWebs != 0 {
		// Cut off any part of the webs that protrude from the top of the pillar
		cut, err := sdf.Cylinder3D(k.PillarHeight, k.WebDiameter, 0)
		if err != nil {
			return nil, err
		}
		s = sdf.Intersect3D(s, cut)
	}
	// Add the pillar hole/stub
	hole, err := pillarHole(k)
	if err != nil {
		return nil, err
	}
	if k.HoleDepth >= 0.0 {
		s = sdf.Difference3D(s, hole)
	} else {
		// support stub
		s = sdf.Union3D(s, hole)
	}
	return s, nil
}

//-----------------------------------------------------------------------------
