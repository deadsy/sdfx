//-----------------------------------------------------------------------------
/*

Servo Models

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

func servos() error {

	names := []string{
		"nano",
		"submicro",
		"micro",
		"mini",
		"standard",
		"large",
		"giant",
	}

	var s sdf.SDF3
	yOfs := 0.0

	for _, n := range names {
		k, err := obj.ServoLookup(n)
		if err != nil {
			return err
		}

		yOfs += 0.5*k.Body.Y + 10.0

		servo, err := obj.Servo3D(k)
		if err != nil {
			return err
		}
		servo = sdf.Transform3D(servo, sdf.Translate3d(v3.Vec{0, yOfs, 20}))

		outline2, err := obj.Servo2D(k, -1)
		if err != nil {
			return err
		}
		outline := sdf.Extrude3D(outline2, 5)
		outline = sdf.Transform3D(outline, sdf.Translate3d(v3.Vec{0, yOfs, 0}))

		s = sdf.Union3D(s, servo, outline)
		yOfs += 0.5 * k.Body.Y
	}

	render.ToSTL(s, "servos.stl", render.NewMarchingCubesOctree(300))
	return nil
}

//-----------------------------------------------------------------------------

func main() {
	err := servos()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

//-----------------------------------------------------------------------------
