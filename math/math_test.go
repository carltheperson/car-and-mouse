package math_test

import (
	"testing"

	"github.com/carltheperson/car-and-mouse/math"
)

func isVectorsSame(v1 math.Vector, v2 math.Vector) bool {
	return v1.A == v2.A && v1.B == v2.B
}

func TestAdd(t *testing.T) {
	v1 := math.Vector{A: 3, B: 6}
	v2 := math.Vector{A: -2, B: 4}
	got := v1.Add(v2)

	want := math.Vector{A: 1, B: 10}

	if !isVectorsSame(got, want) {
		t.Errorf("Got %s wanted %s", got.ToString(), want.ToString())
	}
}

func TestMultiply(t *testing.T) {
	v := math.Vector{A: 3, B: 4}
	got := v.Multiply(2)

	want := math.Vector{A: 6, B: 8}

	if !isVectorsSame(got, want) {
		t.Errorf("Got %s wanted %s", got.ToString(), want.ToString())
	}
}

func TestGetLength(t *testing.T) {
	v := math.Vector{A: 3, B: 4}
	got := v.GetLength()

	want := 5.0

	if got != want {
		t.Errorf("Got %f wanted %f", got, want)
	}
}
