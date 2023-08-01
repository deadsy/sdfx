package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deadsy/sdfx/obj"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	// read the stl file.
	file, err := os.OpenFile("../../files/teapot.stl", os.O_RDONLY, 0400)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// create the SDF from the STL mesh
	inSdf, err := obj.ImportSTL(file, 20, 3, 5)
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
