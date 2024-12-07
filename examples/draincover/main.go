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

func vent2() (sdf.SDF3, error) {
	k := &obj.DrainCoverParms{
		WallDiameter:   1.9 * sdf.MillimetresPerInch,
		WallHeight:     0.75 * sdf.MillimetresPerInch,
		WallThickness:  0.125 * sdf.MillimetresPerInch,
		WallDraft:      0,
		OuterWidth:     0.2 * sdf.MillimetresPerInch,
		InnerWidth:     0.18 * sdf.MillimetresPerInch,
		CoverThickness: 0.125 * sdf.MillimetresPerInch,
		GrateNumber:    8,
		GrateWidth:     1.1,
		GrateDraft:     0,
		CrossBarWidth:  0,
		CrossBarWeb:    false,
	}
	return obj.DrainCover(k)
}

func drain4() (sdf.SDF3, error) {
	k := &obj.DrainCoverParms{
		WallDiameter:   3.9 * sdf.MillimetresPerInch,
		WallHeight:     0.8 * sdf.MillimetresPerInch,
		WallThickness:  0.2 * sdf.MillimetresPerInch,
		WallDraft:      sdf.DtoR(2.0),
		OuterWidth:     0.4 * sdf.MillimetresPerInch,
		InnerWidth:     0.3 * sdf.MillimetresPerInch,
		CoverThickness: 0.2 * sdf.MillimetresPerInch,
		GrateNumber:    8,
		GrateWidth:     1.1,
		GrateDraft:     sdf.DtoR(8.0),
		CrossBarWidth:  0.8,
		CrossBarWeb:    false,
	}
	return obj.DrainCover(k)
}

func drain6() (sdf.SDF3, error) {
	k := &obj.DrainCoverParms{
		WallDiameter:   5.8 * sdf.MillimetresPerInch,
		WallHeight:     0.8 * sdf.MillimetresPerInch,
		WallThickness:  0.2 * sdf.MillimetresPerInch,
		WallDraft:      sdf.DtoR(2.0),
		OuterWidth:     0.4 * sdf.MillimetresPerInch,
		InnerWidth:     0.3 * sdf.MillimetresPerInch,
		CoverThickness: 0.3 * sdf.MillimetresPerInch,
		GrateNumber:    9,
		GrateWidth:     1.0,
		GrateDraft:     sdf.DtoR(8.0),
		CrossBarWidth:  1.8,
		CrossBarWeb:    true,
	}
	return obj.DrainCover(k)
}

func drain12() (sdf.SDF3, error) {
	k := &obj.DrainCoverParms{
		WallDiameter:   11.8 * sdf.MillimetresPerInch,
		WallHeight:     1.0 * sdf.MillimetresPerInch,
		WallThickness:  0.3 * sdf.MillimetresPerInch,
		WallDraft:      sdf.DtoR(2.0),
		OuterWidth:     0.8 * sdf.MillimetresPerInch,
		InnerWidth:     0.5 * sdf.MillimetresPerInch,
		CoverThickness: 0.3 * sdf.MillimetresPerInch,
		GrateNumber:    10,
		GrateWidth:     1.0,
		GrateDraft:     sdf.DtoR(8.0),
		CrossBarWidth:  1.5,
		CrossBarWeb:    true,
	}
	return obj.DrainCover(k)
}

//-----------------------------------------------------------------------------

func main() {

	s, err := vent2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "vent2.stl", render.NewMarchingCubesOctree(300))

	s, err = drain4()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "drain4.stl", render.NewMarchingCubesOctree(300))

	s, err = drain6()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "drain6.stl", render.NewMarchingCubesOctree(300))

	s, err = drain12()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "drain12.stl", render.NewMarchingCubesOctree(400))
}

//-----------------------------------------------------------------------------
