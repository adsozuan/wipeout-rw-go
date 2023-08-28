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

	var err error
	platform, err := wipeout.NewPlatformSdl("Wipeout", 0, 0, 800, 600)
	if err != nil {
		panic(err)
	}
	defer platform.Destroy()

	platform.FindGamepad()
	platform.VideoInit()

	running := true
	for running {
		platform.PumpEvents()
		platform.EndFrame()
	}

	platform.VideoCleanup()

}
