package main

import "math"

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
