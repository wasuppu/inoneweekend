package main

import (
	"math/rand/v2"
)

func main() {
	world := HittableList{}

	groundMaterial := Lambertian{RGB{0.5, 0.5, 0.5}}
	world.Add(Sphere{Point3{0, -1000, 0}, 1000, groundMaterial})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := Point3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.Sub(Point3{4, 0.2, 0}).Length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					albedo := RandomVec3().Mul(RandomVec3())
					sphereMaterial := Lambertian{albedo}
					world.Add(Sphere{center, 0.2, sphereMaterial})
				} else if chooseMat < 0.95 {
					// metal
					albedo := RandomVec3Range(0.5, 1)
					fuzz := RandomRange(0, 0.5)
					sphereMaterial := Metal{albedo, fuzz}
					world.Add(Sphere{center, 0.2, sphereMaterial})
				} else {
					// glass
					sphereMaterial := Dielectric{1.5}
					world.Add(Sphere{center, 0.2, sphereMaterial})
				}
			}
		}
	}

	material1 := Dielectric{1.50}
	material2 := Lambertian{RGB{0.4, 0.2, 0.1}}
	material3 := Metal{RGB{0.7, 0.6, 0.5}, 0.0}

	world.Add(Sphere{Point3{0, 1, 0}, 1.0, material1})
	world.Add(Sphere{Point3{-4, 1, 0}, 1.0, material2})
	world.Add(Sphere{Point3{4, 1, 0}, 1.0, material3})

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 1200
	cam.samplesPerPixel = 10
	cam.maxDepth = 50

	cam.vfov = 20
	cam.lookfrom = Point3{13, 2, 3}
	cam.lookat = Point3{0, 0, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0.6
	cam.focusDist = 10.0

	cam.Render(world)
}
