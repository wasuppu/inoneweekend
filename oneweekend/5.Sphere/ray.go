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
	if r.hitSphere(Point3{0, 0, -1}, 0.5) {
		return Vec3{1, 0, 0}.Color()
	}
	unitDirction := r.dir.Normalize()
	a := 0.5 * (unitDirction.Y() + 1.0)
	return Vec3{1, 1, 1}.Muln(1 - a).Add(Vec3{0.5, 0.7, 1.0}.Muln(a)).Color()
}

func (r Ray) hitSphere(center Point3, radius float64) bool {
	oc := center.Sub(r.orig)
	a := r.dir.Dot(r.dir)
	b := -2.0 * r.dir.Dot(oc)
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant >= 0
}
