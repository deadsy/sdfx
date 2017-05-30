package main

import (
	"fmt"

	. "github.com/deadsy/sdfx/sdf"
)

const CUBIC_INCHES_PER_GALLON = 231

var pool_w = 119.0 + 115.0
var pool_l = 101.0 + 101.0 + 96.0 + 96.0 + 83.0

var pool_depth = []V2{
	V2{0, 43},
	V2{101, 46},
	V2{101 + 101, 58},
	V2{101 + 101 + 96, 83},
	V2{101 + 101 + 96 + 96, 96},
	V2{101 + 101 + 96 + 96 + 83, 96},
}

func main() {

	pool_d := (43.0 + 96.0) / 2

	vol := pool_d * pool_w * pool_l

	fmt.Printf("vol %f in3\n", vol)
	fmt.Printf("vol %f gallons\n", vol/CUBIC_INCHES_PER_GALLON)

}
