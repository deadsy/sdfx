//-----------------------------------------------------------------------------
/*

Arrows and Coordinate Axes

*/
//-----------------------------------------------------------------------------

package obj

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func arrowStyle3D(style byte, size [2]float64, tail bool) (sdf.SDF3, error) {
	// nothing
	if style == '.' {
		return nil, nil
	}
	// ball
	if style == 'b' {
		return sdf.Sphere3D(size[1])
	}
	// cone
	if style == 'c' {
		cone, err := sdf.Cone3D(size[0], size[1], 0, 0)
		if err != nil {
			return nil, err
		}
		cone = sdf.Transform3D(cone, sdf.Translate3d(v3.Vec{0, 0, size[0] * 0.5}))
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
		head = sdf.Transform3D(head, sdf.Translate3d(v3.Vec{0, 0, zOfs}))
	}
	if tail != nil {
		tail = sdf.Transform3D(tail, sdf.Translate3d(v3.Vec{0, 0, -zOfs}))
	}
	return sdf.Union3D(axis, head, tail), nil
}

//-----------------------------------------------------------------------------

func axis3D(a, b, r float64) (sdf.SDF3, error) {
	if a == b {
		return nil, nil
	}
	// ensure a < b
	if a > b {
		a, b = b, a
	}
	style := "cc"
	if a == 0 {
		style = "c."
	}
	if b == 0 {
		style = ".c"
	}
	l0 := b - a
	r0 := r
	r1 := r * 1.5
	l1 := r * 3
	k := ArrowParms{
		Axis:  [2]float64{l0, r0},
		Head:  [2]float64{l1, r1},
		Tail:  [2]float64{l1, r1},
		Style: style,
	}
	s, err := Arrow3D(&k)
	if err != nil {
		return nil, err
	}
	ofs := (a + b) * 0.5
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, ofs})), nil
}

// Axes3D returns a set of axes for a 1, 2 or 3d coordinate systems.
func Axes3D(p0, p1 v3.Vec) (sdf.SDF3, error) {
	// work out the common axis radius
	r := p0.Sub(p1).Abs().MaxComponent() * 0.025
	// x-axis
	x, err := axis3D(p0.X, p1.X, r)
	if err != nil {
		return nil, err
	}
	if x != nil {
		x = sdf.Transform3D(x, sdf.RotateY(sdf.DtoR(90)))
	}
	// y-axis
	y, err := axis3D(p0.Y, p1.Y, r)
	if err != nil {
		return nil, err
	}
	if y != nil {
		y = sdf.Transform3D(y, sdf.RotateX(sdf.DtoR(-90)))
	}
	// z-axis
	z, err := axis3D(p0.Z, p1.Z, r)
	if err != nil {
		return nil, err
	}
	return sdf.Union3D(x, y, z), nil
}

//-----------------------------------------------------------------------------

// DirectedArrow3D returns an arrow between points head/tail.
func DirectedArrow3D(k *ArrowParms, head, tail v3.Vec) (sdf.SDF3, error) {
	v := head.Sub(tail)
	l := v.Length()
	k.Axis[0] = l
	arrow, err := Arrow3D(k)
	if err != nil {
		return nil, err
	}
	// position the arrow
	ofs := head.Add(tail).MulScalar(0.5)
	m := sdf.Translate3d(ofs).Mul(sdf.RotateToVector(v3.Vec{0, 0, 1}, v))
	return sdf.Transform3D(arrow, m), nil
}

//-----------------------------------------------------------------------------
