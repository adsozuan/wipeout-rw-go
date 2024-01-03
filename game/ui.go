package game

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/adsozuan/wipeout-rw-go/engine"

	gl "github.com/chsc/gogl/gl33"
)

type UITextSize int

const (
	UITextSize16 UITextSize = iota
	UITextSize12
	UITextSize8
	UITextSizeMax
)

var (
	UIColorAccent  = engine.RGBA{R: 123, G: 98, B: 12, A: 255}
	UIColorDefault = engine.RGBA{R: 128, G: 128, B: 128, A: 255}
)

type UIIconType int

const (
	UIIconHand UIIconType = iota
	UIIconConfirm
	UIIconCancel
	UIIconEnd
	UIIconDelete
	UIIconStar
	UIIconMax
)

type UIPos int

const (
	UIPosLeft UIPos = 1 << iota
	UIPosCenter
	UIPosRight
	UIPosTop
	UIPosMiddle
	UIPosBottom
)

type Glyph struct {
	Offset engine.Vec2i
	Width  uint16
}

type CharSet struct {
	Texture uint16
	Height  uint16
	Glyphs  [40]Glyph
}

type UI struct {
	charSet      [UITextSizeMax]CharSet
	scale        int
	iconTextures [UIIconMax]uint16
	render       *engine.Render
}

func NewUI(render *engine.Render) *UI {
	return &UI{
		charSet: charSet,
		scale:   2,
		iconTextures: [UIIconMax]uint16{
			0,
			0,
			0,
			0,
			0,
			0,
		},
		render: render,
	}
}

func (ui *UI) Load() error {
	tl, err := ImageGetCompressedTexture("data/textures/drfonts.cmp", ui.render)
	if err != nil {
		return errors.New("UI could not load font textures")
	}
	charSet[UITextSize16].Texture = uint16(TextureFromList(tl, 0))
	charSet[UITextSize12].Texture = uint16(TextureFromList(tl, 1))
	charSet[UITextSize8].Texture = uint16(TextureFromList(tl, 2))
	ui.iconTextures[UIIconHand] = uint16(TextureFromList(tl, 3))
	ui.iconTextures[UIIconConfirm] = uint16(TextureFromList(tl, 5))
	ui.iconTextures[UIIconCancel] = uint16(TextureFromList(tl, 6))
	ui.iconTextures[UIIconEnd] = uint16(TextureFromList(tl, 6))

	return nil
}

func (ui *UI) GetScale() int {
	return ui.scale
}

func (ui *UI) SetScale(scale int) {
	ui.scale = scale
}

func (ui *UI) Scaled(size engine.Vec2i) engine.Vec2i {
	return engine.Vec2i{X: size.X * int32(ui.scale), Y: size.Y * int32(ui.scale)}
}

func (ui *UI) ScaledScreen() engine.Vec2i {
	return engine.Vec2iMulF(ui.render.Size(), gl.Float(ui.scale))
}

func (ui *UI) ScaledPos(anchor UIPos, offset engine.Vec2i) engine.Vec2i {
	var pos engine.Vec2i
	screenSize := ui.render.Size()

	if anchor&UIPosLeft != 0 {
		pos.X = offset.X * int32(ui.scale)
	} else if anchor&UIPosCenter != 0 {
		pos.X = (screenSize.X >> 1) + offset.X*int32(ui.scale)
	} else if anchor&UIPosRight != 0 {
		pos.X = screenSize.X + offset.X*int32(ui.scale)
	}

	if anchor&UIPosTop != 0 {
		pos.Y = offset.Y * int32(ui.scale)
	} else if anchor&UIPosMiddle != 0 {
		pos.Y = (screenSize.Y >> 1) + offset.Y*int32(ui.scale)
	} else if anchor&UIPosBottom != 0 {
		pos.Y = screenSize.Y + offset.Y*int32(ui.scale)
	}

	return pos
}

func charToGlyphIndex(c rune) int {
	if c >= '0' && c <= '9' {
		return int(c - '0' + 26)
	}
	return int(c - 'A')
}

func charWidth(c rune, size UITextSize) int {
	if c == ' ' {
		return 8
	}
	return int(charSet[size].Glyphs[charToGlyphIndex(c)].Width)
}

func textWidth(text string, size UITextSize) int {
	var width uint16 = 0
	cs := &charSet[size]

	for _, ch := range text {
		if ch != ' ' {
			width += cs.Glyphs[charToGlyphIndex(ch)].Width
		} else {
			width += 8
		}

	}

	return int(width)
}

func numberWidth(num int, size UITextSize) int {
	text := strconv.Itoa(num)
	return textWidth(text, size)
}

func (ui *UI) DrawTime(time float32, pos engine.Vec2i, size UITextSize, color engine.RGBA) {
	msec := int(time * 1000)
	tenths := (msec / 100) % 10
	secs := (msec / 1000) % 60
	mins := msec / (60 * 1000)

	textBuffer := fmt.Sprintf("%02d:%02d.%d", mins, secs, tenths)
	ui.DrawText(textBuffer, pos, size, color)
}

func (ui *UI) DrawNumber(num int, pos engine.Vec2i, size UITextSize, color engine.RGBA) {
	textBuffer := strconv.Itoa(num)
	ui.DrawText(textBuffer, pos, size, color)
}

// DrawTextCentered renders centered text on the UI.
func (ui *UI) DrawTextCentered(text string, pos engine.Vec2i, size UITextSize, color engine.RGBA) {
	textWidth := textWidth(text, size) * ui.scale
	pos.X -= int32(textWidth >> 1)
	ui.DrawText(text, pos, size, color)
}

// DrawText renders text on the UI.
func (ui *UI) DrawText(text string, pos engine.Vec2i, size UITextSize, color engine.RGBA) {
	cs := &charSet[size]

	for _, char := range text {
		if char != ' ' {
			glyph := &cs.Glyphs[ui.charToGlyphIndex(char)]
			glyphOffset := engine.Vec2i{X: int32(glyph.Offset.X), Y: int32(glyph.Offset.Y)}
			glyphSize := engine.Vec2i{X: int32(glyph.Width), Y: int32(cs.Height)}
			ui.render.Push2dTitle(pos, glyphOffset, glyphSize, ui.Scaled(glyphSize), color, int(cs.Texture))
			pos.X += int32(int(glyph.Width) * ui.scale)
		} else {
			pos.X += int32(8 * ui.scale)
		}
	}
}

func (ui *UI) DrawImage(pos engine.Vec2i, texture int) error {
	scale, err := ui.render.TextureSize(texture)
	if err != nil {
		return err
	}

	scaledSize := ui.Scaled(scale)
	err = ui.render.Push2d(pos, scaledSize, engine.RGBA{128, 128, 128, 255}, texture)
	if err != nil {
		return err
	}

	return nil
}

func (ui *UI) DrawIcon(icon UIIconType, pos engine.Vec2i, color engine.RGBA) error {
	size, err := ui.render.TextureSize(int(ui.iconTextures[icon]))
	if err != nil {
		return err
	}
	scaledSize := ui.Scaled(size)
	ui.render.Push2d(pos, scaledSize, color, int(ui.iconTextures[icon]))

	return nil
}

// charToGlyphIndex converts a character to a glyph index.
func (ui *UI) charToGlyphIndex(char rune) int {
	if char >= '0' && char <= '9' {
		return int(char-'0') + 26
	}
	return int(char-'A') - 1
}

var charSet [UITextSizeMax]CharSet = [UITextSizeMax]CharSet{
	UITextSize16: {
		Texture: 0,
		Height:  16,
		Glyphs: [40]Glyph{
			{engine.Vec2i{0, 0}, 25}, {engine.Vec2i{25, 0}, 24}, {engine.Vec2i{49, 0}, 17}, {engine.Vec2i{66, 0}, 24}, {engine.Vec2i{90, 0}, 24}, {engine.Vec2i{114, 0}, 17}, {engine.Vec2i{131, 0}, 25}, {engine.Vec2i{156, 0}, 18},
			{engine.Vec2i{174, 0}, 7}, {engine.Vec2i{181, 0}, 17}, {engine.Vec2i{0, 16}, 17}, {engine.Vec2i{17, 16}, 17}, {engine.Vec2i{34, 16}, 28}, {engine.Vec2i{62, 16}, 17}, {engine.Vec2i{79, 16}, 24}, {engine.Vec2i{103, 16}, 24},
			{engine.Vec2i{127, 16}, 26}, {engine.Vec2i{153, 16}, 24}, {engine.Vec2i{177, 16}, 18}, {engine.Vec2i{195, 16}, 17}, {engine.Vec2i{0, 32}, 17}, {engine.Vec2i{17, 32}, 17}, {engine.Vec2i{34, 32}, 29}, {engine.Vec2i{63, 32}, 24},
			{engine.Vec2i{87, 32}, 17}, {engine.Vec2i{104, 32}, 18}, {engine.Vec2i{122, 32}, 24}, {engine.Vec2i{146, 32}, 10}, {engine.Vec2i{156, 32}, 18}, {engine.Vec2i{174, 32}, 17}, {engine.Vec2i{191, 32}, 18}, {engine.Vec2i{0, 48}, 18},
			{engine.Vec2i{18, 48}, 18}, {engine.Vec2i{36, 48}, 18}, {engine.Vec2i{54, 48}, 22}, {engine.Vec2i{76, 48}, 25}, {engine.Vec2i{101, 48}, 7}, {engine.Vec2i{108, 48}, 7}, {engine.Vec2i{198, 0}, 0}, {engine.Vec2i{198, 0}, 0},
		},
	},
	UITextSize12: {
		Texture: 0,
		Height:  12,
		Glyphs: [40]Glyph{
			{engine.Vec2i{0, 0}, 19}, {engine.Vec2i{19, 0}, 19}, {engine.Vec2i{38, 0}, 14}, {engine.Vec2i{52, 0}, 19}, {engine.Vec2i{71, 0}, 19}, {engine.Vec2i{90, 0}, 13}, {engine.Vec2i{103, 0}, 19}, {engine.Vec2i{122, 0}, 14},
			{engine.Vec2i{136, 0}, 6}, {engine.Vec2i{142, 0}, 13}, {engine.Vec2i{155, 0}, 14}, {engine.Vec2i{169, 0}, 14}, {engine.Vec2i{0, 12}, 22}, {engine.Vec2i{22, 12}, 14}, {engine.Vec2i{36, 12}, 19}, {engine.Vec2i{55, 12}, 18},
			{engine.Vec2i{73, 12}, 20}, {engine.Vec2i{93, 12}, 19}, {engine.Vec2i{112, 12}, 15}, {engine.Vec2i{127, 12}, 14}, {engine.Vec2i{141, 12}, 13}, {engine.Vec2i{154, 12}, 13}, {engine.Vec2i{167, 12}, 22}, {engine.Vec2i{0, 24}, 19},
			{engine.Vec2i{19, 24}, 13}, {engine.Vec2i{32, 24}, 14}, {engine.Vec2i{46, 24}, 19}, {engine.Vec2i{65, 24}, 8}, {engine.Vec2i{73, 24}, 15}, {engine.Vec2i{88, 24}, 13}, {engine.Vec2i{101, 24}, 14}, {engine.Vec2i{115, 24}, 15},
			{engine.Vec2i{130, 24}, 14}, {engine.Vec2i{144, 24}, 15}, {engine.Vec2i{159, 24}, 18}, {engine.Vec2i{177, 24}, 19}, {engine.Vec2i{196, 24}, 5}, {engine.Vec2i{201, 24}, 5}, {engine.Vec2i{183, 0}, 0}, {engine.Vec2i{183, 0}, 0},
		},
	},
	UITextSize8: {
		Texture: 0,
		Height:  8,
		Glyphs: [40]Glyph{
			{engine.Vec2i{0, 0}, 13}, {engine.Vec2i{13, 0}, 13}, {engine.Vec2i{26, 0}, 10}, {engine.Vec2i{36, 0}, 13}, {engine.Vec2i{49, 0}, 13}, {engine.Vec2i{62, 0}, 9}, {engine.Vec2i{71, 0}, 13}, {engine.Vec2i{84, 0}, 10},
			{engine.Vec2i{94, 0}, 4}, {engine.Vec2i{98, 0}, 9}, {engine.Vec2i{107, 0}, 10}, {engine.Vec2i{117, 0}, 10}, {engine.Vec2i{127, 0}, 16}, {engine.Vec2i{143, 0}, 10}, {engine.Vec2i{153, 0}, 13}, {engine.Vec2i{166, 0}, 13},
			{engine.Vec2i{179, 0}, 14}, {engine.Vec2i{0, 8}, 13}, {engine.Vec2i{13, 8}, 10}, {engine.Vec2i{23, 8}, 9}, {engine.Vec2i{32, 8}, 9}, {engine.Vec2i{41, 8}, 9}, {engine.Vec2i{50, 8}, 16}, {engine.Vec2i{66, 8}, 14},
			{engine.Vec2i{80, 8}, 9}, {engine.Vec2i{89, 8}, 10}, {engine.Vec2i{99, 8}, 13}, {engine.Vec2i{112, 8}, 6}, {engine.Vec2i{118, 8}, 11}, {engine.Vec2i{129, 8}, 10}, {engine.Vec2i{139, 8}, 10}, {engine.Vec2i{149, 8}, 11},
			{engine.Vec2i{160, 8}, 10}, {engine.Vec2i{170, 8}, 10}, {engine.Vec2i{180, 8}, 12}, {engine.Vec2i{192, 8}, 14}, {engine.Vec2i{206, 8}, 4}, {engine.Vec2i{210, 8}, 4}, {engine.Vec2i{193, 0}, 0}, {engine.Vec2i{193, 0}, 0},
		},
	},
}
