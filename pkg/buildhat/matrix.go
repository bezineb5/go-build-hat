package buildhat

import (
	"fmt"
)

// Matrix creates a matrix interface for the specified port
func (b *Brick) Matrix(port string) *Matrix {
	portNum := int(port[0] - 'A')
	return &Matrix{
		brick:  b,
		port:   portNum,
		pixels: [3][3]Pixel{},
	}
}

// Pixel represents a single LED pixel with color and brightness
type Pixel struct {
	Color      int // 0-10
	Brightness int // 0-10
}

// Matrix provides a Python-like LED matrix interface
type Matrix struct {
	brick  *Brick
	port   int
	pixels [3][3]Pixel
}

// SetPixel sets a single pixel at position (x, y)
func (m *Matrix) SetPixel(x, y, color, brightness int) error {
	if x < 0 || x > 2 || y < 0 || y > 2 {
		return fmt.Errorf("pixel coordinates must be 0-2")
	}
	if color < 0 || color > 10 {
		return fmt.Errorf("color must be 0-10")
	}
	if brightness < 0 || brightness > 10 {
		return fmt.Errorf("brightness must be 0-10")
	}

	m.pixels[x][y] = Pixel{Color: color, Brightness: brightness}
	return m.display()
}

// SetAll sets all pixels to the same color and brightness
func (m *Matrix) SetAll(color, brightness int) error {
	if color < 0 || color > 10 {
		return fmt.Errorf("color must be 0-10")
	}
	if brightness < 0 || brightness > 10 {
		return fmt.Errorf("brightness must be 0-10")
	}

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			m.pixels[x][y] = Pixel{Color: color, Brightness: brightness}
		}
	}
	return m.display()
}

// SetRow sets an entire row to the same color and brightness
func (m *Matrix) SetRow(row, color, brightness int) error {
	if row < 0 || row > 2 {
		return fmt.Errorf("row must be 0-2")
	}
	if color < 0 || color > 10 {
		return fmt.Errorf("color must be 0-10")
	}
	if brightness < 0 || brightness > 10 {
		return fmt.Errorf("brightness must be 0-10")
	}

	for y := 0; y < 3; y++ {
		m.pixels[row][y] = Pixel{Color: color, Brightness: brightness}
	}
	return m.display()
}

// SetColumn sets an entire column to the same color and brightness
func (m *Matrix) SetColumn(col, color, brightness int) error {
	if col < 0 || col > 2 {
		return fmt.Errorf("column must be 0-2")
	}
	if color < 0 || color > 10 {
		return fmt.Errorf("color must be 0-10")
	}
	if brightness < 0 || brightness > 10 {
		return fmt.Errorf("brightness must be 0-10")
	}

	for x := 0; x < 3; x++ {
		m.pixels[x][col] = Pixel{Color: color, Brightness: brightness}
	}
	return m.display()
}

// Clear turns off all pixels
func (m *Matrix) Clear() error {
	return m.SetAll(0, 0)
}

// display sends the current pixel data to the matrix
func (m *Matrix) display() error {
	// Build the data packet
	// Format: 0xc2 followed by 9 bytes, each containing brightness (high nibble) and color (low nibble)
	data := make([]byte, 10)
	data[0] = 0xc2

	idx := 1
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			// Pack brightness and color into single byte
			data[idx] = byte((m.pixels[x][y].Brightness << 4) | m.pixels[x][y].Color)
			idx++
		}
	}

	// Convert to hex string for write1 command
	hexStr := ""
	for i, b := range data {
		if i > 0 {
			hexStr += " "
		}
		hexStr += fmt.Sprintf("%x", b)
	}

	// Send using write1 command
	cmd := fmt.Sprintf("port %d ; write1 %s", m.port, hexStr)
	return m.brick.writeCommand(cmd)
}
