package game

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adsozuan/wipeout-rw-go/engine"
	"github.com/blacktop/lzss"
)

type TimTypePalette int32

const (
	TimTypePaletted4BPP   TimTypePalette = 0x08
	TimTypePaletted8BPP   TimTypePalette = 0x09
	TimTypeTrueColor16BPP TimTypePalette = 0x02
)

type Image struct {
	Width, Height uint32
	Pixels        []engine.RGBA
}

type PsxTim struct {
	Magic         []byte
	Type          TimTypePalette
	HeaderSize    int32
	PaletteX      int16
	PaletteY      int16
	PaletteColors int16
	Palettes      int16
	DataSize      int32
	SkipX         int16
	SkipY         int16
	EntriesPerRow int16
	Rows          int16
}

type TextureList struct {
	start int
	len   int
}

func (t *TextureList) Get(index int) (uint16, error) {
	if index < 0 || index >= t.len {

		return 0, fmt.Errorf("Texture index out of range: %d", index)
	}
	return uint16(t.start + index), nil
}

func ImageAlloc(width, height uint32) *Image {
	image := &Image{
		Width:  width,
		Height: height,
		Pixels: make([]engine.RGBA, width*height),
	}
	return image
}

func ImageLoadFromBytes(bytes []byte, transparent bool) *Image {
	p := uint32(0)

	_ = engine.GetU32LE(bytes, &p)
	imgType := engine.GetU32LE(bytes, &p)
	var palette []uint16

	if imgType == uint32(TimTypePaletted4BPP) || imgType == uint32(TimTypePaletted8BPP) {
		_ = engine.GetU32LE(bytes, &p)
		_ = engine.GetU16LE(bytes, &p)
		_ = engine.GetU16LE(bytes, &p)
		paletteColors := engine.GetU16LE(bytes, &p)
		_ = engine.GetU16LE(bytes, &p)

		palette = make([]uint16, paletteColors)
		for i := uint16(0); i < paletteColors; i++ {
			palette[i] = engine.GetU16LE(bytes, &p)
		}
	}

	_ = engine.GetU32LE(bytes, &p)

	pixelsPer16bit := 1
	if imgType == uint32(TimTypePaletted8BPP) {
		pixelsPer16bit = 2
	} else if imgType == uint32(TimTypePaletted4BPP) {
		pixelsPer16bit = 4
	}

	_ = engine.GetU16LE(bytes, &p)
	_ = engine.GetU16LE(bytes, &p)
	entriesPerRow := engine.GetU16LE(bytes, &p)
	rows := engine.GetU16LE(bytes, &p)

	width := int32(entriesPerRow) * int32(pixelsPer16bit)
	height := int32(rows)
	entries := int32(entriesPerRow) * int32(rows)

	image := &Image{
		Width:  uint32(width),
		Height: uint32(height),
		Pixels: make([]engine.RGBA, width*height),
	}

	pixelPos := 0
	if imgType == uint32(TimTypeTrueColor16BPP) {
		for i := int32(0); i < entries; i++ {
			image.Pixels[pixelPos] = tim16BitToRGBA(engine.GetU16LE(bytes, &p), transparent)
			pixelPos++
		}
	} else if imgType == uint32(TimTypePaletted8BPP) {
		for i := int32(0); i < entries; i++ {
			palettePos := engine.GetU16LE(bytes, &p)
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>0)&0xFF], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>8)&0xFF], transparent)
			pixelPos++
		}
	} else if imgType == uint32(TimTypePaletted4BPP) {
		for i := int32(0); i < entries; i++ {
			palettePos := engine.GetU16LE(bytes, &p)
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>0)&0xF], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>4)&0xF], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>8)&0xF], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>12)&0xF], transparent)
			pixelPos++
		}
	}


	return image
}


func ImageGetTexture(name string) uint16 {
	currentDir, err := os.Getwd()
	if err != nil {
		Logger.Printf("ImageGetTexture-Getwd: %s", err)
		return 0
	}
	Logger.Printf("ImageGetTexture-Loading... %s", name)
	filePath := filepath.Join(currentDir, name)
	data, err := engine.LoadBinaryFile(filePath)
	if err != nil {
		Logger.Printf("ImageGetTexture-LoadBinaryFile: %s", err)
		return 0
	}
	image := ImageLoadFromBytes(data, false)
	texture, err := engine.RenderInstance.TextureCreate(int(image.Width), int(image.Height), image.Pixels)
	if err != nil {
		Logger.Printf("ImageGetTexture: %s", err)
		return 0
	}

	return uint16(texture)
}

func ImageGetCompressedTexture(name string, render *engine.Render) (TextureList, error) {
	currentDir, _ := os.Getwd()
	filePath := filepath.Join(currentDir, name)
	cmp, err := imageLoadCompressed(filePath)
	if err != nil {
		return TextureList{}, err
	}
	tl := TextureList{
		start: render.TexturesLen(),
		len:   int(cmp.Len),
	}

	for i := 0; i < int(cmp.Len); i++ {
		image := ImageLoadFromBytes(cmp.Entries[i], false)
		render.TextureCreate(int(image.Width), int(image.Height), image.Pixels)
	}

	return tl, nil

}

// tim16BitToRGBA converts a 16-bit TIM pixel to RGBA
func tim16BitToRGBA(c uint16, transparentBit bool) engine.RGBA {
	r := byte((c >> 0) & 0x1f) << 3
	g := byte((c >> 5) & 0x1f) << 3
	b := byte((c >> 10) & 0x1f) << 3

	var a byte
	if c == 0 {
		a = 0x00
	} else if transparentBit && (c&0x7FFF) == 0 {
		a = 0x00
	} else {
		a = 0xFF
	}

	return engine.RGBA{R: r, G: g, B: b, A: a}
}


func TextureFromList(tl TextureList, index int) int {
	if index >= tl.len {
		Logger.Printf("texture %d not in list of len %d", index, tl.len)
	}

	return tl.start + index
}

type cmpT struct {
	Len     uint32
	Entries [][]byte
}

func imageLoadCompressed(name string) (*cmpT, error) {
	Logger.Printf("load cmp %s\n", name)

	// Load compressed bytes from the file
	compressedBytes, err := engine.LoadBinaryFile(name)
	if err != nil {
		return nil, err
	}

	decompressedBytes := lzss.Decompress(compressedBytes)

	// Initialize variables
	var p uint32
	var decompressedSize int32

	// Read the number of entries (Len) from data
	imageCount := engine.GetI32LE(compressedBytes, &p)

	for i := 0; i < int(imageCount); i++ {
		decompressedSize += engine.GetI32LE(compressedBytes, &p)
	}

	var cmp cmpT
	cmp.Len = uint32(imageCount)

	// Create a slice to hold bytes
	entries := make([][]byte, imageCount)

	p = 4
	var offset int32 = 0
	var end int32 = 0

	// Iterate through the entries and store their pointers
	for i := uint32(0); i < uint32(imageCount); i++ {
		end = engine.GetI32LE(compressedBytes, &p)
		entries[i] = decompressedBytes[offset : offset+end]
		offset += end
	}

	cmp.Entries = entries

	// Create and return the cmpT struct
	return &cmp, nil
}

