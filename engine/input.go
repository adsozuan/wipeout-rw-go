package engine

import gl "github.com/chsc/gogl/gl33"

// Button enumerations
type Button int

const (
	InputInvalid            Button = 0
	InputKeyA               Button = 4
	InputKeyB               Button = 5
	InputKeyC               Button = 6
	InputKeyD               Button = 7
	InputKeyE               Button = 8
	InputKeyF               Button = 9
	InputKeyG               Button = 10
	InputKeyH               Button = 11
	InputKeyI               Button = 12
	InputKeyJ               Button = 13
	InputKeyK               Button = 14
	InputKeyL               Button = 15
	InputKeyM               Button = 16
	InputKeyN               Button = 17
	InputKeyO               Button = 18
	InputKeyP               Button = 19
	InputKeyQ               Button = 20
	InputKeyR               Button = 21
	InputKeyS               Button = 22
	InputKeyT               Button = 23
	InputKeyU               Button = 24
	InputKeyV               Button = 25
	InputKeyW               Button = 26
	InputKeyX               Button = 27
	InputKeyY               Button = 28
	InputKeyZ               Button = 29
	InputKey1               Button = 30
	InputKey2               Button = 31
	InputKey3               Button = 32
	InputKey4               Button = 33
	InputKey5               Button = 34
	InputKey6               Button = 35
	InputKey7               Button = 36
	InputKey8               Button = 37
	InputKey9               Button = 38
	InputKey0               Button = 39
	InputKeyReturn          Button = 40
	InputKeyEscape          Button = 41
	InputKeyBackspace       Button = 42
	InputKeyTab             Button = 43
	InputKeySpace           Button = 44
	InputKeyMinus           Button = 45
	InputKeyEquals          Button = 46
	InputKeyLeftBrace       Button = 47
	InputKeyRightBrace      Button = 48
	InputKeyBackslash       Button = 49
	InputKeyHash            Button = 50
	InputKeySemicolon       Button = 51
	InputKeyApostrophe      Button = 52
	InputKeyTilde           Button = 53
	InputKeyComma           Button = 54
	InputKeyPeriod          Button = 55
	InputKeySlash           Button = 56
	InputKeyCapsLock        Button = 57
	InputKeyF1              Button = 58
	InputKeyF2              Button = 59
	InputKeyF3              Button = 60
	InputKeyF4              Button = 61
	InputKeyF5              Button = 62
	InputKeyF6              Button = 63
	InputKeyF7              Button = 64
	InputKeyF8              Button = 65
	InputKeyF9              Button = 66
	InputKeyF10             Button = 67
	InputKeyF11             Button = 68
	InputKeyF12             Button = 69
	InputKeyPrintScreen     Button = 70
	InputKeyScrollLock      Button = 71
	InputKeyPause           Button = 72
	InputKeyInsert          Button = 73
	InputKeyHome            Button = 74
	InputKeyPageUp          Button = 75
	InputKeyDelete          Button = 76
	InputKeyEnd             Button = 77
	InputKeyPageDown        Button = 78
	InputKeyRight           Button = 79
	InputKeyLeft            Button = 80
	InputKeyDown            Button = 81
	InputKeyUp              Button = 82
	InputKeyNumlock         Button = 83
	InputKeyKpDivide        Button = 84
	InputKeyKpMultiply      Button = 85
	InputKeyKpMinus         Button = 86
	InputKeyKpPlus          Button = 87
	InputKeyKpEnter         Button = 88
	InputKeyKp1             Button = 89
	InputKeyKp2             Button = 90
	InputKeyKp3             Button = 91
	InputKeyKp4             Button = 92
	InputKeyKp5             Button = 93
	InputKeyKp6             Button = 94
	InputKeyKp7             Button = 95
	InputKeyKp8             Button = 96
	InputKeyKp9             Button = 97
	InputKeyKp0             Button = 98
	InputKeyKpPeriod        Button = 99
	InputKeyLCtrl           Button = 100
	InputKeyLShift          Button = 101
	InputKeyLAlt            Button = 102
	InputKeyLGui            Button = 103
	InputKeyRCtrl           Button = 104
	InputKeyRShift          Button = 105
	InputKeyRAlt            Button = 106
	InputKeyMax             Button = 107
	InputGamepadA           Button = 108
	InputGamepadY           Button = 109
	InputGamepadB           Button = 110
	InputGamepadX           Button = 111
	InputGamepadLShoulder   Button = 112
	InputGamepadRShoulder   Button = 113
	InputGamepadLTrigger    Button = 114
	InputGamepadRTrigger    Button = 115
	InputGamepadSelect      Button = 116
	InputGamepadStart       Button = 117
	InputGamepadLStickPress Button = 118
	InputGamepadRStickPress Button = 119
	InputGamepadDpadUp      Button = 120
	InputGamepadDpadDown    Button = 121
	InputGamepadDpadLeft    Button = 122
	InputGamepadDpadRight   Button = 123
	InputGamepadHome        Button = 124
	InputGamepadLStickUp    Button = 125
	InputGamepadLStickDown  Button = 126
	InputGamepadLStickLeft  Button = 127
	InputGamepadLStickRight Button = 128
	InputGamepadRStickUp    Button = 129
	InputGamepadRStickDown  Button = 130
	InputGamepadRStickLeft  Button = 131
	InputGamepadRStickRight Button = 132
	InputMouseLeft          Button = 134
	InputMouseMiddle        Button = 135
	InputMouseRight         Button = 136
	InputMouseWheelUp       Button = 137
	InputMouseWheelDown     Button = 138
	InputButtonMax          Button = 139
)

// Input layer enumerations
type InputLayer int

const (
	InputLayerSystem InputLayer = 0
	InputLayerUser   InputLayer = 1
	InputLayerMax    InputLayer = 2
)

// Input action constants
const (
	InputActionCommand   = 31
	InputActionMax       = 32
	InputDeadzone        = 0.1
	InputDeadzoneCapture = 0.5
	InputActionNone      = 255
	InputButtonNone      = 0
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
	ActionsState    [InputActionMax]float32
	ActionsPressed  [InputActionMax]bool
	ActionsReleased [InputActionMax]bool
	ExpectedButton  [InputActionMax]byte
	Bindings        [InputLayerMax][InputButtonMax]byte
	CaptureCallback InputCaptureCallback
	CaptureUser     interface{}
	MouseX, MouseY  int32
)

func InputInit() {
	InputUnbindAll(InputLayerSystem)
	InputUnbindAll(InputLayerUser)
}

func InputMousePos() Vec2 {
	return Vec2{X: gl.Float(MouseX), Y: gl.Float(MouseY)}
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
		CaptureCallback(CaptureUser, InputInvalid, asciiChar)
	}
}

func InputCleanUp() {
}

func InputClear() {
	clearBoolSlice(&ActionsPressed)
	clearBoolSlice(&ActionsReleased)
}

func clearBoolSlice(slice *[InputActionMax]bool) {
	for i := range slice {
		slice[i] = false
	}
}

func InputSetLayerButtonState(layer InputLayer, button Button, state float32) {
	if layer < 0 || layer >= InputLayerMax {
		return
	}

	action := Bindings[layer][button]
	if action == InputActionNone {
		return
	}

	expected := ExpectedButton[action]
	if expected == 0 || expected == byte(button) {
		if state > InputDeadzone {
			state = state
		} else {
			state = 0
		}

		if state > 0 && ActionsState[action] == 0 {
			ActionsPressed[action] = true
			ExpectedButton[action] = byte(button)
		} else if state == 0 && ActionsState[action] != 0 {
			ActionsReleased[action] = true
			ExpectedButton[action] = InputButtonNone
		}
		ActionsState[action] = state
	}
}

func InputSetButtonState(button Button, state float32) {
	if button < 0 || button >= InputButtonMax {
		return
	}

	InputSetLayerButtonState(InputLayerSystem, button, state)
	InputSetLayerButtonState(InputLayerUser, button, state)

	if CaptureCallback != nil {
		if state > InputDeadzoneCapture {
			CaptureCallback(CaptureUser, button, 0)
		}
	}
}

func InputBind(layer InputLayer, button Button, action byte) {
	if button < 0 || button >= InputButtonMax || action < 0 || action >= InputActionMax || layer < 0 || layer >= InputLayerMax {
		return
	}
	ActionsState[action] = 0
	Bindings[layer][button] = action
}

func InputBoundToAction(button Button) byte {
	if button < 0 || button >= InputButtonMax {
		return 0
	}
	return Bindings[InputLayerUser][button]
}

func InputUnbind(layer InputLayer, button Button) {
	if layer < 0 || layer >= InputLayerMax || button < 0 || button >= InputButtonMax {
		return
	}
	Bindings[layer][button] = InputActionNone
}

func InputUnbindAll(layer InputLayer) {
	if layer < 0 || layer >= InputLayerMax {
		return
	}
	for button := Button(0); button < InputButtonMax; button++ {
		InputUnbind(layer, button)
	}
}

func InputState(action byte) float32 {
	if action < 0 || action >= InputActionMax {
		return 0
	}
	return ActionsState[action]
}

func InputPressed(action byte) bool {
	if action < 0 || action >= InputActionMax {
		return false
	}
	return ActionsPressed[action]
}

func InputReleased(action byte) bool {
	if action < 0 || action >= InputActionMax {
		return false
	}
	return ActionsReleased[action]
}

func InputNameToButton(name string) Button {
	for i := 0; i < int(InputButtonMax); i++ {
		if buttonNames[i] != "" && buttonNames[i] == name {
			return Button(i)
		}
	}
	return InputInvalid
}

func InputButtonToName(button Button) string {
	if button < 0 || button >= InputButtonMax || buttonNames[button] == "" {
		return ""
	}
	return buttonNames[button]
}
