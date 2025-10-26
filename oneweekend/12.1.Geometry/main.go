package main

import "math"

func main() {
	world := HittableList{}

	R := math.Cos(math.Pi / 4)

	materialLeft := Lambertian{RGB{0, 0, 1}}
	materialRight := Lambertian{RGB{1, 0, 0}}

	world.Add(Sphere{Point3{-R, 0, -1}, R, materialLeft})
	world.Add(Sphere{Point3{R, 0, -1}, R, materialRight})

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.vfov = 90
	cam.Render(world)
}
