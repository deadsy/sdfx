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
// Scalar Functions (similar to GLSL counterparts)

// Return 0 if x < edge, else 1
func step(edge, x float64) float64 {
	if x < edge {
		return 0
	}
	return 1
}

// Linear Interpolation
func mix(x, y, a float64) float64 {
	return (x * (1 - a)) + (y * a)
}

func clamp(x, a, b float64) float64 {
	return math.Min(math.Max(x, a), b)
}

func saturate(x float64) float64 {
	return clamp(x, 0, 1)
}

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
	return mix(vec.V2{p[0], p[2]}.Length()-r, vec.V3{p[0], math.Abs(p[1]) - c, p[2]}.Length()-r, step(c, math.Abs(p[1])))
}

// Distance to line segment between <a> and <b>, used for fCapsule() version 2below
func LineSegment(p, a, b vec.V3) float64 {
	ab := b.Sub(a)
	t := saturate(p.Sub(a).Dot(ab) / ab.Dot(ab))
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

// Version with variable exponent.
// This is slow and does not produce correct distances, but allows for bulging of objects.
func GDF1(p vec.V3, r, e float64, begin, end int) float64 {
	var d float64
	for i := begin; i <= end; i++ {
		d += math.Pow(math.Abs(p.Dot(GDFVectors[i])), e)
	}
	return math.Pow(d, 1/e) - r
}

// Version with without exponent, creates objects with sharp edges and flat faces
func GDF2(p vec.V3, r float64, begin, end int) float64 {
	var d float64
	for i := begin; i <= end; i++ {
		d = math.Max(d, math.Abs(p.Dot(GDFVectors[i])))
	}
	return d - r
}

func Octahedron1(p vec.V3, r, e float64) float64 {
	return GDF1(p, r, e, 3, 6)
}

func Dodecahedron1(p vec.V3, r, e float64) float64 {
	return GDF1(p, r, e, 13, 18)
}

func Icosahedron1(p vec.V3, r, e float64) float64 {
	return GDF1(p, r, e, 3, 12)
}

func TruncatedOctahedron1(p vec.V3, r, e float64) float64 {
	return GDF1(p, r, e, 0, 6)
}

func TruncatedIcosahedron1(p vec.V3, r, e float64) float64 {
	return GDF1(p, r, e, 3, 18)
}

func Octahedron2(p vec.V3, r float64) float64 {
	return GDF2(p, r, 3, 6)
}

func Dodecahedron2(p vec.V3, r float64) float64 {
	return GDF2(p, r, 13, 18)
}

func Icosahedron2(p vec.V3, r float64) float64 {
	return GDF2(p, r, 3, 12)
}

func TruncatedOctahedron2(p vec.V3, r float64) float64 {
	return GDF2(p, r, 0, 6)
}

func TruncatedIcosahedron2(p vec.V3, r float64) float64 {
	return GDF2(p, r, 3, 18)
}

//-----------------------------------------------------------------------------
// Domain Manipulation Operators

/*

// Rotate around a coordinate axis (i.e. in a plane perpendicular to that axis) by angle <a>.
// Read like this: R(p.xz, a) rotates "x towards z".
// This is fast if <a> is a compile-time constant and slower (but still practical) if not.
void pR(inout vec2 p, float a) {
	p = cos(a)*p + sin(a)*vec2(p.y, -p.x);
}

// Shortcut for 45-degrees rotation
void pR45(inout vec2 p) {
	p = (p + vec2(p.y, -p.x))*sqrt(0.5);
}

// Repeat space along one axis. Use like this to repeat along the x axis:
// <float cell = pMod1(p.x,5);> - using the return value is optional.
float pMod1(inout float p, float size) {
	float halfsize = size*0.5;
	float c = floor((p + halfsize)/size);
	p = mod(p + halfsize, size) - halfsize;
	return c;
}

// Same, but mirror every second cell so they match at the boundaries
float pModMirror1(inout float p, float size) {
	float halfsize = size*0.5;
	float c = floor((p + halfsize)/size);
	p = mod(p + halfsize,size) - halfsize;
	p *= mod(c, 2.0)*2 - 1;
	return c;
}

// Repeat the domain only in positive direction. Everything in the negative half-space is unchanged.
float pModSingle1(inout float p, float size) {
	float halfsize = size*0.5;
	float c = floor((p + halfsize)/size);
	if (p >= 0)
		p = mod(p + halfsize, size) - halfsize;
	return c;
}

// Repeat only a few times: from indices <start> to <stop> (similar to above, but more flexible)
float pModInterval1(inout float p, float size, float start, float stop) {
	float halfsize = size*0.5;
	float c = floor((p + halfsize)/size);
	p = mod(p+halfsize, size) - halfsize;
	if (c > stop) { //yes, this might not be the best thing numerically.
		p += size*(c - stop);
		c = stop;
	}
	if (c <start) {
		p += size*(c - start);
		c = start;
	}
	return c;
}


// Repeat around the origin by a fixed angle.
// For easier use, num of repetitions is use to specify the angle.
float pModPolar(inout vec2 p, float repetitions) {
	float angle = 2*PI/repetitions;
	float a = atan(p.y, p.x) + angle/2.;
	float r = length(p);
	float c = floor(a/angle);
	a = mod(a,angle) - angle/2.;
	p = vec2(cos(a), sin(a))*r;
	// For an odd number of repetitions, fix cell index of the cell in -x direction
	// (cell index would be e.g. -5 and 5 in the two halves of the cell):
	if (abs(c) >= (repetitions/2)) c = abs(c);
	return c;
}

// Repeat in two dimensions
vec2 pMod2(inout vec2 p, vec2 size) {
	vec2 c = floor((p + size*0.5)/size);
	p = mod(p + size*0.5,size) - size*0.5;
	return c;
}

// Same, but mirror every second cell so all boundaries match
vec2 pModMirror2(inout vec2 p, vec2 size) {
	vec2 halfsize = size*0.5;
	vec2 c = floor((p + halfsize)/size);
	p = mod(p + halfsize, size) - halfsize;
	p *= mod(c,vec2(2))*2 - vec2(1);
	return c;
}

// Same, but mirror every second cell at the diagonal as well
vec2 pModGrid2(inout vec2 p, vec2 size) {
	vec2 c = floor((p + size*0.5)/size);
	p = mod(p + size*0.5, size) - size*0.5;
	p *= mod(c,vec2(2))*2 - vec2(1);
	p -= size/2;
	if (p.x > p.y) p.xy = p.yx;
	return floor(c/2);
}

// Repeat in three dimensions
vec3 pMod3(inout vec3 p, vec3 size) {
	vec3 c = floor((p + size*0.5)/size);
	p = mod(p + size*0.5, size) - size*0.5;
	return c;
}

// Mirror at an axis-aligned plane which is at a specified distance <dist> from the origin.
float pMirror (inout float p, float dist) {
	float s = sgn(p);
	p = abs(p)-dist;
	return s;
}

// Mirror in both dimensions and at the diagonal, yielding one eighth of the space.
// translate by dist before mirroring.
vec2 pMirrorOctant (inout vec2 p, vec2 dist) {
	vec2 s = sgn(p);
	pMirror(p.x, dist.x);
	pMirror(p.y, dist.y);
	if (p.y > p.x)
		p.xy = p.yx;
	return s;
}

// Reflect space at a plane
float pReflect(inout vec3 p, vec3 planeNormal, float offset) {
	float t = dot(p, planeNormal)+offset;
	if (t < 0) {
		p = p - (2*t)*planeNormal;
	}
	return sgn(t);
}

*/

//-----------------------------------------------------------------------------
