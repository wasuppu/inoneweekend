package main

import (
	"math"
	"math/rand/v2"
)

type Material interface {
	Scatter(in Ray, rec HitRecord) (bool, RGB, Ray)
	Emitted(u, v float64, p Point3) RGB
	ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64
}

type Lambertian struct {
	tex Texture
}

func (m Lambertian) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	scatterDirection := RandomOnHemisphere(rec.normal)

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.normal
	}

	scattered := Ray{rec.p, scatterDirection, in.tm}
	attenuation := m.tex.Value(rec.u, rec.v, rec.p)
	return true, attenuation, scattered
}

func (m Lambertian) Emitted(u, v float64, p Point3) RGB {
	return RGB{}
}

func (m Lambertian) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 1 / (2 * math.Pi)
}

type Metal struct {
	albedo RGB
	fuzz   float64
}

func (m Metal) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	reflected := Reflect(in.dir, rec.normal)
	reflected = reflected.Normalize().Add(RandomUnitVector().Muln(m.fuzz))
	scattered := Ray{rec.p, reflected, in.tm}
	attenuation := m.albedo
	return scattered.dir.Dot(rec.normal) > 0, attenuation, scattered
}

func (m Metal) Emitted(u, v float64, p Point3) RGB {
	return RGB{}
}

func (m Metal) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 0
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
	if cannotRefract || Reflectance(cosTheta, ri) > rand.Float64() {
		direction = Reflect(unitDirection, rec.normal)
	} else {
		direction = Refract(unitDirection, rec.normal, ri)
	}

	scattered := Ray{rec.p, direction, in.tm}
	return true, attenuation, scattered
}

func (m Dielectric) Emitted(u, v float64, p Point3) RGB {
	return RGB{}
}

func (m Dielectric) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 0
}

func Reflectance(cosine, refractionIndex float64) float64 {
	// Use Schlick's approximation for reflectance.
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

type DiffuseLight struct {
	tex Texture
}

func (m DiffuseLight) Emitted(u, v float64, p Point3) RGB {
	return m.tex.Value(u, v, p)
}

func (m DiffuseLight) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	return false, RGB{}, Ray{}
}

func (m DiffuseLight) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 0
}

type Isotropic struct {
	tex Texture
}

func (m Isotropic) Emitted(u, v float64, p Point3) RGB {
	return RGB{}
}

func (m Isotropic) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	attenuation := m.tex.Value(rec.u, rec.v, rec.p)
	scattered := Ray{rec.p, RandomUnitVector(), in.tm}
	return true, attenuation, scattered
}

func (m Isotropic) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 0
}
