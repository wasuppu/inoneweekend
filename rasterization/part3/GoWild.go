package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const (
	WIDTH  = 1280
	HEIGHT = 720
)

type Texture struct {
	Data        []color.RGBA
	WIDTH       int
	Height      int
	NumChannels int
}

type VertexInput struct {
	Pos       Vec3f
	Normal    Vec3f
	TexCoords Vec2f
}

type FragmentInput struct {
	Normal    Vec3f
	TexCoords Vec2f
}

type Mesh struct {
	Offset      int
	Count       int
	DiffuseName string
}

func toRaster(v Vec4f) Vec4f {
	return Vec4f{WIDTH * (v.X() + v.W()) / 2, HEIGHT * (v.W() - v.Y()) / 2, v.Z(), v.W()}
}

func VS(input VertexInput, MVP Mat4) (Vec4f, FragmentInput) {
	var output FragmentInput
	output.Normal = input.Normal
	output.TexCoords = input.TexCoords
	return input.Pos.V4(1).Mulm(MVP), output
}

func FS(input FragmentInput, pTexture Texture) color.RGBA {
	idxS := int(input.TexCoords.S() - float64(int(input.TexCoords.S())*pTexture.WIDTH) - 0.5)
	idxT := int(input.TexCoords.T() - float64(int(input.TexCoords.T())*pTexture.Height) - 0.5)
	idx := (idxT*pTexture.WIDTH + idxS) * pTexture.NumChannels

	return ScaleColorRGB(pTexture.Data[idx], (1.0 / 255))
}

func EvaluateEdgeFunction(e Vec3f, sample Vec2f) bool {
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

func initializeSceneObjects(filename string, meshBuffer []Mesh, vertexBuffer []VertexInput, indexBuffer []int, textures map[string]Texture) {

}

func DrawIndexed(frameBuffer []color.Color, depthBuffer []float64, vertexBuffer []VertexInput, indexBuffer []int, mesh Mesh, MVP Mat4, pTexture Texture) {
	triCount := mesh.Count / 3

	for idx := range triCount {
		vi0 := vertexBuffer[indexBuffer[mesh.Offset+(idx*3)]]
		vi1 := vertexBuffer[indexBuffer[mesh.Offset+(idx*3+1)]]
		vi2 := vertexBuffer[indexBuffer[mesh.Offset+(idx*3+2)]]

		v0Clip, fi0 := VS(vi0, MVP)
		v1Clip, fi1 := VS(vi1, MVP)
		v2Clip, fi2 := VS(vi2, MVP)

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

		e0 := m[0].Divn(math.Abs(m[0].X()) + math.Abs(m[0].Y()))
		e1 := m[1].Divn(math.Abs(m[1].X()) + math.Abs(m[1].Y()))
		e2 := m[2].Divn(math.Abs(m[2].X()) + math.Abs(m[2].Y()))

		C := Vec3f{1, 1, 1}.Mulm(m)

		Z := Vec3f{v0Clip.Z(), v1Clip.Z(), v2Clip.Z()}.Mulm(m)

		pnx := Vec3f{fi0.Normal.X(), fi1.Normal.X(), fi2.Normal.X()}.Mulm(m)
		pny := Vec3f{fi0.Normal.Y(), fi1.Normal.Y(), fi2.Normal.Y()}.Mulm(m)
		pnz := Vec3f{fi0.Normal.Z(), fi1.Normal.Z(), fi2.Normal.Z()}.Mulm(m)

		puvs := Vec3f{fi0.TexCoords.S(), fi1.TexCoords.S(), fi2.TexCoords.S()}
		puvt := Vec3f{fi0.TexCoords.T(), fi1.TexCoords.T(), fi2.TexCoords.T()}

		for y := range HEIGHT {
			for x := range WIDTH {
				sample := Vec2f{float64(x) + 0.5, float64(y) + 0.5}

				inside0 := EvaluateEdgeFunction(e0, sample)
				inside1 := EvaluateEdgeFunction(e1, sample)
				inside2 := EvaluateEdgeFunction(e2, sample)

				if inside0 && inside1 && inside2 {
					oneOverW := C.X()*sample.X() + C.Y()*sample.Y() + C.Z()

					w := 1 / oneOverW

					zOverW := Z.X()*sample.X() + Z.Y()*sample.Y() + Z.Z()
					z := zOverW * w

					if z <= depthBuffer[x+y*WIDTH] {
						depthBuffer[x+y*WIDTH] = z

						nxOverW := pnx.X()*sample.X() + pnx.Y()*sample.Y() + pnx.Z()
						nyOverW := pny.X()*sample.X() + pny.Y()*sample.Y() + pny.Z()
						nzOverW := pnz.X()*sample.X() + pnz.Y()*sample.Y() + pnz.Z()

						uOverW := puvs.X()*sample.X() + puvs.Y()*sample.Y() + puvs.Z()
						vOverW := puvt.X()*sample.X() + puvt.Y()*sample.Y() + puvt.Z()

						normal := Vec3f{nxOverW, nyOverW, nzOverW}.Muln(w)
						texCoords := Vec2f{uOverW, vOverW}.Muln(w)

						fsInput := FragmentInput{normal, texCoords}

						outputColor := FS(fsInput, pTexture)

						frameBuffer[x+y*WIDTH] = outputColor
					}
				}
			}
		}
	}
}

func goWild() {
	frameBuffer := make([]color.Color, WIDTH*HEIGHT)
	for i := range frameBuffer {
		frameBuffer[i] = color.Black
	}
	depthBuffer := make([]float64, WIDTH*HEIGHT)
	for i := range depthBuffer {
		depthBuffer[i] = math.MaxFloat64
	}

	vertexBuffer := []VertexInput{}
	indexBuffer := []int{}
	primitives := []Mesh{}
	textures := make(map[string]Texture)

	filename := "./assets/sponza.obj"

	initializeSceneObjects(filename, primitives, vertexBuffer, indexBuffer, textures)

	nearPlane := 0.125
	farPlane := 5000.0

	eye := Vec3f{0, -8.5, -5}
	lookat := Vec3f{20, 5, 1}
	up := Vec3f{0, 1, 0}

	view := LookAtMatrix(eye, lookat, up)
	proj := Perspective(Radians(-30), float64(WIDTH)/float64(HEIGHT), nearPlane, farPlane)

	mvp := view.Mul(proj)

	for i := range primitives {
		DrawIndexed(frameBuffer, depthBuffer, vertexBuffer, indexBuffer, primitives[i], mvp, textures[primitives[i].DiffuseName])
	}
	writePng("goWild", frameBuffer, WIDTH, HEIGHT)
}

func main() {
	goWild()
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

func ScaleColorRGB(c color.RGBA, scale float64) color.RGBA {
	if scale < 0 {
		scale = 0
	} else if scale > 255 {
		scale = 255
	}

	r := uint8(float64(c.R) * scale)
	g := uint8(float64(c.G) * scale)
	b := uint8(float64(c.B) * scale)
	a := c.A

	return color.RGBA{R: r, G: g, B: b, A: a}
}
