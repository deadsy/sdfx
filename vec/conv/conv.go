//-----------------------------------------------------------------------------
/*

Vector Conversions

*/
//-----------------------------------------------------------------------------

package conv

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/p2"
	v2 "github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v2i"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/deadsy/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------
// V2i to X

// V2iToV2 converts a 2D integer vector to a float vector.
func V2iToV2(a v2i.Vec) v2.Vec {
	return v2.Vec{float64(a.X), float64(a.Y)}
}

//-----------------------------------------------------------------------------
// V3i to X

// V3iToV3 converts a 3D integer vector to a float vector.
func V3iToV3(a v3i.Vec) v3.Vec {
	return v3.Vec{float64(a.X), float64(a.Y), float64(a.Z)}
}

//-----------------------------------------------------------------------------
// V2 to X

// V2ToP2 converts a cartesian to a polar coordinate.
func V2ToP2(a v2.Vec) p2.Vec {
	return p2.Vec{a.Length(), math.Atan2(a.Y, a.X)}
}

// V2ToV3 converts a 2D vector to a 3D vector with a specified Z value.
func V2ToV3(a v2.Vec, z float64) v3.Vec {
	return v3.Vec{a.X, a.Y, z}
}

//-----------------------------------------------------------------------------
// V3 to X

// V3ToSDF converts a 3D vector to the legacy sdf type.
func V3ToSDF(a v3.Vec) sdf.V3 {
	return sdf.V3{a.X, a.Y, a.Z}
}

//-----------------------------------------------------------------------------
// P2 to X

// P2ToV2 converts a polar to a cartesian coordinate.
func P2ToV2(a p2.Vec) v2.Vec {
	return v2.Vec{a.R * math.Cos(a.Theta), a.R * math.Sin(a.Theta)}
}

//-----------------------------------------------------------------------------
