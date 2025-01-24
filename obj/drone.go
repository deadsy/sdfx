//-----------------------------------------------------------------------------
/*

Drone Parts

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

const hexRoundFactor = 0.2

//-----------------------------------------------------------------------------
// motor mount arm

// DroneArmParms are drone arm parameters.
type DroneArmParms struct {
	MotorSize     v2.Vec  // motor diameter/height
	MotorMount    v3.Vec  // motor mount l0, l1, diameter
	RotorCavity   v2.Vec  // cavity for bottom of rotor
	WallThickness float64 // wall thickness
	SideClearance float64 // wall to motor clearance
	MountHeight   float64 // height of motor mount wrt motor height
	ArmHeight     float64 // height of arm wrt motor mount height
	ArmLength     float64 // length of rotor arm
}

// motorMountHeight returns the height of the motor mount.
func mountHeight(k *DroneArmParms) float64 {
	return (k.MountHeight * k.MotorSize.Y) + k.WallThickness
}

// armHeight returns the height of the arm.
func armHeight(k *DroneArmParms) float64 {
	return mountHeight(k) * k.ArmHeight
}

func droneArm(k *DroneArmParms, inner bool) (sdf.SDF3, error) {

	h0 := mountHeight(k)
	h1 := armHeight(k)
	zOfs := 0.5 * (h0 - h1)

	if inner {
		h1 -= 2 * k.WallThickness
	}
	r := h1 / math.Sqrt(3)
	round := r * hexRoundFactor

	arm, err := Hex3D(r, k.ArmLength, round)
	if err != nil {
		return nil, err
	}

	arm = sdf.Transform3D(arm, sdf.RotateX(sdf.DtoR(90)))
	arm = sdf.Transform3D(arm, sdf.Translate3d(v3.Vec{0, 0.5 * k.ArmLength, -zOfs}))
	arm = sdf.Transform3D(arm, sdf.RotateZ(sdf.DtoR(135)))

	return arm, nil
}

func droneMotorBase(k *DroneArmParms) (sdf.SDF3, error) {

	// base
	r0 := (0.5 * k.MotorSize.X) + k.SideClearance + 0.1*k.WallThickness
	h0 := k.WallThickness
	base, err := sdf.Cylinder3D(h0, r0, 0)
	if err != nil {
		return nil, err
	}

	// base rotor cavity
	r1 := 0.5 * k.RotorCavity.X
	h1 := k.RotorCavity.Y
	cavity, err := sdf.Cylinder3D(h1, r1, 0)
	if err != nil {
		return nil, err
	}
	zOfs := 0.5 * (h0 - h1)
	cavity = sdf.Transform3D(cavity, sdf.Translate3d(v3.Vec{0, 0, zOfs}))

	// mount holes
	r2 := 0.5 * k.MotorMount.Z
	h2 := k.WallThickness
	mountHole, err := CounterSunkHole3D(h2, r2)
	if err != nil {
		return nil, err
	}
	mountHole = sdf.Transform3D(mountHole, sdf.RotateX(sdf.DtoR(180)))
	mountPositions := v3.VecSet{
		{0.5 * k.MotorMount.X, 0, 0},
		{-0.5 * k.MotorMount.X, 0, 0},
		{0, 0.5 * k.MotorMount.Y, 0},
		{0, -0.5 * k.MotorMount.Y, 0},
	}
	mountHoles := sdf.Multi3D(mountHole, mountPositions)

	// vent holes
	vent := sdf.Extrude3D(sdf.Box2D(v2.Vec{r0, r0}, 0.2*r0), k.WallThickness)
	v0 := sdf.Transform3D(vent, sdf.Translate3d(v3.Vec{0.8 * r0, 0.8 * r0, 0}))
	v1 := sdf.Transform3D(v0, sdf.RotateZ(sdf.DtoR(90)))
	v2 := sdf.Transform3D(v0, sdf.RotateZ(sdf.DtoR(-90)))

	s := sdf.Difference3D(base, sdf.Union3D(cavity, mountHoles, v0, v1, v2))

	zOfs = -0.5 * (mountHeight(k) - k.WallThickness)
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, zOfs})), nil
}

func droneMotorBody(k *DroneArmParms) (sdf.SDF3, error) {
	r := (0.5 * k.MotorSize.X) + k.SideClearance + k.WallThickness
	h := mountHeight(k)
	round := k.WallThickness * 0.5
	return sdf.Cylinder3D(h, r, round)
}

func droneMotorCavity(k *DroneArmParms) (sdf.SDF3, error) {
	r := (0.5 * k.MotorSize.X) + k.SideClearance
	h := mountHeight(k)
	return sdf.Cylinder3D(h, r, 0)
}

// DroneMotorArm returns a drone motor arm.
func DroneMotorArm(k *DroneArmParms) (sdf.SDF3, error) {

	// outer body
	body, err := droneMotorBody(k)
	if err != nil {
		return nil, err
	}

	// inner cavity
	cavity, err := droneMotorCavity(k)
	if err != nil {
		return nil, err
	}

	// motor base
	base, err := droneMotorBase(k)
	if err != nil {
		return nil, err
	}

	// outer arm
	arm0, err := droneArm(k, false)
	if err != nil {
		return nil, err
	}

	// inner arm
	arm1, err := droneArm(k, true)
	if err != nil {
		return nil, err
	}

	// body + arm
	s := sdf.Union3D(body, arm0)
	s.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	// remove the motor cavity
	s = sdf.Difference3D(s, cavity)
	// remove the inner arm
	s = sdf.Difference3D(s, arm1)
	// add the motor mount base
	s = sdf.Union3D(s, base)

	// carve the top and bottom to remove bumps
	h := mountHeight(k) * 0.5
	s = sdf.Cut3D(s, v3.Vec{0, 0, h}, v3.Vec{0, 0, -1})
	s = sdf.Cut3D(s, v3.Vec{0, 0, -h}, v3.Vec{0, 0, 1})

	// x-axis alignment
	s = sdf.Transform3D(s, sdf.RotateZ(sdf.DtoR(-45)))

	return s, nil
}

//-----------------------------------------------------------------------------
// socket to fit the motor arm

// DroneArmSocketParms defines a socket for a drone motor arm.
type DroneArmSocketParms struct {
	Arm       *DroneArmParms // drone arm parameters
	Size      v3.Vec         // body size for socket
	Clearance float64        // clearance between arm and socket
	Stop      float64        // depth of arm stop
}

func socketHeight(k *DroneArmSocketParms) float64 {
	return armHeight(k.Arm) + (2 * k.Clearance)
}

func socketBody(k *DroneArmSocketParms) (sdf.SDF3, error) {
	return sdf.Box3D(k.Size, 0)
}

func socketArmHole(k *DroneArmSocketParms) (sdf.SDF3, error) {

	h := socketHeight(k)
	r := h / math.Sqrt(3)
	round := r * hexRoundFactor

	s, err := Hex3D(r, k.Size.X, round)
	if err != nil {
		return nil, err
	}

	s = sdf.Transform3D(s, sdf.RotateY(sdf.DtoR(90)))
	s = sdf.Transform3D(s, sdf.RotateX(sdf.DtoR(30)))

	return s, nil
}

func socketStop(k *DroneArmSocketParms) (sdf.SDF3, error) {

	h0 := k.Size.Z
	h1 := 0.5*(h0-socketHeight(k)) + k.Arm.WallThickness

	s, err := sdf.Box3D(v3.Vec{k.Arm.WallThickness, k.Size.Y, h1}, 0)
	if err != nil {
		return nil, err
	}

	zOfs := -0.5 * (h0 - h1)
	xOfs := 0.5*(k.Size.X-k.Arm.WallThickness) - k.Stop

	s = sdf.Transform3D(s, sdf.Translate3d(v3.Vec{xOfs, 0, zOfs}))
	return s, nil
}

// DroneMotorArmSocket returns a socket for a drone motor arm.
func DroneMotorArmSocket(k *DroneArmSocketParms) (sdf.SDF3, error) {
	h := socketHeight(k)
	if k.Size.Y <= h || k.Size.Z <= h {
		return nil, sdf.ErrMsg("socket body is too small for arm")
	}
	if k.Stop >= (k.Size.X - k.Arm.WallThickness) {
		return nil, sdf.ErrMsg("socket body is shorter than arm stop depth")
	}

	body, err := socketBody(k)
	if err != nil {
		return nil, err
	}

	hole, err := socketArmHole(k)
	if err != nil {
		return nil, err
	}

	stop, err := socketStop(k)
	if err != nil {
		return nil, err
	}

	s := sdf.Difference3D(body, hole)
	s = sdf.Union3D(s, stop)

	return s, nil
}

//-----------------------------------------------------------------------------
