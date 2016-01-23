package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	check(sdl.Init(sdl.INIT_EVERYTHING))
	defer sdl.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(
		640, 480,
		sdl.WINDOW_RESIZABLE,
	)
	check(err)
	defer renderer.Destroy()
	defer window.Destroy()
	window.SetTitle("Gophette's Adventures")

	running := true
	for running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch event := e.(type) {
			case *sdl.KeyDownEvent:
				if event.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			case *sdl.WindowEvent:
				if event.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					width, height := int(event.Data1), int(event.Data2)
					fmt.Println("size:", width, height)
				}
			case *sdl.QuitEvent:
				running = false
			}
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
