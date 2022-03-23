//-----------------------------------------------------------------------------
/*

Servo Models

See: https://www.servocity.com/servos/

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
	ShaftOffset float64 // x-offset of drive shaft (from mounting hole center to shaft)
	ShaftLength float64
	ShaftRadius float64
	HoleRadius  float64
}

type servoDatabase map[string]ServoParms

var servoDB = initServoLookup()

func (m servoDatabase) Add(name string, k *ServoParms) {
	m[name] = *k
}

// initServoLookup adds a collection of named servos to the database.
func initServoLookup() servoDatabase {
	m := make(servoDatabase)

	k := ServoParms{
		Body:        sdf.V3{20, 8.7, 20.3},
		Mount:       sdf.V3{28, 8.7, 1},
		Hole:        sdf.V2{24, 0},
		MountOffset: 12,
		ShaftOffset: 6.4,
		ShaftLength: 2.8,
		ShaftRadius: 2,
		HoleRadius:  1,
	}
	m.Add("hitec_hs_40", &k)
	m.Add("nano", &k)

	k = ServoParms{
		Body:        sdf.V3{22.6, 11.5, 24.5},
		Mount:       sdf.V3{32.6, 10.4, 1},
		Hole:        sdf.V2{28.5, 0},
		MountOffset: 16.6,
		ShaftOffset: 9,
		ShaftLength: 2.5,
		ShaftRadius: 2,
		HoleRadius:  1,
	}
	m.Add("hitec_hs_55", &k)
	m.Add("submicro", &k)

	//k = ServoParms{}
	//m.Add("hitec_hs_81", &k)
	//m.Add("micro", &k)

	//k = ServoParms{}
	//m.Add("hitec_hs_225bb", &k)
	//m.Add("mini", &k)

	k = ServoParms{
		Body:        sdf.V3{40.2, 20.2, 38.3},
		Mount:       sdf.V3{52.9, 20.2, 2.5},
		Hole:        sdf.V2{47.6, 10.1},
		MountOffset: 26.5,
		ShaftOffset: 13.85,
		ShaftLength: 3.5,
		ShaftRadius: 2,
		HoleRadius:  2.15,
	}
	m.Add("standard", &k)
	m.Add("hitec_hs_311", &k)

	k = ServoParms{
		Body:        sdf.V3{40, 20, 41.5},
		Mount:       sdf.V3{54.38, 20, 1},
		Hole:        sdf.V2{49.5, 10},
		MountOffset: 27.76,
		ShaftOffset: 30,
		ShaftLength: 4,
		ShaftRadius: 2,
		HoleRadius:  2,
	}
	m.Add("annimos_20kg", &k)

	k = ServoParms{
		Body:        sdf.V3{65.9, 29.9, 59.3},
		Mount:       sdf.V3{82.9, 29.9, 4},
		Hole:        sdf.V2{74.9, 17.8},
		MountOffset: 42,
		ShaftOffset: 18.9,
		ShaftLength: 5.4,
		ShaftRadius: 2,
		HoleRadius:  2.8,
	}
	m.Add("hitec_hs_805bb", &k)
	m.Add("large", &k)

	//k = ServoParms{}
	//m.Add("hitec_hs_1005sgt", &k)
	//m.Add("giant", &k)

	return m
}

// ServoLookup returns the parameters for a named servo.
func ServoLookup(name string) (*ServoParms, error) {
	k, ok := servoDB[name]
	if !ok {
		return nil, fmt.Errorf("servo \"%s\" not found", name)
	}
	return &k, nil
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
	xOfs := 0.5*k.Hole.X - k.ShaftOffset
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
	xOfs = 0.5*k.Hole.X - k.ShaftOffset
	zOfs = 0.5 * k.Body.Z
	s = sdf.Transform3D(s, sdf.Translate3d(sdf.V3{xOfs, 0, zOfs}))

	return s, nil
}

//-----------------------------------------------------------------------------
