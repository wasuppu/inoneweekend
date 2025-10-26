package main

import "math"

type ONB [3]Vec3

func NewONB(n Vec3) ONB {
	onb := ONB{}
	onb[2] = n.Normalize()
	a := Vec3{1, 0, 0}
	if math.Abs(onb[2].X()) > 0.9 {
		a = Vec3{0, 1, 0}
	}
	onb[1] = onb[2].Cross(a).Normalize()
	onb[0] = onb[2].Cross(onb[1])
	return onb
}

func (onb ONB) U() Vec3 {
	return onb[0]
}

func (onb ONB) V() Vec3 {
	return onb[1]
}

func (onb ONB) W() Vec3 {
	return onb[2]
}

func (onb ONB) Transform(v Vec3) Vec3 {
	// Transform from basis coordinates to local space.
	return onb[0].Muln(v[0]).Add(onb[1].Muln(v[1])).Add(onb[2].Muln(v[2]))
}
