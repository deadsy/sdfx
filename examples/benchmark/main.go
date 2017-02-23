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
	s := Circle2D(5)
	fmt.Printf("circle SDF2 %s\n", fmt_eps(BenchmarkSDF2(s)))

	s = NewFlatFlankCam(30, 20, 5)
	fmt.Printf("cam1 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s)))

	s = NewThreeArcCam(30, 20, 5, 200)
	fmt.Printf("cam2 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s)))

	s = NewPolySDF2(Nagon(6, 10.0))
	fmt.Printf("poly6 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s)))

	s = NewPolySDF2(Nagon(12, 10.0))
	fmt.Printf("poly12 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s)))

	s = NewPolySDF2(Nagon(18, 10.0))
	fmt.Printf("poly18 SDF2 %s\n", fmt_eps(BenchmarkSDF2(s)))
}
