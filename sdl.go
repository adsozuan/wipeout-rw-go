package wipeout

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	GlContext  sdl.GLContext
	Gamepad    *sdl.GameController
	WantToExit bool = false
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
func Exit() {
	WantToExit = true
}

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
func (sw *SdlWindow) PumpEvents() {
	var event sdl.Event

	// Keyboards inputs
	for !WantToExit {
		event = sdl.PollEvent()

		// Handle Fullscreen with F11
		if event.GetType() == sdl.KEYDOWN && event.(*sdl.KeyboardEvent).Keysym.Scancode == sdl.SCANCODE_F11 {
			fullscreen := !sw.IsFullScreen()
			sw.SetFullscreen(fullscreen)
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
			Gamepad = sdl.GameControllerOpen(int(event.(*sdl.ControllerDeviceEvent).Which))
		} else if event.GetType() == sdl.CONTROLLERDEVICEREMOVED {
			if Gamepad != nil && event.(*sdl.ControllerDeviceEvent).Which == Gamepad.Joystick().InstanceID() {
				Gamepad.Close()
				Gamepad = nil
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

// IsFullScreen returns true if the window is in fullscreen mode
func (sw *SdlWindow) IsFullScreen() bool {
	return sw.window.GetFlags()&sdl.WINDOW_FULLSCREEN != 0
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
