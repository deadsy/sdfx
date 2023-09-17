//-----------------------------------------------------------------------------

/*

Finite elements from triangle mesh.
The result `inp` file is consumable by ABAQUS or CalculiX.

*/

//-----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"math"
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
		pathSpecs  string // File to be created by test.
		specs      Specs
		loads      []Load // If load is zero, gravity would be the dominant force.
		restraints []Restraint
	}{
		{
			skip:      false,
			name:      "teapot",
			pathSpecs: filepath.Join(os.TempDir(), "specs.json"),
			specs: Specs{
				PathResult:             filepath.Join(os.TempDir(), "result.inp"),
				PathReport:             filepath.Join(os.TempDir(), "report.json"),
				PathStl:                filepath.Join("..", "..", "files", "teapot.stl"), // Valid STL, Unit: mm
				PathLoadPoints:         filepath.Join(os.TempDir(), "load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "restraint-points.json"),
				MassDensity:            1130 * math.Pow(10, -12), // (N*s2/mm4) // Assumed: 1.13 g/cm3
				YoungModulus:           1.6 * 1000,               // MPa (N/mm2)
				PoissonRatio:           0.3,                      // Unitless.
				GravityDirectionX:      0,
				GravityDirectionY:      0,
				GravityDirectionZ:      -1,
				GravityMagnitude:       9810, // mm/s2
				GravityIsNeeded:        false,
				Resolution:             60,
				NonlinearConsidered:    false,
				ExactSurfaceConsidered: true,
			},
			loads: []Load{
				{LocX: -7.7018147506213062, LocY: -0.4793329364029888, LocZ: 5.4655784011739659, MagX: -70.381474830032147, MagY: -174.42493975029208, MagZ: 59.390099907428898},
				{LocX: -0.011008272390835461, LocY: -0.7768798803556729, LocZ: 8.0940818810755175, MagX: -7.1696819796276845, MagY: -157.24707657594607, MagZ: -122.86811489950169},
				{LocX: 7.7771501865767299, LocY: -0.44676917365822177, LocZ: 6.1957182567021745, MagX: 8.5251191596139506, MagY: -198.29032531361477, MagZ: 18.196364257032524},
			},
			restraints: []Restraint{
				{LocX: 2.6121906631017695, LocY: 0.20348199936959829, LocZ: 0.050483960817894413, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
				{LocX: -1.3968227044257533, LocY: -2.035934011608322, LocZ: 0.04909315835598238, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
				{LocX: -1.8197506822193277, LocY: 2.2580011513606717, LocZ: 0.064527793304306025, IsFixedX: true, IsFixedY: true, IsFixedZ: true},
			},
		},
		{
			skip:      true,
			name:      "benchmarkSquare",
			pathSpecs: filepath.Join(os.TempDir(), "bms-specs.json"),
			specs: Specs{
				PathResult:             filepath.Join(os.TempDir(), "bms-result.inp"),
				PathReport:             filepath.Join(os.TempDir(), "bms-report.json"),
				PathStl:                filepath.Join("..", "..", "files", "benchmark-square.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bms-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bms-restraint-points.json"),
				MassDensity:            7.85e-9,
				YoungModulus:           210000,
				PoissonRatio:           0.3,
				GravityDirectionX:      0,
				GravityDirectionY:      0,
				GravityDirectionZ:      -1,
				GravityMagnitude:       9810,
				GravityIsNeeded:        true,
				Resolution:             50,
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
			skip:      true,
			name:      "benchmarkCircle",
			pathSpecs: filepath.Join(os.TempDir(), "bmc-specs.json"),
			specs: Specs{
				PathResult:             filepath.Join(os.TempDir(), "bmc-result.inp"),
				PathReport:             filepath.Join(os.TempDir(), "bmc-report.json"),
				PathStl:                filepath.Join("..", "..", "files", "benchmark-circle.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bmc-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bmc-restraint-points.json"),
				MassDensity:            7.85e-9,
				YoungModulus:           210000,
				PoissonRatio:           0.3,
				GravityDirectionX:      0,
				GravityDirectionY:      0,
				GravityDirectionZ:      -1,
				GravityMagnitude:       9810,
				GravityIsNeeded:        true,
				Resolution:             50,
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
			skip:      true,
			name:      "benchmarkPipe",
			pathSpecs: filepath.Join(os.TempDir(), "bmp-specs.json"),
			specs: Specs{
				PathResult:             filepath.Join(os.TempDir(), "bmp-result.inp"),
				PathReport:             filepath.Join(os.TempDir(), "bmp-report.json"),
				PathStl:                filepath.Join("..", "..", "files", "benchmark-pipe.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bmp-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bmp-restraint-points.json"),
				MassDensity:            7.85e-9,
				YoungModulus:           210000,
				PoissonRatio:           0.3,
				GravityDirectionX:      0,
				GravityDirectionY:      0,
				GravityDirectionZ:      -1,
				GravityMagnitude:       9810,
				GravityIsNeeded:        true,
				Resolution:             50,
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
			skip:      true,
			name:      "benchmarkI",
			pathSpecs: filepath.Join(os.TempDir(), "bmi-specs.json"),
			specs: Specs{
				PathResult:             filepath.Join(os.TempDir(), "bmi-result.inp"),
				PathReport:             filepath.Join(os.TempDir(), "bmi-report.json"),
				PathStl:                filepath.Join("..", "..", "files", "benchmark-I.stl"),
				PathLoadPoints:         filepath.Join(os.TempDir(), "bmi-load-points.json"),
				PathRestraintPoints:    filepath.Join(os.TempDir(), "bmi-restraint-points.json"),
				MassDensity:            7.85e-9,
				YoungModulus:           210000,
				PoissonRatio:           0.3,
				GravityDirectionX:      0,
				GravityDirectionY:      0,
				GravityDirectionZ:      -1,
				GravityMagnitude:       9810,
				GravityIsNeeded:        true,
				Resolution:             50,
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
	}

	for _, tt := range tests {
		if tt.skip {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.MarshalIndent(tt.specs, "", "    ")
			if err != nil {
				t.Error(err)
				return
			}
			err = os.WriteFile(tt.pathSpecs, jsonData, 0644)
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
				tt.pathSpecs,
			}
			main()
		})
	}
}
