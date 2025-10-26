package main

import (
	"math"
	"sort"
)

type Hittable interface {
	Hit(r Ray, intvl Interval) (bool, HitRecord)
	BoundingBox() AABB
}

type HitRecord struct {
	p         Point3
	normal    Vec3
	mat       Material
	t         float64
	u         float64
	v         float64
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
	bbox    AABB
}

func (hit *HittableList) Add(object Hittable) {
	hit.objects = append(hit.objects, object)
	hit.bbox = NewAABBBox(hit.bbox, object.BoundingBox())
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

func (hit HittableList) BoundingBox() AABB {
	return hit.bbox
}

type Sphere struct {
	center Ray
	radius float64
	mat    Material
	bbox   AABB
}

func NewSphere(center Point3, radius float64, mat Material) Sphere {
	rvec := Vec3{radius, radius, radius}
	bbox := NewAABBPoint(center.Add(rvec), center.Sub(rvec))
	return Sphere{Ray{center, Vec3{0, 0, 0}, 0}, radius, mat, bbox}
}

func NewMotionSphere(center1 Point3, center2 Point3, radius float64, mat Material) Sphere {
	rvec := Vec3{radius, radius, radius}
	center := Ray{center1, center2.Sub(center1), 0}
	bbox1 := NewAABBPoint(center.At(0).Sub(rvec), center.At(0).Add(rvec))
	bbox2 := NewAABBPoint(center.At(1).Sub(rvec), center.At(1).Add(rvec))
	bbox := NewAABBBox(bbox1, bbox2)

	return Sphere{center, radius, mat, bbox}
}

func (hit Sphere) Hit(r Ray, intvl Interval) (bool, HitRecord) {
	currentCenter := hit.center.At(r.tm)
	oc := currentCenter.Sub(r.orig)
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
	outwardNormal := rec.p.Sub(currentCenter).Divn(hit.radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.mat = hit.mat

	return true, rec
}

func (hit Sphere) BoundingBox() AABB {
	return hit.bbox
}

type BVHNode struct {
	left  Hittable
	right Hittable
	bbox  AABB
}

func NewBVHNode(list HittableList) BVHNode {
	return BVHNodeConstructor(list.objects, 0, len(list.objects))
}

func BVHNodeConstructor(objects []Hittable, start, end int) BVHNode {
	// Build the bounding box of the span of source objects.
	bbox := EmptyAABB
	for i := start; i < end; i++ {
		bbox = NewAABBBox(bbox, objects[i].BoundingBox())
	}
	axis := bbox.LongestAxis()

	comparator := BoxZCompare
	if axis == 0 {
		comparator = BoxXCompare
	}
	if axis == 1 {
		comparator = BoxYCompare
	}

	var left, right Hittable
	objectSpan := end - start
	if objectSpan == 1 {
		left = objects[start]
		right = objects[start]
	} else if objectSpan == 2 {
		left = objects[start]
		right = objects[start+1]
	} else {
		sort.Slice(objects[start:end], func(i, j int) bool {
			return comparator(objects[start+i], objects[start+j])
		})
		mid := start + objectSpan/2
		left = BVHNodeConstructor(objects, start, mid)
		right = BVHNodeConstructor(objects, mid, end)
	}
	return BVHNode{left, right, bbox}
}

func (hit BVHNode) Hit(r Ray, intvl Interval) (bool, HitRecord) {
	if !hit.bbox.Hit(r, intvl) {
		return false, HitRecord{}
	}

	if hitLeft, rec1 := hit.left.Hit(r, intvl); hitLeft {
		hitRight, rec2 := hit.right.Hit(r, Interval{intvl.min, rec1.t})
		if hitRight {
			return true, rec2
		}
		return true, rec1
	} else {
		hitRight, rec2 := hit.right.Hit(r, Interval{intvl.min, intvl.max})
		return hitRight, rec2
	}
}

func (hit BVHNode) BoundingBox() AABB {
	return hit.bbox
}

func BoxCompare(a, b Hittable, axisIndex int) bool {
	aAxisInterval := a.BoundingBox()[axisIndex]
	bAxisInterval := b.BoundingBox()[axisIndex]
	return aAxisInterval.min < bAxisInterval.min
}

func BoxXCompare(a, b Hittable) bool {
	return BoxCompare(a, b, 0)
}

func BoxYCompare(a, b Hittable) bool {
	return BoxCompare(a, b, 1)
}

func BoxZCompare(a, b Hittable) bool {
	return BoxCompare(a, b, 2)
}
