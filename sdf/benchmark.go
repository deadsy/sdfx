package sdf

import "time"

const N_EVALS = 10000000

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

	start := time.Now()
	for _, p := range points {
		s.Evaluate(p)
	}
	elapsed := time.Since(start)

	return float64(N_EVALS) * float64(time.Second) / float64(elapsed)
}
