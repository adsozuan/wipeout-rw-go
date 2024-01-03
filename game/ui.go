package wipeout

import (
	"encoding/binary"

	"github.com/blacktop/lzss"

	"github.com/adsozuan/wipeout-rw-go/engine"
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

func (ui *UI) Load() {

}

func (ui *UI) GetUiScale() int {
	return ui.scale
}

func (ui *UI) SetUiScale(scale int) {
	ui.scale = scale
}

// DrawTextCentered renders centered text on the UI.
func (ui *UI) DrawTextCentered(text string, pos engine.Vec2i, size UITextSize, color engine.RGBA) {
	textWidth := ui.textWidth(text, size) * ui.scale
	pos.X -= int32(textWidth >> 1)
	ui.DrawText(text, pos, size, color)
}

// textWidth calculates the width of the text.
func (ui *UI) textWidth(text string, size UITextSize) int {
	cs := charSet[size]
	width := 0

	for _, char := range text {
		if char != ' ' {
			glyph := &cs.Glyphs[ui.charToGlyphIndex(char)]
			width += int(glyph.Width) * ui.scale
		} else {
			width += 8 * ui.scale
		}
	}

	return width
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

func (ui *UI) Scaled(size engine.Vec2i) engine.Vec2i {
	return engine.Vec2i{X: size.X * int32(ui.scale), Y: size.Y * int32(ui.scale)}
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

func (ui *UI) ScaledPos(anchor UIPos, offset engine.Vec2i) engine.Vec2i {
	var pos engine.Vec2i
	screenSize := ui.render.Size() 

	if anchor&UIPosLeft != 0 {
		pos.X = offset.X * int32(ui.scale)
	} else if anchor&UIPosCenter != 0 {
		pos.X = (screenSize.X >> 1) + offset.X * int32(ui.scale)
	} else if anchor&UIPosRight != 0 {
		pos.X = screenSize.X + offset.X * int32(ui.scale)
	}

	if anchor&UIPosTop != 0 {
		pos.Y = offset.Y * int32(ui.scale)
	} else if anchor&UIPosMiddle != 0 {
		pos.Y = (screenSize.Y >> 1) + offset.Y * int32(ui.scale)
	} else if anchor&UIPosBottom != 0 {
		pos.Y = screenSize.Y + offset.Y * int32(ui.scale)
	}

	return pos
	
}

type cmpT struct {
	Len     uint32
	Entries []*uint8
}

func imageLoadCompressed(name string) (*cmpT, error) {
	Logger.Printf("load cmp %s\n", name)

	// Load compressed bytes from the file
	compressedBytes, err := engine.LoadBinaryFile(name)
	if err != nil {
		return nil, err
	}

	data := lzss.Decompress(compressedBytes)

	// Initialize variables
	var p uint32
	var decompressedBytesOffset uint32

	// Read the number of entries (Len) from data
	imageCount := binary.LittleEndian.Uint32(data[p:])
	p += 4

	// Create a slice to hold pointers to uint8
	entries := make([]*uint8, imageCount)

	// Iterate through the entries and store their pointers
	for i := uint32(0); i < imageCount; i++ {
		offset := binary.LittleEndian.Uint32(data[p:])
		entries[i] = &data[decompressedBytesOffset+offset]
		p += 4
	}

	// Create and return the cmpT struct
	return &cmpT{
		Len:     imageCount,
		Entries: entries,
	}, nil
}
