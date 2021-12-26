package dc

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------
// UTILITIES/MISC
//-----------------------------------------------------------------------------

func dcFlip(t0 *render.Triangle3) *render.Triangle3 {
	t0.V[1], t0.V[2] = t0.V[2], t0.V[1]
	return t0
}

func dcCompGet(v3 sdf.V3, i int) float64 {
	switch i {
	case 0:
		return v3.X
	case 1:
		return v3.Y
	default:
		return v3.Z
	}
}

func dcCompSet(v3 *sdf.V3, i int, val float64) {
	switch i {
	case 0:
		v3.X = val
	case 1:
		v3.Y = val
	default:
		v3.Z = val
	}
}

func dcMaxI(i int, i2 int) int {
	if i >= i2 {
		return i
	}
	return i2
}
