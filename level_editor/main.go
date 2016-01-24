package main

import (
	"bytes"
	"fmt"
	"github.com/gophergala2016/gophette/resource"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"io/ioutil"
	"unsafe"
)

var (
	renderer      *sdl.Renderer
	backColor     = [3]uint8{255, 0, 0}
	cameraX       = 0
	cameraY       = 0
	draggingImage = false
	images        []image
)

func main() {

	fmt.Print()
	check(sdl.Init(sdl.INIT_EVERYTHING))
	defer sdl.Quit()

	window, r, err := sdl.CreateWindowAndRenderer(
		640, 480,
		sdl.WINDOW_RESIZABLE,
	)
	check(err)
	renderer = r
	defer renderer.Destroy()
	defer window.Destroy()
	window.SetTitle("Gophette's Adventures - Level Editor")
	window.SetPosition(50, 50)
	window.SetSize(1800, 900)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	if len(LevelImages) == 0 {
		for i, id := range []string{
			"grass left",
			"grass right",
			"grass center 1",
			"grass center 2",
			"grass center 3",
		} {
			images = append(images, image{id, loadImage(id), 0, i * 50})
		}
	} else {
		for i := range LevelImages {
			id := LevelImages[i].ID
			x, y := LevelImages[i].X, LevelImages[i].Y
			img := loadImage(id)
			images = append(images, image{id, img, x, y})
		}
	}

	var objects []sdl.Rect

	leftDown := false
	middleDown := false
	selectedImage := -1
	var lastX, lastY int

	moveImage := func(dx, dy int) {
		if selectedImage != -1 {
			images[selectedImage].x += dx
			images[selectedImage].y += dy
		}
	}

	running := true
	for running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch event := e.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if event.Button == sdl.BUTTON_LEFT {
					leftDown = event.State == sdl.PRESSED
					if !leftDown {
						draggingImage = false
					} else {
						for i := range images {
							if images[i].contains(
								int(event.X)-cameraX,
								int(event.Y)-cameraY,
							) {
								draggingImage = true
								selectedImage = i
							}
						}
					}
				}
				if event.Button == sdl.BUTTON_MIDDLE {
					middleDown = event.State == sdl.PRESSED
				}
			case *sdl.MouseMotionEvent:
				dx, dy := int(event.X)-lastX, int(event.Y)-lastY
				if selectedImage != -1 && draggingImage {
					img := &images[selectedImage]
					img.x += dx
					img.y += dy
				}
				lastX, lastY = int(event.X), int(event.Y)

				if middleDown {
					cameraX += dx
					cameraY += dy
				}
			case *sdl.KeyDownEvent:
				switch event.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_LEFT:
					cameraX += 100
				case sdl.K_RIGHT:
					cameraX -= 100
				case sdl.K_UP:
					cameraY += 100
				case sdl.K_DOWN:
					cameraY -= 100
				case sdl.K_a:
					moveImage(-1, 0)
				case sdl.K_d:
					moveImage(1, 0)
				case sdl.K_w:
					moveImage(0, -1)
				case sdl.K_s:
					moveImage(0, 1)
				case sdl.K_c:
					if selectedImage != -1 {
						copy := images[selectedImage]
						copy.x += 10
						copy.y += 10
						images = append(images, copy)
					}
				case sdl.K_DELETE:
					if selectedImage != -1 {
						images = append(images[:selectedImage], images[selectedImage+1:]...)
						selectedImage = -1
					}
				case sdl.K_F3:
					saveImages()
				}
			}
		}

		renderer.SetDrawColor(backColor[0], backColor[1], backColor[2], 255)
		renderer.Clear()

		for i, img := range images {
			img.render(i == selectedImage)
		}

		renderer.SetDrawColor(0, 0, 0, 64)
		for _, obj := range objects {
			obj.X += int32(cameraX)
			obj.Y += int32(cameraY)
			renderer.FillRect(&obj)
		}

		renderer.Present()

		sdl.Delay(10)
	}
}

type image struct {
	id      string
	texture *sdl.Texture
	x, y    int
}

func (img image) contains(x, y int) bool {
	_, _, w, h, _ := img.texture.Query()
	return x >= img.x && y >= img.y && x < img.x+int(w) && y < img.y+int(h)
}

func (img image) render(isSelected bool) {
	_, _, w, h, _ := img.texture.Query()
	x, y := img.x, img.y
	x += cameraX
	y += cameraY
	dest := &sdl.Rect{int32(x), int32(y), w, h}
	renderer.Copy(img.texture, nil, dest)

	if isSelected {
		renderer.SetDrawColor(0, 255, 0, 64)
		renderer.FillRect(dest)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadImage(id string) *sdl.Texture {
	data := resource.Resources[id]
	rwOps := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
	surface, err := img.Load_RW(rwOps, false)
	check(err)
	defer surface.Free()
	texture, err := renderer.CreateTextureFromSurface(surface)
	check(err)
	return texture
}

func saveImages() {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(`package main

type LevelImage struct {
	ID   string
	X, Y int
}

var LevelImages = []LevelImage{
` + imagesToString() + `
}
`)

	ioutil.WriteFile("./images.go", buffer.Bytes(), 0777)
}

func imagesToString() string {
	buffer := bytes.NewBuffer(nil)

	for _, img := range images {
		buffer.WriteString(fmt.Sprintf(`	{"%v", %v, %v},
`, img.id, img.x, img.y))
	}

	return string(buffer.Bytes())
}
