//-----------------------------------------------------------------------------
/*

drain covers

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func drain12() (sdf.SDF3, error) {
	k := &obj.DrainCoverParms{
		WallDiameter:   11.75 * sdf.MillimetresPerInch,
		WallHeight:     1.0 * sdf.MillimetresPerInch,
		WallThickness:  0.25 * sdf.MillimetresPerInch,
		WallDraft:      sdf.DtoR(2.0),
		OuterWidth:     0.75 * sdf.MillimetresPerInch,
		InnerWidth:     0.5 * sdf.MillimetresPerInch,
		CoverThickness: 0.25 * sdf.MillimetresPerInch,
		GrateNumber:    10,
		GrateWidth:     0.6,
		CrossBarWidth:  1.0 * sdf.MillimetresPerInch,
		GrateDraft:     sdf.DtoR(10.0),
	}
	return obj.DrainCover(k)
}

//-----------------------------------------------------------------------------

func main() {
	s, err := drain12()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "drain12.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
