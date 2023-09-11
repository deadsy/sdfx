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

func Test_main(t *testing.T) {
	tests := []struct {
		skip      bool
		name      string
		pathSpecs string // File to be created by test.
		specs     Specs
	}{
		{
			skip:      false,
			name:      "teapot",
			pathSpecs: filepath.Join(os.TempDir(), "teapot-specs.json"),
			specs: Specs{
				PathStl:                   filepath.Join("..", "..", "files", "teapot.stl"),
				PathResultWithPlaceholder: filepath.Join(os.TempDir(), "teapot-result-layer0-to-layer#.inp"),
				PathResultInfo:            filepath.Join(os.TempDir(), "teapot-result-info.json"),
				LayerToStartFea:           3,
				MassDensity:               7.85e-9,
				YoungModulus:              210000,
				PoissonRatio:              0.3,
				GravityDirectionX:         0,
				GravityDirectionY:         0,
				GravityDirectionZ:         +1,   // SLA 3D print is usually done upside down.
				GravityMagnitude:          9810, // mm unit.
				Resolution:                50,
				NonlinearConsidered:       false,
				ExactSurfaceConsidered:    true,
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
			err = os.WriteFile(tt.pathSpecs, jsonData, 0644)
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
