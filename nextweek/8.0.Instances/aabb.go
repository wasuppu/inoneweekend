package main

type AABB [3]Interval

var EmptyAABB = AABB{EmptyInterval, EmptyInterval, EmptyInterval}
var UniverseAABB = AABB{UniverseInterval, UniverseInterval, UniverseInterval}

func NewAABBB(x, y, z Interval) AABB {
	aabb := AABB{x, y, z}
	aabb.PadToMinimums()
	return aabb
}

func NewAABBPoint(a, b Point3) AABB {
	aabb := AABB{}
	if a[0] <= b[0] {
		aabb[0] = Interval{a[0], b[0]}
	} else {
		aabb[0] = Interval{b[0], a[0]}
	}

	if a[1] <= b[1] {
		aabb[1] = Interval{a[1], b[1]}
	} else {
		aabb[1] = Interval{b[1], a[1]}
	}

	if a[2] <= b[2] {
		aabb[2] = Interval{a[2], b[2]}
	} else {
		aabb[2] = Interval{b[2], a[2]}
	}

	aabb.PadToMinimums()
	return aabb
}

func NewAABBBox(box0, box1 AABB) AABB {
	x := NewInterval(box0[0], box1[0])
	y := NewInterval(box0[1], box1[1])
	z := NewInterval(box0[2], box1[2])
	aabb := AABB{x, y, z}
	aabb.PadToMinimums()
	return aabb
}

func (aabb AABB) Hit(r Ray, intvl Interval) bool {
	for axis := range 3 {
		ax := aabb[axis]
		adinv := 1.0 / r.dir[axis]

		t0 := (ax.min - r.orig[axis]) * adinv
		t1 := (ax.max - r.orig[axis]) * adinv

		if t0 < t1 {
			if t0 > intvl.min {
				intvl.max = t0
			}
			if t1 < intvl.max {
				intvl.max = t1
			}
		} else {
			if t1 > intvl.min {
				intvl.min = t1
			}
			if t0 < intvl.max {
				intvl.max = t0
			}
		}

		if intvl.max <= intvl.min {
			return false
		}
	}

	return true
}

func (aabb AABB) LongestAxis() int {
	// Returns the index of the longest axis of the bounding box.
	if aabb[0].Size() > aabb[1].Size() {
		if aabb[0].Size() > aabb[2].Size() {
			return 0
		} else {
			return 2
		}
	} else {
		if aabb[1].Size() > aabb[2].Size() {
			return 1
		} else {
			return 2
		}
	}
}

func (aabb *AABB) PadToMinimums() {
	// Adjust the AABB so that no side is narrower than some delta, padding if necessary.

	delta := 0.0001
	if aabb[0].Size() < delta {
		aabb[0] = aabb[0].Expand(delta)
	}
	if aabb[1].Size() < delta {
		aabb[1] = aabb[1].Expand(delta)
	}
	if aabb[2].Size() < delta {
		aabb[2] = aabb[2].Expand(delta)
	}
}
