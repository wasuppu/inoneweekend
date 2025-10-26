package main

import (
	"image/color"
	"math"
)

type Ray struct {
	orig Point3
	dir  Vec3
}

func (r Ray) At(t float64) Vec3 {
	return r.orig.Add(r.dir.Muln(t))
}

func (r Ray) Color(world Hittable) color.Color {
	if hitAnything, rec := world.Hit(r, 0, math.MaxFloat64); hitAnything {
		return rec.normal.Add(RGB{1, 1, 1}).Muln(0.5).Color()
	}

	unitDirection := r.dir.Normalize()
	a := 0.5 * (unitDirection.Y() + 1.0)
	return RGB{1, 1, 1}.Muln(1.0 - a).Add(RGB{0.5, 0.7, 1.0}.Muln(a)).Color()
}
