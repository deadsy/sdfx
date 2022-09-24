package main

import (
	"log"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	s2d, err := sdf.Circle2D(5)
	if err != nil {
		log.Fatal(err)
	}
	sdf.BenchmarkSDF2("circle SDF2", s2d)

	s2d, err = sdf.FlatFlankCam2D(30, 20, 5)
	if err != nil {
		log.Fatal(err)
	}
	sdf.BenchmarkSDF2("cam1 SDF2", s2d)

	s2d, err = sdf.ThreeArcCam2D(30, 20, 5, 200)
	if err != nil {
		log.Fatal(err)
	}
	sdf.BenchmarkSDF2("cam2 SDF2", s2d)

	s2d, err = sdf.Polygon2D(sdf.Nagon(6, 10.0))
	if err != nil {
		log.Fatal(err)
	}
	sdf.BenchmarkSDF2("poly6 SDF2", s2d)

	s2d, err = sdf.Polygon2D(sdf.Nagon(12, 10.0))
	if err != nil {
		log.Fatal(err)
	}
	sdf.BenchmarkSDF2("poly12 SDF2", s2d)

	s3d, err := sdf.Box3D(v3.Vec{10, 20, 30}, 1)
	if err != nil {
		log.Fatal(err)
	}
	sdf.BenchmarkSDF3("box SDF3", s3d)
}
