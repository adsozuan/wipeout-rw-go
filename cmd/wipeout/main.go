package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/adsozuan/wipeout-rw-go/engine"
	"github.com/adsozuan/wipeout-rw-go/system"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		return err
	}
	// if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
	// 	return err
	// }
	defer sdl.Quit()

	var err error
	platform, err := engine.NewPlatformSdl(system.WindowName, 0, 0,
		system.WindowWidth, system.WindowHeight)
	if err != nil {
		return err
	}
	defer platform.Destroy()

	platform.FindGamepad()
	err = platform.VideoInit()
	if err != nil {
		return err
	}

	system, err := system.New(platform)
	if err != nil {
		return err
	}


	for !platform.ExitWanted() {
		err := platform.PumpEvents()
		if err != nil {
			return err
		}
		err = platform.PrepareFrame()
		if err != nil {
			return err
		}
		system.Update()
		err = platform.EndFrame()
		if err != nil {
			return err
		}
	}
	system.Cleanup()

	platform.VideoCleanup()

	return nil
}
