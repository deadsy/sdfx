//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package sdf

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/deadsy/pt/pt"
	"github.com/yofu/dxf"
)

//-----------------------------------------------------------------------------

func RenderPNG(s SDF3, render_floor bool) {

	scene := pt.Scene{}

	light := pt.LightMaterial(pt.White, 180)

	d := 4.0
	scene.Add(pt.NewSphere(pt.V(-1, -1, 0.5).Normalize().MulScalar(d), 0.25, light))
	scene.Add(pt.NewSphere(pt.V(0, -1, 0.25).Normalize().MulScalar(d), 0.25, light))
	scene.Add(pt.NewSphere(pt.V(-1, 1, 0).Normalize().MulScalar(d), 0.25, light))

	material := pt.GlossyMaterial(pt.HexColor(0x468966), 1.2, pt.Radians(20))

	s0 := NewPtSDF(s)
	//s0 = pt.NewTransformSDF(s0, pt.Translate(pt.V(0, 0, 0.2)))
	//s0 = pt.NewTransformSDF(s0, pt.Rotate(pt.V(0, 0, 1), pt.Radians(30)))

	scene.Add(pt.NewSDFShape(s0, material))

	if render_floor {
		bb := s0.BoundingBox()
		z_min := bb.Min.Z
		z_height := bb.Max.Z - bb.Min.Z
		z_gap := z_height * 0.1 // 10% of height

		floor := pt.GlossyMaterial(pt.HexColor(0xFFF0A5), 1.2, pt.Radians(20))
		floor_plane := pt.V(0, 0, z_min-z_gap)
		floor_normal := pt.V(0, 0, 1)

		scene.Add(pt.NewPlane(floor_plane, floor_normal, floor))
	}

	camera := pt.LookAt(pt.V(-3, 0, 1), pt.V(0, 0, 0), pt.V(0, 0, 1), 35)
	sampler := pt.NewSampler(4, 4)
	sampler.LightMode = pt.LightModeAll
	sampler.SpecularMode = pt.SpecularModeAll
	renderer := pt.NewRenderer(&scene, &camera, sampler, 800, 600)
	renderer.IterativeRender("out%03d.png", 10)
}

//-----------------------------------------------------------------------------

// Render an SDF3 as an STL triangle mesh file.
func RenderSTL(s SDF3, path string) {

	mesh_cells := 200.0 // number of mesh cells on major axis of SDF3 bounding box

	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0_size := bb0.Size()
	mesh_inc := bb0_size.MaxComponent() / mesh_cells
	bb1_size := bb0_size.DivScalar(mesh_inc)
	bb1_size = bb1_size.Ceil().AddScalar(3)
	cells := bb1_size.ToV3i()
	bb1_size = bb1_size.MulScalar(mesh_inc)
	bb := NewBox3(bb0.Center(), bb1_size)

	fmt.Printf("rendering %s (%dx%dx%d)\n", path, cells[0], cells[1], cells[2])

	m := NewSDFMesh(s, bb, mesh_inc)
	err := SaveSTL(path, m)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------

// Render an SDF2 as a PNG image file.
func SDF2_RenderPNG(s SDF2, path string) {

	region_size := 400.0 // number of pixels on major axis of SDF2 bounding box
	border_size := 40.0  // border pixels on either side of bounding box

	// work out the region we will sample
	bb0 := s.BoundingBox()
	bb0_size := bb0.Size()
	inc := bb0_size.MaxComponent() / region_size
	bb1_size := bb0_size.DivScalar(inc)
	bb1_size = bb1_size.Ceil().AddScalar(2 * border_size)
	pixels := bb1_size.ToV2i()
	bb1_size = bb1_size.MulScalar(inc)
	bb := NewBox2(bb0.Center(), bb1_size)

	fmt.Printf("rendering %s (%dx%d)\n", path, pixels[0], pixels[1])

	// sample the distance field
	var dmax, dmin float64
	distance := make([]float64, pixels[0]*pixels[1])
	dx := inc / 2
	dy := inc / 2
	xofs := 0
	for x := 0; x < pixels[0]; x++ {
		for y := 0; y < pixels[1]; y++ {
			d := s.Evaluate(bb.Min.Add(V2{dx, dy}))
			dmax = Max(dmax, d)
			dmin = Min(dmin, d)
			distance[xofs+y] = d
			dy += inc
		}
		dy = inc / 2
		dx += inc
		xofs += pixels[1]
	}

	img := image.NewRGBA(image.Rect(0, 0, pixels[0]-1, pixels[1]-1))

	xofs = 0
	for x := 0; x < pixels[0]; x++ {
		for y := 0; y < pixels[1]; y++ {
			val := 255.0 * ((distance[xofs+y] - dmin) / (dmax - dmin))
			img.Set(x, y, color.Gray{uint8(val)})
		}
		xofs += pixels[1]
	}

	outpng, err := os.Create(path)
	if err != nil {
		panic("Error storing png: " + err.Error())
	}
	defer outpng.Close()
	png.Encode(outpng, img)
}

//-----------------------------------------------------------------------------

func RenderDXF(path string, vlist []V2) {
	d := dxf.NewDrawing()

	for i := 0; i < len(vlist)-1; i++ {
		p0 := vlist[i]
		p1 := vlist[i+1]
		d.Line(p0.X, p0.Y, 0, p1.X, p1.Y, 0)
	}

	err := d.SaveAs(path)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

//-----------------------------------------------------------------------------
