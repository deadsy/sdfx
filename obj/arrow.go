//-----------------------------------------------------------------------------
/*

Arrows

*/
//-----------------------------------------------------------------------------

package obj

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func arrowStyle3D(style byte, size [2]float64, tail bool) (sdf.SDF3, error) {
	// nothing
	if style == '.' {
		return nil, nil
	}
	// ball
	if style == 'b' {
		return sdf.Sphere3D(size[1]), nil
	}
	// cone
	if style == 'c' {
		cone, err := sdf.Cone3D(size[0], size[1], 0, 0)
		if err != nil {
			return nil, err
		}
		cone = sdf.Transform3D(cone, sdf.Translate3d(sdf.V3{0, 0, size[0] * 0.5}))
		if tail {
      // flip it
			cone = sdf.Transform3D(cone, sdf.RotateX(sdf.Pi))
		}
		return cone, nil
	}
	return nil, sdf.ErrMsg(fmt.Sprintf("bad style character '%c'", style))
}

//-----------------------------------------------------------------------------

// ArrowParms defines the parameters for an arrow.
type ArrowParms struct {
	Axis  [2]float64 // length/radius of arrow axis
	Head  [2]float64 // length/radius of arrow head
	Tail  [2]float64 // length/radius of arrow tail
	Style string     // head, tail "c" = cone, "b" == ball (else nothing)
}

// Arrow3D returns an arrow.
func Arrow3D(k *ArrowParms) (sdf.SDF3, error) {
	if k == nil {
		return nil, sdf.ErrMsg("k == nil")
	}

	// decode the head/tail style
	var head, tail sdf.SDF3
	var err error
	switch len(k.Style) {
	case 0:
		// no style
	case 1:
		head, err = arrowStyle3D(k.Style[0], k.Head, false)
		if err != nil {
			return nil, err
		}
	case 2:
		head, err = arrowStyle3D(k.Style[0], k.Head, false)
		if err != nil {
			return nil, err
		}
		tail, err = arrowStyle3D(k.Style[1], k.Tail, true)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdf.ErrMsg("style string is too long")
	}

	// build the axis
	axis, err := sdf.Capsule3D(k.Axis[0]+(2.0*k.Axis[1]), k.Axis[1])
	if err != nil {
		return nil, err
	}

	zOfs := k.Axis[0] * 0.5
	if head != nil {
		head = sdf.Transform3D(head, sdf.Translate3d(sdf.V3{0, 0, zOfs}))
	}
	if tail != nil {
		tail = sdf.Transform3D(tail, sdf.Translate3d(sdf.V3{0, 0, -zOfs}))
	}
	return sdf.Union3D(axis, head, tail), nil
}

//-----------------------------------------------------------------------------
