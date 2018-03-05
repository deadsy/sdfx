//-----------------------------------------------------------------------------
/*

Text Operations

Convert a string and a font specification into an SDF2

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

//-----------------------------------------------------------------------------

const POINT_PER_INCH = 72.0

//-----------------------------------------------------------------------------

func printBounds(b fixed.Rectangle26_6) {
	fmt.Printf("Min.X:%d Min.Y:%d Max.X:%d Max.Y:%d\n", b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
}

//-----------------------------------------------------------------------------

func Test_Text() error {

	// get the font data
	fontfile := "/usr/share/fonts/truetype/msttcorefonts/Arial_Black.ttf"
	b, err := ioutil.ReadFile(fontfile)
	if err != nil {
		return err
	}

	f, err := truetype.Parse(b)
	if err != nil {
		return err
	}

	fupe := fixed.Int26_6(f.FUnitsPerEm())
	printBounds(f.Bounds(fupe))
	fmt.Printf("FUnitsPerEm:%d\n\n", fupe)

	return nil
}

//-----------------------------------------------------------------------------
