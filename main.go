package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	s := sdf.NewRoundedBoxSDF3(sdf.V3{0.4, 0.8, 1.2}, 0.05)
	sdf.Render(&s)
}
