package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_mixer"
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

	game := NewGame(dummyAssetLoader{})

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
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type dummyAssetLoader struct{}

func (dummyAssetLoader) LoadImage(string) Image {
	return nil
}
