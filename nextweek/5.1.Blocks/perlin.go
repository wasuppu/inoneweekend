package main

import "math/rand/v2"

const PointCount = 256

type Perlin struct {
	randfloat [PointCount]float64
	permX     [PointCount]int
	permY     [PointCount]int
	permZ     [PointCount]int
}

func NewPerlin() Perlin {
	perlin := Perlin{}
	for i := range PointCount {
		perlin.randfloat[i] = rand.Float64()
	}

	PerlinGeneratePerm(perlin.permX[:])
	PerlinGeneratePerm(perlin.permY[:])
	PerlinGeneratePerm(perlin.permZ[:])

	return perlin
}

func (pl Perlin) Noise(p Point3) float64 {
	i := int(4*p.X()) & 255
	j := int(4*p.Y()) & 255
	k := int(4*p.Z()) & 255

	return pl.randfloat[pl.permX[i]^pl.permY[j]^pl.permZ[k]]
}

func PerlinGeneratePerm(p []int) {
	for i := range PointCount {
		p[i] = i
	}
	Permute(p, PointCount)
}

func Permute(p []int, n int) {
	for i := range n - 1 {
		target := int(RandomRange(0, float64(i)))
		tmp := p[i]
		p[i] = p[target]
		p[target] = tmp
	}
}
