package sdf

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV3CompareToZero(t *testing.T) {
	tests := []struct {
		name   string
		test   func(V3) bool
		got    V3
		expect bool
	}{
		{"LTZero", (V3).LTZero, V3{1.0, 2.0, 3.0}, false},
		{"LTZero", (V3).LTZero, V3{0.0, 2.0, 3.0}, false},
		{"LTZero", (V3).LTZero, V3{1.0, 0.0, 3.0}, false},
		{"LTZero", (V3).LTZero, V3{1.0, 2.0, 0.0}, false},
		{"LTZero", (V3).LTZero, V3{-1.0, 2.0, 3.0}, true},
		{"LTZero", (V3).LTZero, V3{1.0, -2.0, 3.0}, true},
		{"LTZero", (V3).LTZero, V3{1.0, 2.0, -3.0}, true},

		{"LTEZero", (V3).LTEZero, V3{1.0, 2.0, 3.0}, false},
		{"LTEZero", (V3).LTEZero, V3{0.0, 2.0, 3.0}, true},
		{"LTEZero", (V3).LTEZero, V3{1.0, 0.0, 3.0}, true},
		{"LTEZero", (V3).LTEZero, V3{1.0, 2.0, 0.0}, true},
		{"LTEZero", (V3).LTEZero, V3{-1.0, 2.0, 3.0}, true},
		{"LTEZero", (V3).LTEZero, V3{1.0, -2.0, 3.0}, true},
		{"LTEZero", (V3).LTEZero, V3{1.0, 2.0, -3.0}, true},
	}

	i := 0
	var last string
	for _, test := range tests {
		if last != test.name {
			i = 0
		}

		assert.Equalf(t, test.expect, test.test(test.got), "%s test #%d", test.name, i)
		last = test.name
	}
}

func TestV2CompareToZero(t *testing.T) {
	tests := []struct {
		name   string
		test   func(V2) bool
		got    V2
		expect bool
	}{
		{"LTZero", (V2).LTZero, V2{1.0, 2.0}, false},
		{"LTZero", (V2).LTZero, V2{0.0, 2.0}, false},
		{"LTZero", (V2).LTZero, V2{1.0, 0.0}, false},
		{"LTZero", (V2).LTZero, V2{-1.0, 2.0}, true},
		{"LTZero", (V2).LTZero, V2{1.0, -2.0}, true},

		{"LTEZero", (V2).LTEZero, V2{1.0, 2.0}, false},
		{"LTEZero", (V2).LTEZero, V2{0.0, 2.0}, true},
		{"LTEZero", (V2).LTEZero, V2{1.0, 0.0}, true},
		{"LTEZero", (V2).LTEZero, V2{-1.0, 2.0}, true},
		{"LTEZero", (V2).LTEZero, V2{1.0, -2.0}, true},
	}

	i := 0
	var last string
	for _, test := range tests {
		if last != test.name {
			i = 0
		}

		assert.Equalf(t, test.expect, test.test(test.got), "%s test #%d", test.name, i)
		last = test.name
	}
}

func TestV3Clamp(t *testing.T) {
	a := V3{12.3, 45.6, 78.9}
	b := V3{123.4, 156.7, 189.0}
	tests := []struct {
		got    V3
		expect V3
	}{
		{V3{0.0, 0.0, 0.0}, a},
		{V3{200.0, 200.0, 200.0}, b},
		{a, a},
		{b, b},
	}

	for i, test := range tests {
		assert.Equalf(t, test.expect, test.got.Clamp(a, b), "test #%d", i)
	}
}

func TestV2Clamp(t *testing.T) {
	a := V2{12.3, 45.6}
	b := V2{123.4, 156.7}
	tests := []struct {
		got    V2
		expect V2
	}{
		{V2{0.0, 0.0}, a},
		{V2{200.0, 200.0}, b},
		{a, a},
		{b, b},
	}

	for i, test := range tests {
		assert.Equalf(t, test.expect, test.got.Clamp(a, b), "test #%d", i)
	}
}

func TestBox3RandomSet(t *testing.T) {
	a := V3{12.3, 45.6, 78.9}
	b := V3{123.4, 156.7, 189.0}
	bb := Box3{a, b}
	v3s := bb.RandomSet(42)
	assert.Equal(t, 42, len(v3s), "RandomSet generated correct number of points")

	for _, v := range v3s {
		assert.Equal(t, v, v.Clamp(a, b), "Clamping should never change RandomSet points")
	}
}

func TestBox2RandomSet(t *testing.T) {
	a := V2{12.3, 45.6}
	b := V2{123.4, 156.7}
	bb := Box2{a, b}
	v3s := bb.RandomSet(42)
	assert.Equal(t, 42, len(v3s), "RandomSet generated correct number of points")

	for _, v := range v3s {
		assert.Equal(t, v, v.Clamp(a, b), "Clamping should never change RandomSet points")
	}
}

func TestV3MatrixOps(t *testing.T) {
	a := V3{3.0, 5.0, 7.0}
	b := V3{11.0, 13.0, 17.0}
	assert.Equal(t, 3.0*11.0+5.0*13.0+7.0*17.0, a.Dot(b), "a.b works")
	assert.Equal(t, V3{
		5.0*17.0 - 7.0*13.0,
		7.0*11.0 - 3.0*17.0,
		3.0*13.0 - 5.0*11.0,
	}, a.Cross(b), "axb works")
}

func TestV2MatrixOps(t *testing.T) {
	a := V2{3.0, 5.0}
	b := V2{11.0, 13.0}
	assert.Equal(t, 3.0*11.0+5.0*13.0, a.Dot(b), "a.b works")
	assert.Equal(t, 3.0*13.0-5.0*11.0, a.Cross(b), "axb works")
}

func TestV2Colinearity(t *testing.T) {
	a := V2{37.4, 88.8}
	m := V2{3.0, 5.0}
	b := a.Add(m.MulScalar(16.0))
	c := a.Sub(m.MulScalar(7.0))
	d := V2{55.5, 66.6}

	assert.True(t, colinearFast(a, b, c, 0.0001), "ABC are colinear fast")
	assert.True(t, colinearFast(a, c, b, 0.0001), "ACB are colienar fast")
	assert.True(t, colinearFast(b, a, c, 0.0001), "BAC are colinear fast")
	assert.True(t, colinearFast(b, c, a, 0.0001), "BCA are colienar fast")
	assert.True(t, colinearFast(c, a, b, 0.0001), "CAB are colinear fast")
	assert.True(t, colinearFast(c, b, a, 0.0001), "CBA are colinear fast")

	assert.False(t, colinearFast(a, b, d, 0.0001), "ABD are not colinear fast")
	assert.False(t, colinearFast(a, c, d, 0.0001), "ACD are not colienar fast")
	assert.False(t, colinearFast(b, a, d, 0.0001), "BAD are not colinear fast")
	assert.False(t, colinearFast(b, c, d, 0.0001), "BCD are not colienar fast")
	assert.False(t, colinearFast(c, a, d, 0.0001), "CAD are not colinear fast")
	assert.False(t, colinearFast(c, b, d, 0.0001), "CBD are not colinear fast")

	assert.True(t, colinearSlow(a, b, c, 0.0001), "ABC are colinear slow")
	assert.True(t, colinearSlow(a, c, b, 0.0001), "ACB are colienar slow")
	assert.True(t, colinearSlow(b, a, c, 0.0001), "BAC are colinear slow")
	assert.True(t, colinearSlow(b, c, a, 0.0001), "BCA are colienar slow")
	assert.True(t, colinearSlow(c, a, b, 0.0001), "CAB are colinear slow")
	assert.True(t, colinearSlow(c, b, a, 0.0001), "CBA are colinear slow")

	assert.False(t, colinearSlow(a, b, d, 0.0001), "ABD are not colinear slow")
	assert.False(t, colinearSlow(a, c, d, 0.0001), "ACD are not colienar slow")
	assert.False(t, colinearSlow(b, a, d, 0.0001), "BAD are not colinear slow")
	assert.False(t, colinearSlow(b, c, d, 0.0001), "BCD are not colienar slow")
	assert.False(t, colinearSlow(c, a, d, 0.0001), "CAD are not colinear slow")
	assert.False(t, colinearSlow(c, b, d, 0.0001), "CBD are not colinear slow")
}

func TestV3ScalarOps(t *testing.T) {
	a := 42.0
	v := V3{0.0, 1.0, 2.0}
	assert.Equal(t, V3{0.0 + a, 1.0 + a, 2.0 + a}, v.AddScalar(a), "v+a works")
	assert.Equal(t, V3{0.0 - a, 1.0 - a, 2.0 - a}, v.SubScalar(a), "v-a works")
	assert.Equal(t, V3{0.0 * a, 1.0 * a, 2.0 * a}, v.MulScalar(a), "v*a works")
	assert.Equal(t, V3{0.0 / a, 1.0 / a, 2.0 / a}, v.DivScalar(a), "v/a works")
}

func TestV2ScalarOps(t *testing.T) {
	a := 42.0
	v := V2{0.0, 1.0}
	assert.Equal(t, V2{0.0 + a, 1.0 + a}, v.AddScalar(a), "v+a works")
	assert.Equal(t, V2{0.0 - a, 1.0 - a}, v.SubScalar(a), "v-a works")
	assert.Equal(t, V2{0.0 * a, 1.0 * a}, v.MulScalar(a), "v*a works")
	assert.Equal(t, V2{0.0 / a, 1.0 / a}, v.DivScalar(a), "v/a works")
}

func TestV3Abs(t *testing.T) {
	assert.Equal(t, V3{1.0, 2.0, 3.0}, V3{-1.0, -2.0, -3.0}.Abs(), "abs(v) works")
}

func TestV2Abs(t *testing.T) {
	assert.Equal(t, V2{1.0, 2.0}, V2{-1.0, -2.0}.Abs(), "abs(v) works")
}

func TestV3Ceil(t *testing.T) {
	assert.Equal(t, V3{math.Ceil(1.1), math.Ceil(2.2), math.Ceil(3.3)}, V3{1.1, 2.2, 3.3}.Ceil(), "ceil(v) works")
}

func TestV2Ceil(t *testing.T) {
	assert.Equal(t, V2{math.Ceil(1.1), math.Ceil(2.2)}, V2{1.1, 2.2}.Ceil(), "ceil(v) works")
}

func TestV3Ops(t *testing.T) {
	a := V3{2.0, 11.0, 5.0}
	b := V3{7.0, 3.0, 13.0}

	assert.Equal(t, V3{2.0, 3.0, 5.0}, a.Min(b), "min(a, b) works")
	assert.Equal(t, V3{2.0, 3.0, 5.0}, b.Min(a), "min(b, a) works")

	assert.Equal(t, V3{7.0, 11.0, 13.0}, a.Max(b), "max(a, b) works")
	assert.Equal(t, V3{7.0, 11.0, 13.0}, b.Max(a), "max(b, a) works")

	assert.Equal(t, V3{2.0 + 7.0, 11.0 + 3.0, 5.0 + 13.0}, a.Add(b), "a+b works")
	assert.Equal(t, V3{7.0 + 2.0, 3.0 + 11.0, 13.0 + 5.0}, b.Add(a), "b+a works")

	assert.Equal(t, V3{2.0 - 7.0, 11.0 - 3.0, 5.0 - 13.0}, a.Sub(b), "a-b works")
	assert.Equal(t, V3{7.0 - 2.0, 3.0 - 11.0, 13.0 - 5.0}, b.Sub(a), "b-a works")

	assert.Equal(t, V3{2.0 * 7.0, 11.0 * 3.0, 5.0 * 13.0}, a.Mul(b), "a*b works")
	assert.Equal(t, V3{7.0 * 2.0, 3.0 * 11.0, 13.0 * 5.0}, b.Mul(a), "b*a works")

	assert.Equal(t, V3{2.0 / 7.0, 11.0 / 3.0, 5.0 / 13.0}, a.Div(b), "a/b works")
	assert.Equal(t, V3{7.0 / 2.0, 3.0 / 11.0, 13.0 / 5.0}, b.Div(a), "b/a works")

	assert.Equal(t, V3{-2.0, -11.0, -5.0}, a.Neg(), "-a works")
	assert.Equal(t, V3{-7.0, -3.0, -13.0}, b.Neg(), "-b works")
}

func TestV2Ops(t *testing.T) {
	a := V2{2.0, 11.0}
	b := V2{7.0, 3.0}

	assert.Equal(t, V2{2.0, 3.0}, a.Min(b), "min(a, b) works")
	assert.Equal(t, V2{2.0, 3.0}, b.Min(a), "min(b, a) works")

	assert.Equal(t, V2{7.0, 11.0}, a.Max(b), "max(a, b) works")
	assert.Equal(t, V2{7.0, 11.0}, b.Max(a), "max(b, a) works")

	assert.Equal(t, V2{2.0 + 7.0, 11.0 + 3.0}, a.Add(b), "a+b works")
	assert.Equal(t, V2{7.0 + 2.0, 3.0 + 11.0}, b.Add(a), "b+a works")

	assert.Equal(t, V2{2.0 - 7.0, 11.0 - 3.0}, a.Sub(b), "a-b works")
	assert.Equal(t, V2{7.0 - 2.0, 3.0 - 11.0}, b.Sub(a), "b-a works")

	assert.Equal(t, V2{2.0 * 7.0, 11.0 * 3.0}, a.Mul(b), "a*b works")
	assert.Equal(t, V2{7.0 * 2.0, 3.0 * 11.0}, b.Mul(a), "b*a works")

	assert.Equal(t, V2{2.0 / 7.0, 11.0 / 3.0}, a.Div(b), "a/b works")
	assert.Equal(t, V2{7.0 / 2.0, 3.0 / 11.0}, b.Div(a), "b/a works")

	assert.Equal(t, V2{-2.0, -11.0}, a.Neg(), "-a works")
	assert.Equal(t, V2{-7.0, -3.0}, b.Neg(), "-b works")
}

func TestV3SetOps(t *testing.T) {
	v3s := V3Set{
		{1.0, 99.0, 55.0},
		{95.0, 44.0, 3.0},
		{66.0, 7.0, 88.0},
	}

	assert.Equal(t, V3{1.0, 7.0, 3.0}, v3s.Min(), "min(vs) works")
	assert.Equal(t, V3{95.0, 99.0, 88.0}, v3s.Max(), "max(vs) works")
}

func TestV2SetOps(t *testing.T) {
	v2s := V2Set{
		{1.0, 99.0},
		{95.0, 44.0},
		{66.0, 7.0},
	}

	assert.Equal(t, V2{1.0, 7.0}, v2s.Min(), "min(vs) works")
	assert.Equal(t, V2{95.0, 99.0}, v2s.Max(), "max(vs) works")
}

func TestV3VectorOps(t *testing.T) {
	d := math.Sqrt(2.0*2.0 + 3.0*3.0 + 5.0*5.0)
	assert.Equal(t, d, V3{2.0, 3.0, 5.0}.Length(), "length(v) works")
	assert.Equal(t, 2.0*2.0+3.0*3.0+5.0*5.0, V3{2.0, 3.0, 5.0}.Length2(), "length(v)^2 works")

	assert.Equal(t, 2.0, V3{2.0, 3.0, 5.0}.MinComponent(), "min(v.x, v.y, v.z) works")
	assert.Equal(t, 5.0, V3{2.0, 3.0, 5.0}.MaxComponent(), "max(v.x, v.y, v.z) works")

	assert.Equal(t, V3{2.0 / d, 3.0 / d, 5.0 / d}, V3{2.0, 3.0, 5.0}.Normalize(), "normalize(v) works")
	assert.InDelta(t, V3{2.0, 3.0, 5.0}.Normalize().Length(), 1.0, 0.0001, "length(normalize(v)) == 1")
}

func TestV2VectorOps(t *testing.T) {
	d := math.Sqrt(2.0*2.0 + 3.0*3.0)
	assert.Equal(t, d, V2{2.0, 3.0}.Length(), "length(v) works")
	assert.Equal(t, 2.0*2.0+3.0*3.0, V2{2.0, 3.0}.Length2(), "length(v)^2 works")

	assert.Equal(t, 2.0, V2{2.0, 3.0}.MinComponent(), "min(v.x, v.y) works")
	assert.Equal(t, 3.0, V2{2.0, 3.0}.MaxComponent(), "max(v.x, v.y) works")

	assert.Equal(t, V2{2.0 / d, 3.0 / d}, V2{2.0, 3.0}.Normalize(), "normalize(v) works")
	assert.InDelta(t, V2{2.0, 3.0}.Normalize().Length(), 1.0, 0.0001, "length(normalize(v)) == 1")
}
