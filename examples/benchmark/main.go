package main

import (
	"log"

	. "github.com/deadsy/sdfx/sdf"
)

func main() {
	s2d, err := Circle2D(5)
	if err != nil {
		log.Fatal(err)
	}
	BenchmarkSDF2("circle SDF2", s2d)

	s2d, err = FlatFlankCam2D(30, 20, 5)
	if err != nil {
		log.Fatal(err)
	}
	BenchmarkSDF2("cam1 SDF2", s2d)

	s2d, err = ThreeArcCam2D(30, 20, 5, 200)
	if err != nil {
		log.Fatal(err)
	}
	BenchmarkSDF2("cam2 SDF2", s2d)

	s2d, err = Polygon2D(Nagon(6, 10.0))
	if err != nil {
		log.Fatal(err)
	}
	BenchmarkSDF2("poly6 SDF2", s2d)

	s2d, err = Polygon2D(Nagon(12, 10.0))
	if err != nil {
		log.Fatal(err)
	}
	BenchmarkSDF2("poly12 SDF2", s2d)

	s3d, err := Box3D(V3{10, 20, 30}, 1)
	if err != nil {
		log.Fatal(err)
	}
	BenchmarkSDF3("box SDF3", s3d)
}
