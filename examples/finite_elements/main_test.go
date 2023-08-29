//-----------------------------------------------------------------------------

/*

Finite elements from triangle mesh.
Output `inp` file is consumable by ABAQUS or CalculiX.

*/

//-----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
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
	setup()
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

func setup() {
	bmsSpecs()
	bmcSpecs()
	bmpSpecs()
	bmiSpecs()
	teapotSpecs()
}

func bmsSpecs() {
	specs := map[string]float64{
		"massDensity":       7.85e-9,
		"youngModulus":      210000,
		"poissonRatio":      0.3,
		"gravityConsidered": 1, // 0 or 1
		"gravityDirectionX": 0,
		"gravityDirectionY": 0,
		"gravityDirectionZ": -1,
		"gravityMagnitude":  9810,
		"resolution":        50,
		"layerStart":        -1, // Negative means all layers.
		"layerEnd":          -1, // Negative means all layers.
		"Tet4":              0,
		"Tet10":             0,
		"Hex8":              0,
		"Hex20":             0,
		"Hex8Tet4":          1,
		"Hex20Tet10":        0,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	err = os.WriteFile(bmsSpecsPth, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON file:", err)
		return
	}
}

func bmcSpecs() {
	specs := map[string]float64{
		"massDensity":       7.85e-9,
		"youngModulus":      210000,
		"poissonRatio":      0.3,
		"gravityConsidered": 1, // 0 or 1
		"gravityDirectionX": 0,
		"gravityDirectionY": 0,
		"gravityDirectionZ": -1,
		"gravityMagnitude":  9810,
		"resolution":        50,
		"layerStart":        -1, // Negative means all layers.
		"layerEnd":          -1, // Negative means all layers.
		"Tet4":              0,
		"Tet10":             0,
		"Hex8":              0,
		"Hex20":             0,
		"Hex8Tet4":          1,
		"Hex20Tet10":        0,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	err = os.WriteFile(bmcSpecsPth, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON file:", err)
		return
	}
}

func bmpSpecs() {
	specs := map[string]float64{
		"massDensity":       7.85e-9,
		"youngModulus":      210000,
		"poissonRatio":      0.3,
		"gravityConsidered": 1, // 0 or 1
		"gravityDirectionX": 0,
		"gravityDirectionY": 0,
		"gravityDirectionZ": -1,
		"gravityMagnitude":  9810,
		"resolution":        50,
		"layerStart":        -1, // Negative means all layers.
		"layerEnd":          -1, // Negative means all layers.
		"Tet4":              0,
		"Tet10":             0,
		"Hex8":              0,
		"Hex20":             0,
		"Hex8Tet4":          1,
		"Hex20Tet10":        0,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	err = os.WriteFile(bmpSpecsPth, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON file:", err)
		return
	}
}

func bmiSpecs() {
	specs := map[string]float64{
		"massDensity":       7.85e-9,
		"youngModulus":      210000,
		"poissonRatio":      0.3,
		"gravityConsidered": 1, // 0 or 1
		"gravityDirectionX": 0,
		"gravityDirectionY": 0,
		"gravityDirectionZ": -1,
		"gravityMagnitude":  9810,
		"resolution":        50,
		"layerStart":        -1, // Negative means all layers.
		"layerEnd":          -1, // Negative means all layers.
		"Tet4":              0,
		"Tet10":             0,
		"Hex8":              0,
		"Hex20":             0,
		"Hex8Tet4":          1,
		"Hex20Tet10":        0,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	err = os.WriteFile(bmiSpecsPth, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON file:", err)
		return
	}
}

func teapotSpecs() {
	specs := map[string]float64{
		"massDensity":       7.85e-9,
		"youngModulus":      210000,
		"poissonRatio":      0.3,
		"gravityConsidered": 1, // 0 or 1
		"gravityDirectionX": 0,
		"gravityDirectionY": 0,
		"gravityDirectionZ": -1,
		"gravityMagnitude":  9810,
		"resolution":        50,
		"layerStart":        -1, // Negative means all layers.
		"layerEnd":          -1, // Negative means all layers.
		"Tet4":              0,
		"Tet10":             0,
		"Hex8":              0,
		"Hex20":             0,
		"Hex8Tet4":          1,
		"Hex20Tet10":        0,
	}

	jsonData, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	err = os.WriteFile(teapotSpecsPth, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON file:", err)
		return
	}
}
