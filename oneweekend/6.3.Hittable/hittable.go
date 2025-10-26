package main

import "math"

type Hittable interface {
	Hit(r Ray, tmin, tmax float64) (bool, HitRecord)
}

type HitRecord struct {
	p      Point3
	normal Vec3
	t      float64
}

type Sphere struct {
	center Point3
	radius float64
}

func (hit Sphere) Hit(r Ray, tmin, tmax float64) (bool, HitRecord) {
	oc := hit.center.Sub(r.orig)
	a := r.dir.Dot(r.dir)
	h := r.dir.Dot(oc)
	c := oc.Dot(oc) - hit.radius*hit.radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false, HitRecord{}
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range.
	root := (h - sqrtd) / a
	if tmin >= root || tmax <= root {
		root = (h + sqrtd) / a
		if tmin >= root || tmax <= root {
			return false, HitRecord{}
		}
	}

	rec := HitRecord{}
	rec.t = root
	rec.p = r.At(rec.t)
	rec.normal = rec.p.Sub(hit.center).Divn(hit.radius)

	return true, rec
}
