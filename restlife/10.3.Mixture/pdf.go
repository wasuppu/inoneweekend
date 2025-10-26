package main

import (
	"math"
	"math/rand/v2"
)

type PDF interface {
	Value(direction Vec3) float64
	Generate() Vec3
}

type SpherePDF struct{}

func (pdf SpherePDF) Value(direction Vec3) float64 {
	return 1 / (4 * math.Pi)
}

func (pdf SpherePDF) Generate() Vec3 {
	return RandomUnitVector()
}

type CosinePDF struct {
	uvw ONB
}

func NewCosinePDF(w Vec3) CosinePDF {
	return CosinePDF{NewONB(w)}
}

func (pdf CosinePDF) Value(direction Vec3) float64 {
	cosineTheta := direction.Normalize().Dot(pdf.uvw.W())
	return math.Max(0, cosineTheta/math.Pi)
}

func (pdf CosinePDF) Generate() Vec3 {
	return pdf.uvw.Transform(RandomCosineDirection())
}

type HittablePDF struct {
	objects Hittable
	origin  Point3
}

func (pdf HittablePDF) Value(direction Vec3) float64 {
	return pdf.objects.PDFValue(pdf.origin, direction)
}

func (pdf HittablePDF) Generate() Vec3 {
	return pdf.objects.Random(pdf.origin)
}

type MixturePDF [2]PDF

func (pdf MixturePDF) Value(direction Vec3) float64 {
	return 0.5*pdf[0].Value(direction) + 0.5*pdf[1].Value(direction)
}

func (pdf MixturePDF) Generate() Vec3 {
	if rand.Float64() < 0.5 {
		return pdf[0].Generate()
	} else {
		return pdf[1].Generate()
	}
}
