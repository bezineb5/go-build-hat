package buildhat

import (
	"fmt"
)

// ColorDistanceSensor creates a color distance sensor interface for the specified port
func (b *Brick) ColorDistanceSensor(port Port) *ColorDistanceSensor {
	return &ColorDistanceSensor{
		brick: b,
		port:  port,
	}
}

// ColorDistanceSensor provides a Python-like color distance sensor interface
type ColorDistanceSensor struct {
	brick *Brick
	port  Port
}

// GetColor gets the current color reading as RGBA
func (s *ColorDistanceSensor) GetColor() (Color, error) {
	// Set to color mode (mode 0)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(0))); err != nil {
		return Color{}, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port.Int())
	if err != nil {
		return Color{}, err
	}

	if len(data) < 4 {
		return Color{}, fmt.Errorf("insufficient color data received")
	}

	// Convert from raw values to RGB
	r, _ := data[0].(int)
	g, _ := data[1].(int)
	b, _ := data[2].(int)
	a, _ := data[3].(int)

	return Color{
		R: clamp8(r),
		G: clamp8(g),
		B: clamp8(b),
		A: clamp8(a),
	}, nil
}

// GetDistance gets the distance reading
func (s *ColorDistanceSensor) GetDistance() (int, error) {
	// Set to distance mode (mode 1)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(1))); err != nil {
		return 0, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port.Int())
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no distance data received")
	}

	// Distance is the first value
	if distance, ok := data[0].(int); ok {
		return distance, nil
	}

	return 0, fmt.Errorf("invalid distance data type")
}

// GetReflectedLight gets the reflected light reading
func (s *ColorDistanceSensor) GetReflectedLight() (int, error) {
	// Set to reflected light mode (mode 2)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(2))); err != nil {
		return 0, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port.Int())
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
