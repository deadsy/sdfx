//-----------------------------------------------------------------------------

/*

Finite elements from triangle mesh.
The result `inp` file is consumable by ABAQUS or CalculiX.

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
		pthStl        string // Input STL file.
		pthSpecs      string // To be created by test.
		pthLoads      string // To be created by test.
		pthRestraints string // To be created by test.
		pthResult     string // Result file, consumable by ABAQUS or CalculiX.
		specs         Specs
		loads         []Load
		restraints    []Restraint
	}{
		{
			skip:          false,
			name:          "benchmarkSquare",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-square.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bms-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bms-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bms-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bms-result.inp"),
			specs: Specs{
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
			},
			loads: []Load{
				{
					LocX: 0,
					LocY: 0,
					LocZ: 0,
					MagX: 0,
					MagY: 0,
					MagZ: 0,
				},
			},
			restraints: func() []Restraint {
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
				return restraints
			}(),
		},
		{
			skip:          false,
			name:          "benchmarkCircle",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-circle.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bmc-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bmc-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bmc-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bmc-result.inp"),
			specs: Specs{
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
			},
			loads: []Load{
				{
					LocX: 0,
					LocY: 0,
					LocZ: 0,
					MagX: 0,
					MagY: 0,
					MagZ: 0,
				},
			},
			restraints: []Restraint{
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
			},
		},
		{
			skip:          false,
			name:          "benchmarkPipe",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-pipe.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bmp-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bmp-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bmp-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bmp-result.inp"),
			specs: Specs{
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
			},
			loads: []Load{
				{
					LocX: 0,
					LocY: 0,
					LocZ: 0,
					MagX: 0,
					MagY: 0,
					MagZ: 0,
				},
			},
			restraints: []Restraint{
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
			},
		},
		{
			skip:          false,
			name:          "benchmarkI",
			pthStl:        filepath.Join("..", "..", "files", "benchmark-I.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "bmi-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "bmi-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "bmi-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "bmi-result.inp"),
			specs: Specs{
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
			},
			loads: []Load{
				{
					LocX: 0,
					LocY: 0,
					LocZ: 0,
					MagX: 0,
					MagY: 0,
					MagZ: 0,
				},
			},
			restraints: func() []Restraint {
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

				return restraints
			}(),
		},
		{
			skip:          false,
			name:          "teapot",
			pthStl:        filepath.Join("..", "..", "files", "teapot.stl"),
			pthSpecs:      filepath.Join(os.TempDir(), "teapot-specs.json"),
			pthLoads:      filepath.Join(os.TempDir(), "teapot-loads.json"),
			pthRestraints: filepath.Join(os.TempDir(), "teapot-restraints.json"),
			pthResult:     filepath.Join(os.TempDir(), "teapot-result.inp"),
			specs: Specs{
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
			},
			loads: []Load{
				{
					LocX: 0,
					LocY: 0,
					LocZ: 8.0,
					MagX: 0,
					MagY: 0,
					MagZ: -10,
				},
			},
			restraints: []Restraint{
				{LocX: -2.5, LocY: 2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
				{LocX: 2.5, LocY: 2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
				{LocX: 2.5, LocY: -2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
				{LocX: -2.5, LocY: -2.5, LocZ: 0.3, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.MarshalIndent(tt.specs, "", "    ")
			if err != nil {
				t.Error(err)
				return
			}
			err = os.WriteFile(tt.pthSpecs, jsonData, 0644)
			if err != nil {
				t.Error(err)
				return
			}

			jsonData, err = json.MarshalIndent(tt.loads, "", "    ")
			if err != nil {
				t.Error(err)
				return
			}

			err = os.WriteFile(tt.pthLoads, jsonData, 0644)
			if err != nil {
				t.Error(err)
				return
			}

			jsonData, err = json.MarshalIndent(tt.restraints, "", "    ")
			if err != nil {
				t.Error(err)
				return
			}

			err = os.WriteFile(tt.pthRestraints, jsonData, 0644)
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
