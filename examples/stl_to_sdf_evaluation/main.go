package main

import (
	"fmt"
	"log"

	"github.com/deadsy/sdfx/obj"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	// create the SDF from the STL file mesh
	inSdf, err := obj.ImportSTL("../../files/teapot.stl", 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// This point is definitely inside the teapot model,
	// so SDF value should be negative.
	value := inSdf.Evaluate(v3.Vec{X: -0.8164382918936324, Y: 2.542909114087213, Z: 5.006102143191411})
	if value >= 0 {
		fmt.Println("not expected")
	} else {
		fmt.Println("as expected")
	}
}
