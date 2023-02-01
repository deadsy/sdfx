//-----------------------------------------------------------------------------
/*

Drain Covers

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// DrainCoverParms defines a grated drain pipe cover.
type DrainCoverParms struct {
	WallDiameter   float64 // outer diameter of wall
	WallHeight     float64 // height of wall
	WallThickness  float64 // thickness of wall
	WallDraft      float64 // draft angle of walls
	OuterWidth     float64 // extra width beyond the wall
	InnerWidth     float64 // width between inner wall and grate field
	CoverThickness float64 // thickness of the drain cover
	GrateNumber    int     // number of grate holes
	GrateWidth     float64 // width of the grate hole (0..1)
	CrossBarWidth  float64
	GrateDraft     float64 // draft anfle of grates (radians)
}

//-----------------------------------------------------------------------------

func drainCoverBody(k *DrainCoverParms) (sdf.SDF3, error) {

	// x drafts
	dx0 := 0.5 * k.CoverThickness * math.Tan(k.WallDraft)
	dx1 := 0.5 * k.WallHeight * math.Tan(k.WallDraft)

	// x radii
	r0 := 0.5 * k.WallDiameter
	r1 := r0 + k.OuterWidth
	r2 := r0 - k.WallThickness

	// y thicknesses
	t0 := k.CoverThickness
	t1 := t0 + k.WallHeight

	// build the 2d profile
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(r1+dx0, 0)
	p.Add(r1-dx0, t0)
	p.Add(r0+dx1, t0)
	p.Add(r0-dx1, t1)
	p.Add(r2+dx1, t1)
	p.Add(r2-dx1, t0)
	p.Add(0, t0)

	s, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}

	// return the revolved profile
	return sdf.Revolve3D(s)
}

// DrainCover returns a grated drain pipe cover.
func DrainCover(k *DrainCoverParms) (sdf.SDF3, error) {
	return drainCoverBody(k)
}

//-----------------------------------------------------------------------------
