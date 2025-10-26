package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

type RTWImage struct {
	width  int
	height int
	raster [][]RGB
}

func NewRTWImage(filename string) RTWImage {
	img := RTWImage{}
	err := img.Load(filename)
	if err != nil {
		log.Fatalf("%#v", err)
	}
	return img
}

func (im *RTWImage) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()

	im.width, im.height = img.Bounds().Dx(), img.Bounds().Dy()
	im.raster = make([][]RGB, im.width)
	for i := range im.raster {
		im.raster[i] = make([]RGB, im.height)
	}
	gamma := 2.2
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			i := x - bounds.Min.X
			j := y - bounds.Min.Y
			im.raster[i][j] = RGB{math.Pow(float64(r>>8)/255.0, gamma), math.Pow(float64(g>>8)/255.0, gamma), math.Pow(float64(b>>8)/255.0, gamma)}
		}
	}

	return nil
}

func (im RTWImage) Get(x, y int) RGB {
	x = Clamp(x, 0, im.width-1)
	y = Clamp(y, 0, im.height-1)

	return im.raster[x][y]
}

func (im *RTWImage) Set(x, y int, rgb RGB) bool {
	if x < 0 || x > im.width {
		return false
	}
	if y < 0 || y > im.height {
		return false
	}

	im.raster[x][y] = rgb
	return true
}
