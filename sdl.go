package wipeout

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	GlContext sdl.GLContext
	Gamepad   *sdl.GameController
)

type SdlWindow struct {
	window *sdl.Window
}

// FindGamepad returns the first gamepad found
func FindGamepad() *sdl.GameController {
	for i := 0; i < sdl.NumJoysticks(); i++ {
		if sdl.IsGameController(i) {
			return sdl.GameControllerOpen(i)
		}
	}

	return nil
}

// PumpEvents pumps events from SDL
func PumpEvents() {
	var event sdl.Event

	// Keyboards inputs
	for {
		event = sdl.PollEvent()
		if event.GetType() == sdl.KEYDOWN || event.GetType() == sdl.KEYUP {
			code := event.(*sdl.KeyboardEvent).Keysym.Scancode
			var state float32
			if event.GetType() == sdl.KEYDOWN {
				state = 0
			} else {
				state = 1
			}
			if code >= sdl.SCANCODE_LCTRL && code <= sdl.SCANCODE_RALT {
				codeInternal := code - sdl.SCANCODE_LCTRL + sdl.Scancode(INPUT_KEY_LCTRL)
				InputSetButtonState(Button(codeInternal), state)
			} else if code > 0 && code < sdl.Scancode(INPUT_KEY_MAX) {
				InputSetButtonState(Button(code), state)
			}
		}
	}
}

// NewWindow creates a window
func NewWindow(title string, x, y, w, h int32) (*SdlWindow, error) {

	window, err := sdl.CreateWindow("Wipeout",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	if err = gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)

	return &SdlWindow{
		window: window,
	}, nil

}

// ScreenSize returns the size of the screen
func (sw *SdlWindow) ScreenSize() Vec2i {
	w, h := sw.window.GLGetDrawableSize()
	return Vec2i{w, h}
}

// SetFullscreen sets the window to fullscreen mode
func (sw *SdlWindow) SetFullscreen(fullscreen bool) error {
	if fullscreen {

		display, err := sw.window.GetDisplayIndex()
		if err != nil {
			return err
		}
		mode, err := sdl.GetDesktopDisplayMode(display)
		if err != nil {
			return err
		}
		sw.window.SetDisplayMode(&mode)
		sw.window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
		sdl.ShowCursor(sdl.DISABLE)
	} else {
		sw.window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
		sdl.ShowCursor(sdl.ENABLE)
	}

	return nil
}

func (sw *SdlWindow) VideoInit() {
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_ES)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	var err error
	GlContext, err = sw.window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	sdl.GLSetSwapInterval(1)
}

func (sw *SdlWindow) VideoCleanup() {
	sdl.GLDeleteContext(GlContext)
}

func (sw *SdlWindow) EndFrame() {
	sw.window.GLSwap()
}

func (sw *SdlWindow) Destroy() {
	sw.window.Destroy()
}
