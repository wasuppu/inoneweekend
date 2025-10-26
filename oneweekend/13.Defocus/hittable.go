package main

import "math"

type Hittable interface {
	Hit(r Ray, intvl Interval) (bool, HitRecord)
}

type HitRecord struct {
	p         Point3
	normal    Vec3
	mat       Material
	t         float64
	frontFace bool
}

func (hit *HitRecord) SetFaceNormal(r Ray, outwardNormal Vec3) {
	// Sets the hit record normal vector.
	// NOTE: the parameter `outward_normal` is assumed to have unit length.
	hit.frontFace = r.dir.Dot(outwardNormal) < 0
	if hit.frontFace {
		hit.normal = outwardNormal
	} else {
		hit.normal = outwardNormal.Muln(-1)
	}
}

type HittableList struct {
	objects []Hittable
}

func (hit *HittableList) Add(object Hittable) {
	hit.objects = append(hit.objects, object)
}

func (hit HittableList) Hit(r Ray, intvl Interval) (bool, HitRecord) {
	hitAnything := false
	closestSoFar := intvl.max
	rec := HitRecord{}

	for _, object := range hit.objects {
		if ok, tempRec := object.Hit(r, Interval{intvl.min, closestSoFar}); ok {
			hitAnything = true
			closestSoFar = tempRec.t
			rec = tempRec
		}
	}
	return hitAnything, rec
}

type Sphere struct {
	center Point3
	radius float64
	mat    Material
}

func (hit Sphere) Hit(r Ray, intvl Interval) (bool, HitRecord) {
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
	if !intvl.Surrounds(root) {
		root = (h + sqrtd) / a
		if !intvl.Surrounds(root) {
			return false, HitRecord{}
		}
	}

	rec := HitRecord{}
	rec.t = root
	rec.p = r.At(rec.t)
	outwardNormal := rec.p.Sub(hit.center).Divn(hit.radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.mat = hit.mat

	return true, rec
}
