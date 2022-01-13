//-----------------------------------------------------------------------------
/*

Convert an SDF2 boundary to a set of line segments, using the SURREAL algorithm.

*/
//-----------------------------------------------------------------------------

package render

import (
	"fmt"
	"github.com/Yeicor/surreal/surreal2"
	"github.com/deadsy/sdfx/sdf"
)

// Surreal2 is a Render2 implementation based on the 2D SURREAL algorithm. Advantages:
// - Fast
// - Sharp edges
// - Planar simplification
type Surreal2 struct {
	*surreal2.Algorithm
}

func (s *Surreal2) Render(sdf2 sdf.SDF2, _ int, output chan<- *Line) {
	for _, line := range s.Algorithm.Run(sdf2) {
		output <- &Line{line[0], line[1]}
	}
}

func (s *Surreal2) Info(_ sdf.SDF2, _ int) string {
	return fmt.Sprintf("%+v", s.Algorithm)
}

// TODO: Add the Render3 implementation once it is properly implemented on github.com/Yeicor/surreal
