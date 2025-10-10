package buildhat

import (
	"fmt"
)

// ColorSensor creates a color sensor interface for the specified port
func (b *Brick) ColorSensor(port Port) *ColorSensor {
	return &ColorSensor{
		brick: b,
		port:  port.Int(),
	}
}

// ColorSensor provides a Python-like color sensor interface
type ColorSensor struct {
	brick *Brick
	port  int
}

// GetColor gets the current color reading as RGBA
func (s *ColorSensor) GetColor() (Color, error) {
	// Set to color RGB mode (mode 5 - RGBI)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(5))); err != nil {
		return Color{}, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return Color{}, err
	}

	if len(data) < 4 {
		return Color{}, fmt.Errorf("insufficient color data received")
	}

	// Convert from raw values (0-1024) to 0-255
	r, _ := data[0].(int)
	g, _ := data[1].(int)
	b, _ := data[2].(int)
	i, _ := data[3].(int)

	return Color{
		R: clamp8((r * 255) / 1024),
		G: clamp8((g * 255) / 1024),
		B: clamp8((b * 255) / 1024),
		A: clamp8((i * 255) / 1024),
	}, nil
}

// GetReflectedLight gets the reflected light reading (0-100%)
func (s *ColorSensor) GetReflectedLight() (int, error) {
	// Set to reflected light mode (mode 1)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(1))); err != nil {
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
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(2))); err != nil {
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
