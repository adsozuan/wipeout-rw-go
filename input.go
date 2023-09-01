package wipeout

// Button enumerations
type Button int

const (
	INPUT_INVALID               Button = 0
	INPUT_KEY_A                 Button = 4
	INPUT_KEY_B                 Button = 5
	INPUT_KEY_C                 Button = 6
	INPUT_KEY_D                 Button = 7
	INPUT_KEY_E                 Button = 8
	INPUT_KEY_F                 Button = 9
	INPUT_KEY_G                 Button = 10
	INPUT_KEY_H                 Button = 11
	INPUT_KEY_I                 Button = 12
	INPUT_KEY_J                 Button = 13
	INPUT_KEY_K                 Button = 14
	INPUT_KEY_L                 Button = 15
	INPUT_KEY_M                 Button = 16
	INPUT_KEY_N                 Button = 17
	INPUT_KEY_O                 Button = 18
	INPUT_KEY_P                 Button = 19
	INPUT_KEY_Q                 Button = 20
	INPUT_KEY_R                 Button = 21
	INPUT_KEY_S                 Button = 22
	INPUT_KEY_T                 Button = 23
	INPUT_KEY_U                 Button = 24
	INPUT_KEY_V                 Button = 25
	INPUT_KEY_W                 Button = 26
	INPUT_KEY_X                 Button = 27
	INPUT_KEY_Y                 Button = 28
	INPUT_KEY_Z                 Button = 29
	INPUT_KEY_1                 Button = 30
	INPUT_KEY_2                 Button = 31
	INPUT_KEY_3                 Button = 32
	INPUT_KEY_4                 Button = 33
	INPUT_KEY_5                 Button = 34
	INPUT_KEY_6                 Button = 35
	INPUT_KEY_7                 Button = 36
	INPUT_KEY_8                 Button = 37
	INPUT_KEY_9                 Button = 38
	INPUT_KEY_0                 Button = 39
	INPUT_KEY_RETURN            Button = 40
	INPUT_KEY_ESCAPE            Button = 41
	INPUT_KEY_BACKSPACE         Button = 42
	INPUT_KEY_TAB               Button = 43
	INPUT_KEY_SPACE             Button = 44
	INPUT_KEY_MINUS             Button = 45
	INPUT_KEY_EQUALS            Button = 46
	INPUT_KEY_LEFTBRACKET       Button = 47
	INPUT_KEY_RIGHTBRACKET      Button = 48
	INPUT_KEY_BACKSLASH         Button = 49
	INPUT_KEY_HASH              Button = 50
	INPUT_KEY_SEMICOLON         Button = 51
	INPUT_KEY_APOSTROPHE        Button = 52
	INPUT_KEY_TILDE             Button = 53
	INPUT_KEY_COMMA             Button = 54
	INPUT_KEY_PERIOD            Button = 55
	INPUT_KEY_SLASH             Button = 56
	INPUT_KEY_CAPSLOCK          Button = 57
	INPUT_KEY_F1                Button = 58
	INPUT_KEY_F2                Button = 59
	INPUT_KEY_F3                Button = 60
	INPUT_KEY_F4                Button = 61
	INPUT_KEY_F5                Button = 62
	INPUT_KEY_F6                Button = 63
	INPUT_KEY_F7                Button = 64
	INPUT_KEY_F8                Button = 65
	INPUT_KEY_F9                Button = 66
	INPUT_KEY_F10               Button = 67
	INPUT_KEY_F11               Button = 68
	INPUT_KEY_F12               Button = 69
	INPUT_KEY_PRINTSCREEN       Button = 70
	INPUT_KEY_SCROLLLOCK        Button = 71
	INPUT_KEY_PAUSE             Button = 72
	INPUT_KEY_INSERT            Button = 73
	INPUT_KEY_HOME              Button = 74
	INPUT_KEY_PAGEUP            Button = 75
	INPUT_KEY_DELETE            Button = 76
	INPUT_KEY_END               Button = 77
	INPUT_KEY_PAGEDOWN          Button = 78
	INPUT_KEY_RIGHT             Button = 79
	INPUT_KEY_LEFT              Button = 80
	INPUT_KEY_DOWN              Button = 81
	INPUT_KEY_UP                Button = 82
	INPUT_KEY_NUMLOCK           Button = 83
	INPUT_KEY_KP_DIVIDE         Button = 84
	INPUT_KEY_KP_MULTIPLY       Button = 85
	INPUT_KEY_KP_MINUS          Button = 86
	INPUT_KEY_KP_PLUS           Button = 87
	INPUT_KEY_KP_ENTER          Button = 88
	INPUT_KEY_KP_1              Button = 89
	INPUT_KEY_KP_2              Button = 90
	INPUT_KEY_KP_3              Button = 91
	INPUT_KEY_KP_4              Button = 92
	INPUT_KEY_KP_5              Button = 93
	INPUT_KEY_KP_6              Button = 94
	INPUT_KEY_KP_7              Button = 95
	INPUT_KEY_KP_8              Button = 96
	INPUT_KEY_KP_9              Button = 97
	INPUT_KEY_KP_0              Button = 98
	INPUT_KEY_KP_PERIOD         Button = 99
	INPUT_KEY_LCTRL             Button = 100
	INPUT_KEY_LSHIFT            Button = 101
	INPUT_KEY_LALT              Button = 102
	INPUT_KEY_LGUI              Button = 103
	INPUT_KEY_RCTRL             Button = 104
	INPUT_KEY_RSHIFT            Button = 105
	INPUT_KEY_RALT              Button = 106
	INPUT_KEY_MAX               Button = 107
	INPUT_GAMEPAD_A             Button = 108
	INPUT_GAMEPAD_Y             Button = 109
	INPUT_GAMEPAD_B             Button = 110
	INPUT_GAMEPAD_X             Button = 111
	INPUT_GAMEPAD_L_SHOULDER    Button = 112
	INPUT_GAMEPAD_R_SHOULDER    Button = 113
	INPUT_GAMEPAD_L_TRIGGER     Button = 114
	INPUT_GAMEPAD_R_TRIGGER     Button = 115
	INPUT_GAMEPAD_SELECT        Button = 116
	INPUT_GAMEPAD_START         Button = 117
	INPUT_GAMEPAD_L_STICK_PRESS Button = 118
	INPUT_GAMEPAD_R_STICK_PRESS Button = 119
	INPUT_GAMEPAD_DPAD_UP       Button = 120
	INPUT_GAMEPAD_DPAD_DOWN     Button = 121
	INPUT_GAMEPAD_DPAD_LEFT     Button = 122
	INPUT_GAMEPAD_DPAD_RIGHT    Button = 123
	INPUT_GAMEPAD_HOME          Button = 124
	INPUT_GAMEPAD_L_STICK_UP    Button = 125
	INPUT_GAMEPAD_L_STICK_DOWN  Button = 126
	INPUT_GAMEPAD_L_STICK_LEFT  Button = 127
	INPUT_GAMEPAD_L_STICK_RIGHT Button = 128
	INPUT_GAMEPAD_R_STICK_UP    Button = 129
	INPUT_GAMEPAD_R_STICK_DOWN  Button = 130
	INPUT_GAMEPAD_R_STICK_LEFT  Button = 131
	INPUT_GAMEPAD_R_STICK_RIGHT Button = 132
	INPUT_MOUSE_LEFT            Button = 134
	INPUT_MOUSE_MIDDLE          Button = 135
	INPUT_MOUSE_RIGHT           Button = 136
	INPUT_MOUSE_WHEEL_UP        Button = 137
	INPUT_MOUSE_WHEEL_DOWN      Button = 138
	INPUT_BUTTON_MAX            Button = 139
)

// Input layer enumerations
type InputLayer int

const (
	INPUT_LAYER_SYSTEM InputLayer = 0
	INPUT_LAYER_USER   InputLayer = 1
	INPUT_LAYER_MAX    InputLayer = 2
)

// Input action constants
const (
	INPUT_ACTION_COMMAND   = 31
	INPUT_ACTION_MAX       = 32
	INPUT_DEADZONE         = 0.1
	INPUT_DEADZONE_CAPTURE = 0.5
	INPUT_ACTION_NONE      = 255
	INPUT_BUTTON_NONE      = 0
)

var buttonNames = [...]string{
	"",
	"",
	"",
	"",
	"A",
	"B",
	"C",
	"D",
	"E",
	"F",
	"G",
	"H",
	"I",
	"J",
	"K",
	"L",
	"M",
	"N",
	"O",
	"P",
	"Q",
	"R",
	"S",
	"T",
	"U",
	"V",
	"W",
	"X",
	"Y",
	"Z",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"0",
	"RETURN",
	"ESCAPE",
	"BACKSP",
	"TAB",
	"SPACE",
	"MINUS",
	"EQUALS",
	"LBRACKET",
	"RBRACKET",
	"BSLASH",
	"HASH",
	"SMICOL",
	"APO",
	"TILDE",
	"COMMA",
	"PERIOD",
	"SLASH",
	"CAPS",
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"F8",
	"F9",
	"F10",
	"F11",
	"F12",
	"PRTSC",
	"SCRLK",
	"PAUSE",
	"INSERT",
	"HOME",
	"PG UP",
	"DELETE",
	"END",
	"PG DOWN",
	"RIGHT",
	"LEFT",
	"DOWN",
	"UP",
	"NLOCK",
	"KPDIV",
	"KPMUL",
	"KPMINUS",
	"KPPLUS",
	"KPENTER",
	"KP1",
	"KP2",
	"KP3",
	"KP4",
	"KP5",
	"KP6",
	"KP7",
	"KP8",
	"KP9",
	"KP0",
	"KPPERIOD",
	"LCTRL",
	"LSHIFT",
	"LALT",
	"LGUI",
	"RCTRL",
	"RSHIFT",
	"RALT",
	"",
	"A",
	"Y",
	"B",
	"X",
	"LSHLDR",
	"RSHLDR",
	"LTRIG",
	"RTRIG",
	"SELECT",
	"START",
	"LSTK",
	"RSTK",
	"DPUP",
	"DPDOWN",
	"DPLEFT",
	"DPRIGHT",
	"HOME",
	"LSTKUP",
	"LSTKDOWN",
	"LSTKLEFT",
	"LSTKRIGHT",
	"RSTKUP",
	"RSTKDOWN",
	"RSTKLEFT",
	"RSTKRIGHT",
	"",
	"MLEFT",
	"MMIDDLE",
	"MRIGHT",
	"MWUP",
	"MWDOWN",
}

type InputCaptureCallback func(user interface{}, button Button, asciiChar int32)

var (
	ActionsState    [INPUT_ACTION_MAX]float32
	ActionsPressed  [INPUT_ACTION_MAX]bool
	ActionsReleased [INPUT_ACTION_MAX]bool
	ExpectedButton  [INPUT_ACTION_MAX]uint8
	Bindings        [INPUT_LAYER_MAX][INPUT_BUTTON_MAX]uint8
	CaptureCallback InputCaptureCallback
	CaptureUser     interface{}
	MouseX, MouseY  int32
)

func InputInit() {
	InputUnbindAll(INPUT_LAYER_SYSTEM)
	InputUnbindAll(INPUT_LAYER_USER)
}

func InputMousePos() Vec2 {
	return Vec2{X: float32(MouseX), Y: float32(MouseY)}
}

func InputSetMousePos(x, y int32) {
	MouseX = x
	MouseY = y
}

func InputCapture(cb InputCaptureCallback, user interface{}) {
	CaptureCallback = cb
	CaptureUser = user
	InputClear()
}

func InputTextInput(asciiChar int32) {
	if CaptureCallback != nil {
		CaptureCallback(CaptureUser, INPUT_INVALID, asciiChar)
	}
}

func InputCleanUp() {
}

func InputClear() {
	clearBoolSlice(&ActionsPressed)
	clearBoolSlice(&ActionsReleased)
}

func clearBoolSlice(slice *[INPUT_ACTION_MAX]bool) {
	for i := range slice {
		slice[i] = false
	}
}

func InputSetLayerButtonState(layer InputLayer, button Button, state float32) {
	if layer < 0 || layer >= INPUT_LAYER_MAX {
		return
	}

	action := Bindings[layer][button]
	if action == INPUT_ACTION_NONE {
		return
	}

	expected := ExpectedButton[action]
	if expected == 0 || expected == uint8(button) {
		if state > INPUT_DEADZONE {
			state = state
		} else {
			state = 0
		}

		if state > 0 && ActionsState[action] == 0 {
			ActionsPressed[action] = true
			ExpectedButton[action] = uint8(button)
		} else if state == 0 && ActionsState[action] != 0 {
			ActionsReleased[action] = true
			ExpectedButton[action] = INPUT_BUTTON_NONE
		}
		ActionsState[action] = state
	}
}

func InputSetButtonState(button Button, state float32) {
	if button < 0 || button >= INPUT_BUTTON_MAX {
		return
	}

	InputSetLayerButtonState(INPUT_LAYER_SYSTEM, button, state)
	InputSetLayerButtonState(INPUT_LAYER_USER, button, state)

	if CaptureCallback != nil {
		if state > INPUT_DEADZONE_CAPTURE {
			CaptureCallback(CaptureUser, button, 0)
		}
	}
}

func InputBind(layer InputLayer, button Button, action uint8) {
	if button < 0 || button >= INPUT_BUTTON_MAX || action < 0 || action >= INPUT_ACTION_MAX || layer < 0 || layer >= INPUT_LAYER_MAX {
		return
	}
	ActionsState[action] = 0
	Bindings[layer][button] = action
}

func InputBoundToAction(button Button) uint8 {
	if button < 0 || button >= INPUT_BUTTON_MAX {
		return 0
	}
	return Bindings[INPUT_LAYER_USER][button]
}

func InputUnbind(layer InputLayer, button Button) {
	if layer < 0 || layer >= INPUT_LAYER_MAX || button < 0 || button >= INPUT_BUTTON_MAX {
		return
	}
	Bindings[layer][button] = INPUT_ACTION_NONE
}

func InputUnbindAll(layer InputLayer) {
	if layer < 0 || layer >= INPUT_LAYER_MAX {
		return
	}
	for button := Button(0); button < INPUT_BUTTON_MAX; button++ {
		InputUnbind(layer, button)
	}
}

func InputState(action uint8) float32 {
	if action < 0 || action >= INPUT_ACTION_MAX {
		return 0
	}
	return ActionsState[action]
}

func InputPressed(action uint8) bool {
	if action < 0 || action >= INPUT_ACTION_MAX {
		return false
	}
	return ActionsPressed[action]
}

func InputReleased(action uint8) bool {
	if action < 0 || action >= INPUT_ACTION_MAX {
		return false
	}
	return ActionsReleased[action]
}

func InputNameToButton(name string) Button {
	for i := 0; i < int(INPUT_BUTTON_MAX); i++ {
		if buttonNames[i] != "" && buttonNames[i] == name {
			return Button(i)
		}
	}
	return INPUT_INVALID
}

func InputButtonToName(button Button) string {
	if button < 0 || button >= INPUT_BUTTON_MAX || buttonNames[button] == "" {
		return ""
	}
	return buttonNames[button]
}
