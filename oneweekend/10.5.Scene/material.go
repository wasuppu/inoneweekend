package main

type Material interface {
	Scatter(in Ray, rec HitRecord) (bool, RGB, Ray)
}

type Lambertian struct {
	albedo RGB
}

func (m Lambertian) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	scatterDirection := rec.normal.Add(RandomUnitVector())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.normal
	}

	scattered := Ray{rec.p, scatterDirection}
	attenuation := m.albedo
	return true, attenuation, scattered
}

type Metal struct {
	albedo RGB
}

func (m Metal) Scatter(in Ray, rec HitRecord) (bool, RGB, Ray) {
	reflected := Reflect(in.dir, rec.normal)
	scattered := Ray{rec.p, reflected}
	attenuation := m.albedo
	return true, attenuation, scattered
}
