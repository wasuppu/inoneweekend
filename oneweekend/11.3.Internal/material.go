package main

import "math"

type Material interface {
	Scatter(in Ray, rec HitRecord) (bool, RGB, Ray)
}

type Lambertian struct {
	albedo RGB
}

func (m Lambertian) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	scatterDirection := rec.normal.Add(RandomUnitVector())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.normal
	}

	scattered := Ray{rec.p, scatterDirection}
	attenuation := m.albedo
	return true, attenuation, scattered
}

type Metal struct {
	albedo RGB
	fuzz   float64
}

func (m Metal) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	reflected := Reflect(in.dir, rec.normal)
	reflected = reflected.Normalize().Add(RandomUnitVector().Muln(m.fuzz))
	scattered := Ray{rec.p, reflected}
	attenuation := m.albedo
	return scattered.dir.Dot(rec.normal) > 0, attenuation, scattered
}

type Dielectric struct {
	// Refractive index in vacuum or air, or the ratio of the material's refractive index over
	// the refractive index of the enclosing media
	refractionIndex float64
}

func (m Dielectric) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	attenuation := RGB{1, 1, 1}
	ri := m.refractionIndex
	if rec.frontFace {
		ri = 1 / m.refractionIndex
	}

	unitDirection := in.dir.Normalize()
	cosTheta := math.Min(unitDirection.Muln(-1).Dot(rec.normal), 1)
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	cannotRefract := ri*sinTheta > 1.0
	var direction Vec3
	if cannotRefract {
		direction = Reflect(unitDirection, rec.normal)
	} else {
		direction = Refract(unitDirection, rec.normal, ri)
	}

	scattered := Ray{rec.p, direction}
	return true, attenuation, scattered
}
