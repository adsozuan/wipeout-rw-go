package engine



func GetU8(bytes []byte, p *uint32) byte {
    v := bytes[*p]
    *p++
    return v
}



func GetU16LE(bytes []byte, p *uint32) uint16 {
    v := uint16(bytes[*p]) | uint16(bytes[*p+1])<<8
    *p += 2
    return v
}

func GetU32LE(bytes []byte, p *uint32) uint32 {
    v := uint32(bytes[*p]) | uint32(bytes[*p+1])<<8 | uint32(bytes[*p+2])<<16 | uint32(bytes[*p+3])<<24
    *p += 4
    return v
}


func GetU16(bytes []byte, p *uint32) uint16 {
	v := uint16(bytes[*p]) << 8
	*p++
	v |= uint16(bytes[*p]) << 0
	*p++
	return v
}

func GetU32(bytes []byte, p *uint32) uint32 {
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

func GetI8(bytes []byte, p *uint32) int8 {
	return int8(GetU8(bytes, p))
}

func GetI16(bytes []byte, p *uint32) int16 {
	return int16(GetU16(bytes, p))
}

func GetI16LE(bytes []byte, p *uint32) int16 {
	return int16(GetU16LE(bytes, p))
}

func GetI32(bytes []byte, p *uint32) int32 {
	return int32(GetU32(bytes, p))
}

func GetI32LE(bytes []byte, p *uint32) int32 {
	return int32(GetU32LE(bytes, p))
}

