package main

import (
	"math"
)

type Texture interface {
	Value(u, v float64, p Point3) RGB
}

type SolidColor struct {
	albedo RGB
}

func NewSolidColor(red, green, blue float64) SolidColor {
	return SolidColor{RGB{red, green, blue}}
}

func (t SolidColor) Value(u, v float64, p Point3) RGB {
	return t.albedo
}

type CheckerTexture struct {
	invScale float64
	even     Texture
	odd      Texture
}

func NewCheckerTexture(scale float64, even, odd Texture) CheckerTexture {
	return CheckerTexture{1 / scale, even, odd}
}

func (t CheckerTexture) Value(u, v float64, p Point3) RGB {
	x := int(math.Floor(t.invScale * p.X()))
	y := int(math.Floor(t.invScale * p.Y()))
	z := int(math.Floor(t.invScale * p.Z()))

	isEven := (x+y+z)%2 == 0

	if isEven {
		return t.even.Value(u, v, p)
	} else {
		return t.odd.Value(u, v, p)
	}
}

type ImageTexture struct {
	rtwImage RTWImage
}

func NewImageTexture(filename string) ImageTexture {
	return ImageTexture{NewRTWImage(filename)}
}

func (t ImageTexture) Value(u, v float64, p Point3) RGB {
	// If we have no texture data, then return solid cyan as a debugging aid.
	if t.rtwImage.height <= 0 {
		return RGB{0, 1, 1}
	}

	// Clamp input texture coordinates to [0,1] x [1,0]
	u = Interval{0, 1}.Clamp(u)
	v = 1 - Interval{0, 1}.Clamp(v) // Flip V to image coordinates

	i := int(u * float64(t.rtwImage.width))
	j := int(v * float64(t.rtwImage.height))

	return t.rtwImage.Get(i, j)
}

type NoiseTexture struct {
	noise Perlin
	scale float64
}

func NewNoiseTexture(scale float64) NoiseTexture {
	return NoiseTexture{NewPerlin(), scale}
}

func (t NoiseTexture) Value(u, v float64, p Point3) RGB {
	return RGB{1, 1, 1}.Muln(0.5 * (1.0 + t.noise.Noise(p.Muln(t.scale))))
}
