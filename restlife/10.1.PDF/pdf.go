package main

import "math"

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
