package engine

import (
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var Logger *log.Logger

const (
	PLATFORM_WINDOW_FLAGS = sdl.WINDOW_OPENGL
)

var (
	ScreenBuffer       *sdl.Texture
	ScreenBufferPixels *sdl.Rect
	ScreenBufferPitch  int32
	ScreenBufferSize   Vec2i = Vec2i{0, 0}
	ScreenSize         Vec2i = Vec2i{0, 0}
)

var GamepadMap = map[sdl.GameControllerButton]Button{
	sdl.CONTROLLER_BUTTON_A:             INPUT_GAMEPAD_A,
	sdl.CONTROLLER_BUTTON_B:             INPUT_GAMEPAD_B,
	sdl.CONTROLLER_BUTTON_X:             INPUT_GAMEPAD_X,
	sdl.CONTROLLER_BUTTON_Y:             INPUT_GAMEPAD_Y,
	sdl.CONTROLLER_BUTTON_BACK:          INPUT_GAMEPAD_SELECT,
	sdl.CONTROLLER_BUTTON_START:         INPUT_GAMEPAD_START,
	sdl.CONTROLLER_BUTTON_GUIDE:         INPUT_GAMEPAD_HOME,
	sdl.CONTROLLER_BUTTON_LEFTSTICK:     INPUT_GAMEPAD_L_STICK_PRESS,
	sdl.CONTROLLER_BUTTON_RIGHTSTICK:    INPUT_GAMEPAD_R_STICK_PRESS,
	sdl.CONTROLLER_BUTTON_LEFTSHOULDER:  INPUT_GAMEPAD_L_SHOULDER,
	sdl.CONTROLLER_BUTTON_RIGHTSHOULDER: INPUT_GAMEPAD_R_SHOULDER,
	sdl.CONTROLLER_BUTTON_DPAD_UP:       INPUT_GAMEPAD_DPAD_UP,
	sdl.CONTROLLER_BUTTON_DPAD_DOWN:     INPUT_GAMEPAD_DPAD_DOWN,
	sdl.CONTROLLER_BUTTON_DPAD_LEFT:     INPUT_GAMEPAD_DPAD_LEFT,
	sdl.CONTROLLER_BUTTON_DPAD_RIGHT:    INPUT_GAMEPAD_DPAD_RIGHT,
	sdl.CONTROLLER_BUTTON_MAX:           INPUT_INVALID,
}

var GamepadAxisMap = map[sdl.GameControllerAxis]Button{
	sdl.CONTROLLER_AXIS_LEFTX:        INPUT_GAMEPAD_L_STICK_LEFT,
	sdl.CONTROLLER_AXIS_LEFTY:        INPUT_GAMEPAD_L_STICK_UP,
	sdl.CONTROLLER_AXIS_RIGHTX:       INPUT_GAMEPAD_R_STICK_LEFT,
	sdl.CONTROLLER_AXIS_RIGHTY:       INPUT_GAMEPAD_R_STICK_UP,
	sdl.CONTROLLER_AXIS_TRIGGERLEFT:  INPUT_GAMEPAD_L_TRIGGER,
	sdl.CONTROLLER_AXIS_TRIGGERRIGHT: INPUT_GAMEPAD_R_TRIGGER,
	sdl.CONTROLLER_AXIS_MAX:          INPUT_INVALID,
}

// Exit() exits the game
func (sw *PlatformSdl) Exit() {
	sw.wantToExit = true
}

// ExitWanted() returns true if the user wants to exit the game
func (sw *PlatformSdl) ExitWanted() bool {
	return sw.wantToExit
}

type PlatformSdl struct {
	window     *sdl.Window
	glContext  sdl.GLContext
	gamepad    *sdl.GameController
	perfFreq   uint64
	wantToExit bool
}

// NewPlatformSdl creates a window
func NewPlatformSdl(title string, x, y, w, h int32) (*PlatformSdl, error) {
	Logger = log.New(os.Stderr, "engine |", log.Ldate|log.Ltime)

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
    sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
    sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	window, err := sdl.CreateWindow("Wipeout",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		w, h, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|PLATFORM_WINDOW_FLAGS|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		return nil, err
	}

	if err != nil {
		panic(err)
	}
	// gl.Enable(gl.DEPTH_TEST)

	return &PlatformSdl{
		window:     window,
		perfFreq:   sdl.GetPerformanceFrequency(),
		wantToExit: false,
	}, nil
}

func (sw *PlatformSdl) Now() float64 {
	perfCounter := sdl.GetPerformanceCounter()

	return float64(perfCounter) / float64(sw.perfFreq)
}

// FindGamepad returns the first gamepad found
func (sw *PlatformSdl) FindGamepad() {
	for i := 0; i < sdl.NumJoysticks(); i++ {
		if sdl.IsGameController(i) {
			sw.gamepad = sdl.GameControllerOpen(i)
		}
	}
}

// PumpEvents pumps events from SDL
func (sw *PlatformSdl) PumpEvents() error {
	// Keyboards inputs
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		// Handle Fullscreen with F11
		if event.GetType() == sdl.KEYDOWN && event.(*sdl.KeyboardEvent).Keysym.Scancode == sdl.SCANCODE_F11 {
			Logger.Println("full screen")
			fullscreen := !sw.IsFullScreen()
			err := sw.SetFullscreen(fullscreen)
			if err != nil {
				return err
			}

		} else if event.GetType() == sdl.KEYDOWN || event.GetType() == sdl.KEYUP {
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
		} else if event.GetType() == sdl.TEXTINPUT {
			InputTextInput(int32(event.(*sdl.TextInputEvent).Text[0]))

			// Gamepad connected/disconnected
		} else if event.GetType() == sdl.CONTROLLERDEVICEADDED {
			sw.gamepad = sdl.GameControllerOpen(int(event.(*sdl.ControllerDeviceEvent).Which))
		} else if event.GetType() == sdl.CONTROLLERDEVICEREMOVED {
			if sw.gamepad != nil && event.(*sdl.ControllerDeviceEvent).Which == sw.gamepad.Joystick().InstanceID() {
				sw.gamepad.Close()
				sw.gamepad = nil
			}
			// Input Gamepad buttons
		} else if event.GetType() == sdl.CONTROLLERBUTTONDOWN || event.GetType() == sdl.CONTROLLERBUTTONUP {
			button := GamepadMap[sdl.GameControllerButton(event.(*sdl.ControllerButtonEvent).Button)]
			if button != INPUT_INVALID {
				var state float32
				if event.GetType() == sdl.CONTROLLERBUTTONDOWN {
					state = 0
				} else {
					state = 1
				}
				InputSetButtonState(button, state)
			}
		} else if event.GetType() == sdl.CONTROLLERAXISMOTION {
			var state float32 = float32(event.(*sdl.ControllerAxisEvent).Value) / 32767.0

			if event.(*sdl.ControllerAxisEvent).Axis < sdl.CONTROLLER_AXIS_MAX {
				code := GamepadAxisMap[sdl.GameControllerAxis(event.(*sdl.ControllerAxisEvent).Axis)]
				if code == INPUT_GAMEPAD_L_TRIGGER || code == INPUT_GAMEPAD_R_TRIGGER {
					InputSetButtonState(code, state)
				} else if state > 0 {
					InputSetButtonState(code, 0.0)
					InputSetButtonState(code+1, state)
				} else {
					InputSetButtonState(code, -state)
					InputSetButtonState(code+1, 0.0)
				}
			}
		} else if event.GetType() == sdl.QUIT {
			sw.wantToExit = true
			Logger.Println("quit wanted")
		}
	}

	return nil
}

// ScreenSize returns the size of the screen
func (sw *PlatformSdl) ScreenSize() Vec2i {
	w, h := sw.window.GLGetDrawableSize()
	return Vec2i{w, h}
}

// SetFullscreen sets the window to fullscreen mode
func (sw *PlatformSdl) SetFullscreen(fullscreen bool) error {
	if fullscreen {

		display, err := sw.window.GetDisplayIndex()
		if err != nil {
			return err
		}
		mode, err := sdl.GetDesktopDisplayMode(display)
		if err != nil {
			return err
		}
		err = sw.window.SetDisplayMode(&mode)
		if err != nil {
			return err
		}
		err = sw.window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
		if err != nil {
			return err
		}
		_, err = sdl.ShowCursor(sdl.DISABLE)
		if err != nil {
			return err
		}
	} else {
		err := sw.window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
		if err != nil {
			return err
		}
		_, err = sdl.ShowCursor(sdl.ENABLE)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsFullScreen returns true if the window is in fullscreen mode
func (sw *PlatformSdl) IsFullScreen() bool {
	return sw.window.GetFlags()&sdl.WINDOW_FULLSCREEN != 0
}

func (sw *PlatformSdl) VideoInit() error {
	Logger.Printf("Video init")
	var err error
	// err = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	// if err != nil {
	// 	return err
	// }
	// err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	// if err != nil {
	// 	return err
	// }
	// err = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	// if err != nil {
	// 	return err
	// }
	sw.glContext, err = sw.window.GLCreateContext()
	Logger.Printf("GlContext=%v", sw.glContext)
	if err != nil {
		return err
	}
	err = sdl.GLSetSwapInterval(1)
	if err != nil {
		return err
	}

	return nil
}

func (sw *PlatformSdl) VideoCleanup() {
	sdl.GLDeleteContext(sw.glContext)
}

func (sw *PlatformSdl) PrepareFrame() error {
	// if ScreenBufferSize.X != ScreenSize.X || ScreenBufferSize.Y != ScreenSize.Y {
	// 	if ScreenBuffer != nil {
	// 		err := ScreenBuffer.Destroy()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	ScreenBuffer, _ = Renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, ScreenSize.X, ScreenSize.Y)
	// 	ScreenBufferSize = ScreenSize
	// }
	// _, _, err := ScreenBuffer.Lock(ScreenBufferPixels)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (sw *PlatformSdl) EndFrame() error {
	// ScreenBufferPixels = nil
	// ScreenBuffer.Unlock()
	// err := Renderer.Copy(ScreenBuffer, nil, nil)
	// if err != nil {
	// 	return err
	// }
	// Renderer.Present()
	sw.window.GLSwap()

	return nil
}

func (sw *PlatformSdl) GetScreenBuffer() *sdl.Rect {
	return ScreenBufferPixels
}

func (sw *PlatformSdl) GetScreenSize() Vec2i {
	ScreenSize.X, ScreenSize.Y = sw.window.GetSize()
	return ScreenSize
}

func (sw *PlatformSdl) Destroy() error {
	return sw.window.Destroy()
}
