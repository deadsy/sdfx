//-----------------------------------------------------------------------------
/*

Report benchmarking results for evaluations on SDF2/SDF3 objects.

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"time"
)

//-----------------------------------------------------------------------------

const nEvals = 10000000

//-----------------------------------------------------------------------------

// fmtEPS returns a string with a formatted evaluations per second.
func fmtEPS(eps float64) string {
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

// BenchmarkSDF2 reports the evaluation speed for an SDF2.
func BenchmarkSDF2(description string, s SDF2) {
	// sample over a region larger than the bounding box
	box := NewBox2(s.BoundingBox().Center(), s.BoundingBox().Size().MulScalar(1.2))
	points := box.RandomSet(nEvals)

	start := time.Now()
	for _, p := range points {
		s.Evaluate(p)
	}
	elapsed := time.Since(start)

	eps := float64(nEvals) * float64(time.Second) / float64(elapsed)
	fmt.Printf("%s %s\n", description, fmtEPS(eps))
}

//-----------------------------------------------------------------------------

// BenchmarkSDF3 reports the evaluation speed for an SDF3.
func BenchmarkSDF3(description string, s SDF3) {
	// sample over a region larger than the bounding box
	box := NewBox3(s.BoundingBox().Center(), s.BoundingBox().Size().MulScalar(1.2))
	points := box.RandomSet(nEvals)

	start := time.Now()
	for _, p := range points {
		s.Evaluate(p)
	}
	elapsed := time.Since(start)

	eps := float64(nEvals) * float64(time.Second) / float64(elapsed)
	fmt.Printf("%s %s\n", description, fmtEPS(eps))
}

//-----------------------------------------------------------------------------
