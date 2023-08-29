package main

import (
	"fmt"
	"os"

	"github.com/adsozuan/wipeout-rw-go"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		return err
	}
	defer sdl.Quit()

	var err error
	platform, err := wipeout.NewPlatformSdl("Wipeout", 0, 0, 800, 600)
	if err != nil {
		return err
	}
	defer platform.Destroy()

	platform.FindGamepad()
	err = platform.VideoInit()
	if err != nil {
		return err
	}

	system := wipeout.NewSystem(platform)

	for !platform.ExitWanted() {
		platform.PumpEvents()
		platform.PrepareFrame()
		system.Update()
		platform.EndFrame()
	}

	platform.VideoCleanup()

	return nil
}
