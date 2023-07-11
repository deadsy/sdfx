package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
)

func bmSquareRestraint(x, y, z float64) (bool, bool, bool) {
	return false, false, false
}

func bmSquareLoad(x, y, z float64) (float64, float64, float64) {
	return 0, 0, 0
}

// Benchmark reference:
// https://github.com/calculix/CalculiX-Examples/tree/master/NonLinear/Sections
func benchmarkSquare() {
	prg := "openscad"

	arg1 := "-o"
	arg2_stl := "3d-beam-square.stl"
	arg3_cad := "../../files/beam-square.scad"

	cmd := exec.Command(prg, arg1, arg2_stl, arg3_cad)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	fmt.Println(string(stdout))

	// read the stl file.
	file, err := os.OpenFile(arg2_stl, os.O_RDONLY, 0400)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// create the SDF from the STL mesh
	inSdf, err := obj.ImportSTL(file, 20, 3, 5)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet4 i.e. 4-node tetrahedron
	err = fe(inSdf, 50, render.Linear, render.Tetrahedral, "tet4.inp", bmSquareRestraint, bmSquareLoad)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet4 i.e. 4-node tetrahedron
	err = fePartial(inSdf, 50, render.Linear, render.Tetrahedral, "partial-tet4.inp", bmSquareRestraint, bmSquareLoad, 0, 3)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fe(inSdf, 50, render.Quadratic, render.Tetrahedral, "tet10.inp", bmSquareRestraint, bmSquareLoad)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// tet10 i.e. 10-node tetrahedron
	err = fePartial(inSdf, 50, render.Quadratic, render.Tetrahedral, "partial-tet10.inp", bmSquareRestraint, bmSquareLoad, 0, 3)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fe(inSdf, 50, render.Linear, render.Hexahedral, "hex8.inp", bmSquareRestraint, bmSquareLoad)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 i.e. 8-node hexahedron
	err = fePartial(inSdf, 50, render.Linear, render.Hexahedral, "partial-hex8.inp", bmSquareRestraint, bmSquareLoad, 0, 3)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fe(inSdf, 50, render.Quadratic, render.Hexahedral, "hex20.inp", bmSquareRestraint, bmSquareLoad)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 i.e. 20-node hexahedron
	err = fePartial(inSdf, 50, render.Quadratic, render.Hexahedral, "partial-hex20.inp", bmSquareRestraint, bmSquareLoad, 0, 3)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fe(inSdf, 50, render.Linear, render.Both, "hex8tet4.inp", bmSquareRestraint, bmSquareLoad)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex8 and tet4
	err = fePartial(inSdf, 50, render.Linear, render.Both, "partial-hex8tet4.inp", bmSquareRestraint, bmSquareLoad, 0, 3)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fe(inSdf, 50, render.Quadratic, render.Both, "hex20tet10.inp", bmSquareRestraint, bmSquareLoad)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// hex20 and tet10
	err = fePartial(inSdf, 50, render.Quadratic, render.Both, "partial-hex20tet10.inp", bmSquareRestraint, bmSquareLoad, 0, 3)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
