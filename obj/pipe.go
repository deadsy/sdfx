//-----------------------------------------------------------------------------
/*

Standard Pipes

*/
//-----------------------------------------------------------------------------

package obj

import (
	"fmt"
	"log"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// PipeParameters stores the parameters that define pipe.
type PipeParameters struct {
	Name  string  // name
	Outer float64 // outer radius
	Inner float64 // inner radius
	Wall  float64 // wall thickness
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
		Wall:  (outer - inner) * 0.5,
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
		Wall:  k.Wall * scale,
		Units: units,
	}
	return &k0, nil
}

//-----------------------------------------------------------------------------

// Pipe3D returns a length of the named standard pipe.
func Pipe3D(name string, length float64, units string) (sdf.SDF3, error) {
	if length < 0 {
		return nil, sdf.ErrMsg("length < 0")
	}
	if length == 0 {
		return nil, nil
	}
	k, err := PipeLookup(name, units)
	if err != nil {
		return nil, err
	}
	outer, err := sdf.Cylinder3D(length, k.Outer, 0)
	if err != nil {
		return nil, err
	}
	inner, err := sdf.Cylinder3D(length, k.Inner, 0)
	if err != nil {
		return nil, err
	}
	return sdf.Difference3D(outer, inner), nil
}

//-----------------------------------------------------------------------------
