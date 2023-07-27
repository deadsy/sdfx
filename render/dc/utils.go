package dc

import (
	"github.com/deadsy/sdfx/render"
)

//-----------------------------------------------------------------------------
// UTILITIES/MISC
//-----------------------------------------------------------------------------

func dcFlip(t *render.Triangle3) *render.Triangle3 {
	t[1], t[2] = t[2], t[1]
	return t
}

func dcMaxI(i int, i2 int) int {
	if i >= i2 {
		return i
	}
	return i2
}
