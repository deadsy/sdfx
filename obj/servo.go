//-----------------------------------------------------------------------------
/*

Servo Models

See: https://www.servocity.com/servos/

Note: Servos fall into several well known size categories. In general you could
design for this nominal size but you need to check the final fit against the specific
servo being used. There is dimensional variance within and across manufacturers.

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
		ShaftRadius: 1.4,
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
		ShaftRadius: 1.25,
		HoleRadius:  0.95,
	}
	m.Add("hitec_hs_55", &k)
	m.Add("submicro", &k)

	k = ServoParms{
		Body:        sdf.V3{29.1, 13, 30.4},
		Mount:       sdf.V3{40, 12, 2},
		Hole:        sdf.V2{35.6, 0},
		MountOffset: 19,
		ShaftOffset: 9.8,
		ShaftLength: 3.8,
		ShaftRadius: 1.9,
		HoleRadius:  2.25,
	}
	m.Add("hitec_hs_85bb", &k)
	m.Add("micro", &k)

	k = ServoParms{
		Body:        sdf.V3{32.3, 16.8, 33},
		Mount:       sdf.V3{44.3, 16, 2.2},
		Hole:        sdf.V2{39.6, 7.9},
		MountOffset: 23.5,
		ShaftOffset: 12.2,
		ShaftLength: 3.3,
		ShaftRadius: 1.65,
		HoleRadius:  2.25,
	}
	m.Add("hitec_hs_225bb", &k)
	m.Add("mini", &k)

	k = ServoParms{
		Body:        sdf.V3{40.2, 20.2, 38.3},
		Mount:       sdf.V3{52.9, 20.2, 2.5},
		Hole:        sdf.V2{47.6, 10.1},
		MountOffset: 26.5,
		ShaftOffset: 13.85,
		ShaftLength: 3.5,
		ShaftRadius: 1.75,
		HoleRadius:  2.15,
	}
	m.Add("hitec_hs_311", &k)
	m.Add("standard", &k)

	k = ServoParms{
		Body:        sdf.V3{40, 20, 41.5},
		Mount:       sdf.V3{54.2, 18.5, 3},
		Hole:        sdf.V2{49.5, 10},
		MountOffset: 28,
		ShaftOffset: 14.75,
		ShaftLength: 4.2,
		ShaftRadius: 2.1,
		HoleRadius:  2.15,
	}
	m.Add("annimos_ds3218", &k)

	k = ServoParms{
		Body:        sdf.V3{65.9, 29.9, 59.3},
		Mount:       sdf.V3{82.9, 29.9, 4},
		Hole:        sdf.V2{74.9, 17.8},
		MountOffset: 42,
		ShaftOffset: 18.9,
		ShaftLength: 5.4,
		ShaftRadius: 2.7,
		HoleRadius:  2.8,
	}
	m.Add("hitec_hs_805bb", &k)
	m.Add("large", &k)

	k = ServoParms{
		Body:        sdf.V3{64, 33, 73.3},
		Mount:       sdf.V3{88, 33, 4},
		Hole:        sdf.V2{76, 21},
		MountOffset: 53.3,
		ShaftOffset: 20.6,
		ShaftLength: 7.6,
		ShaftRadius: 3.8,
		HoleRadius:  3,
	}
	m.Add("hitec_hs_1005sgt", &k)
	m.Add("giant", &k)

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

// Servo2D returns a 2D cutout model for servo mounting.
func Servo2D(k *ServoParms, holeRadius float64) (sdf.SDF2, error) {

	if holeRadius < 0 {
		holeRadius = k.HoleRadius
	}

	const clearance = 1.0 // mounting hole clearance

	// servo body
	body := sdf.Box2D(sdf.V2{k.Body.X + clearance, k.Body.Y + clearance}, 0)

	// holes
	hole, err := sdf.Circle2D(holeRadius)
	if err != nil {
		return nil, err
	}
	xOfs := 0.5 * k.Hole.X
	yOfs := 0.5 * k.Hole.Y
	holes := sdf.Multi2D(hole, []sdf.V2{{xOfs, yOfs}, {-xOfs, yOfs}, {xOfs, -yOfs}, {-xOfs, -yOfs}})

	s := sdf.Union2D(body, holes)

	// position the shaft at the origin
	xOfs = 0.5*k.Hole.X - k.ShaftOffset
	s = sdf.Transform2D(s, sdf.Translate2d(sdf.V2{xOfs, 0}))

	return s, nil
}

//-----------------------------------------------------------------------------

// ServoHornParms stores the parameters that define a servo horn.
type ServoHornParms struct {
	CenterRadius float64 // radius of center hole
	NumHoles     int     // numer of mount holes
	CircleRadius float64 // radius of bolt circle
	HoleRadius   float64 // radius of mount hole
}

// ServoHorn returns a 2D cutout model for a servo horn nount.
func ServoHorn(k *ServoHornParms) (sdf.SDF2, error) {
	if k.CenterRadius < 0 {
		return nil, sdf.ErrMsg("CenterRadius < 0")
	}
	if k.NumHoles < 0 {
		return nil, sdf.ErrMsg("NumHoles < 0")
	}
	if k.CircleRadius < 0 {
		return nil, sdf.ErrMsg("CircleRadius < 0")
	}
	if k.HoleRadius < 0 {
		return nil, sdf.ErrMsg("HoleRadius < 0")
	}

	var s sdf.SDF2

	if k.CenterRadius > 0 {
		h, err := sdf.Circle2D(k.CenterRadius)
		if err != nil {
			return nil, err
		}
		s = sdf.Union2D(s, h)
	}

	if k.NumHoles > 0 && k.CircleRadius > 0 && k.HoleRadius > 0 {
		h, err := BoltCircle2D(k.HoleRadius, k.CircleRadius, k.NumHoles)
		if err != nil {
			return nil, err
		}
		s = sdf.Union2D(s, h)
	}

	return s, nil
}

//-----------------------------------------------------------------------------
