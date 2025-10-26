package main

import (
	"math"
	"math/rand/v2"
	"sort"
)

type Hittable interface {
	Hit(r Ray, intvl Interval) (bool, HitRecord)
	BoundingBox() AABB
	PDFValue(origin Point3, direction Vec3) float64
	Random(origin Point3) Vec3
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

func (hit HittableList) PDFValue(origin Point3, direction Vec3) float64 {
	weight := 1.0 / float64(len(hit.objects))
	sum := 0.0

	for _, object := range hit.objects {
		sum += float64(weight) * object.PDFValue(origin, direction)
	}

	return sum
}

func (hit HittableList) Random(origin Point3) Vec3 {
	intSize := len(hit.objects)
	return hit.objects[int(RandomRange(0, float64(intSize-1)))].Random(origin)
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
	rec.u, rec.v = GetSphereUV(outwardNormal)
	rec.mat = hit.mat

	return true, rec
}

func (hit Sphere) BoundingBox() AABB {
	return hit.bbox
}

func (hit Sphere) PDFValue(origin Point3, direction Vec3) float64 {
	// This method only works for stationary spheres.

	hitAnything, _ := hit.Hit(Ray{origin, direction, 0}, Interval{0.001, math.MaxFloat64})
	if !hitAnything {
		return 0
	}

	dist := hit.center.At(0).Sub(origin)
	distSquared := dist.Dot(dist)
	cosThetaMax := math.Sqrt(1 - hit.radius*hit.radius/distSquared)
	solidAngle := 2 * math.Pi * (1 - cosThetaMax)

	return 1 / solidAngle
}

func (hit Sphere) Random(origin Point3) Vec3 {
	direction := hit.center.At(0).Sub(origin)
	distanceSquared := direction.Dot(direction)
	uvw := NewONB(direction)
	return uvw.Transform(RandomToSphere(hit.radius, distanceSquared))
}

func GetSphereUV(p Point3) (float64, float64) {
	// p: a given point on the sphere of radius one, centered at the origin.
	// u: returned value [0,1] of angle around the Y axis from X=-1.
	// v: returned value [0,1] of angle from Y=-1 to Y=+1.
	//     <1 0 0> yields <0.50 0.50>       <-1  0  0> yields <0.00 0.50>
	//     <0 1 0> yields <0.50 1.00>       < 0 -1  0> yields <0.50 0.00>
	//     <0 0 1> yields <0.25 0.50>       < 0  0 -1> yields <0.75 0.50>

	theta := math.Acos(-p.Y())
	phi := math.Atan2(-p.Z(), p.X()) + math.Pi
	u := phi / (2 * math.Pi)
	v := theta / math.Pi

	return u, v
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

func (hit BVHNode) PDFValue(origin Point3, direction Vec3) float64 {
	return 0
}

func (hit BVHNode) Random(origin Point3) Vec3 {
	return Vec3{1, 0, 0}
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

type Quad struct {
	q      Point3
	u, v   Vec3
	w      Vec3
	mat    Material
	bbox   AABB
	normal Vec3
	d      float64
	area   float64
}

func NewQuad(q Point3, u, v Vec3, mat Material) Quad {
	n := u.Cross(v)
	normal := n.Normalize()
	d := normal.Dot(q)
	w := n.Divn(n.Dot(n))
	area := n.Length()
	quad := Quad{q: q, u: u, v: v, w: w, mat: mat, normal: normal, d: d, area: area}
	quad.SetBoundingBox()
	return quad
}

func (hit *Quad) SetBoundingBox() {
	// Compute the bounding box of all four vertices.
	bboxDiagonal1 := NewAABBPoint(hit.q, hit.q.Add(hit.u).Add(hit.v))
	bboxDiagonal2 := NewAABBPoint(hit.q.Add(hit.u), hit.q.Add(hit.v))
	hit.bbox = NewAABBBox(bboxDiagonal1, bboxDiagonal2)
}

func (hit Quad) Hit(r Ray, intvl Interval) (bool, HitRecord) {
	denom := hit.normal.Dot(r.dir)

	// No hit if the ray is parallel to the plane.
	if math.Abs(denom) < 1e-8 {
		return false, HitRecord{}
	}

	// Return false if the hit point parameter t is outside the ray interval.
	t := (hit.d - hit.normal.Dot(r.orig)) / denom
	if !intvl.Contains(t) {
		return false, HitRecord{}
	}

	// Determine if the hit point lies within the planar shape using its plane coordinates.
	intersection := r.At(t)
	planarHitptVector := intersection.Sub(hit.q)
	alpha := hit.w.Dot(planarHitptVector.Cross(hit.v))
	beta := hit.w.Dot(hit.u.Cross(planarHitptVector))

	ok, rec := IsInterior(alpha, beta)
	if !ok {
		return false, rec
	}

	// Ray hits the 2D shape; set the rest of the hit record and return true.
	rec.t = t
	rec.p = intersection
	rec.mat = hit.mat
	rec.SetFaceNormal(r, hit.normal)

	return true, rec
}

func (hit Quad) BoundingBox() AABB {
	return hit.bbox
}

func (hit Quad) PDFValue(origin Point3, direction Vec3) float64 {
	hitAnything, rec := hit.Hit(Ray{origin, direction, 0}, Interval{0.001, math.MaxFloat64})
	if !hitAnything {
		return 0
	}

	distanceSquared := rec.t * rec.t * direction.Dot(direction)
	cosine := math.Abs(direction.Dot(rec.normal) / direction.Length())

	return distanceSquared / (cosine * hit.area)
}

func (hit Quad) Random(origin Point3) Vec3 {
	p := hit.q.Add(hit.u.Muln(rand.Float64())).Add(hit.v.Muln(rand.Float64()))
	return p.Sub(origin)
}

func IsInterior(a, b float64) (bool, HitRecord) {
	unitInterval := Interval{0, 1}
	// Given the hit point in plane coordinates, return false if it is outside the
	// primitive, otherwise set the hit record UV coordinates and return true.

	if !unitInterval.Contains(a) || !unitInterval.Contains(b) {
		return false, HitRecord{}
	}

	rec := HitRecord{}
	rec.u = a
	rec.v = b
	return true, rec
}

func Box(a, b Point3, mat Material) HittableList {
	// Returns the 3D box (six sides) that contains the two opposite vertices a & b.
	sides := HittableList{}

	// Construct the two opposite vertices with the minimum and maximum coordinates.
	min := Point3{math.Min(a.X(), b.X()), math.Min(a.Y(), b.Y()), math.Min(a.Z(), b.Z())}
	max := Point3{math.Max(a.X(), b.X()), math.Max(a.Y(), b.Y()), math.Max(a.Z(), b.Z())}

	dx := Vec3{max.X() - min.X(), 0, 0}
	dy := Vec3{0, max.Y() - min.Y(), 0}
	dz := Vec3{0, 0, max.Z() - min.Z()}

	sides.Add(NewQuad(Point3{min.X(), min.Y(), max.Z()}, dx, dy, mat))          // front
	sides.Add(NewQuad(Point3{max.X(), min.Y(), max.Z()}, dz.Muln(-1), dy, mat)) // right
	sides.Add(NewQuad(Point3{max.X(), min.Y(), min.Z()}, dx.Muln(-1), dy, mat)) // back
	sides.Add(NewQuad(Point3{min.X(), min.Y(), min.Z()}, dz, dy, mat))          // left
	sides.Add(NewQuad(Point3{min.X(), max.Y(), max.Z()}, dx, dz.Muln(-1), mat)) // top
	sides.Add(NewQuad(Point3{min.X(), min.Y(), min.Z()}, dx, dz, mat))          // bottom

	return sides
}

type Translate struct {
	object Hittable
	offset Vec3
	bbox   AABB
}

func NewTranslate(object Hittable, offset Vec3) Translate {
	bbox := object.BoundingBox().AddVec3(offset)
	return Translate{object, offset, bbox}
}

func (hit Translate) Hit(r Ray, intvl Interval) (bool, HitRecord) {
	// Move the ray backwards by the offset
	offsetR := Ray{r.orig.Sub(hit.offset), r.dir, r.tm}

	// Determine whether an intersection exists along the offset ray (and if so, where)
	hitAnything, rec := hit.object.Hit(offsetR, intvl)
	if !hitAnything {
		return false, HitRecord{}
	}

	// Move the intersection point forwards by the offset
	rec.p = rec.p.Add(hit.offset)

	return true, rec
}

func (hit Translate) BoundingBox() AABB {
	return hit.bbox
}

func (hit Translate) PDFValue(origin Point3, direction Vec3) float64 {
	return 0
}

func (hit Translate) Random(origin Point3) Vec3 {
	return Vec3{1, 0, 0}
}

type RotateY struct {
	object   Hittable
	sinTheta float64
	cosTheta float64
	bbox     AABB
}

func NewRotateY(object Hittable, angle float64) RotateY {
	radians := Radians(angle)
	sinTheta := math.Sin(radians)
	cosTheta := math.Cos(radians)
	bbox := object.BoundingBox()

	min := Point3{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	max := Point3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

	for i := range 2 {
		for j := range 2 {
			for k := range 2 {
				x := float64(i)*bbox[0].max + (1-float64(i))*bbox[0].min
				y := float64(j)*bbox[1].max + (1-float64(j))*bbox[1].min
				z := float64(k)*bbox[2].max + (1-float64(k))*bbox[2].min

				newx := cosTheta*x + sinTheta*z
				newz := -sinTheta*x + cosTheta*z

				tester := Vec3{newx, y, newz}

				for c := range 3 {
					min[c] = math.Min(min[c], tester[c])
					max[c] = math.Max(max[c], tester[c])
				}
			}
		}
	}

	bbox = NewAABBPoint(min, max)

	return RotateY{object, sinTheta, cosTheta, bbox}
}

func (hit RotateY) Hit(r Ray, intvl Interval) (bool, HitRecord) {
	// Transform the ray from world space to object space.
	origin := Point3{hit.cosTheta*r.orig.X() - hit.sinTheta*r.orig.Z(),
		r.orig.Y(),
		hit.sinTheta*r.orig.X() + hit.cosTheta*r.orig.Z()}

	direction := Vec3{hit.cosTheta*r.dir.X() - hit.sinTheta*r.dir.Z(),
		r.dir.Y(),
		hit.sinTheta*r.dir.X() + hit.cosTheta*r.dir.Z()}

	roratedR := Ray{origin, direction, r.tm}

	// Determine whether an intersection exists in object space (and if so, where).
	hitAnything, rec := hit.object.Hit(roratedR, intvl)
	if !hitAnything {
		return false, HitRecord{}
	}

	// Transform the intersection from object space back to world space.
	rec.p = Point3{hit.cosTheta*rec.p.X() + hit.sinTheta*rec.p.Z(),
		rec.p.Y(),
		-hit.sinTheta*rec.p.X() + hit.cosTheta*rec.p.Z()}

	rec.normal = Vec3{hit.cosTheta*rec.normal.X() + hit.sinTheta*rec.normal.Z(),
		rec.normal.Y(),
		-hit.sinTheta*rec.normal.X() + hit.cosTheta*rec.normal.Z()}

	return true, rec
}

func (hit RotateY) BoundingBox() AABB {
	return hit.bbox
}

func (hit RotateY) PDFValue(origin Point3, direction Vec3) float64 {
	return 0
}

func (hit RotateY) Random(origin Point3) Vec3 {
	return Vec3{1, 0, 0}
}

type ConstantMedium struct {
	boundary      Hittable
	negInvDensity float64
	phaseFunction Material
}

func NewConstantMedium(boundary Hittable, density float64, tex Texture) ConstantMedium {
	return ConstantMedium{boundary, -1 / density, Isotropic{tex}}
}

func (hit ConstantMedium) Hit(r Ray, intvl Interval) (bool, HitRecord) {
	hitAnything1, rec1 := hit.boundary.Hit(r, UniverseInterval)
	if !hitAnything1 {
		return false, HitRecord{}
	}
	hitAnything2, rec2 := hit.boundary.Hit(r, Interval{rec1.t + 0.0001, math.MaxFloat64})
	if !hitAnything2 {
		return false, HitRecord{}
	}

	if rec1.t < intvl.min {
		rec1.t = intvl.min
	}

	if rec2.t > intvl.max {
		rec2.t = intvl.max
	}

	if rec1.t >= rec2.t {
		return false, HitRecord{}
	}

	if rec1.t < 0 {
		rec1.t = 0
	}

	rayLength := r.dir.Length()
	distanceInsideBoundary := (rec2.t - rec1.t) * rayLength
	hitDistance := hit.negInvDensity * math.Log(rand.Float64())

	if hitDistance > distanceInsideBoundary {
		return false, HitRecord{}
	}

	rec := HitRecord{}
	rec.t = rec1.t + hitDistance/rayLength
	rec.p = r.At(rec.t)

	rec.normal = Vec3{1, 0, 0} // arbitrary
	rec.frontFace = true       // also arbitrary
	rec.mat = hit.phaseFunction

	return true, rec
}

func (hit ConstantMedium) BoundingBox() AABB {
	return hit.boundary.BoundingBox()
}

func (hit ConstantMedium) PDFValue(origin Point3, direction Vec3) float64 {
	return 0
}

func (hit ConstantMedium) Random(origin Point3) Vec3 {
	return Vec3{1, 0, 0}
}
