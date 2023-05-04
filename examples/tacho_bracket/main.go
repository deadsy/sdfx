//-----------------------------------------------------------------------------
/*

Tachometer Holding Bracket

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

const tachoRadius = 0.5 * 3.5 * sdf.MillimetresPerInch
const bracketHeight = 10.0
const bracketWidth = 10.0
const tabWidth = 30.0
const tabLength = 30.0

//-----------------------------------------------------------------------------

func tachoBracket() (sdf.SDF3, error) {

	const outerRadius = tachoRadius + bracketWidth

	s0, err := sdf.Circle2D(outerRadius)
	if err != nil {
		return nil, err
	}

	s1, err := sdf.Circle2D(tachoRadius)
	if err != nil {
		return nil, err
	}

	s2 := sdf.Box2D(v2.Vec{2.0 * (outerRadius + tabLength), tabWidth}, 0.1*(tabWidth+tabLength))

	s3 := sdf.Difference2D(sdf.Union2D(s0, s2), s1)

	return sdf.Extrude3D(s3, bracketHeight), nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := tachoBracket()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "tacho.stl", render.NewMarchingCubesOctree(300))

}

//-----------------------------------------------------------------------------
