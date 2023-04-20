//-----------------------------------------------------------------------------
/*

Drone Parts

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// material shrinkage

const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const wallThickness = 3.0

func tab(upper bool) (sdf.SDF3, sdf.SDF3, error) {
	k := obj.TabParms{
		Size:      v3.Vec{4.0 * wallThickness, 0.7 * wallThickness, wallThickness}, // size of tab
		Clearance: 0.1,                                                             // clearance between male and female elements
		Angled:    false,                                                           // tab at 45 degrees (in x-direction)
	}
	return obj.Tab(&k, upper)
}

func tabbox(mode bool) (sdf.SDF3, error) {

	round := 0.5 * wallThickness
	oSize := v3.Vec{40, 40, 30}
	iSize := oSize.SubScalar(2.0 * wallThickness)

	outer, err := sdf.Box3D(oSize, round)
	if err != nil {
		return nil, err
	}
	inner, err := sdf.Box3D(iSize, round)
	if err != nil {
		return nil, err
	}

	box := sdf.Difference3D(outer, inner)
	lidHeight := oSize.Z * 0.25

	if mode == true {
		// upper
		s := sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, 1})
		body, env, err := tab(true)
		if err != nil {
			return nil, err
		}
		s = sdf.Union3D(sdf.Difference3D(s, env), body)
		return s, nil
	}

	// lower
	s := sdf.Cut3D(box, v3.Vec{0, 0, lidHeight}, v3.Vec{0, 0, -1})
	body, env, err := tab(false)
	if err != nil {
		return nil, err
	}
	s = sdf.Union3D(sdf.Difference3D(s, env), body)
	return s, nil
}

//-----------------------------------------------------------------------------

func main() {

	upper, err := tabbox(true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(upper, shrink), "upper.stl", render.NewMarchingCubesOctree(300))

	lower, err := tabbox(false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(lower, shrink), "lower.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
