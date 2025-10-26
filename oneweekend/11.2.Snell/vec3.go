package main

import (
	"image/color"
	"math"
	"math/rand/v2"
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

func (v Vec3) Mul(o Vec3) Vec3 {
	return Vec3{v[0] * o[0], v[1] * o[1], v[2] * o[2]}
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

func (v Vec3) NearZero() bool {
	// Return true if the vector is close to zero in all dimensions.
	s := 1e-8
	return math.Abs(v[0]) < s && math.Abs(v[1]) < s && math.Abs(v[2]) < s
}

func Clamp[T int | float64](val T, min T, max T) T {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

func (v Vec3) Clamp() Vec3 {
	return Vec3{Clamp(v[0], 0, 1), Clamp(v[1], 0, 1), Clamp(v[2], 0, 1)}
}

func RandomRange(min, max float64) float64 {
	// Returns a random real in [min,max).
	return min + (max-min)*rand.Float64()
}

func RandomVec3() Vec3 {
	return Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
}

func RandomVec3Range(min, max float64) Vec3 {
	return Vec3{RandomRange(min, max), RandomRange(min, max), RandomRange(min, max)}
}

func RandomUnitVector() Vec3 {
	for {
		p := RandomVec3Range(-1, 1)
		lensq := p.Dot(p)
		if lensq <= 1 && lensq > 1e-160 {
			return p.Divn(math.Sqrt(lensq))
		}
	}
}

func RandomOnHemisphere(normal Vec3) Vec3 {
	onUnitSphere := RandomUnitVector()
	// In the same hemisphere as the normal
	if onUnitSphere.Dot(normal) > 0.0 {
		return onUnitSphere
	} else {
		return onUnitSphere.Muln(-1)
	}
}

func (v RGB) Color() color.RGBA {
	intensity := Interval{0.000, 0.999}
	for i := range v {
		v[i] = 256 * intensity.Clamp(LinearToGamma(v[i]))
	}
	return color.RGBA{uint8(v[0]), uint8(v[1]), uint8(v[2]), 0xff}
}

func LinearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}
	return 0
}

func Reflect(v Vec3, n Vec3) Vec3 {
	return v.Sub(n.Muln(2 * v.Dot(n)))
}

func Refract(uv Vec3, n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := math.Min(uv.Muln(-1).Dot(n), 1)
	rOutPerp := uv.Add(n.Muln(cosTheta)).Muln(etaiOverEtat)
	rOutParallel := n.Muln(-math.Sqrt(math.Abs(1 - rOutPerp.Dot(rOutPerp))))
	return rOutPerp.Add(rOutParallel)
}
