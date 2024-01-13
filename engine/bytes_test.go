package engine

import "testing"


func TestGetU8(t *testing.T) {
    tests := []struct {
        bytes []byte
        pos   uint32
        want  byte
    }{
        {[]byte{0x01, 0x02, 0x03}, 0, 0x01},
        {[]byte{0x01, 0x02, 0x03}, 1, 0x02},
        // More test cases, especially edge cases
    }

    for _, tt := range tests {
        pos := tt.pos // Copy position
        got := GetU8(tt.bytes, &pos)
        if got != tt.want || pos != tt.pos+1 {
            t.Errorf("GetU8(%v, %v) = %v, pos %v; want %v, pos %v", tt.bytes, tt.pos, got, pos, tt.want, tt.pos+1)
        }
    }
}