package main

type Ray struct {
	orig Point3
	dir  Vec3
}

func (r Ray) At(t float64) Vec3 {
	return r.orig.Add(r.dir.Muln(t))
}
