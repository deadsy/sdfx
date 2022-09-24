//-----------------------------------------------------------------------------
/*

Keyways in Shafts

*/
//-----------------------------------------------------------------------------

package obj

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

// KeywayParameters defines the parameters for a keyway and shaft.
type KeywayParameters struct {
	ShaftRadius float64 // shaft radius
	KeyRadius   float64 // shaft center to bottom/top of key
	KeyWidth    float64 // key width
	ShaftLength float64 // shaft length (3d only)
}

// Keyway2D returns the 2d profile of a shaft and keyway.
func Keyway2D(k *KeywayParameters) (sdf.SDF2, error) {
	if k.ShaftRadius <= 0 {
		return nil, sdf.ErrMsg("k.ShaftRadius <= 0")
	}
	if k.KeyRadius < 0 {
		return nil, sdf.ErrMsg("k.KeyRadius < 0")
	}
	if k.KeyWidth < 0 {
		return nil, sdf.ErrMsg("k.KeyWidth < 0")
	}
	shaft, err := sdf.Circle2D(k.ShaftRadius)
	if err != nil {
		return nil, err
	}
	var s sdf.SDF2
	if k.KeyRadius < k.ShaftRadius {
		// The key is cut into the shaft (shaft profile)
		l := k.ShaftRadius - k.KeyRadius
		key := sdf.Box2D(v2.Vec{l, k.KeyWidth}, 0)
		key = sdf.Transform2D(key, sdf.Translate2d(v2.Vec{k.ShaftRadius - l*0.5, 0}))
		s = sdf.Difference2D(shaft, key)
	} else {
		// The key is proud of the shaft (bore profile)
		key := sdf.Box2D(v2.Vec{k.KeyRadius, k.KeyWidth}, 0)
		key = sdf.Transform2D(key, sdf.Translate2d(v2.Vec{k.KeyRadius * 0.5, 0}))
		s = sdf.Union2D(shaft, key)
	}
	return s, nil
}

// Keyway3D returns a shaft and keyway.
func Keyway3D(k *KeywayParameters) (sdf.SDF3, error) {
	if k.ShaftLength <= 0 {
		return nil, sdf.ErrMsg("k.ShaftLength <= 0")
	}
	s, err := Keyway2D(k)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, k.ShaftLength), nil
}

//-----------------------------------------------------------------------------
