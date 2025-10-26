package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

type Camera struct {
	aspectRadio       float64 // Ratio of image width over height
	imageWidth        int     // Rendered image width in pixel count
	samplesPerPixel   int     // Count of random samples for each pixel
	maxDepth          int     // Maximum number of ray bounces into scene
	vfov              float64 // Vertical view angle (field of view)
	lookfrom          Point3  // Point camera is looking from
	lookat            Point3  // Point camera is looking at
	vup               Vec3    // Camera-relative "up" direction
	imageHeight       int     // Rendered image height
	pixelSamplesScale float64 // Color scale factor for a sum of pixel samples
	center            Point3  // Camera center
	pixel00Loc        Point3  // Location of pixel 0, 0
	pixelDeltaU       Vec3    // Offset to pixel to the right
	pixelDeltaV       Vec3    // Offset to pixel below
	u, v, w           Vec3    // Camera frame basis vectors
}

func DefaultCamera() Camera {
	return Camera{
		aspectRadio:     1.0,
		imageWidth:      100,
		samplesPerPixel: 10,
		maxDepth:        10,
		vfov:            90,
		lookfrom:        Point3{0, 0, 0},
		lookat:          Point3{0, 0, -1},
		vup:             Vec3{0, 1, 0},
	}
}

func (c *Camera) Initialize() {
	// Calculate the image height, and ensure that it's at least 1.
	c.imageHeight = max(int(float64(c.imageWidth)/c.aspectRadio), 1)

	c.pixelSamplesScale = 1.0 / float64(c.samplesPerPixel)

	c.center = c.lookfrom

	// Determine viewport dimensions.
	focalLength := c.lookfrom.Sub(c.lookat).Length()
	theta := Radians(c.vfov)
	h := math.Tan(theta / 2)
	viewportHeight := 2 * h * focalLength
	viewportWidth := viewportHeight * (float64(c.imageWidth) / float64(c.imageHeight))

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	c.w = c.lookfrom.Sub(c.lookat).Normalize()
	c.u = c.vup.Cross(c.w).Normalize()
	c.v = c.w.Cross(c.u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := c.u.Muln(viewportWidth)   // Vector across viewport horizontal edge
	viewportV := c.v.Muln(-viewportHeight) // Vector down viewport vertical edge

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	c.pixelDeltaU = viewportU.Divn(float64(c.imageWidth))
	c.pixelDeltaV = viewportV.Divn(float64(c.imageHeight))

	// Calculate the location of the upper left pixel.
	viewportUpperLeft := c.center.Sub(c.w.Muln(focalLength)).Sub(viewportU.Divn(2)).Sub(viewportV.Divn(2))
	c.pixel00Loc = viewportUpperLeft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Muln(0.5))
}

func (c *Camera) Render(world Hittable) {
	c.Initialize()

	framebuffer := make([]color.Color, c.imageWidth*c.imageHeight)
	for j := range c.imageHeight {
		log.Printf("\rScanlines remaining: %d", c.imageHeight-j)
		for i := range c.imageWidth {
			pixelColor := RGB{0, 0, 0}
			for range c.samplesPerPixel {
				r := c.GetRay(i, j)
				pixelColor = pixelColor.Add(c.RayColor(r, c.maxDepth, world))
			}

			framebuffer[i+j*c.imageWidth] = pixelColor.Muln(c.pixelSamplesScale).Color()
		}
	}
	log.Println("Done.")

	WritePng("out", framebuffer, c.imageWidth, c.imageHeight)
}

func (c *Camera) GetRay(i, j int) Ray {
	// Construct a Camera ray Originating from the Origin and Directed at randomly sampled point around the pixel location i, j.
	offset := SampleSquare()
	pixelSample := c.pixel00Loc.Add(c.pixelDeltaU.Muln(float64(i) + offset.X())).Add(c.pixelDeltaV.Muln(float64(j) + offset.Y()))

	orig := c.center
	dir := pixelSample.Sub(orig)
	return Ray{orig, dir}
}

func (c Camera) RayColor(r Ray, depth int, world Hittable) RGB {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return RGB{0, 0, 0}
	}

	if hitAnything, rec := world.Hit(r, Interval{0.001, math.MaxFloat64}); hitAnything {
		if ok, attenuation, scattered := rec.mat.Scatter(r, rec); ok {
			return attenuation.Mul(c.RayColor(scattered, depth-1, world))
		}
		return RGB{0, 0, 0}
	}

	unitDirection := r.dir.Normalize()
	a := 0.5 * (unitDirection.Y() + 1.0)
	return RGB{1, 1, 1}.Muln(1.0 - a).Add(RGB{0.5, 0.7, 1.0}.Muln(a))
}

func SampleSquare() Vec3 {
	// Returns the vector to a random point in the [-.5,-.5]-[+.5,+.5] unit square.
	return Vec3{rand.Float64() - 0.5, rand.Float64() - 0.5, 0}
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
