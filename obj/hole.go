//-----------------------------------------------------------------------------
/*

Holes

*/
//-----------------------------------------------------------------------------

package obj

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// CircleGrilleParms defines the parameters for a circular grille pattern.
type CircleGrilleParms struct {
	HoleDiameter      float64 // diameter of the grille holes
	GrilleDiameter    float64 // diameter of grille
	RadialSpacing     float64 // radial spacing of holes
	TangentialSpacing float64 // tangential spacing of holes
	Thickness         float64 // hole thickness (3d only)
}

// CircleGrille2D returns the SDF2 for a circular grille.
func CircleGrille2D(k *CircleGrilleParms) (sdf.SDF2, error) {
	var s []sdf.SDF2

	holeRadius := k.HoleDiameter * 0.5
	grilleRadius := k.GrilleDiameter * 0.5
	rSpacing := 1.0 + k.RadialSpacing
	tSpacing := 1.0 + k.TangentialSpacing

	// central base hole
	c, err := sdf.Circle2D(holeRadius)
	if err != nil {
		return nil, err
	}
	s = append(s, c)

	steps := int(math.Floor((grilleRadius / (rSpacing * k.HoleDiameter)) - 0.5))
	rDelta := grilleRadius / (float64(steps) + 0.5)
	theta := 0.0

	for i := 1; i <= steps; i++ {
		r := rDelta * float64(i)
		k := int(math.Floor((sdf.Tau * r) / (k.HoleDiameter * tSpacing)))
		dTheta := sdf.Tau / float64(k)
		c0 := sdf.Transform2D(c, sdf.Translate2d(v2.Vec{0, r}))
		for j := 0; j < k; j++ {
			c1 := sdf.Transform2D(c0, sdf.Rotate2d(theta+float64(j)*dTheta))
			s = append(s, c1)
		}
		theta += 0.5 * dTheta
	}
	return sdf.Union2D(s...), nil
}

// CircleGrille3D returns the SDF3 for a circular grille.
func CircleGrille3D(k *CircleGrilleParms) (sdf.SDF3, error) {
	s, err := CircleGrille2D(k)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, k.Thickness), nil
}

//-----------------------------------------------------------------------------

// CounterBoredHole3D returns the SDF3 for a counterbored hole.
func CounterBoredHole3D(
	l float64, // total length (includes counterbore)
	r float64, // hole radius
	cbRadius float64, // counter bore radius
	cbDepth float64, // counter bore depth
) (sdf.SDF3, error) {
	s0, err := sdf.Cylinder3D(l, r, 0)
	if err != nil {
		return nil, err
	}
	s1, err := sdf.Cylinder3D(cbDepth, cbRadius, 0)
	if err != nil {
		return nil, err
	}
	s1 = sdf.Transform3D(s1, sdf.Translate3d(v3.Vec{0, 0, (l - cbDepth) * 0.5}))
	return sdf.Union3D(s0, s1), nil
}

// ChamferedHole3D returns the SDF3 for a chamfered hole (45 degrees).
func ChamferedHole3D(
	l float64, // total length (includes chamfer)
	r float64, // hole radius
	chRadius float64, // chamfer radius
) (sdf.SDF3, error) {
	s0, err := sdf.Cylinder3D(l, r, 0)
	if err != nil {
		return nil, err
	}
	s1, err := sdf.Cone3D(chRadius, r, r+chRadius, 0)
	if err != nil {
		return nil, err
	}
	s1 = sdf.Transform3D(s1, sdf.Translate3d(v3.Vec{0, 0, (l - chRadius) * 0.5}))
	return sdf.Union3D(s0, s1), nil
}

// CounterSunkHole3D returns the SDF3 for a countersunk hole (45 degrees).
func CounterSunkHole3D(
	l float64, // total length
	r float64, // hole radius
) (sdf.SDF3, error) {
	return ChamferedHole3D(l, r, r)
}

//-----------------------------------------------------------------------------

// BoltCircle2D returns a 2D profile for a flange bolt circle.
func BoltCircle2D(
	holeRadius float64, // radius of bolt holes
	circleRadius float64, // radius of bolt circle
	numHoles int, // number of bolts
) (sdf.SDF2, error) {
	s, err := sdf.Circle2D(holeRadius)
	if err != nil {
		return nil, err
	}
	s = sdf.Transform2D(s, sdf.Translate2d(v2.Vec{circleRadius, 0}))
	s = sdf.RotateCopy2D(s, numHoles)
	return s, nil
}

// BoltCircle3D returns a 3D object for a flange bolt circle.
func BoltCircle3D(
	holeDepth float64, // depth of bolt holes
	holeRadius float64, // radius of bolt holes
	circleRadius float64, // radius of bolt circle
	numHoles int, // number of bolts
) (sdf.SDF3, error) {
	s, err := BoltCircle2D(holeRadius, circleRadius, numHoles)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, holeDepth), nil
}

//-----------------------------------------------------------------------------

// KeyedHoleParms defines the parameters for a keyed hole.
type KeyedHoleParms struct {
	Diameter  float64 // diameter of hole
	KeySize   float64 // key size / hole diameter, [0..1]
	NumKeys   int     // number of key flats (1 or 2)
	Thickness float64 // hole thickness (3d only)
}

// KeyedHole2D retuns a 2D object for a round hole with a flat section.
func KeyedHole2D(k *KeyedHoleParms) (sdf.SDF2, error) {
	s, err := sdf.Circle2D(k.Diameter * 0.5)
	if err != nil {
		return nil, err
	}
	if k.NumKeys == 1 {
		yOfs := k.Diameter * (k.KeySize - 0.5)
		return sdf.Cut2D(s, v2.Vec{0, yOfs}, v2.Vec{1, 0}), nil
	} else if k.NumKeys == 2 {
		yOfs := 0.5 * k.Diameter * k.KeySize
		s = sdf.Cut2D(s, v2.Vec{0, yOfs}, v2.Vec{1, 0})
		return sdf.Cut2D(s, v2.Vec{0, -yOfs}, v2.Vec{-1, 0}), nil
	}
	return nil, sdf.ErrMsg("NumKeys must be 1 or 2")
}

// KeyedHole3D retuns a 3D object for a round hole with a flat section.
func KeyedHole3D(k *KeyedHoleParms) (sdf.SDF3, error) {
	s, err := KeyedHole2D(k)
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(s, k.Thickness), nil
}

//-----------------------------------------------------------------------------
