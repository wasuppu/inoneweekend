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

func toRaster(v Vec4f) Vec4f {
	return Vec4f{WIDTH * (v.X() + v.W()) / 2, HEIGHT * (v.W() - v.Y()) / 2, v.Z(), v.W()}
}

func initializeSceneObjects() []Mat4 {
	identity := ID4()
	M0 := identity.Translate(Vec3f{0, 0, 2})
	M0 = M0.Rotate(Radians(45), Vec3f{0, 1, 0})
	M1 := identity.Translate(Vec3f{-3.75, 0, 0})
	M1 = M1.Rotate(Radians(30), Vec3f{1, 0, 0})
	M2 := identity.Translate(Vec3f{3.75, 0, 0})
	M2 = M2.Rotate(Radians(60), Vec3f{0, 1, 0})
	M3 := identity.Translate(Vec3f{0, 0, -2})
	M3 = M3.Rotate(Radians(90), Vec3f{0, 0, 1})
	return []Mat4{M0, M1, M2, M3}
}

func VS(pos Vec3f, M, V, P Mat4) Vec4f {
	return pos.V4(1).Mulm(M).Mulm(V).Mulm(P)
}

func evaluateEdgeFunction(e Vec3f, sample Vec3f) bool {
	result := e.X()*sample.X() + e.Y()*sample.Y() + e.Z()

	if result > 0 {
		return true
	} else if result < 0 {
		return false
	}

	if e.X() > 0 {
		return true
	} else if e.X() < 0 {
		return false
	}

	if e.X() == 0 && e.Y() < 0 {
		return false
	} else {
		return true
	}
}

func Go3d() {
	vertices := []Vec3f{
		{1.0, -1.0, -1.0},
		{1.0, -1.0, 1.0},
		{-1.0, -1.0, 1.0},
		{-1.0, -1.0, -1.0},
		{1.0, 1.0, -1.0},
		{1.0, 1.0, 1.0},
		{-1.0, 1.0, 1.0},
		{-1.0, 1.0, -1.0},
	}

	indices := []int{1, 3, 0, 7, 5, 4, 4, 1, 0, 5, 2, 1, 2, 7, 3, 0, 7, 4, 1, 2, 3, 7, 6, 5, 4, 5, 1, 5, 6, 2, 2, 6, 7, 0, 3, 7}

	colors := []Vec3f{
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 1},
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}

	frameBuffer := make([]color.Color, WIDTH*HEIGHT)
	for y := range HEIGHT {
		for x := range WIDTH {
			frameBuffer[x+y*WIDTH] = color.Black
		}
	}
	depthBuffer := make([]float64, WIDTH*HEIGHT)

	objects := initializeSceneObjects()

	nearPlane := 0.1
	farPlane := 100.0

	eye := Vec3f{0, 3.75, 6.5}
	lookat := Vec3f{0, 0, 0}
	up := Vec3f{0, 1, 0}

	view := LookAtMatrix(eye, lookat, up)
	proj := Perspective(Radians(60), float64(WIDTH)/float64(HEIGHT), nearPlane, farPlane)

	for n := range objects {
		for idx := range len(indices) / 3 {
			v0 := vertices[indices[idx*3]]
			v1 := vertices[indices[idx*3+1]]
			v2 := vertices[indices[idx*3+2]]

			v0Clip := VS(v0, objects[n], view, proj)
			v1Clip := VS(v1, objects[n], view, proj)
			v2Clip := VS(v2, objects[n], view, proj)

			v0Homogen := toRaster(v0Clip)
			v1Homogen := toRaster(v1Clip)
			v2Homogen := toRaster(v2Clip)

			m := Mat3{
				{v0Homogen.X(), v1Homogen.X(), v2Homogen.X()},
				{v0Homogen.Y(), v1Homogen.Y(), v2Homogen.Y()},
				{v0Homogen.W(), v1Homogen.W(), v2Homogen.W()},
			}

			det := m.Determinant()
			if det > 0 {
				continue
			}

			m = m.Inverse()

			e0 := Vec3f(m[0])
			e1 := Vec3f(m[1])
			e2 := Vec3f(m[2])

			c := Vec3f{1, 1, 1}.Mulm(m)

			for y := range HEIGHT {
				for x := range WIDTH {
					sample := Vec3f{float64(x) + 0.5, float64(y) + 0.5, 1}

					inside0 := evaluateEdgeFunction(e0, sample)
					inside1 := evaluateEdgeFunction(e1, sample)
					inside2 := evaluateEdgeFunction(e2, sample)

					if inside0 && inside1 && inside2 {
						oneOverW := c.X()*sample.X() + c.Y()*sample.Y() + c.Z()

						if oneOverW >= depthBuffer[x+y*WIDTH] {
							depthBuffer[x+y*WIDTH] = oneOverW
							frameBuffer[x+y*WIDTH] = toColor(colors[indices[3*idx]%6])
						}
					}
				}
			}
		}
	}

	writePng("Go3d", frameBuffer, WIDTH, HEIGHT)
}

func toColor(v Vec3f) color.RGBA {
	v = Vec3f{Clamp(v.X(), 0, 1), Clamp(v.Y(), 0, 1), Clamp(v.Z(), 0, 1)}
	return color.RGBA{uint8(255.999 * v[0]), uint8(255.999 * v[1]), uint8(255.999 * v[2]), 0xff}
}

func writePng(name string, pixels []color.Color, WIDTH, HEIGHT int) {
	f, _ := os.Create(name + ".png")
	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	for j := range HEIGHT {
		for i := range WIDTH {
			img.Set(i, j, pixels[i+j*WIDTH])
		}
	}
	png.Encode(f, img)
}

func main() {
	Go3d()
}
