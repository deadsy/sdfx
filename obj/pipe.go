//-----------------------------------------------------------------------------
/*

Standard Pipes

*/
//-----------------------------------------------------------------------------

package obj

import (
	"fmt"
	"log"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// PipeParameters stores the parameters that define pipe.
type PipeParameters struct {
	Name  string  // name
	Outer float64 // outer radius
	Inner float64 // inner radius
	Units string  // "inch" or "mm"
}

type pipeDatabase map[string]*PipeParameters

var pipeDB = initPipeLookup()

func (m pipeDatabase) Sch40Add(name string, outer, inner float64) {
	if inner >= outer {
		log.Panicf("inner >= outer for \"sch40:%s\"", name)
	}
	name = "sch40:" + name
	k := PipeParameters{
		Name:  name,
		Outer: outer * 0.5,
		Inner: inner * 0.5,
		Units: "inch",
	}
	m[name] = &k
}

// initPipeLookup adds a collection of standard pipes to the pipe database.
func initPipeLookup() pipeDatabase {
	m := make(pipeDatabase)

	// schedule 40 PVC
	m.Sch40Add("1/8", 0.405, 0.249)
	m.Sch40Add("1/4", 0.540, 0.344)
	m.Sch40Add("3/8", 0.675, 0.473)
	m.Sch40Add("1/2", 0.840, 0.602)
	m.Sch40Add("3/4", 1.050, 0.804)
	m.Sch40Add("1", 1.315, 1.029)
	m.Sch40Add("1-1/4", 1.660, 1.360)
	m.Sch40Add("1-1/2", 1.900, 1.590)
	m.Sch40Add("2", 2.375, 2.047)
	m.Sch40Add("2-1/2", 2.875, 2.445)
	m.Sch40Add("3", 3.500, 3.042)
	m.Sch40Add("3-1/2", 4.000, 3.521)
	m.Sch40Add("4", 4.500, 3.998)
	m.Sch40Add("5", 5.563, 5.016)
	m.Sch40Add("6", 6.625, 6.031)
	m.Sch40Add("8", 8.625, 7.942)
	m.Sch40Add("10", 10.750, 9.976)
	m.Sch40Add("12", 12.750, 11.889)
	m.Sch40Add("14", 14.000, 13.073)
	m.Sch40Add("16", 16.000, 14.940)
	m.Sch40Add("18", 18.000, 16.809)
	m.Sch40Add("20", 20.000, 18.743)
	m.Sch40Add("24", 24.000, 22.544)

	return m
}

// PipeLookup returns the parameters for a named pipe.
func PipeLookup(name, units string) (*PipeParameters, error) {
	if units != "mm" && units != "inch" {
		return nil, sdf.ErrMsg("units must be mm/inch")
	}

	k, ok := pipeDB[name]
	if !ok {
		return nil, fmt.Errorf("pipe \"%s\" not found", name)
	}
	// handle scale conversion
	scale := 1.0
	if units != k.Units {
		if units == "mm" && k.Units == "inch" {
			scale = sdf.MillimetresPerInch
		}
		if units == "inch" && k.Units == "mm" {
			scale = 1.0 / sdf.MillimetresPerInch
		}
	}
	k0 := PipeParameters{
		Outer: k.Outer * scale,
		Inner: k.Inner * scale,
		Units: units,
	}
	return &k0, nil
}

//-----------------------------------------------------------------------------

// Pipe3D returns a length of pipe.
func Pipe3D(oRadius, iRadius, length float64) (sdf.SDF3, error) {
	if oRadius <= 0 {
		return nil, sdf.ErrMsg("oRadius <= 0")
	}
	if iRadius <= 0 {
		return nil, sdf.ErrMsg("iRadius <= 0")
	}
	if oRadius <= iRadius {
		return nil, sdf.ErrMsg("oRadius <= iRadius")
	}
	if length < 0 {
		return nil, sdf.ErrMsg("length < 0")
	}
	if length == 0 {
		return nil, nil
	}
	outer, err := sdf.Cylinder3D(length, oRadius, 0)
	if err != nil {
		return nil, err
	}
	inner, err := sdf.Cylinder3D(length, iRadius, 0)
	if err != nil {
		return nil, err
	}
	return sdf.Difference3D(outer, inner), nil
}

//-----------------------------------------------------------------------------

// StdPipe3D returns a length of the named standard pipe.
func StdPipe3D(name, units string, length float64) (sdf.SDF3, error) {
	k, err := PipeLookup(name, units)
	if err != nil {
		return nil, err
	}
	return Pipe3D(k.Outer, k.Inner, length)
}

//-----------------------------------------------------------------------------

func connectorArm(radius, length float64) (sdf.SDF3, error) {
	if radius <= 0 {
		return nil, sdf.ErrMsg("radius <= 0")
	}
	if length < radius {
		return nil, sdf.ErrMsg("length < radius")
	}
	s, err := sdf.Cylinder3D(length+(2*radius), radius, radius)
	if err != nil {
		return nil, err
	}
	s = sdf.Cut3D(s, sdf.V3{0, 0, 0.5 * length}, sdf.V3{0, 0, -1})
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, 0, length * 0.5})), nil
}

func pipeConnector1(radius, length float64, cfg [6]bool) (sdf.SDF3, error) {
	var dirn []sdf.V3
	if cfg[0] {
		dirn = append(dirn, sdf.V3{1, 0, 0})
	}
	if cfg[1] {
		dirn = append(dirn, sdf.V3{-1, 0, 0})
	}
	if cfg[2] {
		dirn = append(dirn, sdf.V3{0, 1, 0})
	}
	if cfg[3] {
		dirn = append(dirn, sdf.V3{0, -1, 0})
	}
	if cfg[4] {
		dirn = append(dirn, sdf.V3{0, 0, 1})
	}
	if cfg[5] {
		dirn = append(dirn, sdf.V3{0, 0, -1})
	}
	if len(dirn) < 1 {
		return nil, sdf.ErrMsg("no connectors")
	}
	s, err := connectorArm(radius, length)
	if err != nil {
		return nil, err
	}
	return sdf.Orient3D(s, sdf.V3{0, 0, 1}, dirn), nil
}

func pipeConnector2(outer, inner, length float64, cfg [6]bool) (sdf.SDF3, error) {
	// outer
	s, err := pipeConnector1(outer, length, cfg)
	if err != nil {
		return nil, err
	}
	// inner
	if inner > 0 {
		inner, err := pipeConnector1(inner, length, cfg)
		if err != nil {
			return nil, err
		}
		s = sdf.Difference3D(s, inner)
	}
	return s, nil
}

// PipeConnectorParms defines an n-way female pipe connector.
type PipeConnectorParms struct {
	Length        float64 // length of connector arm to center
	OuterRadius   float64 // outer radius of connector arm
	InnerRadius   float64 // inner radius of connector arm
	RecessDepth   float64 // depth for recessed stop
	RecessWidth   float64 // width of internal recess step
	Configuration [6]bool // position of arms. +x,-x,+y,-y,+z,-z
}

// PipeConnector3D returns an n-way female pipe connector.
func PipeConnector3D(k *PipeConnectorParms) (sdf.SDF3, error) {

	if k.Length <= 0 {
		return nil, sdf.ErrMsg("k.Length <= 0")
	}
	if k.OuterRadius <= 0 {
		return nil, sdf.ErrMsg("k.OuterRadius <= 0")
	}
	if k.InnerRadius < 0 {
		return nil, sdf.ErrMsg("k.InnerRadius < 0")
	}
	if k.RecessDepth < 0 {
		return nil, sdf.ErrMsg("k.RecessDepth < 0")
	}
	if k.RecessWidth < 0 {
		return nil, sdf.ErrMsg("k.RecessWidth < 0")
	}
	if k.InnerRadius >= k.OuterRadius {
		return nil, sdf.ErrMsg("k.InnerRadius >= k.OuterRadius")
	}
	if k.RecessDepth >= k.Length {
		return nil, sdf.ErrMsg("k.RecessDepth >= k.Length")
	}
	if k.RecessWidth >= k.InnerRadius {
		return nil, sdf.ErrMsg("k.RecessWidth >= k.InnerRadius")
	}

	// outer surface
	s, err := pipeConnector2(k.OuterRadius, k.InnerRadius, k.Length, k.Configuration)
	if err != nil {
		return nil, err
	}

	// recessed stop
	if k.RecessWidth > 0 {
		length := k.Length - k.RecessDepth
		inner := k.InnerRadius - k.RecessWidth
		recess, err := pipeConnector2(k.InnerRadius, inner, length, k.Configuration)
		if err != nil {
			return nil, err
		}
		s = sdf.Union3D(s, recess)
	}

	return s, nil
}

//-----------------------------------------------------------------------------

// StdPipeConnector3D returns an n-way female pipe connector for a standard pipe size.
func StdPipeConnector3D(name, units string, length float64, cfg [6]bool) (sdf.SDF3, error) {
	p, err := PipeLookup(name, units)
	if err != nil {
		return nil, err
	}
	wall := p.Outer - p.Inner
	k := PipeConnectorParms{
		Length:        length,
		OuterRadius:   p.Outer + wall,
		InnerRadius:   p.Outer,
		RecessDepth:   math.Min(2*p.Outer, length-p.Outer-(0.5*wall)),
		RecessWidth:   wall,
		Configuration: cfg,
	}
	return PipeConnector3D(&k)
}

//-----------------------------------------------------------------------------
