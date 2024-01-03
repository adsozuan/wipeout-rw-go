package wipeout

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/adsozuan/wipeout-rw-go/engine"
)

const (
	TimTypePaletted4BPP   = 0x08
	TimTypePaletted8BPP   = 0x09
	TimTypeTrueColor16BPP = 0x02
)

type Image struct {
	Width, Height uint32
	Pixels        []engine.RGBA
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

func ImageLoadFromBytes(bytes []uint8, transparent bool) *Image {
	var p uint32 = 0
	_ = engine.GetI32LE(bytes, &p) // magic
	typ := engine.GetI32LE(bytes, &p)
	var palette []uint16

	if typ == TimTypePaletted4BPP || typ == TimTypePaletted8BPP {
		_ = engine.GetI32LE(bytes, &p) // header size
		_ = engine.GetI16LE(bytes, &p) // paletteX
		_ = engine.GetI16LE(bytes, &p) // paletteY
		paletteColors := engine.GetI16LE(bytes, &p)
		_ = engine.GetI16LE(bytes, &p) // palettes
		palette = *(*[]uint16)(unsafe.Pointer(&bytes[p]))
		p += uint32(paletteColors * 2)
	}

	_ = engine.GetI32LE(bytes, &p) // data size

	pixelsPer16Bit := 1
	if typ == TimTypePaletted8BPP {
		pixelsPer16Bit = 2
	} else if typ == TimTypePaletted4BPP {
		pixelsPer16Bit = 4
	}

	_ = engine.GetI16LE(bytes, &p) // skipX
	_ = engine.GetI16LE(bytes, &p) // skipY
	entriesPerRow := engine.GetI16LE(bytes, &p)
	rows := engine.GetI16LE(bytes, &p)

	width := int32(entriesPerRow * int16(pixelsPer16Bit))
	height := int32(rows)
	entries := int32(entriesPerRow * rows)

	image := ImageAlloc(uint32(width), uint32(height))
	pixelPos := int32(0)

	if typ == TimTypeTrueColor16BPP {
		for i := int32(0); i < entries; i++ {
			image.Pixels[pixelPos] = tim16BitToRGBA(uint16(engine.GetI16LE(bytes, &p)), transparent)
			pixelPos++
		}
	} else if typ == TimTypePaletted8BPP {
		for i := int32(0); i < entries; i++ {
			palettePos := int32(engine.GetI16LE(bytes, &p))
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>0)&0xff], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>8)&0xff], transparent)
			pixelPos++
		}
	} else if typ == TimTypePaletted4BPP {
		for i := int32(0); i < entries; i++ {
			palettePos := int32(engine.GetI16LE(bytes, &p))
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>0)&0xf], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>4)&0xf], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>8)&0xf], transparent)
			pixelPos++
			image.Pixels[pixelPos] = tim16BitToRGBA(palette[(palettePos>>12)&0xf], transparent)
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

// tim16BitToRGBA converts a 16-bit TIM pixel to RGBA
func tim16BitToRGBA(value uint16, transparent bool) engine.RGBA {
	r := uint8((value >> 11) & 0x1F)
	g := uint8((value >> 5) & 0x3F)
	b := uint8(value & 0x1F)

	if transparent {
		a := uint8(0x1F)
		return engine.RGBA{R: r, G: g, B: b, A: a}
	}

	return engine.RGBA{R: r, G: g, B: b, A: 0xFF}
}