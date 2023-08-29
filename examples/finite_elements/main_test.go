//-----------------------------------------------------------------------------

/*

Finite elements from triangle mesh.
Output `inp` file is consumable by ABAQUS or CalculiX.

*/

//-----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var bmsSpecsPth string = filepath.Join(os.TempDir(), "bms-specs.json")
var bmcSpecsPth string = filepath.Join(os.TempDir(), "bmc-specs.json")
var bmpSpecsPth string = filepath.Join(os.TempDir(), "bmp-specs.json")
var bmiSpecsPth string = filepath.Join(os.TempDir(), "bmi-specs.json")
var teapotSpecsPth string = filepath.Join(os.TempDir(), "teapot-specs.json")

var bmsLoadsPth string = filepath.Join(os.TempDir(), "bms-loads.json")
var bmcLoadsPth string = filepath.Join(os.TempDir(), "bmc-loads.json")
var bmpLoadsPth string = filepath.Join(os.TempDir(), "bmp-loads.json")
var bmiLoadsPth string = filepath.Join(os.TempDir(), "bmi-loads.json")
var teapotLoadsPth string = filepath.Join(os.TempDir(), "teapot-loads.json")

var bmsRestraintsPth string = filepath.Join(os.TempDir(), "bms-restraints.json")
var bmcRestraintsPth string = filepath.Join(os.TempDir(), "bmc-restraints.json")
var bmpRestraintsPth string = filepath.Join(os.TempDir(), "bmp-restraints.json")
var bmiRestraintsPth string = filepath.Join(os.TempDir(), "bmi-restraints.json")
var teapotRestraintsPth string = filepath.Join(os.TempDir(), "teapot-restraints.json")

func Test_main(t *testing.T) {
	err := setup()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		skip          bool
		name          string
		pthStl        string
		pthSpecs      string
		pthLoads      string
		pthRestraints string
	}{
		{
			skip:          false,
			name:          "benchmarkSquare",
			pthStl:        "../../files/benchmark-square.stl",
			pthSpecs:      bmsSpecsPth,
			pthLoads:      bmsLoadsPth,
			pthRestraints: bmsRestraintsPth,
		},
		{
			skip:          false,
			name:          "benchmarkCircle",
			pthStl:        "../../files/benchmark-circle.stl",
			pthSpecs:      bmcSpecsPth,
			pthLoads:      bmcLoadsPth,
			pthRestraints: bmcRestraintsPth,
		},
		{
			skip:          false,
			name:          "benchmarkPipe",
			pthStl:        "../../files/benchmark-pipe.stl",
			pthSpecs:      bmpSpecsPth,
			pthLoads:      bmpLoadsPth,
			pthRestraints: bmpRestraintsPth,
		},
		{
			skip:          false,
			name:          "benchmarkI",
			pthStl:        "../../files/benchmark-I.stl",
			pthSpecs:      bmiSpecsPth,
			pthLoads:      bmiLoadsPth,
			pthRestraints: bmiRestraintsPth,
		},
		{
			skip:          false,
			name:          "teapot",
			pthStl:        "../../files/teapot.stl",
			pthSpecs:      teapotSpecsPth,
			pthLoads:      teapotLoadsPth,
			pthRestraints: teapotRestraintsPth,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{
				"executable-name-dummy",
				tt.pthStl,
				tt.pthSpecs,
				tt.pthLoads,
				tt.pthRestraints,
			}
			main()
		})
	}
}

func setup() error {
	err := bmsSpecs()
	if err != nil {
		return err
	}
	err = bmcSpecs()
	if err != nil {
		return err
	}
	err = bmpSpecs()
	if err != nil {
		return err
	}
	err = bmiSpecs()
	if err != nil {
		return err
	}
	err = teapotSpecs()
	if err != nil {
		return err
	}
	err = bmsRestraints()
	if err != nil {
		return err
	}
	err = bmcRestraints()
	if err != nil {
		return err
	}
	err = bmpRestraints()
	if err != nil {
		return err
	}
	err = bmiRestraints()
	if err != nil {
		return err
	}
	err = teapotRestraints()
	if err != nil {
		return err
	}
	err = bmsLoads()
	if err != nil {
		return err
	}
	err = bmcLoads()
	if err != nil {
		return err
	}
	err = bmpLoads()
	if err != nil {
		return err
	}
	err = bmiLoads()
	if err != nil {
		return err
	}
	return teapotLoads()
}

func bmsSpecs() error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
		GravityConsidered:      true,
		GravityDirectionX:      0,
		GravityDirectionY:      0,
		GravityDirectionZ:      -1,
		GravityMagnitude:       9810,
		Resolution:             50,
		LayersAllConsidered:    true,
		LayerStart:             -1, // Negative means all layers.
		LayerEnd:               -1, // Negative means all layers.
		NonlinearConsidered:    false,
		ExactSurfaceConsidered: true,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmsSpecsPth, jsonData, 0644)
}

func bmcSpecs() error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
		GravityConsidered:      true,
		GravityDirectionX:      0,
		GravityDirectionY:      0,
		GravityDirectionZ:      -1,
		GravityMagnitude:       9810,
		Resolution:             50,
		LayersAllConsidered:    true,
		LayerStart:             -1, // Negative means all layers.
		LayerEnd:               -1, // Negative means all layers.
		NonlinearConsidered:    false,
		ExactSurfaceConsidered: true,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmcSpecsPth, jsonData, 0644)
}

func bmpSpecs() error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
		GravityConsidered:      true,
		GravityDirectionX:      0,
		GravityDirectionY:      0,
		GravityDirectionZ:      -1,
		GravityMagnitude:       9810,
		Resolution:             50,
		LayersAllConsidered:    true,
		LayerStart:             -1, // Negative means all layers.
		LayerEnd:               -1, // Negative means all layers.
		NonlinearConsidered:    false,
		ExactSurfaceConsidered: true,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmpSpecsPth, jsonData, 0644)
}

func bmiSpecs() error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
		GravityConsidered:      true,
		GravityDirectionX:      0,
		GravityDirectionY:      0,
		GravityDirectionZ:      -1,
		GravityMagnitude:       9810,
		Resolution:             50,
		LayersAllConsidered:    true,
		LayerStart:             -1, // Negative means all layers.
		LayerEnd:               -1, // Negative means all layers.
		NonlinearConsidered:    false,
		ExactSurfaceConsidered: true,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmiSpecsPth, jsonData, 0644)
}

func teapotSpecs() error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
		GravityConsidered:      true,
		GravityDirectionX:      0,
		GravityDirectionY:      0,
		GravityDirectionZ:      -1,
		GravityMagnitude:       9810,
		Resolution:             50,
		LayersAllConsidered:    true,
		LayerStart:             -1, // Negative means all layers.
		LayerEnd:               -1, // Negative means all layers.
		NonlinearConsidered:    false,
		ExactSurfaceConsidered: true,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(teapotSpecsPth, jsonData, 0644)
}

func bmsRestraints() error {
	restraints := make([]Restraint, 0)

	gap := 1.0
	var y float64
	for y <= 17.32 {
		restraint := Restraint{
			LocX:     0,
			LocY:     y,
			LocZ:     0,
			IsFixedX: true,
			IsFixedY: true,
			IsFixedZ: true,
		}
		restraints = append(restraints, restraint)
		y += gap
	}

	y = 0
	for y <= 17.32 {
		restraint := Restraint{
			LocX:     200,
			LocY:     y,
			LocZ:     0,
			IsFixedX: false,
			IsFixedY: true,
			IsFixedZ: true,
		}
		restraints = append(restraints, restraint)
		y += gap
	}

	jsonData, err := json.MarshalIndent(restraints, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmsRestraintsPth, jsonData, 0644)
}

func bmcRestraints() error {
	restraints := []Restraint{
		{LocX: 0, LocY: 0, LocZ: 0, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: -2.0313, LocZ: 0.213498, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: -3.97382, LocZ: 0.844661, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: 2.0313, LocZ: 0.213498, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: 3.97382, LocZ: 0.844661, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: 0, LocZ: 0, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: -2.0313, LocZ: 0.213498, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: -3.97382, LocZ: 0.844661, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: 2.0313, LocZ: 0.213498, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: 3.97382, LocZ: 0.844661, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
	}

	jsonData, err := json.MarshalIndent(restraints, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmcRestraintsPth, jsonData, 0644)
}

func bmpRestraints() error {
	restraints := []Restraint{
		{LocX: 0, LocY: 0, LocZ: 0, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: -2.0313, LocZ: 0.213498, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: -3.97382, LocZ: 0.844661, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: 2.0313, LocZ: 0.213498, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 0, LocY: 3.97382, LocZ: 0.844661, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: 0, LocZ: 0, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: -2.0313, LocZ: 0.213498, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: -3.97382, LocZ: 0.844661, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: 2.0313, LocZ: 0.213498, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
		{LocX: 200, LocY: 3.97382, LocZ: 0.844661, IsFixedX: false, IsFixedY: true, IsFixedZ: true},
	}

	jsonData, err := json.MarshalIndent(restraints, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmpRestraintsPth, jsonData, 0644)
}

func bmiRestraints() error {
	restraints := make([]Restraint, 0)

	gap := 1.0
	var y float64
	for y <= 25 {
		restraints = append(restraints, Restraint{LocX: 0, LocY: y, LocZ: 0, IsFixedX: true, IsFixedY: true, IsFixedZ: true})
		y += gap
	}

	y = 0
	for y <= 25 {
		restraints = append(restraints, Restraint{LocX: 200, LocY: y, LocZ: 0, IsFixedX: false, IsFixedY: true, IsFixedZ: true})
		y += gap
	}

	jsonData, err := json.MarshalIndent(restraints, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmiRestraintsPth, jsonData, 0644)
}

func teapotRestraints() error {
	restraints := []Restraint{
		{LocX: -2.5, LocY: 2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 2.5, LocY: 2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: 2.5, LocY: -2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
		{LocX: -2.5, LocY: -2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
	}

	jsonData, err := json.MarshalIndent(restraints, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(teapotRestraintsPth, jsonData, 0644)
}

func bmsLoads() error {
	loads := []Load{
		{
			LocX: 0,
			LocY: 0,
			LocZ: 0,
			MagX: 0,
			MagY: 0,
			MagZ: 0,
		},
	}

	jsonData, err := json.MarshalIndent(loads, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmsLoadsPth, jsonData, 0644)
}

func bmcLoads() error {
	loads := []Load{
		{
			LocX: 0,
			LocY: 0,
			LocZ: 0,
			MagX: 0,
			MagY: 0,
			MagZ: 0,
		},
	}

	jsonData, err := json.MarshalIndent(loads, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmcLoadsPth, jsonData, 0644)
}

func bmpLoads() error {
	loads := []Load{
		{
			LocX: 0,
			LocY: 0,
			LocZ: 0,
			MagX: 0,
			MagY: 0,
			MagZ: 0,
		},
	}

	jsonData, err := json.MarshalIndent(loads, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmpLoadsPth, jsonData, 0644)
}

func bmiLoads() error {
	loads := []Load{
		{
			LocX: 0,
			LocY: 0,
			LocZ: 0,
			MagX: 0,
			MagY: 0,
			MagZ: 0,
		},
	}

	jsonData, err := json.MarshalIndent(loads, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmiLoadsPth, jsonData, 0644)
}

func teapotLoads() error {
	loads := []Load{
		{
			LocX: 0,
			LocY: 0,
			LocZ: 8.0,
			MagX: 0,
			MagY: 0,
			MagZ: -10,
		},
	}

	jsonData, err := json.MarshalIndent(loads, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(teapotLoadsPth, jsonData, 0644)
}
