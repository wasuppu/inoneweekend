package main

import "image/color"

type Ray struct {
	orig Point3
	dir  Vec3
}

func (r Ray) At(t float64) Vec3 {
	return r.orig.Add(r.dir.Muln(t))
}

func (r Ray) Color() color.Color {
	unitDirction := r.dir.Normalize()
	a := 0.5 * (unitDirction.Y() + 1.0)
	return Vec3{1, 1, 1}.Muln(1 - a).Add(Vec3{0.5, 0.7, 1.0}.Muln(a)).Color()
}
