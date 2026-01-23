//-----------------------------------------------------------------------------
/*

Basic Polygon Shapes

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// IsocelesTrapezoid2D
func IsocelesTrapezoid2D(base0, base1, height float64) (sdf.SDF2, error) {
	b0 := 0.5 * base0
	b1 := 0.5 * base1
	h := 0.5 * height
	p0 := v2.Vec{b0, -h}
	p1 := v2.Vec{b1, h}
	p2 := v2.Vec{-b1, h}
	p3 := v2.Vec{-b0, -h}
	return sdf.Polygon2D([]v2.Vec{p0, p1, p2, p3})
}

//-----------------------------------------------------------------------------

// IsocelesTriangle2D
func IsocelesTriangle2D(base, height float64) (sdf.SDF2, error) {
	b := 0.5 * base
	h := 0.5 * height
	p0 := v2.Vec{b, -h}
	p1 := v2.Vec{0, h}
	p2 := v2.Vec{-b, -h}
	return sdf.Polygon2D([]v2.Vec{p0, p1, p2})
}

//-----------------------------------------------------------------------------
