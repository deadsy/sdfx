//-----------------------------------------------------------------------------
/*

Gyroid Cubes
Gyroid Teapot

*/
//-----------------------------------------------------------------------------

package main

import (
	"errors"
	"log"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func gyroidCube() (sdf.SDF3, error) {

	l := 100.0   // cube side
	k := l * 0.2 // 5 cycles per side

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

// gyroidSurface - suitable for printing
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

func gyroidTeapot(cyclesPerSide int) (sdf.SDF3, error) {
	if cyclesPerSide < 1 {
		return nil, errors.New("cycles per side should not be <= 0")
	}

	// create the SDF from the STL mesh
	teapot, err := obj.ImportSTL("../../files/teapot.stl", 20, 3, 5)
	if err != nil {
		return nil, err
	}

	min := teapot.BoundingBox().Min
	max := teapot.BoundingBox().Max

	kX := (max.X - min.X) / float64(cyclesPerSide)
	kY := (max.Y - min.Y) / float64(cyclesPerSide)
	kZ := (max.Z - min.Z) / float64(cyclesPerSide)

	gyroid, err := sdf.Gyroid3D(v3.Vec{kX, kY, kZ})
	if err != nil {
		return nil, err
	}

	return sdf.Intersect3D(teapot, gyroid), nil
}

//-----------------------------------------------------------------------------

func main() {

	s0, err := gyroidCube()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s0, "gyroid_cube.stl", render.NewMarchingCubesUniform(300))

	s1, err := gyroidSurface()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s1, "gyroid_surface.stl", render.NewMarchingCubesUniform(150))

	s2, err := gyroidTeapot(10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	// Rendering to STL is fine:
	render.ToSTL(s2, "gyroid_teapot.stl", render.NewMarchingCubesUniform(200))

	// Rendering to triangles and then saving the STL, the output file is not as expected:
	m2 := render.ToTriangles(s2, render.NewMarchingCubesUniform(200))
	render.SaveSTL("gyroid_teapot_mesh.stl", m2)
}

//-----------------------------------------------------------------------------
