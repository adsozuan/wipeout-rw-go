package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/adsozuan/wipeout-rw-go/engine"
	wipeout "github.com/adsozuan/wipeout-rw-go/game"
	gl "github.com/chsc/gogl/gl33"

	//gl "github.com/chsc/gogl/gl33"
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
	defer sdl.Quit()

	var err error
	platform, err := engine.NewPlatformSdl(engine.SystemWindowName, 0, 0,
		engine.SystemWindowHeight, engine.SystemWindowWidth)
	if err != nil {
		return err
	}
	defer platform.Destroy()

	platform.FindGamepad()
	err = platform.VideoInit()
	if err != nil {
		return err
	}
	if err = gl.Init(); err != nil {
		panic(err)
	}

	gl.Viewport(0, 0, gl.Sizei(engine.SystemWindowWidth), gl.Sizei(engine.SystemWindowHeight))
	// OPENGL FLAGS
	gl.ClearColor(0.0, 0.1, 0.0, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	system := engine.NewSystem(platform)

	title := wipeout.NewTitle(float32(system.Time()), system.Render)

	for !platform.ExitWanted() {
		platform.PumpEvents()
		platform.PrepareFrame()
		title.Update()
		system.Update()
		platform.EndFrame()
	}

	platform.VideoCleanup()

	return nil
}
