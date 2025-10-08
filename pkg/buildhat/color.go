package buildhat

// Color represents an RGBA color value
type Color struct {
	R uint8 // Red component (0-255)
	G uint8 // Green component (0-255)
	B uint8 // Blue component (0-255)
	A uint8 // Alpha/Intensity component (0-255)
}

// clamp8 ensures a value is within 0-255 range
func clamp8(val int) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}
