//-----------------------------------------------------------------------------
/*

Pillar Holder

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

var wallThickness = 2.5
var wallHeight = 15.0
var pillarWidth = 33.0
var pillarRadius = 4.0
var feetWidth = 6.0
var baseThickness = 3.0

//-----------------------------------------------------------------------------

func base() sdf.SDF3 {
	w := pillarWidth + 2.0*(feetWidth+wallThickness)
	h := pillarWidth + 2.0*wallThickness
	r := pillarRadius + wallThickness
	base2d := sdf.Box2D(sdf.V2{w, h}, r)
	return sdf.Extrude3D(base2d, baseThickness)
}

func wall(w, r float64) sdf.SDF3 {
	base := sdf.Box2D(sdf.V2{w, w}, r)
	s := sdf.Extrude3D(base, wallHeight)
	ofs := 0.5 * (wallHeight - baseThickness)
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{0, 0, ofs}))
}

func holder() sdf.SDF3 {
	base := base()
	outer := wall(pillarWidth+2.0*wallThickness, pillarRadius+wallThickness)
	inner := wall(pillarWidth, pillarRadius)
	return sdf.Difference3D(sdf.Union3D(base, outer), inner)
}

//-----------------------------------------------------------------------------

func main() {
	s := holder()
	render.RenderSTL(sdf.ScaleUniform3D(s, shrink), 300, "holder.stl")
}

//-----------------------------------------------------------------------------
