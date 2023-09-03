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
		skip       bool
		name       string
		pthSpecs   string // File to be created by test.
		specs      Specs
		loads      []Load // If load is zero, gravity would be the dominant force.
		restraints []Restraint
	}{
		{
			skip:     false,
			name:     "benchmarkSquare",
			pthSpecs: filepath.Join(os.TempDir(), "bms-specs.json"),
			specs: Specs{
				PathStl:                filepath.Join("..", "..", "files", "benchmark-square.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bms-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bms-restraint-points.json"),
				PathResult:             filepath.Join(os.TempDir(), "bms-result.inp"),
				PathResultInfo:         filepath.Join(os.TempDir(), "bms-result-info.json"),
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
			skip:     false,
			name:     "benchmarkCircle",
			pthSpecs: filepath.Join(os.TempDir(), "bmc-specs.json"),
			specs: Specs{
				PathStl:                filepath.Join("..", "..", "files", "benchmark-circle.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bmc-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bmc-restraint-points.json"),
				PathResult:             filepath.Join(os.TempDir(), "bmc-result.inp"),
				PathResultInfo:         filepath.Join(os.TempDir(), "bmc-result-info.json"),
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
			skip:     false,
			name:     "benchmarkPipe",
			pthSpecs: filepath.Join(os.TempDir(), "bmp-specs.json"),
			specs: Specs{
				PathStl:                filepath.Join("..", "..", "files", "benchmark-pipe.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bmp-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bmp-restraint-points.json"),
				PathResult:             filepath.Join(os.TempDir(), "bmp-result.inp"),
				PathResultInfo:         filepath.Join(os.TempDir(), "bmp-result-info.json"),
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
			skip:     false,
			name:     "benchmarkI",
			pthSpecs: filepath.Join(os.TempDir(), "bmi-specs.json"),
			specs: Specs{
				PathStl:                filepath.Join("..", "..", "files", "benchmark-I.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bmi-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bmi-restraint-points.json"),
				PathResult:             filepath.Join(os.TempDir(), "bmi-result.inp"),
				PathResultInfo:         filepath.Join(os.TempDir(), "bmi-result-info.json"),
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
			skip:     false,
			name:     "teapot",
			pthSpecs: filepath.Join(os.TempDir(), "teapot-specs.json"),
			specs: Specs{
				PathStl:                filepath.Join("..", "..", "files", "teapot.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "teapot-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "teapot-restraint-points.json"),
				PathResult:             filepath.Join(os.TempDir(), "teapot-result.inp"),
				PathResultInfo:         filepath.Join(os.TempDir(), "teapot-result-info.json"),
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

			err = os.WriteFile(tt.specs.PathLoadPoints, jsonData, 0644)
			if err != nil {
				t.Error(err)
				return
			}

			jsonData, err = json.MarshalIndent(tt.restraints, "", "    ")
			if err != nil {
				t.Error(err)
				return
			}

			err = os.WriteFile(tt.specs.PathRestraintPoints, jsonData, 0644)
			if err != nil {
				t.Error(err)
				return
			}

			os.Args = []string{
				"executable-name-dummy",
				tt.pthSpecs,
			}
			main()
		})
	}
}
