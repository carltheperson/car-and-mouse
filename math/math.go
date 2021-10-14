package math

import (
	"fmt"
	"math"
)

type Vector struct {
	A, B float64
}

func (v *Vector) ToString() string {
	return fmt.Sprintf("[%f, %f]", v.A, v.B)
}

func (v1 *Vector) Add(v2 Vector) Vector {
	return Vector{
		A: v1.A + v2.A,
		B: v1.B + v2.B,
	}
}

func (v *Vector) Multiply(n float64) Vector {
	return Vector{
		A: v.A * n,
		B: v.B * n,
	}
}

func (v *Vector) GetLength() float64 {
	return math.Sqrt(v.A*v.A + v.B*v.B)
}

func (v *Vector) GetUnitVector() Vector {
	length := v.GetLength()
	if length == 0 {
		length = 0.01
	}
	return Vector{
		A: v.A / length,
		B: v.B / length,
	}
}

func ConvertDirectionVectorToRadians(vector Vector) float64 {
	return math.Atan2(vector.A, vector.B) * -1
}

func GetRadiansBetweenTwoVectors(v1 Vector, v2 Vector) float64 {
	return math.Atan2(v1.A-v2.A, v1.B-v2.B)
}

func ConvertRadiansToDirectionVector(radians float64) Vector {
	return Vector{A: math.Cos(radians), B: math.Sin(radians)}
}

func RotatePoint(point Vector, center Vector, angle float64) Vector {
	newPoint := Vector{A: point.A - center.A, B: point.B - center.B}

	return Vector{
		A: math.Cos(angle)*newPoint.A - math.Sin(angle)*newPoint.B + center.A,
		B: math.Sin(angle)*newPoint.A + math.Cos(angle)*newPoint.B + center.B}
}

// GetDirectionDifference gets the difference between angels (in radians)
//
// Important context: Radians range from 0 to 2PI
//
// There should be a small difference between a value close to 2PI and a value close to 0
// because if drawn on a circle they would be close.
//
// Examples: 1PI - 1.5PI = 0.5PI, 1.75PI - 0.25PI = 0.5PI, 1.99PI - 0.01 = 0.02PI
func GetDirectionDifference(d1 float64, d2 float64) float64 {
	return math.Atan2(math.Sin(d1-d2), math.Cos(d1-d2))
}

func GetDistanceBetweenTwoPoints(point1 Vector, point2 Vector) float64 {
	return math.Sqrt(math.Pow(point1.A-point2.A, 2) + math.Pow(point1.B-point2.B, 2))
}
