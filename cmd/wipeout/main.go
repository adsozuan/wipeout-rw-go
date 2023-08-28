package main

import (
	"github.com/adsozuan/wipeout-rw-go"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	wipeout.Gamepad = wipeout.FindGamepad()

	var err error
	window, err := wipeout.NewWindow("Wipeout", 0, 0, 800, 600)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.VideoInit()

	running := true
	for running {
		wipeout.PumpEvents()
		window.EndFrame()
	}

	window.VideoCleanup()

}
