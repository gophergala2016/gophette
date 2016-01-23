package main

import (
	"github.com/gonutz/xcf"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"os"
)

func main() {
	gophette, err := xcf.LoadFromFile("./gophette.xcf")
	check(err)
	for i := range gophette.Layers {
		if gophette.Layers[i].Visible {
			file, err := os.Create("./gophette_" + gophette.Layers[i].Name + ".png")
			check(err)
			defer file.Close()
			check(png.Encode(file, scaleImage(gophette.Layers[i])))
		}
	}
}

func scaleImage(img image.Image) image.Image {
	const factor = 0.33
	return resize.Resize(
		uint(0.5+factor*float64(img.Bounds().Dx())),
		uint(0.5+factor*float64(img.Bounds().Dy())),
		img,
		resize.Bicubic,
	)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
