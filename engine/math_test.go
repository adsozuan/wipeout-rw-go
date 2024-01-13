package engine

import "testing"

func TestScale(t *testing.T) {
	tests := []struct {
		v, inMin, inMax, outMin, outMax, want float64
	}{
		{10, 0, 100, 0, 1, 0.1},
		// Add more test cases here
	}

	for _, tt := range tests {
		got := Scale(tt.v, tt.inMin, tt.inMax, tt.outMin, tt.outMax)
		if got != tt.want {
			t.Errorf("Scale(%v, %v, %v, %v, %v) = %v; want %v", tt.v, tt.inMin, tt.inMax, tt.outMin, tt.outMax, got, tt.want)
		}
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		a, b, t, want float64
	}{
		{0, 10, 0.5, 5},
		// Add more test cases here
	}

	for _, tt := range tests {
		got := Lerp(tt.a, tt.b, tt.t)
		if got != tt.want {
			t.Errorf("Lerp(%v, %v, %v) = %v; want %v", tt.a, tt.b, tt.t, got, tt.want)
		}
	}
}