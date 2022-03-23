//-----------------------------------------------------------------------------
/*

Servo Models

*/
//-----------------------------------------------------------------------------

package obj

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// ServoParms stores the parameters that define the servo.
type ServoParms struct {
	Body        sdf.V3  // body size
	Mount       sdf.V3  // mounting lugs size
	Hole        sdf.V2  // hole layout
	MountOffset float64 // z-offset of mounting lugs (from base of servo)
	ShaftOffset float64 // x-offset of drive shaft (from side of servo)
	ShaftLength float64
	ShaftRadius float64
	HoleRadius  float64
}

type servoDatabase map[string]*ServoParms

var servoDB = initServoLookup()

func (m servoDatabase) Add(name string, k *ServoParms) {
	m[name] = k
}

// initServoLookup adds a collection of named servos to the database.
func initServoLookup() servoDatabase {
	m := make(servoDatabase)

	k := ServoParms{
		Body:        sdf.V3{40, 20, 41.5},
		Mount:       sdf.V3{54.38, 20, 2.6},
		Hole:        sdf.V2{49.5, 10},
		MountOffset: 27.76,
		ShaftOffset: 30,
		ShaftLength: 4,
		ShaftRadius: 2,
		HoleRadius:  2,
	}
	m.Add("annimos_20kg", &k)

	return m
}

// ServoLookup returns the parameters for a named servo.
func ServoLookup(name string) (*ServoParms, error) {
	k, ok := servoDB[name]
	if !ok {
		return nil, fmt.Errorf("servo \"%s\" not found", name)
	}
	return k, nil
}

//-----------------------------------------------------------------------------

// Servo3D returns a 3D model for a servo.
func Servo3D(k *ServoParms) (sdf.SDF3, error) {

	// servo body
	body, err := sdf.Box3D(k.Body, 0.06*k.Body.Y)
	if err != nil {
		return nil, err
	}

	// mounting lugs
	m := sdf.Box2D(sdf.V2{k.Mount.X, k.Mount.Y}, 0.1*k.Mount.Y)
	mount := sdf.Extrude3D(m, k.Mount.Z)
	zOfs := k.MountOffset - 0.5*(k.Body.Z-k.Mount.Z)
	mount = sdf.Transform3D(mount, sdf.Translate3d(sdf.V3{0, 0, zOfs}))

	// output shaft
	shaft, err := sdf.Cylinder3D(k.ShaftLength, k.ShaftRadius, 0)
	if err != nil {
		return nil, err
	}
	xOfs := k.ShaftOffset - 0.5*k.Body.X
	zOfs = 0.5 * (k.Body.Z + k.ShaftLength)
	shaft = sdf.Transform3D(shaft, sdf.Translate3d(sdf.V3{-xOfs, 0, zOfs}))

	// holes
	hole, err := sdf.Cylinder3D(k.Body.Z, k.HoleRadius, 0)
	if err != nil {
		return nil, err
	}
	xOfs = 0.5 * k.Hole.X
	yOfs := 0.5 * k.Hole.Y
	holes := sdf.Multi3D(hole, []sdf.V3{{xOfs, yOfs, 0}, {-xOfs, yOfs, 0}, {xOfs, -yOfs, 0}, {-xOfs, -yOfs, 0}})

	s := sdf.Difference3D(sdf.Union3D(body, mount, shaft), holes)

	// position the shaft on the z-axis and the bottom of the servo at z=0
	xOfs = k.ShaftOffset - 0.5*k.Body.X
	zOfs = 0.5 * k.Body.Z
	s = sdf.Transform3D(s, sdf.Translate3d(sdf.V3{xOfs, 0, zOfs}))

	return s, nil
}

//-----------------------------------------------------------------------------
