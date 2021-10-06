package math

import (
	"fmt"
	"math"
)

type Vector2D struct {
	A, B float64
}

func (v *Vector2D) ToString() string {
	return fmt.Sprintf("[%f, %f]", v.A, v.B)
}

func (v1 *Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{
		A: v1.A + v2.A,
		B: v1.B + v2.B,
	}
}

func (v *Vector2D) Multiply(n float64) Vector2D {
	return Vector2D{
		A: v.A * n,
		B: v.B * n,
	}
}

func (v *Vector2D) GetLength() float64 {
	return math.Sqrt(v.A*v.A + v.B*v.B)
}

func (v *Vector2D) GetUnitVector() Vector2D {
	length := v.GetLength()
	if length == 0 {
		length = 0.01
	}
	return Vector2D{
		A: v.A / length,
		B: v.B / length,
	}
}

func ConvertDirectionVectorToRadians(vector Vector2D) float64 {
	return math.Atan2(vector.A, vector.B) * -1
}

func GetRadiansBetweenTwoVectors(v1 Vector2D, v2 Vector2D) float64 {
	return math.Atan2(v1.A-v2.A, v1.B-v2.B)
}

func ConvertRadiansToDirectionVector(radians float64) Vector2D {
	return Vector2D{A: math.Cos(radians), B: math.Sin(radians)}
}

func RotatePoint(point Vector2D, center Vector2D, angle float64) Vector2D {
	newPoint := Vector2D{A: point.A - center.A, B: point.B - center.B}

	return Vector2D{
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
