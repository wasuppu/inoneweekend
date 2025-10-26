package main

import (
	"image/color"
	"math"
)

type Vec3 [3]float64
type Point3 = Vec3
type RGB = Vec3

func (v Vec3) X() float64 {
	return v[0]
}

func (v Vec3) Y() float64 {
	return v[1]
}

func (v Vec3) Z() float64 {
	return v[2]
}

func (v Vec3) R() float64 {
	return v[0]
}

func (v Vec3) G() float64 {
	return v[1]
}

func (v Vec3) B() float64 {
	return v[2]
}

func (v Vec3) Add(o Vec3) Vec3 {
	return Vec3{v[0] + o[0], v[1] + o[1], v[2] + o[2]}
}

func (v Vec3) Sub(o Vec3) Vec3 {
	return Vec3{v[0] - o[0], v[1] - o[1], v[2] - o[2]}
}

func (v Vec3) Muln(t float64) Vec3 {
	return Vec3{v[0] * t, v[1] * t, v[2] * t}
}

func (v Vec3) Divn(t float64) Vec3 {
	return Vec3{v[0] / t, v[1] / t, v[2] / t}
}

func (v Vec3) Dot(o Vec3) float64 {
	return v[0]*o[0] + v[1]*o[1] + v[2]*o[2]
}

func (v Vec3) Cross(o Vec3) Vec3 {
	x := v[1]*o[2] - v[2]*o[1]
	y := v[2]*o[0] - v[0]*o[2]
	z := v[0]*o[1] - v[1]*o[0]

	return Vec3{x, y, z}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(float64(v.Dot(v)))
}

func (v Vec3) Normalize() Vec3 {
	return v.Muln(1 / v.Length())
}

func (v Vec3) Clamp() Vec3 {
	return Vec3{Clamp(v[0], 0, 1), Clamp(v[1], 0, 1), Clamp(v[2], 0, 1)}
}

func Clamp[T int | float64](val T, min T, max T) T {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

func (v RGB) Color() color.RGBA {
	return color.RGBA{uint8(255.999 * v[0]), uint8(255.999 * v[1]), uint8(255.999 * v[2]), 0xff}
}
