package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

type Camera struct {
	aspectRadio float64 // Ratio of image width over height
	imageWidth  int     // Rendered image width in pixel count
	imageHeight int     // Rendered image height
	center      Point3  // Camera center
	pixel00Loc  Point3  // Location of pixel 0, 0
	pixelDeltaU Vec3    // Offset to pixel to the right
	pixelDeltaV Vec3    // Offset to pixel below
}

func DefaultCamera() Camera {
	return Camera{aspectRadio: 1.0, imageWidth: 100}
}

func (c *Camera) Initialize() {
	// Calculate the image height, and ensure that it's at least 1.
	c.imageHeight = max(int(float64(c.imageWidth)/c.aspectRadio), 1)

	c.center = Point3{0, 0, 0}

	// Determine viewport dimensions.
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(c.imageWidth) / float64(c.imageHeight))

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := Vec3{viewportWidth, 0, 0}
	viewportV := Vec3{0, -viewportHeight, 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	c.pixelDeltaU = viewportU.Divn(float64(c.imageWidth))
	c.pixelDeltaV = viewportV.Divn(float64(c.imageHeight))

	// Calculate the location of the upper left pixel.
	viewportUpperLeft := c.center.Sub(Vec3{0, 0, focalLength}).Sub(viewportU.Divn(2)).Sub(viewportV.Divn(2))
	c.pixel00Loc = viewportUpperLeft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Muln(0.5))
}

func (c *Camera) Render(world Hittable) {
	c.Initialize()

	framebuffer := make([]color.Color, c.imageWidth*c.imageHeight)
	for j := range c.imageHeight {
		log.Printf("\rScanlines remaining: %d", c.imageHeight-j)
		for i := range c.imageWidth {
			pixelCenter := c.pixel00Loc.Add(c.pixelDeltaU.Muln(float64(i))).Add(c.pixelDeltaV.Muln(float64(j)))
			rayDirection := pixelCenter.Sub(c.center)
			r := Ray{c.center, rayDirection}
			framebuffer[i+j*c.imageWidth] = c.RayColor(r, world)
		}
	}
	log.Println("Done.")

	WritePng("out", framebuffer, c.imageWidth, c.imageHeight)
}

func (c Camera) RayColor(r Ray, world Hittable) color.Color {
	if hitAnything, rec := world.Hit(r, Interval{0, math.MaxFloat64}); hitAnything {
		return rec.normal.Add(RGB{1, 1, 1}).Muln(0.5).Color()
	}

	unitDirection := r.dir.Normalize()
	a := 0.5 * (unitDirection.Y() + 1.0)
	return RGB{1, 1, 1}.Muln(1.0 - a).Add(RGB{0.5, 0.7, 1.0}.Muln(a)).Color()
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
