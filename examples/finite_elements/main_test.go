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
		pthLoadPoints string
		pthLoadDirs   string
		pthLoadMags   string
		pthRestraints string
	}{
		{
			skip:          false,
			name:          "benchmarkSquare",
			pthStl:        "../../files/benchmark-square.stl",
			pthSpecs:      bmsSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: bmsRestraintsPth,
		},
		{
			skip:          false,
			name:          "benchmarkCircle",
			pthStl:        "../../files/benchmark-circle.stl",
			pthSpecs:      bmcSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: bmcRestraintsPth,
		},
		{
			skip:          false,
			name:          "benchmarkPipe",
			pthStl:        "../../files/benchmark-pipe.stl",
			pthSpecs:      bmpSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: bmpRestraintsPth,
		},
		{
			skip:          false,
			name:          "benchmarkI",
			pthStl:        "../../files/benchmark-I.stl",
			pthSpecs:      bmiSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: bmiRestraintsPth,
		},
		{
			skip:          false,
			name:          "teapot",
			pthStl:        "../../files/teapot.stl",
			pthSpecs:      teapotSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: teapotRestraintsPth,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	return err
}

type Specs struct {
	MassDensity            float64
	YoungModulus           float64
	PoissonRatio           float64
	GravityConsidered      bool
	GravityDirectionX      float64
	GravityDirectionY      float64
	GravityDirectionZ      float64
	GravityMagnitude       float64
	Resolution             int
	LayersAllConsidered    bool
	LayerStart             int
	LayerEnd               int
	NonlinearConsidered    bool
	ExactSurfaceConsidered bool
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

type Restraint struct {
	LocationX float64
	LocationY float64
	LocationZ float64
	IsFixedX  bool
	IsFixedY  bool
	IsFixedZ  bool
}

func bmsRestraints() error {
	restraints := make([]Restraint, 0)

	gap := 1.0
	var y float64
	for y <= 17.32 {
		restraint := Restraint{
			LocationX: 0,
			LocationY: y,
			LocationZ: 0,
			IsFixedX:  true,
			IsFixedY:  true,
			IsFixedZ:  true,
		}
		restraints = append(restraints, restraint)
		y += gap
	}

	y = 0
	for y <= 17.32 {
		restraint := Restraint{
			LocationX: 200,
			LocationY: y,
			LocationZ: 0,
			IsFixedX:  false,
			IsFixedY:  true,
			IsFixedZ:  true,
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
		{
			LocationX: 0,
			LocationY: 0,
			LocationZ: 0,
			IsFixedX:  true,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 0,
			LocationY: -2.0313,
			LocationZ: 0.213498,
			IsFixedX:  true,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 0,
			LocationY: -3.97382,
			LocationZ: 0.844661,
			IsFixedX:  true,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 0,
			LocationY: 2.0313,
			LocationZ: 0.213498,
			IsFixedX:  true,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 0,
			LocationY: 3.97382,
			LocationZ: 0.844661,
			IsFixedX:  true,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 200,
			LocationY: 0,
			LocationZ: 0,
			IsFixedX:  false,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 200,
			LocationY: -2.0313,
			LocationZ: 0.213498,
			IsFixedX:  false,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 200,
			LocationY: -3.97382,
			LocationZ: 0.844661,
			IsFixedX:  false,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 200,
			LocationY: 2.0313,
			LocationZ: 0.213498,
			IsFixedX:  false,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
		{
			LocationX: 200,
			LocationY: 3.97382,
			LocationZ: 0.844661,
			IsFixedX:  false,
			IsFixedY:  true,
			IsFixedZ:  true,
		},
	}

	jsonData, err := json.MarshalIndent(restraints, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bmcRestraintsPth, jsonData, 0644)
}
