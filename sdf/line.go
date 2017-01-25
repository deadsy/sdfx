//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import ()

//-----------------------------------------------------------------------------
// 2D Line Segment

type Line2 struct {
	A, B V2 // line start/end points
	V    V2 // vector in line direction x = a + tv, where t = [0,1]
	N    V2 // normal to line
}

func NewLine2(a, b V2) Line2 {
	l := Line2{}
	l.A = a
	l.B = b
	ba := b.Sub(a)
	v := ba.Normalize()
	l.N = V2{v.Y, -v.X}
	l.V = ba.MulScalar(1 / ba.Dot(ba))
	return l
}

// return the distance to the line, +ve implies same side as line normal
func (l *Line2) Distance(p V2) float64 {
	pa := p.Sub(l.A)
	t := pa.Dot(l.V)  // t-parameter of projection onto line
	dn := pa.Dot(l.N) // distance normal to line
	var d float64
	if t < 0 {
		d = pa.Length()
		if dn < 0 {
			d = -d
		}
	} else if t > 1 {
		d = p.Sub(l.B).Length()
		if dn < 0 {
			d = -d
		}
	} else {
		d = dn
	}
	return d
}

//-----------------------------------------------------------------------------
