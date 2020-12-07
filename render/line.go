//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package render

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// Line is a 2d line segment defined with 2 points.
type Line [2]sdf.V2

// Degenerate returns true if the line is degenerate.
func (l Line) Degenerate(tolerance float64) bool {
	// check for identical vertices
	return l[0].Equals(l[1], tolerance)
}

//-----------------------------------------------------------------------------
