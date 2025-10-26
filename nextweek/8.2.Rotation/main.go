package main

import (
	"math/rand/v2"
	"path"
	"path/filepath"
	"runtime"
)

var (
	basepath string
	rootpath string
)

func init() {
	_, exepath, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(exepath)
	rootpath = filepath.Dir(filepath.Dir(basepath))
}

func BouncingSpheres() {
	world := HittableList{}

	checker := NewCheckerTexture(0.32, SolidColor{RGB{0.2, 0.3, 0.1}}, SolidColor{RGB{0.9, 0.9, 0.9}})
	groundMaterial := Lambertian{checker}
	world.Add(NewSphere(Point3{0, -1000, 0}, 1000, groundMaterial))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := Point3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.Sub(Point3{4, 0.2, 0}).Length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					albedo := RandomVec3().Mul(RandomVec3())
					sphereMaterial := Lambertian{SolidColor{albedo}}
					center2 := center.Add(Vec3{0, RandomRange(0, 0.5), 0})
					world.Add(NewMotionSphere(center, center2, 0.2, sphereMaterial))
				} else if chooseMat < 0.95 {
					// metal
					albedo := RandomVec3Range(0.5, 1)
					fuzz := RandomRange(0, 0.5)
					sphereMaterial := Metal{albedo, fuzz}
					world.Add(NewSphere(center, 0.2, sphereMaterial))
				} else {
					// glass
					sphereMaterial := Dielectric{1.5}
					world.Add(NewSphere(center, 0.2, sphereMaterial))
				}
			}
		}
	}

	material1 := Dielectric{1.50}
	material2 := Lambertian{SolidColor{RGB{0.4, 0.2, 0.1}}}
	material3 := Metal{RGB{0.7, 0.6, 0.5}, 0.0}

	world.Add(NewSphere(Point3{0, 1, 0}, 1.0, material1))
	world.Add(NewSphere(Point3{-4, 1, 0}, 1.0, material2))
	world.Add(NewSphere(Point3{4, 1, 0}, 1.0, material3))

	bvh := NewBVHNode(world)

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.background = RGB{0.70, 0.80, 1.00}

	cam.vfov = 20
	cam.lookfrom = Point3{13, 2, 3}
	cam.lookat = Point3{0, 0, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0.6
	cam.focusDist = 10.0

	cam.Render(bvh)
}

func CheckeredSpheres() {
	world := HittableList{}

	checker := NewCheckerTexture(0.32, SolidColor{RGB{0.2, 0.3, 0.1}}, SolidColor{RGB{0.9, 0.9, 0.9}})
	world.Add(NewSphere(Point3{0, -10, 0}, 10, Lambertian{checker}))
	world.Add(NewSphere(Point3{0, 10, 0}, 10, Lambertian{checker}))

	bvh := NewBVHNode(world)

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.background = RGB{0.70, 0.80, 1.00}

	cam.vfov = 20
	cam.lookfrom = Point3{13, 2, 3}
	cam.lookat = Point3{0, 0, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0

	cam.Render(bvh)
}

func Earth() {
	world := HittableList{}

	earthTexture := NewImageTexture(path.Join(rootpath, "textures", "earthmap.jpg"))
	earthSurface := Lambertian{earthTexture}
	globe := NewSphere(Point3{0, 0, 0}, 2, earthSurface)

	world.Add(globe)

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.background = RGB{0.70, 0.80, 1.00}

	cam.vfov = 20
	cam.lookfrom = Point3{0, 0, 12}
	cam.lookat = Point3{0, 0, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0

	cam.Render(world)
}

func PerlinSpheres() {
	world := HittableList{}

	pertext := NewNoiseTexture(4)
	world.Add(NewSphere(Point3{0, -1000, 0}, 1000, Lambertian{pertext}))
	world.Add(NewSphere(Point3{0, 2, 0}, 2, Lambertian{pertext}))

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.background = RGB{0.70, 0.80, 1.00}

	cam.vfov = 20
	cam.lookfrom = Point3{13, 2, 3}
	cam.lookat = Point3{0, 0, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0

	cam.Render(world)
}

func Quads() {
	world := HittableList{}

	// Materials
	leftRed := Lambertian{NewSolidColor(1.0, 0.2, 0.2)}
	backGreen := Lambertian{NewSolidColor(0.2, 1.0, 0.2)}
	rightBlue := Lambertian{NewSolidColor(0.2, 0.2, 1.0)}
	upperOrange := Lambertian{NewSolidColor(1.0, 0.5, 0.0)}
	lowerTeal := Lambertian{NewSolidColor(0.2, 0.8, 0.8)}

	// Quads
	world.Add(NewQuad(Point3{-3, -2, 5}, Vec3{0, 0, -4}, Vec3{0, 4, 0}, leftRed))
	world.Add(NewQuad(Point3{-2, -2, 0}, Vec3{4, 0, 0}, Vec3{0, 4, 0}, backGreen))
	world.Add(NewQuad(Point3{3, -2, 1}, Vec3{0, 0, 4}, Vec3{0, 4, 0}, rightBlue))
	world.Add(NewQuad(Point3{-2, 3, 1}, Vec3{4, 0, 0}, Vec3{0, 0, 4}, upperOrange))
	world.Add(NewQuad(Point3{-2, -3, 5}, Vec3{4, 0, 0}, Vec3{0, 0, -4}, lowerTeal))

	cam := DefaultCamera()
	cam.aspectRadio = 1.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.background = RGB{0.70, 0.80, 1.00}

	cam.vfov = 80
	cam.lookfrom = Point3{0, 0, 9}
	cam.lookat = Point3{0, 0, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0

	cam.Render(world)
}

func SimpleLight() {
	world := HittableList{}

	pertext := NewNoiseTexture(4)
	world.Add(NewSphere(Point3{0, -1000, 0}, 1000, Lambertian{pertext}))
	world.Add(NewSphere(Point3{0, 2, 0}, 2, Lambertian{pertext}))

	difflight := DiffuseLight{NewSolidColor(4, 4, 4)}
	world.Add(NewSphere(Point3{0, 7, 0}, 2, difflight))
	world.Add(NewQuad(Point3{3, 1, -2}, Vec3{2, 0, 0}, Vec3{0, 2, 0}, difflight))

	cam := DefaultCamera()
	cam.aspectRadio = 16.0 / 9.0
	cam.imageWidth = 400
	cam.samplesPerPixel = 100
	cam.maxDepth = 50
	cam.background = RGB{0, 0, 0}

	cam.vfov = 20
	cam.lookfrom = Point3{26, 3, 6}
	cam.lookat = Point3{0, 2, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0

	cam.Render(world)
}

func CornellBox() {
	world := HittableList{}

	red := Lambertian{NewSolidColor(0.65, 0.05, 0.05)}
	white := Lambertian{NewSolidColor(0.73, 0.73, 0.73)}
	green := Lambertian{NewSolidColor(0.12, 0.45, 0.15)}
	light := DiffuseLight{NewSolidColor(15, 15, 15)}

	world.Add(NewQuad(Point3{555, 0, 0}, Vec3{0, 555, 0}, Vec3{0, 0, 555}, green))
	world.Add(NewQuad(Point3{0, 0, 0}, Vec3{0, 555, 0}, Vec3{0, 0, 555}, red))
	world.Add(NewQuad(Point3{343, 554, 332}, Vec3{-130, 0, 0}, Vec3{0, 0, -105}, light))
	world.Add(NewQuad(Point3{0, 0, 0}, Vec3{555, 0, 0}, Vec3{0, 0, 555}, white))
	world.Add(NewQuad(Point3{555, 555, 555}, Vec3{-555, 0, 0}, Vec3{0, 0, -555}, white))
	world.Add(NewQuad(Point3{0, 0, 555}, Vec3{555, 0, 0}, Vec3{0, 555, 0}, white))

	world.Add(NewTranslate(NewRotateY(Box(Point3{0, 0, 0}, Point3{165, 330, 165}, white), 15), Vec3{265, 0, 295}))
	world.Add(NewTranslate(NewRotateY(Box(Point3{0, 0, 0}, Point3{165, 165, 165}, white), -18), Vec3{130, 0, 65}))

	cam := DefaultCamera()
	cam.aspectRadio = 1.0
	cam.imageWidth = 600
	cam.samplesPerPixel = 200
	cam.maxDepth = 50
	cam.background = RGB{0, 0, 0}

	cam.vfov = 40
	cam.lookfrom = Point3{278, 278, -800}
	cam.lookat = Point3{278, 278, 0}
	cam.vup = Vec3{0, 1, 0}

	cam.defocusAngle = 0

	cam.Render(world)
}

func main() {
	switch 7 {
	case 1:
		BouncingSpheres()
	case 2:
		CheckeredSpheres()
	case 3:
		Earth()
	case 4:
		PerlinSpheres()
	case 5:
		Quads()
	case 6:
		SimpleLight()
	case 7:
		CornellBox()
	}
}
