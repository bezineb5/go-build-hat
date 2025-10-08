package buildhat

import (
	"fmt"
)

// ColorSensor creates a color sensor interface for the specified port
func (b *Brick) ColorSensor(port string) *ColorSensor {
	portNum := int(port[0] - 'A')
	return &ColorSensor{
		brick: b,
		port:  portNum,
	}
}

// ColorSensor provides a Python-like color sensor interface
type ColorSensor struct {
	brick *Brick
	port  int
}

// GetColor gets the current color reading as RGBA
func (s *ColorSensor) GetColor() (struct{ R, G, B, A uint8 }, error) {
	// Set to color RGB mode (mode 5 - RGBI)
	cmd := fmt.Sprintf("port %d ; select 5", s.port)
	if err := s.brick.writeCommand(cmd); err != nil {
		return struct{ R, G, B, A uint8 }{}, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return struct{ R, G, B, A uint8 }{}, err
	}

	if len(data) < 4 {
		return struct{ R, G, B, A uint8 }{}, fmt.Errorf("insufficient color data received")
	}

	// Convert from raw values (0-1024) to 0-255
	r, _ := data[0].(int)
	g, _ := data[1].(int)
	b, _ := data[2].(int)
	i, _ := data[3].(int)

	return struct{ R, G, B, A uint8 }{
		R: clamp8((r * 255) / 1024),
		G: clamp8((g * 255) / 1024),
		B: clamp8((b * 255) / 1024),
		A: clamp8((i * 255) / 1024),
	}, nil
}

// GetReflectedLight gets the reflected light reading (0-100%)
func (s *ColorSensor) GetReflectedLight() (int, error) {
	// Set to reflected light mode (mode 1)
	cmd := fmt.Sprintf("port %d ; select 1", s.port)
	if err := s.brick.writeCommand(cmd); err != nil {
		return 0, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no reflected light data received")
	}

	// Reflected light is the first value
	if light, ok := data[0].(int); ok {
		return light, nil
	}

	return 0, fmt.Errorf("invalid reflected light data type")
}

// GetAmbientLight gets the ambient light reading (0-100%)
func (s *ColorSensor) GetAmbientLight() (int, error) {
	// Set to ambient light mode (mode 2)
	cmd := fmt.Sprintf("port %d ; select 2", s.port)
	if err := s.brick.writeCommand(cmd); err != nil {
		return 0, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no ambient light data received")
	}

	// Ambient light is the first value
	if light, ok := data[0].(int); ok {
		return light, nil
	}

	return 0, fmt.Errorf("invalid ambient light data type")
}
