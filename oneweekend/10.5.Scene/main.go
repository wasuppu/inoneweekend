package main

func main() {
	world := HittableList{}

	materialGround := Lambertian{RGB{0.8, 0.8, 0.0}}
	materialCenter := Lambertian{RGB{0.1, 0.2, 0.5}}
	materialLeft := Metal{RGB{0.8, 0.8, 0.8}}
	materialRight := Metal{RGB{0.8, 0.6, 0.2}}

	world.Add(Sphere{Point3{0, -100.5, -1}, 100, materialGround})
	world.Add(Sphere{Point3{0, 0, -1.2}, 0.5, materialCenter})
	world.Add(Sphere{Point3{-1, 0, -1}, 0.5, materialLeft})
	world.Add(Sphere{Point3{1, 0, -1}, 0.5, materialRight})

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.Render(world)
}
