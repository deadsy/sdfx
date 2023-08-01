//-----------------------------------------------------------------------------
/*

STL File Load/Save Testing

*/
//-----------------------------------------------------------------------------

package render

import (
	"testing"
)

//-----------------------------------------------------------------------------

func Test_LoadSTL(t *testing.T) {
	loadTests := []struct {
		path     string
		meshSize int
	}{
		{"../files/bottle.stl", 1240},
		{"../files/monkey.stl", 366},
		{"../files/teapot.stl", 9438},
	}
	for _, test := range loadTests {
		mesh, err := LoadSTL(test.path)
		if err != nil {
			t.Errorf("%s", err)
		}
		if len(mesh) != test.meshSize {
			t.Errorf("%s expected %d triangles (got %d)", test.path, test.meshSize, len(mesh))
		}
	}
}

//-----------------------------------------------------------------------------
