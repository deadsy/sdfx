package main

import (
	. "github.com/deadsy/sdfx/sdf"
)

func main() {
	RenderSTL(cc16a(), 200, "cc16a.stl")
	RenderSTL(cc16b(), 200, "cc16b.stl")
	//	RenderSTL(cc18b(), 200, "cc18b.stl")
}
