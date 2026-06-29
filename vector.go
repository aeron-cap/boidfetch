package main

import "math"

type Vector2D struct {
	X float64
	Y float64
}

func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vector2D) Subtract(other Vector2D) Vector2D {
	return Vector2D{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vector2D) Multiply(scalar float64) Vector2D {
	return Vector2D{X: v.X * scalar, Y: v.Y * scalar}
}

func (v Vector2D) Divide(scalar float64) Vector2D {
	if scalar == 0 {
		return v
	}
	return Vector2D{X: v.X / scalar, Y: v.Y / scalar}
}

func (v Vector2D) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v Vector2D) Normalize() Vector2D {
	mag := v.Magnitude()
	if mag > 0 {
		return v.Divide(mag)
	}
	return v
}
