package main

import (
	"github.com/deadsy/sdfx/render"
)

func main() {
	render.RenderSTL(cc16a(), 200, "cc16a.stl")
	render.RenderSTL(cc16b(), 200, "cc16b.stl")
	cc18a()
	render.RenderSTL(cc18b(), 200, "cc18b.stl")
	render.RenderSTL(cc18c(), 200, "cc18c.stl")
}
