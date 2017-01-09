package main

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec"
)

func main() {
	fmt.Printf("%f\n", sdf.Sphere(vec.V3{0, 0, 1}, 1))
	fmt.Printf("%f\n", sdf.Plane(vec.V3{1, 2, 3}, vec.V3{2, 3, 4}.Normalize(), 3))
	fmt.Printf("%f\n", sdf.BoxCheap(vec.V3{1, 2, 3}, vec.V3{2, 3, 4}))
}
