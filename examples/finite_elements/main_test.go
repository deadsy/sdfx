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

// Benchmark reference:
// https://github.com/calculix/CalculiX-Examples/tree/master/NonLinear/Sections
func Test_main(t *testing.T) {
	tests := []struct {
		skip          bool
		name          string
		pthStl        string
		pthSpecs      string
		pthLoads      string
		pthRestraints string
		pthResult     string
	}{
		{
			skip:          false,
			name:          "benchmarkSquare",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-square.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bms-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bms-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bms-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bms-result.inp"),
		},
		{
			skip:          false,
			name:          "benchmarkCircle",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-circle.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bmc-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bmc-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bmc-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bmc-result.inp"),
		},
		{
			skip:          false,
			name:          "benchmarkPipe",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-pipe.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bmp-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bmp-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bmp-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bmp-result.inp"),
		},
		{
			skip:          false,
			name:          "benchmarkI",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-I.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bmi-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bmi-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bmi-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bmi-result.inp"),
		},
		{
			skip:          false,
			name:          "teapot",
			pthStl:        filepath.Join("..", "..", "files", "teapot.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "teapot-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "teapot-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "teapot-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "teapot-result.inp"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setup()
			if err != nil {
				t.Error(err)
				return
			}

			os.Args = []string{
				"executable-name-dummy",
				tt.pthStl,
				tt.pthSpecs,
				tt.pthLoads,
				tt.pthRestraints,
				tt.pthResult,
			}
			main()
		})
	}
}

func setup() error {
	err := bmsSpecs(bmsSpecsPth)
	if err != nil {
		return err
	}
	err = bmcSpecs(bmcSpecsPth)
	if err != nil {
		return err
	}
	err = bmpSpecs(bmpSpecsPth)
	if err != nil {
		return err
	}
	err = bmiSpecs(bmiSpecsPth)
	if err != nil {
		return err
	}
	err = teapotSpecs(teapotSpecsPth)
	if err != nil {
		return err
	}
	err = bmsRestraints(bmsRestraintsPth)
	if err != nil {
		return err
	}
	err = bmcRestraints(bmcRestraintsPth)
	if err != nil {
		return err
	}
	err = bmpRestraints(bmpRestraintsPth)
	if err != nil {
		return err
	}
	err = bmiRestraints(bmiRestraintsPth)
	if err != nil {
		return err
	}
	err = teapotRestraints(teapotRestraintsPth)
	if err != nil {
		return err
	}
	err = bmsLoads(bmsLoadsPth)
	if err != nil {
		return err
	}
	err = bmcLoads(bmcLoadsPth)
	if err != nil {
		return err
	}
	err = bmpLoads(bmpLoadsPth)
	if err != nil {
		return err
	}
	err = bmiLoads(bmiLoadsPth)
	if err != nil {
		return err
	}
	return teapotLoads(teapotLoadsPth)
}

func bmsSpecs(pth string) error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmcSpecs(pth string) error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmpSpecs(pth string) error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmiSpecs(pth string) error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
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

	return os.WriteFile(pth, jsonData, 0644)
}

func teapotSpecs(pth string) error {
	specs := Specs{
		MassDensity:            7.85e-9,
		YoungModulus:           210000,
		PoissonRatio:           0.3,
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmsRestraints(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmcRestraints(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmpRestraints(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmiRestraints(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func teapotRestraints(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmsLoads(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmcLoads(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmpLoads(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func bmiLoads(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}

func teapotLoads(pth string) error {
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

	return os.WriteFile(pth, jsonData, 0644)
}
