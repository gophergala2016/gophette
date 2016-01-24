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

	camera := newWindowCamera(window.GetSize())

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

		check(renderer.SetDrawColor(255, 255, 255, 255))
		check(renderer.Clear())
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
	camera   *windowCamera
	texture  *sdl.Texture
}

func (img *textureImage) DrawAt(x, y int) {
	_, _, w, h, err := img.texture.Query()
	check(err)
	dx, dy := img.camera.offset()
	dest := sdl.Rect{int32(x + dx), int32(y + dy), w, h}
	check(img.renderer.Copy(img.texture, nil, &dest))
}

func (img *textureImage) Size() (int, int) {
	_, _, w, h, err := img.texture.Query()
	check(err)
	return int(w), int(h)
}

type sdlAssetLoader struct {
	camera   *windowCamera
	renderer *sdl.Renderer
	images   map[string]*textureImage
}

func newSDLAssetLoader(cam *windowCamera, renderer *sdl.Renderer) *sdlAssetLoader {
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
	camera   *windowCamera
}

func (graphics *sdlGraphics) FillRect(rect Rectangle, r, g, b, a uint8) {
	check(graphics.renderer.SetDrawColor(r, g, b, a))
	rect = rect.MoveBy(graphics.camera.offset())
	sdlRect := sdl.Rect{int32(rect.X), int32(rect.Y), int32(rect.W), int32(rect.H)}
	graphics.renderer.FillRect(&sdlRect)
}
