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
			pthRestraints: "",
		},
		{
			skip:          false,
			name:          "benchmarkCircle",
			pthStl:        "../../files/benchmark-circle.stl",
			pthSpecs:      bmcSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: "",
		},
		{
			skip:          false,
			name:          "benchmarkPipe",
			pthStl:        "../../files/benchmark-pipe.stl",
			pthSpecs:      bmpSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: "",
		},
		{
			skip:          false,
			name:          "benchmarkI",
			pthStl:        "../../files/benchmark-I.stl",
			pthSpecs:      bmiSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: "",
		},
		{
			skip:          false,
			name:          "teapot",
			pthStl:        "../../files/teapot.stl",
			pthSpecs:      teapotSpecsPth,
			pthLoadPoints: "",
			pthLoadDirs:   "",
			pthLoadMags:   "",
			pthRestraints: "",
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
	return teapotSpecs()
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
