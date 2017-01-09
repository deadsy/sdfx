/*
Derived from the hg_sdf library
http://mercury.sexy/hg_sdf/
*/

package sdf

import (
	"math"

	"github.com/deadsy/sdfx/vec"
)

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
		vec.V2{p[0], p[2]}.Dot(vec.V2{math.Phi + 1, 1}.Normalize())),
		vec.V2{p[1], p[0]}.Dot(vec.V2{1, math.Phi}.Normalize())),
		vec.V2{p[0], p[2]}.Dot(vec.V2{1, math.Phi}.Normalize()))
	l := p.Length()
	return l - 1.5 - 0.2*(1.5/2)*math.Cos(math.Min(math.Sqrt(1.01-b/l)*(math.Pi/0.25), math.Pi))
}

// Cylinder standing upright on the xz plane
func Cylinder(p vec.V3, r, height float64) float64 {
	d := vec.V2{p[0], p[2]}.Length() - r
	return math.Max(d, math.Abs(p[1])-height)
}

// Capsule: A Cylinder with round caps on both sides
func Capsule(p vec.V3, r, c float64) float64 {
	return vec.Mix(vec.V2{p[0], p[2]}.Length()-r, vec.V3{p[0], math.Abs(p[1]) - c, p[2]}.Length()-r, vec.Step(c, math.Abs(p[1])))
}

/*

// Distance to line segment between <a> and <b>, used for fCapsule() version 2below
func fLineSegment(vec3 p, vec3 a, vec3 b) float64 {
	vec3 ab = b - a;
	float t = saturate(dot(p - a, ab) / dot(ab, ab));
	return length((ab*t + a) - p);
}

// Capsule version 2: between two end points <a> and <b> with radius r
func fCapsule(vec3 p, vec3 a, vec3 b, float r) float64 {
	return fLineSegment(p, a, b) - r;
}

// Torus in the XZ-plane
func fTorus(vec3 p, float smallRadius, float largeRadius) float64 {
	return length(vec2(length(p.xz) - largeRadius, p.y)) - smallRadius;
}

// A circle line. Can also be used to make a torus by subtracting the smaller radius of the torus.
func fCircle(vec3 p, float r) float64 {
	float l = length(p.xz) - r;
	return length(vec2(p.y, l));
}

// A circular disc with no thickness (i.e. a cylinder with no height).
// Subtract some value to make a flat disc with rounded edge.
func fDisc(vec3 p, float r) float64 {
	float l = length(p.xz) - r;
	return l < 0 ? abs(p.y) : length(vec2(p.y, l));
}

// Hexagonal prism, circumcircle variant
func fHexagonCircumcircle(vec3 p, vec2 h) float64 {
	vec3 q = abs(p);
	return max(q.y - h.y, max(q.x*sqrt(3)*0.5 + q.z*0.5, q.z) - h.x);
	//this is mathematically equivalent to this line, but less efficient:
	//return max(q.y - h.y, max(dot(vec2(cos(PI/3), sin(PI/3)), q.zx), q.z) - h.x);
}

// Hexagonal prism, incircle variant
func fHexagonIncircle(vec3 p, vec2 h) float64 {
	return fHexagonCircumcircle(p, vec2(h.x*sqrt(3)*0.5, h.y));
}

// Cone with correct distances to tip and base circle. Y is up, 0 is in the middle of the base.
func fCone(vec3 p, float radius, float height) float64 {
	vec2 q = vec2(length(p.xz), p.y);
	vec2 tip = q - vec2(0, height);
	vec2 mantleDir = normalize(vec2(height, radius));
	float mantle = dot(tip, mantleDir);
	float d = max(mantle, -q.y);
	float projected = dot(tip, vec2(mantleDir.y, -mantleDir.x));

	// distance to tip
	if ((q.y > height) && (projected < 0)) {
		d = max(d, length(tip));
	}

	// distance to base ring
	if ((q.x > radius) && (projected > length(vec2(height, radius)))) {
		d = max(d, length(q - vec2(radius, 0)));
	}
	return d;
}

*/
