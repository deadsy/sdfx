package dc

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------
// UTILITIES/MISC
//-----------------------------------------------------------------------------

func dcFlip(t *sdf.Triangle3) *sdf.Triangle3 {
	t[1], t[2] = t[2], t[1]
	return t
}

func dcMaxI(i int, i2 int) int {
	if i >= i2 {
		return i
	}
	return i2
}
