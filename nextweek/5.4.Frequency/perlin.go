package main

import (
	"math"
	"math/rand/v2"
)

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
	u := p.X() - math.Floor(p.X())
	v := p.Y() - math.Floor(p.Y())
	w := p.Z() - math.Floor(p.Z())
	u = u * u * (3 - 2*u)
	v = v * v * (3 - 2*v)
	w = w * w * (3 - 2*w)

	i := int(math.Floor(p.X()))
	j := int(math.Floor(p.Y()))
	k := int(math.Floor(p.Z()))
	c := [2][2][2]float64{}

	for di := range 2 {
		for dj := range 2 {
			for dk := range 2 {
				c[di][dj][dk] = pl.randfloat[pl.permX[(i+di)&255]^
					pl.permY[(j+dj)&255]^
					pl.permZ[(k+dk)&255]]
			}
		}
	}

	return TrilinearInterp(c, u, v, w)
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

func TrilinearInterp(c [2][2][2]float64, u, v, w float64) float64 {
	accum := 0.0
	for i := range 2 {
		for j := range 2 {
			for k := range 2 {
				accum += (float64(i)*u + (1-float64(i))*(1-u)) *
					(float64(j)*v + (1-float64(j))*(1-v)) *
					(float64(k)*w + (1-float64(k))*(1-w)) *
					c[i][j][k]
			}
		}
	}
	return accum
}
