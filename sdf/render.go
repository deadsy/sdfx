package sdf

import (
	"github.com/deadsy/pt/pt"
)

func Render(s *SDF3) {

	scene := pt.Scene{}

	light := pt.LightMaterial(pt.White, 180)

	d := 4.0
	scene.Add(pt.NewSphere(pt.V(-1, -1, 0.5).Normalize().MulScalar(d), 0.25, light))
	scene.Add(pt.NewSphere(pt.V(0, -1, 0.25).Normalize().MulScalar(d), 0.25, light))
	scene.Add(pt.NewSphere(pt.V(-1, 1, 0).Normalize().MulScalar(d), 0.25, light))

	material := pt.GlossyMaterial(pt.HexColor(0x468966), 1.2, pt.Radians(20))

	s0 := NewPtSDF(s)
	s0 = pt.NewTransformSDF(s0, pt.Translate(pt.V(0, 0, 0.2)))
	s0 = pt.NewTransformSDF(s0, pt.Rotate(pt.V(0, 0, 1), pt.Radians(30)))

	scene.Add(pt.NewSDFShape(s0, material))

	floor := pt.GlossyMaterial(pt.HexColor(0xFFF0A5), 1.2, pt.Radians(20))
	scene.Add(pt.NewPlane(pt.V(0, 0, -0.5), pt.V(0, 0, 1), floor))

	camera := pt.LookAt(pt.V(-3, 0, 1), pt.V(0, 0, 0), pt.V(0, 0, 1), 35)
	sampler := pt.NewSampler(4, 4)
	sampler.LightMode = pt.LightModeAll
	sampler.SpecularMode = pt.SpecularModeAll
	renderer := pt.NewRenderer(&scene, &camera, sampler, 800, 600)
	renderer.IterativeRender("out%03d.png", 10)
}
