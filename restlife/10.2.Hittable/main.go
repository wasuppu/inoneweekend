package main

func main() {
	world := HittableList{}

	red := Lambertian{NewSolidColor(0.65, 0.05, 0.05)}
	white := Lambertian{NewSolidColor(0.73, 0.73, 0.73)}
	green := Lambertian{NewSolidColor(0.12, 0.45, 0.15)}
	light := DiffuseLight{NewSolidColor(15, 15, 15)}

	// Cornell box sides
	world.Add(NewQuad(Point3{555, 0, 0}, Vec3{0, 0, 555}, Vec3{0, 555, 0}, green))
	world.Add(NewQuad(Point3{0, 0, 555}, Vec3{0, 0, -555}, Vec3{0, 555, 0}, red))
	world.Add(NewQuad(Point3{0, 555, 0}, Vec3{555, 0, 0}, Vec3{0, 0, 555}, white))
	world.Add(NewQuad(Point3{0, 0, 555}, Vec3{555, 0, 0}, Vec3{0, 0, -555}, white))
	world.Add(NewQuad(Point3{555, 0, 555}, Vec3{-555, 0, 0}, Vec3{0, 555, 0}, white))

	// Light
	world.Add(NewQuad(Point3{213, 554, 227}, Vec3{130, 0, 0}, Vec3{0, 0, 105}, light))

	// Box
	world.Add(NewTranslate(NewRotateY(Box(Point3{0, 0, 0}, Point3{165, 330, 165}, white), 15), Vec3{265, 0, 295}))
	world.Add(NewTranslate(NewRotateY(Box(Point3{0, 0, 0}, Point3{165, 165, 165}, white), -18), Vec3{130, 0, 65}))

	// Light Sources
	lights := NewQuad(Point3{343, 554, 332}, Vec3{-130, 0, 0}, Vec3{0, 0, -105}, EmptyMaterial{})

	cam := DefaultCamera()
	cam.aspectRadio = 1.0
	cam.imageWidth = 600
	cam.samplesPerPixel = 10
	cam.maxDepth = 50
	cam.background = RGB{0, 0, 0}

	cam.vfov = 40
	cam.lookfrom = Point3{278, 278, -800}
	cam.lookat = Point3{278, 278, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0

	cam.Render(world, lights)
}
