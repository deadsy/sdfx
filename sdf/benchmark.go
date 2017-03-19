//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"time"
)

//-----------------------------------------------------------------------------

const N_EVALS = 10000000

//-----------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------

// Benchmark evaluation speed for an SDF2.
func BenchmarkSDF2(description string, s SDF2) {

	// sample over a region larger than the bounding box
	center := s.BoundingBox().Center()
	size := s.BoundingBox().Size().MulScalar(1.2)

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

	eps := float64(N_EVALS) * float64(time.Second) / float64(elapsed)
	fmt.Printf("%s %s\n", description, fmt_eps(eps))
}

//-----------------------------------------------------------------------------

// Benchmark evaluation speed for an SDF3.

func BenchmarkSDF3(description string, s SDF3) {

	// sample over a region larger than the bounding box
	center := s.BoundingBox().Center()
	size := s.BoundingBox().Size().MulScalar(1.2)

	// build an array of random sample points
	var points [N_EVALS]V3
	for i, _ := range points {
		points[i] = center.Add(size.Random())
	}

	start := time.Now()
	for _, p := range points {
		s.Evaluate(p)
	}
	elapsed := time.Since(start)

	eps := float64(N_EVALS) * float64(time.Second) / float64(elapsed)
	fmt.Printf("%s %s\n", description, fmt_eps(eps))
}

//-----------------------------------------------------------------------------
