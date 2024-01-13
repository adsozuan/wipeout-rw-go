package engine

func Scale(v, inMin, inMax, outMin, outMax float64) float64 {
	return outMin + (outMax-outMin)*((v-inMin)/(inMax-inMin))
}

func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}


// Clamp limits the value between min and max
func Clamp[T Number](value, min, max T) T {
	if value < min {
		return (min)
	}
	if value > max {
		return (max)
	}
	return (value)
}
