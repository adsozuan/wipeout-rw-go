package wipeout

func Scale(v, inMin, inMax, outMin, outMax float64) float64 {
	return outMin + (outMax-outMin)*((v-inMin)/(inMax-inMin))
}

func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func GetU8(bytes []uint8, p *uint32) uint8 {
	v := bytes[*p]
	*p++
	return v
}

func GetU16(bytes []uint8, p *uint32) uint16 {
	v := uint16(bytes[*p]) << 8
	*p++
	v |= uint16(bytes[*p]) << 0
	*p++
	return v
}

func GetU32(bytes []uint8, p *uint32) uint32 {
	v := uint32(bytes[*p]) << 24
	*p++
	v |= uint32(bytes[*p]) << 16
	*p++
	v |= uint32(bytes[*p]) << 8
	*p++
	v |= uint32(bytes[*p]) << 0
	*p++
	return v
}

func GetU16LE(bytes []uint8, p *uint32) uint16 {
	v := uint16(bytes[*p]) << 0
	*p++
	v |= uint16(bytes[*p]) << 8
	*p++
	return v
}

func GetU32LE(bytes []uint8, p *uint32) uint32 {
	v := uint32(bytes[*p]) << 0
	*p++
	v |= uint32(bytes[*p]) << 8
	*p++
	v |= uint32(bytes[*p]) << 16
	*p++
	v |= uint32(bytes[*p]) << 24
	*p++
	return v
}

func GetI8(bytes []uint8, p *uint32) int8 {
	return int8(GetU8(bytes, p))
}

func GetI16(bytes []uint8, p *uint32) int16 {
	return int16(GetU16(bytes, p))
}

func GetI16LE(bytes []uint8, p *uint32) int16 {
	return int16(GetU16LE(bytes, p))
}

func GetI32(bytes []uint8, p *uint32) int32 {
	return int32(GetU32(bytes, p))
}

func GetI32LE(bytes []uint8, p *uint32) int32 {
	return int32(GetU32LE(bytes, p))
}
