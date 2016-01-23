package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gonutz/xcf"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

const scale = 0.33

type ResourceMap map[string][]byte

func main() {
	resources := make(ResourceMap)
	constants := bytes.NewBuffer(nil)

	gophette, err := xcf.LoadFromFile("./gophette.xcf")
	check(err)
	barney, err := xcf.LoadFromFile("./barney.xcf")
	check(err)

	// create the collision information for Gophette and Barney
	addCollisionInfo := func(canvas xcf.Canvas, variable string) {
		collision := canvas.GetLayerByName("collision")
		left, top := findTopLeftNonTransparentPixel(collision)
		right, bottom := findBottomRightNonTransparentPixel(collision)
		// scale the collision rect just like the images
		left = int(0.5 + scale*float64(left))
		top = int(0.5 + scale*float64(top))
		right = int(0.5 + scale*float64(right))
		bottom = int(0.5 + scale*float64(bottom))
		width, height := right-left+1, bottom-top+1
		line := fmt.Sprintf(
			"var %v = Rectangle{%v, %v, %v, %v}\n",
			variable,
			left, top, width, height,
		)
		constants.WriteString(line)
	}
	addCollisionInfo(gophette, "HeroCollisionRect")
	addCollisionInfo(barney, "BarneyCollisionRect")

	// create the image resources
	for _, layer := range []string{
		"jump",
		"run1",
		"run2",
		"run3",
	} {
		small := scaleImage(gophette.GetLayerByName(layer))
		resources["gophette_left_"+layer] = imageToBytes(small)
		resources["gophette_right_"+layer] = imageToBytes(imaging.FlipH(small))
	}

	for _, layer := range []string{
		"stand",
		"jump",
		"run1",
		"run2",
		"run3",
		"run4",
		"run5",
		"run6",
	} {
		small := scaleImage(barney.GetLayerByName(layer))
		resources["barney_left_"+layer] = imageToBytes(small)
		resources["barney_right_"+layer] = imageToBytes(imaging.FlipH(small))
	}

	content := toGoFile(resources, string(constants.Bytes()))
	ioutil.WriteFile("../resources.go", content, 0777)
}

func imageToBytes(img image.Image) []byte {
	buffer := bytes.NewBuffer(nil)
	check(png.Encode(buffer, img))
	return buffer.Bytes()
}

func findTopLeftNonTransparentPixel(img image.Image) (x, y int) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				return x, y
			}
		}
	}
	return -1, -1
}

func findBottomRightNonTransparentPixel(img image.Image) (x, y int) {
	for y := img.Bounds().Max.Y - 1; y >= img.Bounds().Min.Y; y-- {
		for x := img.Bounds().Max.X - 1; x >= img.Bounds().Min.X; x-- {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				return x, y
			}
		}
	}
	return -1, -1
}

func scaleImage(img image.Image) image.Image {
	return resize.Resize(
		uint(0.5+scale*float64(img.Bounds().Dx())),
		uint(0.5+scale*float64(img.Bounds().Dy())),
		img,
		resize.Bicubic,
	)
}

func toGoFile(resources ResourceMap, constants string) []byte {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(`package main

// NOTE this file is generated, do not edit it

` + constants + `
var Resources = map[string][]byte{`)

	var table sortableResourceEntries
	for id, data := range resources {
		table = append(table, resourceEntry{id, data})
	}
	sort.Sort(table)

	for _, entry := range table {
		buffer.WriteString(`
	"` + entry.id + `": []byte{`)
		for i, b := range entry.data {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(strconv.Itoa(int(b)))
		}
		buffer.WriteString("},")
	}

	buffer.WriteString("\n}\n")
	return buffer.Bytes()
}

type resourceEntry struct {
	id   string
	data []byte
}

type sortableResourceEntries []resourceEntry

func (e sortableResourceEntries) Len() int {
	return len(e)
}

func (e sortableResourceEntries) Less(i, j int) bool {
	return strings.Compare(e[i].id, e[j].id) < 0
}

func (e sortableResourceEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
