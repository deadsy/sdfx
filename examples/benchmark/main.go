package main

import (
	"fmt"

	. "github.com/deadsy/sdfx/sdf"
)

func fmt_eps(eps float64) string {
	if eps > 1000000000.0 {
		return fmt.Sprintf("%.2f G evals/sec", eps/1000000000.0)
	} else if eps > 1000000.0 {
		return fmt.Sprintf("%.2f M evals/sec", eps/1000000.0)
	} else if eps > 1000.0 {
		return fmt.Sprintf("%.2f K evals/sec", eps/1000.0)
	}
	return fmt.Sprintf("%.2f evals/sec", eps)
}

func main() {
	s2d := Circle2D(5)
	fmt.Printf("circle SDF2 %s\n", fmt_eps(BenchmarkSDF2(s2d)))

	s2d = FlatFlankCam2D(30, 20, 5)
	fmt.Printf("cam1 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s2d)))

	s2d = ThreeArcCam2D(30, 20, 5, 200)
	fmt.Printf("cam2 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s2d)))

	s2d = Polygon2D(Nagon(6, 10.0))
	fmt.Printf("poly6 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s2d)))

	s2d = Polygon2D(Nagon(12, 10.0))
	fmt.Printf("poly12 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s2d)))

	s3d := Box3D(V3{10, 20, 30}, 1)
	fmt.Printf("box SDF3 %s\n", fmt_eps(BenchmarkSDF3(s3d)))
}
