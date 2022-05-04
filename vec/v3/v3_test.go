package v3

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareToZero(t *testing.T) {
	tests := []struct {
		name   string
		test   func(Vec) bool
		got    Vec
		expect bool
	}{
		{"LTZero", (Vec).LTZero, Vec{1.0, 2.0, 3.0}, false},
		{"LTZero", (Vec).LTZero, Vec{0.0, 2.0, 3.0}, false},
		{"LTZero", (Vec).LTZero, Vec{1.0, 0.0, 3.0}, false},
		{"LTZero", (Vec).LTZero, Vec{1.0, 2.0, 0.0}, false},
		{"LTZero", (Vec).LTZero, Vec{-1.0, 2.0, 3.0}, true},
		{"LTZero", (Vec).LTZero, Vec{1.0, -2.0, 3.0}, true},
		{"LTZero", (Vec).LTZero, Vec{1.0, 2.0, -3.0}, true},

		{"LTEZero", (Vec).LTEZero, Vec{1.0, 2.0, 3.0}, false},
		{"LTEZero", (Vec).LTEZero, Vec{0.0, 2.0, 3.0}, true},
		{"LTEZero", (Vec).LTEZero, Vec{1.0, 0.0, 3.0}, true},
		{"LTEZero", (Vec).LTEZero, Vec{1.0, 2.0, 0.0}, true},
		{"LTEZero", (Vec).LTEZero, Vec{-1.0, 2.0, 3.0}, true},
		{"LTEZero", (Vec).LTEZero, Vec{1.0, -2.0, 3.0}, true},
		{"LTEZero", (Vec).LTEZero, Vec{1.0, 2.0, -3.0}, true},
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

func TestClamp(t *testing.T) {
	a := Vec{12.3, 45.6, 78.9}
	b := Vec{123.4, 156.7, 189.0}
	tests := []struct {
		got    Vec
		expect Vec
	}{
		{Vec{0.0, 0.0, 0.0}, a},
		{Vec{200.0, 200.0, 200.0}, b},
		{a, a},
		{b, b},
	}

	for i, test := range tests {
		assert.Equalf(t, test.expect, test.got.Clamp(a, b), "test #%d", i)
	}
}

func TestScalarOps(t *testing.T) {
	a := 42.0
	v := Vec{0.0, 1.0, 2.0}
	assert.Equal(t, Vec{0.0 + a, 1.0 + a, 2.0 + a}, v.AddScalar(a), "v+a works")
	assert.Equal(t, Vec{0.0 - a, 1.0 - a, 2.0 - a}, v.SubScalar(a), "v-a works")
	assert.Equal(t, Vec{0.0 * a, 1.0 * a, 2.0 * a}, v.MulScalar(a), "v*a works")
	assert.Equal(t, Vec{0.0 / a, 1.0 / a, 2.0 / a}, v.DivScalar(a), "v/a works")
}

func TestAbs(t *testing.T) {
	assert.Equal(t, Vec{1.0, 2.0, 3.0}, Vec{-1.0, -2.0, -3.0}.Abs(), "abs(v) works")
}

func TestCeil(t *testing.T) {
	assert.Equal(t, Vec{math.Ceil(1.1), math.Ceil(2.2), math.Ceil(3.3)}, Vec{1.1, 2.2, 3.3}.Ceil(), "ceil(v) works")
}

func TestOps(t *testing.T) {
	a := Vec{2.0, 11.0, 5.0}
	b := Vec{7.0, 3.0, 13.0}

	assert.Equal(t, Vec{2.0, 3.0, 5.0}, a.Min(b), "min(a, b) works")
	assert.Equal(t, Vec{2.0, 3.0, 5.0}, b.Min(a), "min(b, a) works")

	assert.Equal(t, Vec{7.0, 11.0, 13.0}, a.Max(b), "max(a, b) works")
	assert.Equal(t, Vec{7.0, 11.0, 13.0}, b.Max(a), "max(b, a) works")

	assert.Equal(t, Vec{2.0 + 7.0, 11.0 + 3.0, 5.0 + 13.0}, a.Add(b), "a+b works")
	assert.Equal(t, Vec{7.0 + 2.0, 3.0 + 11.0, 13.0 + 5.0}, b.Add(a), "b+a works")

	assert.Equal(t, Vec{2.0 - 7.0, 11.0 - 3.0, 5.0 - 13.0}, a.Sub(b), "a-b works")
	assert.Equal(t, Vec{7.0 - 2.0, 3.0 - 11.0, 13.0 - 5.0}, b.Sub(a), "b-a works")

	assert.Equal(t, Vec{2.0 * 7.0, 11.0 * 3.0, 5.0 * 13.0}, a.Mul(b), "a*b works")
	assert.Equal(t, Vec{7.0 * 2.0, 3.0 * 11.0, 13.0 * 5.0}, b.Mul(a), "b*a works")

	assert.Equal(t, Vec{2.0 / 7.0, 11.0 / 3.0, 5.0 / 13.0}, a.Div(b), "a/b works")
	assert.Equal(t, Vec{7.0 / 2.0, 3.0 / 11.0, 13.0 / 5.0}, b.Div(a), "b/a works")

	assert.Equal(t, Vec{-2.0, -11.0, -5.0}, a.Neg(), "-a works")
	assert.Equal(t, Vec{-7.0, -3.0, -13.0}, b.Neg(), "-b works")
}

func TestSetOps(t *testing.T) {
	v3s := VecSet{
		{1.0, 99.0, 55.0},
		{95.0, 44.0, 3.0},
		{66.0, 7.0, 88.0},
	}

	assert.Equal(t, Vec{1.0, 7.0, 3.0}, v3s.Min(), "min(vs) works")
	assert.Equal(t, Vec{95.0, 99.0, 88.0}, v3s.Max(), "max(vs) works")
}

func TestVectorOps(t *testing.T) {
	d := math.Sqrt(2.0*2.0 + 3.0*3.0 + 5.0*5.0)
	assert.Equal(t, d, Vec{2.0, 3.0, 5.0}.Length(), "length(v) works")
	assert.Equal(t, 2.0*2.0+3.0*3.0+5.0*5.0, Vec{2.0, 3.0, 5.0}.Length2(), "length(v)^2 works")

	assert.Equal(t, 2.0, Vec{2.0, 3.0, 5.0}.MinComponent(), "min(v.x, v.y, v.z) works")
	assert.Equal(t, 5.0, Vec{2.0, 3.0, 5.0}.MaxComponent(), "max(v.x, v.y, v.z) works")

	assert.Equal(t, Vec{2.0 / d, 3.0 / d, 5.0 / d}, Vec{2.0, 3.0, 5.0}.Normalize(), "normalize(v) works")
	assert.InDelta(t, Vec{2.0, 3.0, 5.0}.Normalize().Length(), 1.0, 0.0001, "length(normalize(v)) == 1")
}
