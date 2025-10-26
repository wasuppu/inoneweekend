package main

import (
	"math"
)

const PointCount = 256

type Perlin struct {
	randvec [PointCount]Vec3
	permX   [PointCount]int
	permY   [PointCount]int
	permZ   [PointCount]int
}

func NewPerlin() Perlin {
	perlin := Perlin{}
	for i := range PointCount {
		perlin.randvec[i] = RandomVec3Range(-1, 1).Normalize()
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

	i := int(math.Floor(p.X()))
	j := int(math.Floor(p.Y()))
	k := int(math.Floor(p.Z()))
	c := [2][2][2]Vec3{}

	for di := range 2 {
		for dj := range 2 {
			for dk := range 2 {
				c[di][dj][dk] = pl.randvec[pl.permX[(i+di)&255]^
					pl.permY[(j+dj)&255]^
					pl.permZ[(k+dk)&255]]
			}
		}
	}

	return PerlinInterp(c, u, v, w)
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

func PerlinInterp(c [2][2][2]Vec3, u, v, w float64) float64 {
	uu := u * u * (3 - 2*u)
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)
	accum := 0.0

	for i := range 2 {
		for j := range 2 {
			for k := range 2 {
				weightV := Vec3{u - float64(i), v - float64(j), w - float64(k)}
				accum += (float64(i)*uu + (1-float64(i))*(1-uu)) *
					(float64(j)*vv + (1-float64(j))*(1-vv)) *
					(float64(k)*ww + (1-float64(k))*(1-ww)) *
					c[i][j][k].Dot(weightV)
			}
		}
	}
	return accum
}
