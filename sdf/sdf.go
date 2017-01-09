//-----------------------------------------------------------------------------
/*
Derived from the hg_sdf library
http://mercury.sexy/hg_sdf/
*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	"github.com/deadsy/sdfx/vec"
)

//-----------------------------------------------------------------------------

const PHI = math.Phi
const PI = math.Pi

//-----------------------------------------------------------------------------
// Primitive Distance Functions

// Sphere
func Sphere(p vec.V3, r float64) float64 {
	return p.Length() - r
}

// Plane with normal n (n is normalized) at some distance from the origin
func Plane(p, n vec.V3, distanceFromOrigin float64) float64 {
	return p.Dot(n) + distanceFromOrigin
}

// Cheap Box: distance to corners is overestimated
func BoxCheap(p, b vec.V3) float64 {
	return p.Abs().Sub(b).Vmax()
}

// Box: correct distance to corners
func Box(p, b vec.V3) float64 {
	d := p.Abs().Sub(b)
	return d.Max(vec.V3{0, 0, 0}).Length() + d.Min(vec.V3{0, 0, 0}).Vmax()
}

// Same as above, but in two dimensions (an endless box)
func Box2Cheap(p, b vec.V2) float64 {
	return p.Abs().Sub(b).Vmax()
}

func Box2(p, b vec.V2) float64 {
	d := p.Abs().Sub(b)
	return d.Max(vec.V2{0, 0}).Length() + d.Min(vec.V2{0, 0}).Vmax()
}

// Endless "corner"
func Corner(p vec.V2) float64 {
	return p.Max(vec.V2{0, 0}).Length() + p.Min(vec.V2{0, 0}).Vmax()
}

// Blobby ball object. You've probably seen it somewhere. This is not a correct distance bound, beware.
func Blob(p vec.V3) float64 {
	p = p.Abs()
	if p[0] < math.Max(p[1], p[2]) {
		p = vec.V3{p[1], p[2], p[0]}
	}
	if p[0] < math.Max(p[1], p[2]) {
		p = vec.V3{p[1], p[2], p[0]}
	}
	b := math.Max(math.Max(math.Max(
		p.Dot(vec.V3{1, 1, 1}.Normalize()),
		vec.V2{p[0], p[2]}.Dot(vec.V2{PHI + 1, 1}.Normalize())),
		vec.V2{p[1], p[0]}.Dot(vec.V2{1, PHI}.Normalize())),
		vec.V2{p[0], p[2]}.Dot(vec.V2{1, PHI}.Normalize()))
	l := p.Length()
	return l - 1.5 - 0.2*(1.5/2)*math.Cos(math.Min(math.Sqrt(1.01-b/l)*(PI/0.25), PI))
}

// Cylinder standing upright on the xz plane
func Cylinder(p vec.V3, r, height float64) float64 {
	d := vec.V2{p[0], p[2]}.Length() - r
	return math.Max(d, math.Abs(p[1])-height)
}

// Capsule version 1: A Cylinder with round caps on both sides
func Capsule1(p vec.V3, r, c float64) float64 {
	return vec.Mix(vec.V2{p[0], p[2]}.Length()-r, vec.V3{p[0], math.Abs(p[1]) - c, p[2]}.Length()-r, vec.Step(c, math.Abs(p[1])))
}

// Distance to line segment between <a> and <b>, used for fCapsule() version 2below
func LineSegment(p, a, b vec.V3) float64 {
	ab := b.Sub(a)
	t := vec.Saturate(p.Sub(a).Dot(ab) / ab.Dot(ab))
	return ab.Scale(t).Sum(a).Sub(p).Length()
}

// Capsule version 2: between two end points <a> and <b> with radius r
func Capsule2(p, a, b vec.V3, r float64) float64 {
	return LineSegment(p, a, b) - r
}

// Torus in the XZ-plane
func Torus(p vec.V3, smallRadius, largeRadius float64) float64 {
	return vec.V2{vec.V2{p[0], p[2]}.Length() - largeRadius, p[1]}.Length() - smallRadius
}

// A circle line. Can also be used to make a torus by subtracting the smaller radius of the torus.
func Circle(p vec.V3, r float64) float64 {
	l := vec.V2{p[0], p[2]}.Length() - r
	return vec.V2{p[1], l}.Length()
}

// A circular disc with no thickness (i.e. a cylinder with no height).
// Subtract some value to make a flat disc with rounded edge.
func Disc(p vec.V3, r float64) float64 {
	l := vec.V2{p[0], p[2]}.Length() - r
	if l < 0 {
		return math.Abs(p[1])
	}
	return vec.V2{p[1], l}.Length()
}

// Hexagonal prism, circumcircle variant
func HexagonCircumcircle(p vec.V3, h vec.V2) float64 {
	q := p.Abs()
	return math.Max(q[1]-h[1], math.Max(q[0]*math.Sqrt(3)*0.5+q[2]*0.5, q[2])-h[0])
}

// Hexagonal prism, incircle variant
func HexagonIncircle(p vec.V3, h vec.V2) float64 {
	return HexagonCircumcircle(p, vec.V2{h[0] * math.Sqrt(3) * 0.5, h[1]})
}

// Cone with correct distances to tip and base circle. Y is up, 0 is in the middle of the base.
func Cone(p vec.V3, radius, height float64) float64 {

	q := vec.V2{vec.V2{p[0], p[2]}.Length(), p[1]}
	tip := q.Sub(vec.V2{0, height})
	mantleDir := vec.V2{height, radius}.Normalize()

	mantle := tip.Dot(mantleDir)
	d := math.Max(mantle, -q[1])
	projected := tip.Dot(vec.V2{mantleDir[1], -mantleDir[0]})

	// distance to tip
	if (q[1] > height) && (projected < 0) {
		d = math.Max(d, tip.Length())
	}

	// distance to base ring
	if (q[0] > radius) && (projected > vec.V2{height, radius}.Length()) {
		d = math.Max(d, q.Sub(vec.V2{radius, 0}).Length())
	}

	return d
}

//-----------------------------------------------------------------------------
// "Generalized Distance Functions" by Akleman and Chen.
// see the Paper at https://www.viz.tamu.edu/faculty/ergun/research/implicitmodeling/papers/sm99.pdf

var GDFVectors = [19]vec.V3{

	vec.V3{1, 0, 0}.Normalize(),
	vec.V3{0, 1, 0}.Normalize(),
	vec.V3{0, 0, 1}.Normalize(),

	vec.V3{1, 1, 1}.Normalize(),
	vec.V3{-1, 1, 1}.Normalize(),
	vec.V3{1, -1, 1}.Normalize(),
	vec.V3{1, 1, -1}.Normalize(),

	vec.V3{0, 1, PHI + 1}.Normalize(),
	vec.V3{0, -1, PHI + 1}.Normalize(),
	vec.V3{PHI + 1, 0, 1}.Normalize(),
	vec.V3{-PHI - 1, 0, 1}.Normalize(),
	vec.V3{1, PHI + 1, 0}.Normalize(),
	vec.V3{-1, PHI + 1, 0}.Normalize(),

	vec.V3{0, PHI, 1}.Normalize(),
	vec.V3{0, -PHI, 1}.Normalize(),
	vec.V3{1, 0, PHI}.Normalize(),
	vec.V3{-1, 0, PHI}.Normalize(),
	vec.V3{PHI, 1, 0}.Normalize(),
	vec.V3{-PHI, 1, 0}.Normalize(),
}

/*

// Version with variable exponent.
// This is slow and does not produce correct distances, but allows for bulging of objects.
float fGDF(vec3 p, float r, float e, int begin, int end) {
	float d = 0;
	for (int i = begin; i <= end; ++i)
		d += pow(abs(dot(p, GDFVectors[i])), e);
	return pow(d, 1/e) - r;
}

// Version with without exponent, creates objects with sharp edges and flat faces
float fGDF(vec3 p, float r, int begin, int end) {
	float d = 0;
	for (int i = begin; i <= end; ++i)
		d = max(d, abs(dot(p, GDFVectors[i])));
	return d - r;
}

// Primitives follow:

float fOctahedron(vec3 p, float r, float e) {
	return fGDF(p, r, e, 3, 6);
}

float fDodecahedron(vec3 p, float r, float e) {
	return fGDF(p, r, e, 13, 18);
}

float fIcosahedron(vec3 p, float r, float e) {
	return fGDF(p, r, e, 3, 12);
}

float fTruncatedOctahedron(vec3 p, float r, float e) {
	return fGDF(p, r, e, 0, 6);
}

float fTruncatedIcosahedron(vec3 p, float r, float e) {
	return fGDF(p, r, e, 3, 18);
}

float fOctahedron(vec3 p, float r) {
	return fGDF(p, r, 3, 6);
}

float fDodecahedron(vec3 p, float r) {
	return fGDF(p, r, 13, 18);
}

float fIcosahedron(vec3 p, float r) {
	return fGDF(p, r, 3, 12);
}

float fTruncatedOctahedron(vec3 p, float r) {
	return fGDF(p, r, 0, 6);
}

float fTruncatedIcosahedron(vec3 p, float r) {
	return fGDF(p, r, 3, 18);
}

*/

//-----------------------------------------------------------------------------
