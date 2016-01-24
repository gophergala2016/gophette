package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"time"
	"unsafe"
)

func main() {
	sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1")

	check(sdl.Init(sdl.INIT_EVERYTHING))
	defer sdl.Quit()

	check(mix.Init(mix.INIT_OGG))
	defer mix.Quit()
	check(mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 1, 512))
	defer mix.CloseAudio()

	if img.Init(img.INIT_PNG)&img.INIT_PNG == 0 {
		panic("error init png")
	}
	defer img.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(
		640, 480,
		sdl.WINDOW_RESIZABLE,
	)
	check(err)
	defer renderer.Destroy()
	defer window.Destroy()
	window.SetTitle("Gophette's Adventures")
	sdl.ShowCursor(0)

	window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	fullscreen := true

	camera := newSDLCamera(window.GetSize())

	assetLoader := newSDLAssetLoader(camera, renderer)
	defer assetLoader.close()

	// charIndex selects which character is being controlled by the user, for
	// the final game this must be 0 but for creating the "AI" for Barney, set
	// this to 1 and delete the recorded inputs so they are not applied
	// additionally to the user controls

	// NOTE either this
	const charIndex = 0
	// NOTE or these
	//const charIndex = 1
	//recordedInputs = recordedInputs[:0]

	game := NewGame(
		assetLoader,
		&sdlGraphics{renderer, camera},
		camera,
		charIndex,
	)

	frameTime := time.Second / 60
	lastUpdate := time.Now().Add(-frameTime)

	for game.Running() {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch event := e.(type) {
			case *sdl.KeyDownEvent:
				if event.Repeat == 0 {
					switch event.Keysym.Sym {
					case sdl.K_LEFT:
						game.HandleInput(InputEvent{GoLeft, true, charIndex})
					case sdl.K_RIGHT:
						game.HandleInput(InputEvent{GoRight, true, charIndex})
					case sdl.K_UP:
						game.HandleInput(InputEvent{Jump, true, charIndex})
					case sdl.K_ESCAPE:
						game.HandleInput(InputEvent{QuitGame, true, charIndex})
					}
				}
			case *sdl.KeyUpEvent:
				switch event.Keysym.Sym {
				case sdl.K_LEFT:
					game.HandleInput(InputEvent{GoLeft, false, charIndex})
				case sdl.K_RIGHT:
					game.HandleInput(InputEvent{GoRight, false, charIndex})
				case sdl.K_UP:
					game.HandleInput(InputEvent{Jump, false, charIndex})
				case sdl.K_F11:
					if fullscreen {
						window.SetFullscreen(0)
					} else {
						window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
					}
					fullscreen = !fullscreen
				case sdl.K_ESCAPE:
					game.HandleInput(InputEvent{QuitGame, false, charIndex})
				}
			case *sdl.WindowEvent:
				if event.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					width, height := int(event.Data1), int(event.Data2)
					camera.setWindowSize(width, height)
				}
			case *sdl.QuitEvent:
				game.HandleInput(InputEvent{QuitGame, true, charIndex})
			}
		}

		now := time.Now()
		dt := now.Sub(lastUpdate)
		// TODO make sure the animations are not all jittery
		_ = dt
		//if dt > frameTime {
		game.Update()
		//lastUpdate = now
		//}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		game.Render()
		renderer.Present()
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type textureImage struct {
	renderer *sdl.Renderer
	camera   *sdlCamera
	texture  *sdl.Texture
}

func (img *textureImage) DrawAt(x, y int) {
	_, _, w, h, _ := img.texture.Query()
	dx, dy := img.camera.offset()
	dest := sdl.Rect{int32(x + dx), int32(y + dy), w, h}
	img.renderer.Copy(img.texture, nil, &dest)
}

func (img *textureImage) Size() (int, int) {
	_, _, w, h, _ := img.texture.Query()
	return int(w), int(h)
}

type sdlAssetLoader struct {
	camera   *sdlCamera
	renderer *sdl.Renderer
	images   map[string]*textureImage
}

func newSDLAssetLoader(cam *sdlCamera, renderer *sdl.Renderer) *sdlAssetLoader {
	return &sdlAssetLoader{
		camera:   cam,
		renderer: renderer,
		images:   make(map[string]*textureImage),
	}
}

func (l *sdlAssetLoader) LoadImage(id string) Image {
	if img, ok := l.images[id]; ok {
		return img
	}
	data := Resources[id]
	if data == nil {
		panic("unknown resource: " + id)
	}

	rwOps := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
	surface, err := img.Load_RW(rwOps, false)
	check(err)
	defer surface.Free()
	texture, err := l.renderer.CreateTextureFromSurface(surface)
	check(err)
	image := &textureImage{l.renderer, l.camera, texture}
	l.images[id] = image

	return image
}

func (l *sdlAssetLoader) close() {
	for _, image := range l.images {
		image.texture.Destroy()
	}
}

type sdlGraphics struct {
	renderer *sdl.Renderer
	camera   *sdlCamera
}

func (graphics *sdlGraphics) FillRect(rect Rectangle, r, g, b, a uint8) {
	graphics.renderer.SetDrawColor(r, g, b, a)
	rect = rect.MoveBy(graphics.camera.offset())
	sdlRect := sdl.Rect{int32(rect.X), int32(rect.Y), int32(rect.W), int32(rect.H)}
	graphics.renderer.FillRect(&sdlRect)
}

type sdlCamera struct {
	position Rectangle
	bounds   Rectangle
}

func newSDLCamera(windowW, windowH int) *sdlCamera {
	cam := &sdlCamera{
		// initially set no bounds (big integers)
		bounds: Rectangle{-999999, -999999, 2 * 999999, 2 * 999999},
	}
	cam.setWindowSize(windowW, windowH)
	return cam
}

func (cam *sdlCamera) setWindowSize(w, h int) {
	cx, cy := cam.position.Center()
	cam.position.W, cam.position.H = w, h
	cam.CenterAround(cx, cy)
}

func (cam *sdlCamera) CenterAround(x, y int) {
	cam.position.X = x - cam.position.W/2
	cam.position.Y = y - cam.position.H/2

	// keep the camera in the bounds
	if cam.position.X < cam.bounds.X {
		cam.position.X = cam.bounds.X
	}
	if cam.position.Y < cam.bounds.Y {
		cam.position.Y = cam.bounds.Y
	}
	if cam.position.X+cam.position.W > cam.bounds.X+cam.bounds.W {
		cam.position.X = cam.bounds.X + cam.bounds.W - cam.position.W
	}
	if cam.position.Y+cam.position.H > cam.bounds.Y+cam.bounds.H {
		cam.position.Y = cam.bounds.Y + cam.bounds.H - cam.position.H
	}
}

func (cam *sdlCamera) SetBounds(bounds Rectangle) {
	cam.bounds = bounds
}

func (cam *sdlCamera) offset() (dx, dy int) {
	return -cam.position.X, -cam.position.Y
}
