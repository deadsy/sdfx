package dc

import (
	"github.com/deadsy/sdfx/render"
)

//-----------------------------------------------------------------------------
// UTILITIES/MISC
//-----------------------------------------------------------------------------

func dcFlip(t0 *render.Triangle3) *render.Triangle3 {
	t0.V[1], t0.V[2] = t0.V[2], t0.V[1]
	return t0
}

func dcMaxI(i int, i2 int) int {
	if i >= i2 {
		return i
	}
	return i2
}
