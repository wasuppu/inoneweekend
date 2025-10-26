package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	WIDTH  = 1280
	HEIGHT = 720
)

func toRaster(v Vec3f) Vec3f {
	return Vec3f{WIDTH * (v.X() + 1) / 2, HEIGHT * (v.Y() + 1) / 2, 1}
}

func HelloTriangle() {
	v0 := Vec3f{-0.5, 0.5, 1.0}
	v1 := Vec3f{0.5, 0.5, 1.0}
	v2 := Vec3f{0.0, -0.5, 1.0}

	v0 = toRaster(v0)
	v1 = toRaster(v1)
	v2 = toRaster(v2)

	c0 := Vec3f{1, 0, 0}
	c1 := Vec3f{0, 1, 0}
	c2 := Vec3f{0, 0, 1}

	m := Mat3{
		{v0.X(), v1.X(), v2.X()},
		{v0.Y(), v1.Y(), v2.Y()},
		{v0.Z(), v1.Z(), v2.Z()},
	}

	m = m.Inverse()

	e0 := Vec3f{1, 0, 0}.Mulm(m)
	e1 := Vec3f{0, 1, 0}.Mulm(m)
	e2 := Vec3f{0, 0, 1}.Mulm(m)

	frameBuffer := make([]color.Color, WIDTH*HEIGHT)
	for y := range HEIGHT {
		for x := range WIDTH {
			frameBuffer[x+y*WIDTH] = color.Black
		}
	}

	for y := range HEIGHT {
		for x := range WIDTH {
			sample := Vec3f{float64(x) + 0.5, float64(y) + 0.5, 1}
			alpha := e0.Dot(sample)
			beta := e1.Dot(sample)
			gamma := e2.Dot(sample)

			if alpha >= 0 && beta >= 0 && gamma >= 0 {
				frameBuffer[x+y*WIDTH] = toColor(c0.Muln(alpha).Add(c1.Muln(beta)).Add(c2.Muln(gamma)))
			}

		}
	}

	writePng("traignale", frameBuffer, WIDTH, HEIGHT)
}

func toColor(v Vec3f) color.RGBA {
	v = Vec3f{Clamp(v.X(), 0, 1), Clamp(v.Y(), 0, 1), Clamp(v.Z(), 0, 1)}
	return color.RGBA{uint8(255.999 * v[0]), uint8(255.999 * v[1]), uint8(255.999 * v[2]), 0xff}
}

func writePng(name string, pixels []color.Color, width, height int) {
	f, _ := os.Create(name + ".png")
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for j := range height {
		for i := range width {
			img.Set(i, j, pixels[i+j*width])
		}
	}
	png.Encode(f, img)
}

func main() {
	HelloTriangle()
}
