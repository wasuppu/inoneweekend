package main

func main() {
	world := HittableList{}

	materialGround := Lambertian{RGB{0.8, 0.8, 0.0}}
	materialCenter := Lambertian{RGB{0.1, 0.2, 0.5}}
	materialLeft := Dielectric{1.50}
	materialBubble := Dielectric{1.00 / 1.50}
	materialRight := Metal{RGB{0.8, 0.6, 0.2}, 1.0}

	world.Add(Sphere{Point3{0, -100.5, -1}, 100, materialGround})
	world.Add(Sphere{Point3{0, 0, -1.2}, 0.5, materialCenter})
	world.Add(Sphere{Point3{-1, 0, -1}, 0.5, materialLeft})
	world.Add(Sphere{Point3{-1, 0, -1}, 0.4, materialBubble})
	world.Add(Sphere{Point3{1, 0, -1}, 0.5, materialRight})

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50

	cam.vfov = 20
	cam.lookfrom = Point3{-2, 2, 1}
	cam.lookat = Point3{0, 0, -1}
	cam.vup = Vec3{0, 1, 0}

	cam.Render(world)
}
