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
		smallLeft := scaleImage(barney.GetLayerByName("left_" + layer))
		smallRight := scaleImage(barney.GetLayerByName("right_" + layer))
		resources["barney_left_"+layer] = imageToBytes(smallLeft)
		resources["barney_right_"+layer] = imageToBytes(smallRight)
	}

	grass, err := xcf.LoadFromFile("./grass.xcf")
	check(err)
	for _, layer := range []string{
		"grass left",
		"grass right",
		"grass center 1",
		"grass center 2",
		"grass center 3",
	} {
		resources[layer] = imageToBytes(grass.GetLayerByName(layer))
	}

	grassLong, err := xcf.LoadFromFile("./grass_long.xcf")
	check(err)
	for _, layer := range []string{
		"grass long 1",
		"grass long 2",
		"grass long 3",
	} {
		resources[layer] = imageToBytes(grassLong.GetLayerByName(layer))
	}

	ground, err := xcf.LoadFromFile("./ground.xcf")
	check(err)
	for _, layer := range []string{
		"ground left",
		"ground right",
		"ground center 1",
		"ground center 2",
		"ground center 3",
	} {
		resources[layer] = imageToBytes(ground.GetLayerByName(layer))
	}

	groundLong, err := xcf.LoadFromFile("./ground_long.xcf")
	check(err)
	for _, layer := range []string{
		"ground long 1",
		"ground long 2",
	} {
		resources[layer] = imageToBytes(groundLong.GetLayerByName(layer))
	}

	rock, err := xcf.LoadFromFile("./rock.xcf")
	check(err)
	resources["square rock"] = imageToBytes(scaleImage(rock.GetLayerByName("rock")))

	tree, err := xcf.LoadFromFile("./tree.xcf")
	check(err)
	smallTree := scaleImage(tree.GetLayerByName("small"))
	resources["small tree"] = imageToBytes(smallTree)

	tree, err = xcf.LoadFromFile("./tree_big.xcf")
	check(err)
	bigTree := scaleImage(tree.GetLayerByName("big"))
	resources["big tree"] = imageToBytes(bigTree)

	tree, err = xcf.LoadFromFile("./tree_huge.xcf")
	check(err)
	hugeTree := scaleImage(tree.GetLayerByName("huge"))
	resources["huge tree"] = imageToBytes(hugeTree)

	content := toGoFile(resources, string(constants.Bytes()))
	ioutil.WriteFile("../resource/resources.go", content, 0777)
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
	buffer.WriteString(`package resource

// NOTE this file is generated, do not edit it

type Rectangle struct{ X, Y, W, H int }

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
