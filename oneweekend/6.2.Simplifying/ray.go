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

func (r Ray) Color() color.Color {
	t := r.hitSphere(Point3{0, 0, -1}, 0.5)
	if t >= 0 {
		n := r.At(t).Sub(Vec3{0, 0, -1}).Normalize()
		return Vec3{n.X() + 1, n.Y() + 1, n.Z() + 1}.Muln(0.5).Color()
	}
	unitDirction := r.dir.Normalize()
	a := 0.5 * (unitDirction.Y() + 1.0)
	return Vec3{1, 1, 1}.Muln(1 - a).Add(Vec3{0.5, 0.7, 1.0}.Muln(a)).Color()
}

func (r Ray) hitSphere(center Point3, radius float64) float64 {
	oc := center.Sub(r.orig)
	a := r.dir.Dot(r.dir)
	h := r.dir.Dot(oc)
	c := oc.Dot(oc) - radius*radius
	discriminant := h*h - a*c
	if discriminant < 0 {
		return -1
	} else {
		return (h - math.Sqrt(discriminant)) / a
	}
}
