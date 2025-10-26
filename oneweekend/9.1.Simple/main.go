package main

func main() {
	world := HittableList{}
	world.Add(Sphere{Point3{0, 0, -1}, 0.5})
	world.Add(Sphere{Point3{0, -100.5, -1}, 100})

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.Render(world)
}
