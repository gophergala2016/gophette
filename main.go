package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"time"
	"unsafe"
)

func main() {
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

	window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	fullscreen := true

	assetLoader := newSDLAssetLoader(renderer)
	defer assetLoader.close()
	game := NewGame(assetLoader)

	frameTime := time.Second / 60
	lastUpdate := time.Now().Add(-frameTime)

	for game.Running() {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch event := e.(type) {
			case *sdl.KeyDownEvent:
				switch event.Keysym.Sym {
				case sdl.K_LEFT:
					game.HandleInput(InputEvent{GoLeft, true})
				case sdl.K_RIGHT:
					game.HandleInput(InputEvent{GoRight, true})
				case sdl.K_UP:
					game.HandleInput(InputEvent{Jump, true})
				case sdl.K_ESCAPE:
					game.HandleInput(InputEvent{QuitGame, true})
				}
			case *sdl.KeyUpEvent:
				switch event.Keysym.Sym {
				case sdl.K_LEFT:
					game.HandleInput(InputEvent{GoLeft, false})
				case sdl.K_RIGHT:
					game.HandleInput(InputEvent{GoRight, false})
				case sdl.K_UP:
					game.HandleInput(InputEvent{Jump, false})
				case sdl.K_F11:
					if fullscreen {
						window.SetFullscreen(0)
					} else {
						window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
					}
					fullscreen = !fullscreen
				case sdl.K_ESCAPE:
					game.HandleInput(InputEvent{QuitGame, false})
				}
			case *sdl.WindowEvent:
				if event.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					width, height := int(event.Data1), int(event.Data2)
					fmt.Println("size:", width, height)
				}
			case *sdl.QuitEvent:
				game.HandleInput(InputEvent{QuitGame, true})
			}
		}

		now := time.Now()
		dt := now.Sub(lastUpdate)
		if dt > frameTime {
			game.Update()
			lastUpdate = now
		}

		renderer.SetDrawColor(0, 0, 0, 255)
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
	texture  *sdl.Texture
}

func (img *textureImage) DrawAt(x, y int) {
	_, _, w, h, _ := img.texture.Query()
	dest := sdl.Rect{int32(x), int32(y), w, h}
	img.renderer.Copy(img.texture, nil, &dest)
}

type sdlAssetLoader struct {
	renderer *sdl.Renderer
	images   map[string]*textureImage
}

func newSDLAssetLoader(renderer *sdl.Renderer) *sdlAssetLoader {
	return &sdlAssetLoader{
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
	image := &textureImage{l.renderer, texture}
	l.images[id] = image

	return image
}

func (l *sdlAssetLoader) close() {
	for _, image := range l.images {
		image.texture.Destroy()
	}
}
