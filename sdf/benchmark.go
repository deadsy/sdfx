package sdf

import "fmt"

const N_EVALS = 10 // 100000

// Benchmark evaluation speed for an SDF2.
// Return evaluations per second.
func BenchmarkSDF2(s SDF2) float64 {

	// sample over a region larger than the bounding box
	center := s.BoundingBox().Center()
	size := s.BoundingBox().Size().MulScalar(1.5)

	// build an array of random sample points
	var points [N_EVALS]V2
	for i, _ := range points {
		points[i] = center.Add(size.Random())
	}

	fmt.Printf("%+v\n", points)

	return 0

}
