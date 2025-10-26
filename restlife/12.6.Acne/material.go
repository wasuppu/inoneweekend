package main

import (
	"math"
	"math/rand/v2"
)

type Material interface {
	Scatter(in Ray, rec HitRecord) (bool, ScatterRecord)
	Emitted(in Ray, rec HitRecord, u, v float64, p Point3) RGB
	ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64
}

type ScatterRecord struct {
	attenuation RGB
	pdf         PDF
	skipPdf     bool
	skipPdfRay  Ray
}

type EmptyMaterial struct{}

func (m EmptyMaterial) Scatter(in Ray, rec HitRecord) (bool, ScatterRecord) {
	return false, ScatterRecord{skipPdf: true}
}

func (m EmptyMaterial) Emitted(in Ray, rec HitRecord, u, v float64, p Point3) RGB {
	return RGB{}
}

func (m EmptyMaterial) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 0
}

type Lambertian struct {
	tex Texture
}

func (m Lambertian) Scatter(in Ray, rec HitRecord) (bool, ScatterRecord) {
	attenuation := m.tex.Value(rec.u, rec.v, rec.p)
	pdf := NewCosinePDF(rec.normal)
	skipPdf := false
	return true, ScatterRecord{attenuation: attenuation, pdf: pdf, skipPdf: skipPdf}
}

func (m Lambertian) Emitted(in Ray, rec HitRecord, u, v float64, p Point3) RGB {
	return RGB{}
}

func (m Lambertian) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	cosTheta := rec.normal.Dot(scattered.dir.Normalize())
	if cosTheta < 0 {
		return 0
	} else {
		return cosTheta / math.Pi
	}
}

type Metal struct {
	albedo RGB
	fuzz   float64
}

func (m Metal) Scatter(in Ray, rec HitRecord) (bool, ScatterRecord) {
	reflected := Reflect(in.dir, rec.normal)
	reflected = reflected.Normalize().Add(RandomUnitVector().Muln(m.fuzz))

	attenuation := m.albedo
	scattered := Ray{rec.p, reflected, in.tm}
	return true, ScatterRecord{attenuation: attenuation, skipPdfRay: scattered, skipPdf: true}
}

func (m Metal) Emitted(in Ray, rec HitRecord, u, v float64, p Point3) RGB {
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

func (m Dielectric) Scatter(in Ray, rec HitRecord) (bool, ScatterRecord) {
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
	return true, ScatterRecord{attenuation: attenuation, skipPdfRay: scattered, skipPdf: true}
}

func (m Dielectric) Emitted(in Ray, rec HitRecord, u, v float64, p Point3) RGB {
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

func (m DiffuseLight) Emitted(in Ray, rec HitRecord, u, v float64, p Point3) RGB {
	if !rec.frontFace {
		return RGB{}
	}
	return m.tex.Value(u, v, p)
}

func (m DiffuseLight) Scatter(in Ray, rec HitRecord) (bool, ScatterRecord) {
	return false, ScatterRecord{}
}

func (m DiffuseLight) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 0
}

type Isotropic struct {
	tex Texture
}

func (m Isotropic) Emitted(in Ray, rec HitRecord, u, v float64, p Point3) RGB {
	return RGB{}
}

func (m Isotropic) Scatter(in Ray, rec HitRecord) (bool, ScatterRecord) {
	attenuation := m.tex.Value(rec.u, rec.v, rec.p)
	pdf := SpherePDF{}
	return true, ScatterRecord{attenuation: attenuation, pdf: pdf, skipPdf: false}
}

func (m Isotropic) ScatteringPdf(in Ray, rec HitRecord, scattered Ray) float64 {
	return 1 / (4 * math.Pi)
}
