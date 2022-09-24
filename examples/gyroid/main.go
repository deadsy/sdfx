//-----------------------------------------------------------------------------
/*

Gyroid Cubes

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func gyroidCube() (sdf.SDF3, error) {

	l := 100.0   // cube side
	k := l * 0.1 // 10 cycles per side

	gyroid, err := sdf.Gyroid3D(v3.Vec{k, k, k})
	if err != nil {
		return nil, err
	}

	box, err := sdf.Box3D(v3.Vec{l, l, l}, 0)
	if err != nil {
		return nil, err
	}

	return sdf.Intersect3D(box, gyroid), nil
}

//-----------------------------------------------------------------------------

func gyroidSurface() (sdf.SDF3, error) {

	l := 60.0    // cube side
	k := l * 0.5 // 2 cycles per side

	s, err := sdf.Gyroid3D(v3.Vec{k, k, k})
	if err != nil {
		return nil, err
	}

	s, err = sdf.Shell3D(s, k*0.025)
	if err != nil {
		return nil, err
	}

	box, err := sdf.Box3D(v3.Vec{l, l, l}, 0)
	if err != nil {
		return nil, err
	}

	s = sdf.Intersect3D(box, s)

	// remove the isolated bits on the cube corners
	sphere, err := sdf.Sphere3D(k * 0.15)
	if err != nil {
		return nil, err
	}
	d := l * 0.5
	s0 := sdf.Transform3D(sphere, sdf.Translate3d(v3.Vec{d, d, d}))
	s1 := sdf.Transform3D(sphere, sdf.Translate3d(v3.Vec{-d, -d, -d}))

	return sdf.Difference3D(s, sdf.Union3D(s0, s1)), nil
}

//-----------------------------------------------------------------------------

func main() {

	s0, err := gyroidCube()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTLSlow(s0, 300, "gyroid_cube.stl")

	s1, err := gyroidSurface()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTLSlow(s1, 150, "gyroid_surface.stl")
}

//-----------------------------------------------------------------------------
