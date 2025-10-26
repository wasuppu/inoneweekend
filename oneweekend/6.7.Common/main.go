package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	// Image
	aspectRadio := 16.0 / 9.0
	imageWidth := 400

	// Calculate the image height, and ensure that it's at least 1.
	imageHeight := max(int(float64(imageWidth)/aspectRadio), 1)

	// World
	world := HittableList{}
	world.Add(Sphere{Point3{0, 0, -1}, 0.5})
	world.Add(Sphere{Point3{0, -100.5, -1}, 100})
	// Camera
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	CameraCenter := Point3{0, 0, 0}

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := Vec3{viewportWidth, 0, 0}
	viewportV := Vec3{0, -viewportHeight, 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	pixelDeltaU := viewportU.Divn(float64(imageWidth))
	pixelDeltaV := viewportV.Divn(float64(imageHeight))

	// Calculate the location of the upper left pixel.
	viewportUpperLeft := CameraCenter.Sub(Vec3{0, 0, focalLength}).Sub(viewportU.Divn(2)).Sub(viewportV.Divn(2))
	pixel00Loc := viewportUpperLeft.Add(pixelDeltaU.Add(pixelDeltaV).Muln(0.5))

	// Render
	framebuffer := make([]color.Color, imageWidth*imageHeight)
	for j := range imageHeight {
		log.Printf("\rScanlines remaining: %d", imageHeight-j)
		for i := range imageWidth {
			pixelCenter := pixel00Loc.Add(pixelDeltaU.Muln(float64(i))).Add(pixelDeltaV.Muln(float64(j)))
			rayDirection := pixelCenter.Sub(CameraCenter)
			r := Ray{CameraCenter, rayDirection}
			framebuffer[i+j*imageWidth] = r.Color(world)
		}
	}
	log.Println("Done.")

	WritePng("out", framebuffer, imageWidth, imageHeight)
}

func WritePng(name string, pixels []color.Color, imageWidth, imageHeight int) {
	f, _ := os.Create(name + ".png")
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	for j := range imageHeight {
		for i := range imageWidth {
			img.Set(i, j, pixels[i+j*imageWidth])
		}
	}

	png.Encode(f, img)
}
