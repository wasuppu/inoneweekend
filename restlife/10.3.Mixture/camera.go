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
	background        RGB     // Scene background color
	vfov              float64 // Vertical view angle (field of view)
	lookfrom          Point3  // Point camera is looking from
	lookat            Point3  // Point camera is looking at
	vup               Vec3    // Camera-relative "up" direction
	defocusAngle      float64 // Variation angle of rays through each pixel
	focusDist         float64 // Distance from camera lookfrom point to plane of perfect focus
	imageHeight       int     // Rendered image height
	pixelSamplesScale float64 // Color scale factor for a sum of pixel samples
	sqrtSpp           int     // Square root of number of samples per pixel
	recipSqrtSpp      float64 // 1 / sqrt_spp
	center            Point3  // Camera center
	pixel00Loc        Point3  // Location of pixel 0, 0
	pixelDeltaU       Vec3    // Offset to pixel to the right
	pixelDeltaV       Vec3    // Offset to pixel below
	u, v, w           Vec3    // Camera frame basis vectors
	defocusDiskU      Vec3    // Defocus disk horizontal radius
	defocusDiskV      Vec3    // Defocus disk vertical radius
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
		defocusAngle:    0,
		focusDist:       10,
	}
}

func (c *Camera) Initialize() {
	// Calculate the image height, and ensure that it's at least 1.
	c.imageHeight = max(int(float64(c.imageWidth)/c.aspectRadio), 1)

	c.sqrtSpp = int(math.Sqrt(float64(c.samplesPerPixel)))
	c.pixelSamplesScale = 1.0 / float64(c.sqrtSpp*c.sqrtSpp)
	c.recipSqrtSpp = 1.0 / float64(c.sqrtSpp)

	c.pixelSamplesScale = 1.0 / float64(c.samplesPerPixel)

	c.center = c.lookfrom

	// Determine viewport dimensions.
	theta := Radians(c.vfov)
	h := math.Tan(theta / 2)
	viewportHeight := 2 * h * c.focusDist
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
	viewportUpperLeft := c.center.Sub(c.w.Muln(c.focusDist)).Sub(viewportU.Divn(2)).Sub(viewportV.Divn(2))
	c.pixel00Loc = viewportUpperLeft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Muln(0.5))

	// Calculate the camera defocus disk basis vectors.
	defocusRadius := c.focusDist * math.Tan(Radians(c.defocusAngle/2))
	c.defocusDiskU = c.u.Muln(defocusRadius)
	c.defocusDiskV = c.v.Muln(defocusRadius)
}

func (c *Camera) Render(world Hittable, lights Hittable) {
	c.Initialize()

	framebuffer := make([]color.Color, c.imageWidth*c.imageHeight)
	for j := range c.imageHeight {
		log.Printf("\rScanlines remaining: %d", c.imageHeight-j)
		for i := range c.imageWidth {
			pixelColor := RGB{0, 0, 0}
			for sj := range c.sqrtSpp {
				for si := range c.sqrtSpp {
					r := c.GetRay(i, j, si, sj)
					pixelColor = pixelColor.Add(c.RayColor(r, c.maxDepth, world, lights))
				}
			}
			framebuffer[i+j*c.imageWidth] = pixelColor.Muln(c.pixelSamplesScale).Color()
		}
	}
	log.Println("Done.")

	WritePng("3-10.3", framebuffer, c.imageWidth, c.imageHeight)
}

func (c *Camera) GetRay(i, j, si, sj int) Ray {
	// Construct a Camera ray Originating from the Origin and Directed at randomly sampled point around the pixel location i, j.
	offset := c.SampleSquareStratified(si, sj)
	pixelSample := c.pixel00Loc.Add(c.pixelDeltaU.Muln(float64(i) + offset.X())).Add(c.pixelDeltaV.Muln(float64(j) + offset.Y()))

	orig := c.DefocusDiskSample()
	if c.defocusAngle <= 0 {
		orig = c.center
	}
	dir := pixelSample.Sub(orig)
	tm := rand.Float64()
	return Ray{orig, dir, tm}
}

func (c Camera) RayColor(r Ray, depth int, world Hittable, lights Hittable) RGB {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return RGB{0, 0, 0}
	}

	// If the ray hits nothing, return the background color.
	hitAnything, rec := world.Hit(r, Interval{0.001, math.MaxFloat64})
	if !hitAnything {
		return c.background
	}

	colorFromEmission := rec.mat.Emitted(r, rec, rec.u, rec.v, rec.p)
	ok, attenuation, _, _ := rec.mat.Scatter(r, rec)
	if !ok {
		return colorFromEmission
	}

	p0 := HittablePDF{lights, rec.p}
	p1 := NewCosinePDF(rec.normal)
	mixedPdf := MixturePDF{p0, p1}

	scattered := Ray{rec.p, mixedPdf.Generate(), r.tm}
	pdfValue := mixedPdf.Value(scattered.dir)

	scatteringPdf := rec.mat.ScatteringPdf(r, rec, scattered)

	sampleColor := c.RayColor(scattered, depth-1, world, lights)
	colorFromScatter := attenuation.Muln(scatteringPdf).Mul(sampleColor).Divn(pdfValue)

	return colorFromEmission.Add(colorFromScatter)
}

func (c Camera) DefocusDiskSample() Point3 {
	// Returns a random point in the camera defocus disk.
	p := RandomInUnitDisk()
	return c.center.Add(c.defocusDiskU.Muln(p[0])).Add(c.defocusDiskV.Muln(p[1]))
}

func (c Camera) SampleSquareStratified(si, sj int) Vec3 {
	// Returns the vector to a random point in the square sub-pixel specified by grid
	// indices s_i and s_j, for an idealized unit square pixel [-.5,-.5] to [+.5,+.5].

	px := (float64(si)+rand.Float64())*c.recipSqrtSpp - 0.5
	py := (float64(sj)+rand.Float64())*c.recipSqrtSpp - 0.5

	return Vec3{px, py, 0}
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
