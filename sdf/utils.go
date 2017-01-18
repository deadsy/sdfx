package sdf

import (
	"math"
)

func DtoR(degrees float64) float64 {
	return math.Pi * degrees / 180
}

func RtoD(radians float64) float64 {
	return 180 * radians / math.Pi
}
