package main

import (
	"bytes"
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

type ResourceMap map[string][]byte

func main() {
	resources := make(ResourceMap)

	gophette, err := xcf.LoadFromFile("./gophette.xcf")
	check(err)
	for i := range gophette.Layers {
		if gophette.Layers[i].Visible {
			small := scaleImage(gophette.Layers[i])

			buffer := bytes.NewBuffer(nil)
			check(png.Encode(buffer, small))
			resources["gophette_left_"+gophette.Layers[i].Name] = buffer.Bytes()

			buffer = bytes.NewBuffer(nil)
			check(png.Encode(buffer, imaging.FlipH(small)))
			resources["gophette_right_"+gophette.Layers[i].Name] = buffer.Bytes()
		}
	}

	ioutil.WriteFile("../resources.go", toGoFile(resources), 0777)
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

func toGoFile(resources ResourceMap) []byte {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(`package main

// NOTE this file is generated, do not edit it

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
