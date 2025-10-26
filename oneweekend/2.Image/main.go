package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

const (
	imageWidth  = 1024
	imageHeight = 768
)

func main() {
	framebuffer := make([]color.Color, imageWidth*imageHeight)
	render(framebuffer, imageWidth, imageHeight)
	WritePng("out", framebuffer, imageWidth, imageHeight)
}

func render(framebuffer []color.Color, imageWidth, imageHeight int) {
	for j := range imageHeight {
		log.Printf("\rScanlines remaining: %d", imageHeight-j)
		for i := range imageWidth {
			framebuffer[i+j*imageWidth] = RGB{float64(i) / float64(imageWidth-1), float64(j) / float64(imageHeight-1), 0}.Color()
		}
	}
	log.Println("Done.")
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
