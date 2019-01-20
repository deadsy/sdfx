// spiral generates a PNG of a spiral.
package main

import (
	"flag"
	"log"
	"math"

	. "github.com/deadsy/sdfx/sdf"
)

var (
	start = flag.Float64("start", 0.0, "Start radius (and angle) in radians of spiral")
	end   = flag.Float64("end", 2*math.Pi, "End radius (and angle) in radians of spiral")
	round = flag.Float64("round", 0.0, "Round radius for spiral")
	size  = flag.Int("size", 800, "Size of output PNG file (width and height)")
	out   = flag.String("out", "spiral.png", "Output PNG filename of spiral")
)

func main() {
	flag.Parse()

	s := Spiral2D(*start, *end, *round)
	png, err := NewPNG(*out, s.BoundingBox(), V2i{*size, *size})
	if err != nil {
		log.Fatalf("NewPNG: %v", err)
	}
	png.RenderSDF2(s)
	if err := png.Save(); err != nil {
		log.Fatalf("Save: %v", err)
	}
}
