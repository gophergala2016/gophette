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
	renderer       *sdl.Renderer
	backColor      = [3]uint8{255, 0, 0}
	cameraX        = 0
	cameraY        = 0
	draggingImage  = false
	draggingObject = false
	images         []image
)

func main() {
	fmt.Print()

	sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1")

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
	window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	if len(LevelImages) == 0 {
		for i, id := range []string{
			"grass left",
			"grass right",
			"grass center 1",
			"grass center 2",
			"grass center 3",
			"small tree",
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

	leftDown := false
	middleDown := false
	rightDown := false
	selectedImage := -1
	selectedObject := -1
	var lastX, lastY int

	moveImage := func(dx, dy int) {
		if selectedImage != -1 {
			images[selectedImage].x += dx
			images[selectedImage].y += dy
		}
	}

	stretchObject := func(dx, dy int) {
		if selectedObject != -1 {
			if sdl.GetKeyboardState()[sdl.SCANCODE_LCTRL] != 0 {
				dx *= 20
				dy *= 20
			}
			obj := &LevelObjects[selectedObject]
			obj.W += dx
			obj.H += dy
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
						draggingObject = false
					} else {
						selectedObject = -1
						selectedImage = -1
						for i := range images {
							if images[i].contains(
								int(event.X)-cameraX,
								int(event.Y)-cameraY,
							) {
								draggingImage = true
								selectedImage = i
							}
						}

						if selectedImage == -1 {
							for i := range LevelObjects {
								if contains(LevelObjects[i],
									int(event.X)-cameraX,
									int(event.Y)-cameraY,
								) {
									draggingObject = true
									selectedObject = i
								}
							}
						}
					}
				}
				if event.Button == sdl.BUTTON_MIDDLE {
					middleDown = event.State == sdl.PRESSED
				}
				if event.Button == sdl.BUTTON_RIGHT {
					rightDown = event.State == sdl.PRESSED
					LevelObjects = append(LevelObjects, LevelObject{
						int(event.X) - cameraX,
						int(event.Y) - cameraY,
						0,
						0,
						true,
					})
					selectedObject = -1
				}
			case *sdl.MouseMotionEvent:
				dx, dy := int(event.X)-lastX, int(event.Y)-lastY
				if selectedImage != -1 && draggingImage {
					img := &images[selectedImage]
					img.x += dx
					img.y += dy
				}
				if selectedObject != -1 && draggingObject {
					obj := &LevelObjects[selectedObject]
					obj.X += dx
					obj.Y += dy
				}
				lastX, lastY = int(event.X), int(event.Y)

				if middleDown {
					cameraX += dx
					cameraY += dy
				}

				if rightDown {
					last := &LevelObjects[len(LevelObjects)-1]
					last.W += dx
					last.H += dy
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
				case sdl.K_j:
					stretchObject(-1, 0)
				case sdl.K_l:
					stretchObject(1, 0)
				case sdl.K_i:
					stretchObject(0, -1)
				case sdl.K_k:
					stretchObject(0, 1)
				case sdl.K_MINUS:
					if selectedImage != -1 && selectedImage != 0 {
						images = append(append(
							[]image{images[selectedImage]},
							images[0:selectedImage]...),
							images[selectedImage+1:]...,
						)
						selectedImage = 0
					}
				case sdl.K_PLUS:
					last := len(images) - 1
					if selectedImage != -1 && selectedImage != last {
						images = append(append(
							images[0:selectedImage],
							images[selectedImage+1:]...),
							images[selectedImage],
						)
						selectedImage = last
					}
				case sdl.K_SPACE:
					if selectedObject != -1 {
						obj := &LevelObjects[selectedObject]
						obj.Solid = !obj.Solid
					}
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
					} else if selectedObject != -1 {
						LevelObjects = append(
							LevelObjects[:selectedObject],
							LevelObjects[selectedObject+1:]...,
						)
						selectedObject = -1
					}
				case sdl.K_F3:
					saveLevel()
				}
			}
		}

		renderer.SetDrawColor(backColor[0], backColor[1], backColor[2], 255)
		renderer.Clear()

		for i, img := range images {
			img.render(i == selectedImage)
		}

		for i, obj := range LevelObjects {
			var g uint8 = 0
			var a uint8 = 100
			if i == selectedObject {
				g = 255
			}
			renderer.SetDrawColor(0, g, 0, a)
			if obj.Solid {
				renderer.SetDrawColor(0, g, 255, a)
			}
			obj.X += cameraX
			obj.Y += cameraY
			r := sdl.Rect{int32(obj.X), int32(obj.Y), int32(obj.W), int32(obj.H)}
			renderer.FillRect(&r)
		}

		renderer.Present()
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
` + imagesToString() + `}
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

func saveObjects() {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(`package main

type LevelObject struct {
	X, Y, W, H int
	Solid      bool
}

var LevelObjects = []LevelObject{
` + objectsToString() + `}
`)

	ioutil.WriteFile("./objects.go", buffer.Bytes(), 0777)
}

func objectsToString() string {
	buffer := bytes.NewBuffer(nil)

	for _, obj := range LevelObjects {
		buffer.WriteString(fmt.Sprintf(`	{%v, %v, %v, %v, %v},
`,
			obj.X, obj.Y, obj.W, obj.H, obj.Solid,
		))
	}

	return string(buffer.Bytes())
}

func contains(obj LevelObject, x, y int) bool {
	return x >= obj.X && y >= obj.Y && x < obj.X+obj.W && y < obj.Y+obj.H
}

func saveLevel() {
	for i := 0; i < len(LevelObjects); i++ {
		if LevelObjects[i].W == 0 || LevelObjects[i].H == 0 {
			LevelObjects = append(LevelObjects[:i], LevelObjects[i+1:]...)
			i--
		}
	}

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(`package main

var level1 = Level{
	[]LevelObject{` + objectsToString() + `
	},
	[]LevelImage{` + imagesToString() + `
	},
}
`)
	ioutil.WriteFile("../level1.go", buffer.Bytes(), 0777)

	saveImages()
	saveObjects()
}
