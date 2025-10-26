package main

type AABB [3]Interval

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

	return aabb
}

func NewAABBBox(box0, box1 AABB) AABB {
	x := NewInterval(box0[0], box1[0])
	y := NewInterval(box0[1], box1[1])
	z := NewInterval(box0[2], box1[2])
	return AABB{x, y, z}
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
