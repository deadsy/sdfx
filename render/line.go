//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package render

import v2 "github.com/deadsy/sdfx/vec/v2"

//-----------------------------------------------------------------------------

// Line is a 2d line segment defined with 2 points.
type Line [2]v2.Vec

// Degenerate returns true if the line is degenerate.
func (l Line) Degenerate(tolerance float64) bool {
	// check for identical vertices
	return l[0].Equals(l[1], tolerance)
}

//-----------------------------------------------------------------------------
