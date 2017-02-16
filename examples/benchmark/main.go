package main

import (
	"fmt"

	. "github.com/deadsy/sdfx/sdf"
)

func main() {
	s := NewCircleSDF2(5)
	eps := BenchmarkSDF2(s)
	fmt.Printf("%f\n", eps)
}
