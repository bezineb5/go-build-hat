package buildhat

import (
	"fmt"
)

// TiltSensor creates a tilt sensor interface for the specified port
func (b *Brick) TiltSensor(port string) *TiltSensor {
	portNum := int(port[0] - 'A')
	return &TiltSensor{
		brick: b,
		port:  portNum,
	}
}

// TiltSensor provides a Python-like tilt sensor interface (WeDo sensor)
type TiltSensor struct {
	brick *Brick
	port  int
}

// GetTilt gets the current tilt reading as X, Y, Z coordinates
func (s *TiltSensor) GetTilt() (struct{ X, Y, Z int }, error) {
	// Set to tilt mode (mode 0)
	cmd := fmt.Sprintf("port %d ; select 0", s.port)
	if err := s.brick.writeCommand(cmd); err != nil {
		return struct{ X, Y, Z int }{}, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return struct{ X, Y, Z int }{}, err
	}

	// Tilt data can be 1 value (direction) or 3 values (X, Y, Z)
	if len(data) == 1 {
		// Single value - direction code
		// Return as a simple representation
		direction, _ := data[0].(int)
		return struct{ X, Y, Z int }{X: direction, Y: 0, Z: 0}, nil
	}

	if len(data) < 3 {
		return struct{ X, Y, Z int }{}, fmt.Errorf("insufficient tilt data received")
	}

	x, _ := data[0].(int)
	y, _ := data[1].(int)
	z, _ := data[2].(int)

	return struct{ X, Y, Z int }{X: x, Y: y, Z: z}, nil
}

// GetDirection returns the tilt direction as a string
func (s *TiltSensor) GetDirection() (string, error) {
	tilt, err := s.GetTilt()
	if err != nil {
		return "", err
	}

	// Simple direction mapping based on tilt values
	// This is a simplified version - actual implementation may vary
	switch {
	case tilt.X > 45:
		return "right", nil
	case tilt.X < -45:
		return "left", nil
	case tilt.Y > 45:
		return "forward", nil
	case tilt.Y < -45:
		return "backward", nil
	default:
		return "level", nil
	}
}
