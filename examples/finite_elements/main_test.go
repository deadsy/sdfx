//-----------------------------------------------------------------------------

/*

Finite elements from triangle mesh.
Output `inp` file is consumable by ABAQUS or CalculiX.

*/

//-----------------------------------------------------------------------------

package main

import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
